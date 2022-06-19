/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createRouter, createWebHistory } from 'vue-router';
import store from '@/store';
import { getGnapKeyPair, gnapRequestAccess } from '@/mixins';
import routes from './routes';
import { computed } from 'vue';

const router = createRouter({
  history: createWebHistory(__webpack_public_path__),
  routes,
});

router.beforeEach(async (to, from, next) => {
  store.dispatch('agent/flushStore');
  if (to.path === '/gnap') {
    const gnapAccessTokenConfig = computed(() => store.getters['getGnapAccessTokenConfig']);
    const gnapAccessTokens = await gnapAccessTokenConfig.value;
    const gnapAuthServerURL = computed(() => store.getters['hubAuthURL']).value;
    const gnapKeyPair = await getGnapKeyPair();
    const signer = { SignatureVal: gnapKeyPair };
    const resp = await gnapRequestAccess(signer, gnapAccessTokens, gnapAuthServerURL);
    // TODO Issue-1699 Save properties from resp to Vuex/local storage and call GNAPClient.continue() in separate func
  }
  if (to.path === 'gnap/redirect') {
    // TODO Issue-1701 HTTP Response GNAP grant response with interact-redirect (to display list of OIDC providers)
  }
  const locale = store.getters.getLocale;
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (store.getters.getCurrentUser) {
      next();
      return;
    } else if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      next();
      return;
    } else {
      const { signin, disableCHAPI } = to.meta;
      router.replace({
        name: signin ? 'signin' : 'signup',
        params: {
          ...router.currentRoute._value.params,
          locale: locale.base,
          redirect: to.name,
          disableCHAPI,
        },
        query: to.query,
      });
      next();
      return;
    }
  } else if (to.matched.some((record) => record.meta.blockNoAuth)) {
    if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      next();
      return;
    } else {
      router.replace({
        name: 'block-no-auth',
        params: {
          ...router.currentRoute._value.params,
          locale: locale.base,
          redirect: { name: 'signup' },
        },
      });
      next();
      return;
    }
  } else {
    if (to.params.locale && to.params.locale !== locale.id) {
      router.replace({
        name: to.params.name,
        params: {
          ...router.currentRoute._value.params,
          ...to.params,
          locale: locale.base,
        },
        query: to.query,
      });
      next();
      return;
    } else {
      next();
      return;
    }
  }
});

export default router;
