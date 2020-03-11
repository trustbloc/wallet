/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue'
import App from './App.vue'
import router from "./router";
import * as polyfill from "credential-handler-polyfill";
import * as webCredentialHandler from "web-credential-handler";
import * as Aries from "@hyperledger/aries-framework-go"

Vue.config.productionTip = false

Vue.prototype.$polyfill = polyfill
Vue.prototype.$webCredentialHandler = webCredentialHandler

async function loadAriesOnce() {
    import('@hyperledger/aries-framework-go/dist/web/aries.js')
        .catch(err => console.log('error importing aries library : errMsg=' + err.message))

    if (!Vue.prototype.$aries) {
        // TODO start up option should be fetched from environment
        await new Aries.Framework({
            assetsPath: "/aries-framework-go/assets",
            "agent-default-label": "dem-js-agent",
            "http-resolver-url": [],
            "auto-accept": true,
            "outbound-transport": ["ws", "http"],
            "transport-return-route": "all",
            "log-level": "debug",
            "db-namespace": "agent"
        }).then(resp => {
            Vue.prototype.$aries = resp
            console.log('aries started successfully')
        }).catch(err => {
            console.error('error starting aries framework : errMsg=', err)
        })
    }

    return Vue.prototype.$aries
}

Vue.prototype.$arieslib = loadAriesOnce()

new Vue({
    router,
    render: h => h(App),
}).$mount('#app')
