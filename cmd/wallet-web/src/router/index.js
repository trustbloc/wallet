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
import i18n, { updateI18nLocale } from '@/plugins/i18n';

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
      next({
        path: `${locale.base}/signup`,
        params: { redirect: to.name },
      });
      return;
    }
  } else if (to.matched.some((record) => record.meta.blockNoAuth)) {
    console.log('blockNoAuth');
    if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      next();
      return;
    } else {
      console.log('other options ');
      next({
        name: 'block-no-auth',
        params: { loginURL: `${locale.base}/signup` },
      });
      return;
    }
  } else {
    if (locale.id !== router.history.pending.params.locale) {
      next({
        path: `${locale.base}${router.history._startLocation}`,
      });
      return;
    } else {
      next();
      return;
    }
  }
});

export default router;
