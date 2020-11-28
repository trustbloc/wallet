/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const axios = require('axios').default;

const agentOptsLocation = l => `${l}/walletconfig/agent`
const credentialMediator = url => url ? `${url}?origin=${encodeURIComponent(window.location.origin)}${__webpack_public_path__}/` : undefined

let defaultAgentStartupOpts = {
    assetsPath: `${__webpack_public_path__}agent-js-worker/assets`,
    'outbound-transport': ['ws', 'http'],
    'transport-return-route': 'all',
    'http-resolver-url': ["trustbloc:testnet.trustbloc.local@http://localhost:8080/1.0/identifiers", "web@http://localhost:8080/1.0/identifiers"],

    'agent-default-label': 'demo-wallet-web',
    'auto-accept': true,
    'log-level': 'debug',
    'indexedDB-namespace': 'agent',
    'edge-agent-server':'',

    blocDomain: 'testnet.trustbloc.local',
    walletMediatorURL: 'https://localhost:10063',
    blindedRouting: false,
    credentialMediatorURL: '',
    storageType: `edv`,
    edvServerURL: '',
    edvVaultID: '',
    edvCapability: '',
    authzKeyStoreURL: '',
    opsKeyStoreURL: '',
    edvOpsKIDURL: '',
    edvHMACKIDURL: '',
    keyServer: {authzKMSURL: '', opsKMSURL: '', keyEDVURL:'', useRemoteKMS: true},
    useEDVCache: false,
    clearCache: ''
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

                agentOpts['http-resolver-url'] = agentOpts['http-resolver-url'].split(',')

                if (agentOpts.storageType === 'edv') {
                    const userInfoURL = agentOpts["edge-agent-server"] + "/oidc/userinfo"

                    console.log("User info URL is: " + userInfoURL)

                    const client = axios.create({
                        withCredentials: true
                    })

                    await client.get(userInfoURL)
                        .then(resp => {
                            const edvVaultURL = resp.data.bootstrap.edvVaultURL

                            console.log("User EDV Vault URL is: " + edvVaultURL)

                            const edvVaultID = edvVaultURL.substring(edvVaultURL.lastIndexOf('/')+1)

                            console.log("User EDV Vault ID is: " + edvVaultID)

                            agentOpts.edvVaultID = edvVaultID
                            agentOpts.edvCapability=resp.data.bootstrap.edvCapability
                            agentOpts.authzKeyStoreURL=resp.data.bootstrap.authzKeyStoreURL
                            agentOpts.opsKeyStoreURL=resp.data.bootstrap.opsKeyStoreURL
                            agentOpts.userConfig=resp.data.userConfig
                            agentOpts.edvOpsKIDURL=resp.data.bootstrap.edvOpsKIDURL
                            agentOpts.edvHMACKIDURL=resp.data.bootstrap.edvHMACKIDURL
                            agentOpts.keyServer=resp.data.keyServer
                        })
                        .catch(err => {
                            console.log("error fetching user info: errMsg=", err);
                            console.log("Note: If you haven't logged in yet and you just got a 403 error, then it's expected")
                        })
                }
                console.log("agent-sdk will be started with:")
                console.log(agentOpts)
            }

            commit('updateAgentOpts', {
                assetsPath: defaultAgentStartupOpts['assetsPath'],
                'outbound-transport': defaultAgentStartupOpts['outbound-transport'],
                'transport-return-route': defaultAgentStartupOpts['transport-return-route'],
                'http-resolver-url': ('http-resolver-url' in agentOpts) ? agentOpts['http-resolver-url'] : defaultAgentStartupOpts['http-resolver-url'],
                'agent-default-label': ('agent-default-label' in agentOpts) ? agentOpts['agent-default-label'] : defaultAgentStartupOpts['agent-default-label'],
                'auto-accept': ('auto-accept' in agentOpts) ? agentOpts['auto-accept'] : defaultAgentStartupOpts['auto-accept'],
                'log-level': ('log-level' in agentOpts) ? agentOpts['log-level'] : defaultAgentStartupOpts['log-level'],
                'indexedDB-namespace': ('indexedDB-namespace' in agentOpts) ? agentOpts['indexedDB-namespace'] : defaultAgentStartupOpts['indexedDB-namespace'],
                'edge-agent-server': ('edge-agent-server' in agentOpts) ? agentOpts['edge-agent-server'] : defaultAgentStartupOpts['edge-agent-server'],
                blocDomain: ('blocDomain' in agentOpts) ? agentOpts['blocDomain'] : defaultAgentStartupOpts['blocDomain'],
                walletMediatorURL: ('walletMediatorURL' in agentOpts) ? agentOpts['walletMediatorURL'] : defaultAgentStartupOpts['walletMediatorURL'],
                credentialMediatorURL: credentialMediator(('credentialMediatorURL' in agentOpts) ? agentOpts['credentialMediatorURL'] : defaultAgentStartupOpts['credentialMediatorURL']),
                blindedRouting: ('blindedRouting' in agentOpts) ? agentOpts['blindedRouting'] : defaultAgentStartupOpts['blindedRouting'],
                storageType: ('storageType' in agentOpts) ? agentOpts['storageType'] : defaultAgentStartupOpts['storageType'],
                edvServerURL: ('edvServerURL' in agentOpts) ? agentOpts['edvServerURL'] : defaultAgentStartupOpts['edvServerURL'],
                edvVaultID: ('edvVaultID' in agentOpts) ? agentOpts['edvVaultID'] : defaultAgentStartupOpts['edvVaultID'],
                edvCapability: ('edvCapability' in agentOpts) ? agentOpts['edvCapability'] : defaultAgentStartupOpts['edvCapability'],
                userConfig: ('userConfig' in agentOpts) ? agentOpts['userConfig'] : defaultAgentStartupOpts['userConfig'],
                authzKeyStoreURL: ('authzKeyStoreURL' in agentOpts) ? agentOpts['authzKeyStoreURL'] : defaultAgentStartupOpts['authzKeyStoreURL'],
                opsKeyStoreURL: ('opsKeyStoreURL' in agentOpts) ? agentOpts['opsKeyStoreURL'] : defaultAgentStartupOpts['opsKeyStoreURL'],
                edvOpsKIDURL: ('edvOpsKIDURL' in agentOpts) ? agentOpts['edvOpsKIDURL'] : defaultAgentStartupOpts['edvOpsKIDURL'],
                edvHMACKIDURL: ('edvHMACKIDURL' in agentOpts) ? agentOpts['edvHMACKIDURL'] : defaultAgentStartupOpts['edvHMACKIDURL'],
                keyServer: ('keyServer' in agentOpts) ? agentOpts['keyServer'] : defaultAgentStartupOpts['keyServer'],
                useEDVCache: ('useEDVCache' in agentOpts) ? agentOpts['useEDVCache'] : defaultAgentStartupOpts['useEDVCache'],
                clearCache: ('clearCache' in agentOpts) ? agentOpts['clearCache'] : defaultAgentStartupOpts['clearCache'],
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
        serverURL(state) {
            return state.agentOpts['edge-agent-server']
        },
    }
}
