/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import store from './store'
import App from './App.vue';
import VueRouter from "vue-router";
import routes from "./router/index";
import * as polyfill from "credential-handler-polyfill";
import * as webCredentialHandler from "web-credential-handler";
import * as Aries from "@trustbloc-cicd/aries-framework-go"
import * as trustblocAgent from "@trustbloc/trustbloc-agent"
import MaterialDashboard from "./material-dashboard";
import {mapActions} from "vuex";

Vue.config.productionTip = false

Vue.prototype.$polyfill = polyfill
Vue.prototype.$webCredentialHandler = webCredentialHandler
Vue.prototype.$trustblocAgent = trustblocAgent

let defaultAriesStartupOpts = {
    assetsPath: '/aries-framework-go/assets',
    'outbound-transport': ['ws', 'http'],
    'transport-return-route': 'all',
    'http-resolver-url': ["trustbloc:testnet.trustbloc.local@http://localhost:8080/1.0/identifiers", "web@http://localhost:8080/1.0/identifiers"],

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
    'log-level': 'debug',
    walletMediatorURL: 'https://localhost:10063',
    credentialMediatorURL: '',
    agentUsername: '',
    sdsServerURL: ''
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

    return {
        assetsPath: defaultTrustBlocStartupOpts['assetsPath'],
        blocDomain: ('blocDomain' in startupOpts) ? startupOpts['blocDomain'] : defaultTrustBlocStartupOpts['blocDomain'],
        walletMediatorURL: ('walletMediatorURL' in startupOpts) ? startupOpts['walletMediatorURL'] : defaultTrustBlocStartupOpts['walletMediatorURL'],
        credentialMediatorURL: credentialMediator(('credentialMediatorURL' in startupOpts) ? startupOpts['credentialMediatorURL'] : defaultTrustBlocStartupOpts['credentialMediatorURL']),
        'log-level': ('log-level' in startupOpts) ? startupOpts['log-level'] : defaultTrustBlocStartupOpts['log-level'],
        agentUsername: ('agentUsername' in startupOpts) ? startupOpts['agentUsername'] : defaultTrustBlocStartupOpts['agentUsername'],
        sdsServerURL: ('sdsServerURL' in startupOpts) ? startupOpts['sdsServerURL'] : defaultTrustBlocStartupOpts['sdsServerURL']
    }
}

var credentialMediator = url => url ? `${url}?origin=${encodeURIComponent(window.location.origin)}` : undefined

// configure router
const router = new VueRouter({
    mode: 'history',
    routes, // short for routes: routes
    linkExactActiveClass: "nav-item active"
});

Vue.use(VueRouter);
Vue.use(MaterialDashboard);

new Vue({
    store,
    el: "#app",
    data: () => ({
        loaded: false,
    }),
    methods: mapActions(['initStore', 'onDidExchangeState', 'onIssueCredentialState']),
    mounted: async function () {
        // gets aries options
        let ariesOpts = await ariesStartupOpts()
        // sets aries instance globally
        window.$aries = await new Aries.Framework(ariesOpts)
        Vue.prototype.$arieslib = window.$aries

        // gets trustbloc options
        let trustblocOpts = await trustblocStartupOpts()
        Vue.prototype.$trustblocStartupOpts = trustblocOpts

        // registers listener which will update connections
        window.$aries.startNotifier(this.onDidExchangeState, ["didexchange_states"])
        // registers listener which will update credentials
        window.$aries.startNotifier(this.onIssueCredentialState, ["issue-credential_states"])
        // inits storage
        await this.initStore({aries: ariesOpts, trustbloc: trustblocOpts})
        // removes spinner
        this.loaded = true
    },
    render: h => h(App),
    router
});

