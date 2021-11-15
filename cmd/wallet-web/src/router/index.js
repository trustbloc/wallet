/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createRouter, createWebHistory } from 'vue-router';
import Root from './Root';
import store from '@/store';
import routes from './routes';
import supportedLocales from '@/config/supportedLocales';

// Creates regex (en|fr)
function getLocaleRegex() {
  let reg = '';
  supportedLocales.forEach((locale, index) => {
    reg = `${reg}${locale.id}${index !== supportedLocales.length - 1 ? '|' : ''}`;
  });
  return `(${reg})`;
}

const router = createRouter({
  history: createWebHistory(__webpack_public_path__),
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
  console.log('to ', to, 'to meta --> ', to.meta);
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
        name: to.meta.signin ? 'signin' : 'signup',
        params: {
          ...router.currentRoute._value.params,
          locale: locale.base,
          redirect: to.name,
        },
        meta: to.meta,
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
