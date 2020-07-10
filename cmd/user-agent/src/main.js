/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import App from './App.vue';
import VueRouter from "vue-router";
import routes from "./router/index";
import * as polyfill from "credential-handler-polyfill";
import * as webCredentialHandler from "web-credential-handler";
import * as Aries from "@trustbloc-cicd/aries-framework-go"
import * as trustblocAgent from "@trustbloc/trustbloc-agent"
import MaterialDashboard from "./material-dashboard";


Vue.config.productionTip = false

Vue.prototype.$polyfill = polyfill
Vue.prototype.$webCredentialHandler = webCredentialHandler
Vue.prototype.$trustblocAgent = trustblocAgent


async function loadAriesOnce() {
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

let defaultTrustBlocStartupOpts = {
    assetsPath: '/trustbloc-agent/assets',
    blocDomain: 'testnet.trustbloc.local',
    walletMediatorURL: '',
    'log-level': 'debug'
}

async function trustblocStartupOpts() {
    let startupOpts = {}
    if (process.env.NODE_ENV === "production") {
        const axios = require('axios').default;

        // call service to get the opts
        await axios.get(window.location.origin + '/trustbloc-agent/jsopts')
            .then(resp => {
                startupOpts = resp.data
                console.log("successfully fetched start up options: resp=" + JSON.stringify(startupOpts));
            })
            .catch(err => {
                console.log("error fetching start up options - using default options : errMsg=", err);
            })
    }

    return  {
        assetsPath: defaultTrustBlocStartupOpts['assetsPath'],
        blocDomain: ('blocDomain' in startupOpts) ? startupOpts['blocDomain'] : defaultTrustBlocStartupOpts['blocDomain'],
        walletMediatorURL: ('walletMediatorURL' in startupOpts) ? startupOpts['walletMediatorURL'] : defaultTrustBlocStartupOpts['walletMediatorURL'],
        'log-level': ('log-level' in startupOpts) ? startupOpts['log-level'] : defaultTrustBlocStartupOpts['log-level']
    }
}

Vue.prototype.$arieslib = loadAriesOnce()
Vue.prototype.$trustblocStartupOpts = trustblocStartupOpts()


// configure router
const router = new VueRouter({
    mode:'history',
    routes, // short for routes: routes
    linkExactActiveClass: "nav-item active"
});

Vue.use(VueRouter);
Vue.use(MaterialDashboard);

new Vue({
    el: "#app",
    render: h => h(App),
    router
});

