/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue'
import App from './App.vue'
import router from "./router";
import invitation from "@trustbloc/invitation";

Vue.config.productionTip = false
Vue.use(invitation);

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
