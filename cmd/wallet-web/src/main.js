/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import store from './store';
import App from './App.vue';
import VueRouter from 'vue-router';
import router from './router/index';
import * as polyfill from 'credential-handler-polyfill';
import * as webCredentialHandler from 'web-credential-handler';
import { mapActions, mapGetters } from 'vuex';
import VueCookies from 'vue-cookies';
import '@/assets/css/tailwind.css';
import i18n from './plugins/i18n';
import supportedLocales from '@/config/supportedLocales';
import { updateI18nLocale } from '@/plugins/i18n';
import getBrowserLocale from '@/utils/i18n/getBrowserLocale';

Vue.config.productionTip = false;

Vue.prototype.$polyfill = polyfill;
Vue.prototype.$webCredentialHandler = webCredentialHandler;

Vue.use(VueRouter);
Vue.use(VueCookies);
Vue.$cookies.config('7d');
Vue.component(
  'ToastNotification',
  require('./components/ToastNotification/ToastNotification').default
);

// Returns locale configuration. By default, try VUE_APP_I18N_LOCALE. As fallback, use 'en'.
function getMappedLocale(locale = process.env.VUE_APP_I18N_LOCALE || 'en') {
  return supportedLocales.find((loc) => loc.id === locale);
}

function getStartingLocale() {
  // Get locale parameter form the URL
  const localeUrlParam = router.history._startLocation
    .replaceAll(/^\//gi, '')
    .replace(/\/.*$/gi, '');
  // If locale parameter is set, check if it is amongst the supported locales and return it.
  if (localeUrlParam && supportedLocales.find((loc) => loc.id === localeUrlParam)) {
    return getMappedLocale(localeUrlParam);
  }
  // If no locale parameter is set in the URL, use the browser default.
  else {
    const browserLocale = getBrowserLocale({ countryCodeOnly: true });
    return getMappedLocale(browserLocale);
  }
}

// Get starting locale, set it in i18n and in the store
const startingLocale = getStartingLocale();
store.dispatch('setLocale', startingLocale);

new Vue({
  store,

  data: () => ({
    loaded: false,
  }),

  beforeCreate: async function () {
    await updateI18nLocale(startingLocale.id);
  },

  mounted: async function () {
    // load opts
    await this.initOpts();

    // load user if already logged in
    this.loadUser();

    // load agent if user already logged in and agent not initialized (scenario: page refresh)
    if (store.getters.getCurrentUser && !this.isAgentInitialized()) {
      await this.initAgent();
    }

    // removes spinner
    this.loaded = true;
  },

  methods: {
    ...mapActions(['initOpts', 'loadUser']),
    ...mapActions('agent', { initAgent: 'init' }),
    ...mapGetters('agent', { isAgentInitialized: 'isInitialized' }),
    ...mapGetters(['getAgentOpts']),
  },
  i18n,
  router,
  render: (h) => h(App),
}).$mount('#app');

window.onbeforeunload = function () {
  if (store) {
    store.dispatch('agent/flushStore');
  }
};
