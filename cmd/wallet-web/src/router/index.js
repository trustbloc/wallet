/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import VueRouter from 'vue-router';
import Root from './Root';
import store from '@/store';
import routes from './routes';
import supportedLocales from '@/config/supportedLocales';

Vue.use(VueRouter);

// Creates regex (en|fr)
function getLocaleRegex() {
  let reg = '';
  supportedLocales.forEach((locale, index) => {
    reg = `${reg}${locale.id}${index !== supportedLocales.length - 1 ? '|' : ''}`;
  });
  return `(${reg})`;
}

const router = new VueRouter({
  mode: 'history',
  base: `${__webpack_public_path__}`,
  routes: [
    {
      path: `/:locale${getLocaleRegex()}?`,
      component: Root,
      children: routes,
    },
  ],
});

router.beforeEach((to, from, next) => {
  store.dispatch('agent/flushStore');
  const locale = store.getters.getLocale;
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (store.getters.getCurrentUser) {
      next();
      return;
    } else if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      next();
      return;
    } else {
      router.replace({
        name: 'signup',
        params: {
          ...router.history.current.params,
          locale: locale.base,
          redirect: { name: to.name },
          query: to.query,
        },
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
          ...router.history.current.params,
          locale: locale.base,
          redirect: { name: 'signup' },
        },
      });
      next();
      return;
    }
  } else {
    if (locale.id !== router.history.pending.params.locale) {
      router.replace({
        name: router.history.pending.name,
        params: {
          ...router.history.pending.params,
          locale: locale.base,
        },
        query: router.history.pending.query,
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
