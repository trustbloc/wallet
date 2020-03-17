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
        let startupOpts = await ariesStartupOpts()
        console.log("aries startup options ", JSON.stringify(startupOpts))

        await new Aries.Framework(startupOpts).then(resp => {
            Vue.prototype.$aries = resp
            console.log('aries started successfully')
        }).catch(err => {
            console.error('error starting aries framework : errMsg=', err)
        })
    }

    return Vue.prototype.$aries
}


let defaultAriesStartupOpts = {
    assetsPath: '/aries-framework-go/assets',
    'outbound-transport': ['ws', 'http'],
    'transport-return-route': 'all',
    'http-resolver-url': [],
    'agent-default-label': 'demo-user-agent',
    'auto-accept': true,
    'log-level': 'debug',
    'db-namespace': 'agent'
}

async function ariesStartupOpts() {
    let startupOpts = {}

    if (process.env.NODE_ENV === "production") {
        const axios = require('axios').default;

        // call service to get the opts
        await axios.get(window.location.origin + '/aries/jsopts')
            .then(resp => {
                startupOpts = resp.data
                console.log("successfully fetched start up options: resp=" + JSON.stringify(startupOpts));
            })
            .catch(err => {
                console.log("error fetching start up options - using default options : errMsg=", err);
            })
    }

    return {
        assetsPath: defaultAriesStartupOpts['assetsPath'],
        'outbound-transport': defaultAriesStartupOpts['outbound-transport'],
        'transport-return-route': defaultAriesStartupOpts['transport-return-route'],
        'http-resolver-url': ('http-resolver-url' in startupOpts) ? startupOpts['http-resolver-url'] : defaultAriesStartupOpts['http-resolver-url'],
        'agent-default-label': ('agent-default-label' in startupOpts) ? startupOpts['agent-default-label'] : defaultAriesStartupOpts['agent-default-label'],
        'auto-accept': ('auto-accept' in startupOpts) ? startupOpts['auto-accept'] : defaultAriesStartupOpts['auto-accept'],
        'log-level': ('log-level' in startupOpts) ? startupOpts['log-level'] : defaultAriesStartupOpts['log-level'],
        'db-namespace': ('db-namespace' in startupOpts) ? startupOpts['db-namespace'] : defaultAriesStartupOpts['db-namespace']
    }
}

Vue.prototype.$arieslib = loadAriesOnce()

new Vue({
    router,
    render: h => h(App),
}).$mount('#app')

