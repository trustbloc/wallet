/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const axios = require('axios').default;

const agentOptsLocation = l => `${l}/walletconfig/agent`
const credentialMediator = url => url ? `${url}?origin=${encodeURIComponent(window.location.origin)}${__webpack_public_path__}/` : undefined

let defaultAgentStartupOpts = {
    assetsPath: `${__webpack_public_path__}/agent-js-worker/assets`,
    'outbound-transport': ['ws', 'http'],
    'transport-return-route': 'all',
    'http-resolver-url': ["trustbloc:testnet.trustbloc.local@http://localhost:8080/1.0/identifiers", "web@http://localhost:8080/1.0/identifiers"],

    'agent-default-label': 'demo-user-agent',
    'auto-accept': true,
    'log-level': 'debug',
    'db-namespace': 'agent',

    blocDomain: 'testnet.trustbloc.local',
    walletMediatorURL: 'https://localhost:10063',
    blindedRouting: false,
    credentialMediatorURL: '',
    sdsServerURL: ''
}

export default {
    actions: {
        async initOpts({commit}, location) {
            location = location ? location : window.location.origin

            let agentOpts = {}

            console.log('process.env.NODE_ENV process.env.NODE_ENV', process.env.NODE_ENV)

            if (process.env.NODE_ENV === "production") {
                // call service to get the agent opts
                await axios.get(agentOptsLocation(location))
                    .then(resp => {
                        agentOpts = resp.data
                        console.log("successfully fetched agent start up options");
                    })
                    .catch(err => {
                        console.log("error fetching start up options - using default options : errMsg=", err);
                    })
            }

            commit('updateAgentOpts', {
                assetsPath: defaultAgentStartupOpts['assetsPath'],
                'outbound-transport': defaultAgentStartupOpts['outbound-transport'],
                'transport-return-route': defaultAgentStartupOpts['transport-return-route'],
                'http-resolver-url': ('http-resolver-url' in agentOpts) ? agentOpts['http-resolver-url'] : defaultAgentStartupOpts['http-resolver-url'],
                'agent-default-label': ('agent-default-label' in agentOpts) ? agentOpts['agent-default-label'] : defaultAgentStartupOpts['agent-default-label'],
                'auto-accept': ('auto-accept' in agentOpts) ? agentOpts['auto-accept'] : defaultAgentStartupOpts['auto-accept'],
                'log-level': ('log-level' in agentOpts) ? agentOpts['log-level'] : defaultAgentStartupOpts['log-level'],
                'db-namespace': ('db-namespace' in agentOpts) ? agentOpts['db-namespace'] : defaultAgentStartupOpts['db-namespace'],

                blocDomain: ('blocDomain' in agentOpts) ? agentOpts['blocDomain'] : defaultAgentStartupOpts['blocDomain'],
                walletMediatorURL: ('walletMediatorURL' in agentOpts) ? agentOpts['walletMediatorURL'] : defaultAgentStartupOpts['walletMediatorURL'],
                credentialMediatorURL: credentialMediator(('credentialMediatorURL' in agentOpts) ? agentOpts['credentialMediatorURL'] : defaultAgentStartupOpts['credentialMediatorURL']),
                blindedRouting: ('blindedRouting' in agentOpts) ? agentOpts['blindedRouting'] : defaultAgentStartupOpts['blindedRouting'],
                sdsServerURL: ('sdsServerURL' in agentOpts) ? agentOpts['sdsServerURL'] : defaultAgentStartupOpts['sdsServerURL']
            })
        },
    },
    mutations: {
        updateAgentOpts(state, opts) {
            state.agentOpts = opts
        }
    },
    state: {
        agentOpts: {},
    },
    getters: {
        getAgentOpts(state) {
            return state.agentOpts
        },
        agentDefaultLabel(state) {
            return state.agentOpts['agent-default-label']
        },
    }
}
