/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createApp } from 'vue';
import store from './store';
import App from './App.vue';
import router from './router/index';
import * as polyfill from 'credential-handler-polyfill';
import * as webCredentialHandler from 'web-credential-handler';
import { mapActions, mapGetters } from 'vuex';
import '@/assets/css/tailwind.css';
import i18n from './plugins/i18n';
import getStartingLocale from '@/utils/i18n/getStartingLocale.js';
import { updateI18nLocale } from '@/plugins/i18n';

// Get starting locale, set it in i18n and in the store
const startingLocale = getStartingLocale();
store.dispatch('setLocale', startingLocale);

const app = createApp({
  ...App,
  data: () => ({
    loaded: false,
  }),

  beforeCreate: async function () {
    await updateI18nLocale(startingLocale.id);
  },

  mounted: async function () {
    try {
      // load opts
      await this.initOpts();

      // load user if already logged in
      this.loadUser();

      // load agent if user already logged in and agent not initialized (scenario: page refresh)
      if (store.getters.getCurrentUser && !this.isAgentInitialized()) {
        await this.initAgent();
      }
    } catch (e) {
      console.log('main.js error:', e);
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
});

app.config.globalProperties.$polyfill = polyfill;
app.config.globalProperties.$webCredentialHandler = webCredentialHandler;
app.use(router);
app.use(store);
app.use(i18n);
app.component(
  'ToastNotification',
  require('./components/ToastNotification/ToastNotification').default
);

window.onbeforeunload = function () {
  if (store) {
    store.dispatch('agent/flushStore');
  }
};

app.mount('#app');
