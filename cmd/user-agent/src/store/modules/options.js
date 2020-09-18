/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const axios = require('axios').default;

const trustblocOptsLocation = l => `${l}/trustbloc-agent/jsopts`

const ariesOptsLocation = l => `${l}/aries/jsopts`

const credentialMediator = url => url ? `${url}?origin=${encodeURIComponent(window.location.origin)}` : undefined

const defaultTrustBlocStartupOpts = {
    assetsPath: '/trustbloc-agent/assets',
    blocDomain: 'testnet.trustbloc.local',
    'log-level': 'debug',
    walletMediatorURL: 'https://localhost:10063',
    credentialMediatorURL: '',
    sdsServerURL: ''
}

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

export default {
    actions: {
        async initOpts({commit}, location) {
            location = location ? location :  window.location.origin

            let tbOpts = {}
            let ariesOpts = {}

            console.log('process.env.NODE_ENV process.env.NODE_ENV', process.env.NODE_ENV)

            if (process.env.NODE_ENV === "production") {
                // call service to get the trustbloc opts
                await axios.get(trustblocOptsLocation(location))
                    .then(resp => {
                        tbOpts = resp.data
                        console.log("successfully fetched trustbloc start up options");
                    })
                    .catch(err => {
                        console.log("error fetching start up options - using default options : errMsg=", err);
                    })

                // call service to get the aries opts
                await axios.get(ariesOptsLocation(location))
                    .then(resp => {
                        ariesOpts = resp.data
                        console.log("successfully fetched aries start up options");
                    })
                    .catch(err => {
                        console.log("error fetching start up options - using default options : errMsg=", err);
                    })
            }

            commit('updateTrustblocOpts', {
                assetsPath: defaultTrustBlocStartupOpts['assetsPath'],
                blocDomain: ('blocDomain' in tbOpts) ? tbOpts['blocDomain'] : defaultTrustBlocStartupOpts['blocDomain'],
                walletMediatorURL: ('walletMediatorURL' in tbOpts) ? tbOpts['walletMediatorURL'] : defaultTrustBlocStartupOpts['walletMediatorURL'],
                credentialMediatorURL: credentialMediator(('credentialMediatorURL' in tbOpts) ? tbOpts['credentialMediatorURL'] : defaultTrustBlocStartupOpts['credentialMediatorURL']),
                'log-level': ('log-level' in tbOpts) ? tbOpts['log-level'] : defaultTrustBlocStartupOpts['log-level'],
                sdsServerURL: ('sdsServerURL' in tbOpts) ? tbOpts['sdsServerURL'] : defaultTrustBlocStartupOpts['sdsServerURL']
            })

            commit('updateAriesOpts', {
                assetsPath: defaultAriesStartupOpts['assetsPath'],
                'outbound-transport': defaultAriesStartupOpts['outbound-transport'],
                'transport-return-route': defaultAriesStartupOpts['transport-return-route'],
                'http-resolver-url': ('http-resolver-url' in ariesOpts) ? ariesOpts['http-resolver-url'] : defaultAriesStartupOpts['http-resolver-url'],
                'agent-default-label': ('agent-default-label' in ariesOpts) ? ariesOpts['agent-default-label'] : defaultAriesStartupOpts['agent-default-label'],
                'auto-accept': ('auto-accept' in ariesOpts) ? ariesOpts['auto-accept'] : defaultAriesStartupOpts['auto-accept'],
                'log-level': ('log-level' in ariesOpts) ? ariesOpts['log-level'] : defaultAriesStartupOpts['log-level'],
                'db-namespace': ('db-namespace' in ariesOpts) ? ariesOpts['db-namespace'] : defaultAriesStartupOpts['db-namespace']
            })
        },
    },
    mutations: {
        updateTrustblocOpts(state, opts) {
            state.trustblocOpts = opts
        },
        updateAriesOpts(state, opts) {
            state.ariesOpts = opts
        }
    },
    state: {
        ariesOpts: {},
        trustblocOpts: {},
    },
    getters: {
        getAriesOpts(state) {
            return state.ariesOpts
        },
        agentDefaultLabel(state) {
            return state.ariesOpts['agent-default-label']
        },
        getTrustblocOpts(state) {
            return state.trustblocOpts
        },
    }
}
