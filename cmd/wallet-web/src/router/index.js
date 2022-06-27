/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createRouter, createWebHistory } from 'vue-router';
import store from '@/store';
import { getGnapKeyPair, gnapContinue, gnapRequestAccess } from '@/mixins';
import routes from './routes';
import { computed } from 'vue';
import { SHA3 } from 'sha3';

const router = createRouter({
  history: createWebHistory(__webpack_public_path__),
  routes,
});

router.beforeEach(async (to, from, next) => {
  store.dispatch('agent/flushStore');
  if (to.path === '/gnap') {
    const gnapAccessTokenConfig = computed(() => store.getters['getGnapAccessTokenConfig']);
    const gnapAccessTokens = await gnapAccessTokenConfig.value;
    const gnapAuthServerURL = computed(() => store.getters['hubAuthURL']);
    const walletWebUrl = computed(() => store.getters['walletWebUrl']);
    const gnapKeyPair = await getGnapKeyPair();
    const signer = { SignatureVal: gnapKeyPair };
    const clientNonceVal = (Math.random() + 1).toString(36).substring(7);
    const resp = await gnapRequestAccess(
      signer,
      gnapAccessTokens,
      gnapAuthServerURL.value,
      walletWebUrl.value,
      clientNonceVal
    );
    // If user have already logged in then just redirect
    if (resp.data.access_token || false) {
      store.dispatch('updateSessionToken', resp.data.access_token);
      router.push({ name: 'vaults' });
      next();
      return;
    }
    const respMetaData = {
      uri: resp.data.continue.uri,
      continue_access_token: resp.data.continue.access_token,
      finish: resp.data.interact.finish,
      clientNonceVal: clientNonceVal,
    };
    store.dispatch('updateGnapReqAccessResp', respMetaData);
    window.location.href = resp.data.interact.redirect;
  }
  if (to.path === '/gnap/redirect') {
    const gnapResp = store.getters.getGnapReqAccessResp;
    const params = new URL(document.location).searchParams;
    const hash = params.get('hash');
    const interactRef = params.get('interact_ref');
    const data = gnapResp.clientNonceVal + '\n' + gnapResp.finish + '\n' + interactRef + '\n';

    const shaHash = new SHA3(512);
    shaHash.update(data);
    let hashB64 = shaHash.digest({ format: 'base64' });
    hashB64 = hashB64.replace(/\+/g, '-').replace(/\//g, '_').replace(/\=+$/, '');
    if (hash === hashB64) {
      const gnapAuthServerURL = computed(() => store.getters['hubAuthURL']);
      const gnapKeyPair = await getGnapKeyPair();
      const signer = { SignatureVal: gnapKeyPair };
      const gnapContinueResp = await gnapContinue(
        signer,
        gnapAuthServerURL.value,
        interactRef,
        gnapResp.continue_access_token.value
      );
      const accessToken = gnapContinueResp.data.access_token[0].value;
      const subjectId = gnapContinueResp.data.subject.sub_ids[0].id;
      store.dispatch('updateSessionToken', gnapContinueResp.data.access_token);
      store.dispatch('agent/init', { accessToken, subjectId });
    }
    // TODO Issue-1744 Fetch user data to continue the wallet dashboard flow - Integrate with agent sdk
    window.top.close();
    router.push({ name: 'vaults' });
    next();
    return;
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
