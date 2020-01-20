/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue'
import App from './App.vue'
import router from "./router";
import invitation from "@trustbloc/invitation";
import * as polyfill from "credential-handler-polyfill";
import * as webCredentialHandler from "web-credential-handler";

Vue.config.productionTip = false
Vue.use(invitation);

Vue.prototype.$polyfill = polyfill
Vue.prototype.$webCredentialHandler = webCredentialHandler


new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
