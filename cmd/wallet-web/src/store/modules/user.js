/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {Messenger, WalletManager} from "../../pages/chapi/wallet";
import * as Agent from "@trustbloc/agent-sdk";

// TODO message type domain needs to be finalized
const msgServices = [
    {name: 'request-peer-did', type: 'https://didcomm.org/peerdidrequest/1.0/message'},
    {name: 'diddoc-res', type: 'https://trustbloc.dev/adapter/1.0/diddoc-resp'},
]

export default {
    state: {
        username: null,
        metadata: null,
        setupStatus: null,
    },
    mutations: {
        setUser(state, val) {
            state.username = val
            localStorage.setItem('user', val)
        },
        setUserMetadata(state, val) {
            state.metadata = val
            localStorage.setItem('metadata', val)
        },
        setUserSetupStatus(state, val){
            state.setupStatus = val
            localStorage.setItem('setupStatus', val)
        },
        clearUser(state) {
            state.username = null
            state.metadata = null
            state.setupStatus= null

            localStorage.removeItem('user')
            localStorage.removeItem('metadata')
            localStorage.removeItem('setupStatus')
        },
        loadUser(state) {
            state.username = localStorage.getItem('user');
            state.metadata = localStorage.getItem('metadata');
            state.setupStatus = localStorage.getItem('setupStatus');
        }
    },
    actions: {
        async refreshUserMetadata({commit, state, rootGetters}) {
            if (!state.username) {
                throw 'invalid operation, user not logged in'
            }

            await new WalletManager(rootGetters['agent/getInstance']).getWalletMetadata(state.username).then(
                async resp => {
                    commit('setUserMetadata', JSON.stringify(resp))
                }
            )
        },
        async loadOIDCUser({commit, dispatch, getters}) {
            console.log("getting from server URL", getters, getters.serverURL)
            let userInfo = await fetch(getters.serverURL + "/oidc/userinfo", {
                method: 'GET', credentials: 'include'
            })

            if (userInfo.ok) {
                let profile = await userInfo.json()
                console.log("received user data: " + JSON.stringify(profile, null, 2))

                commit('setUser', profile.sub)

                await dispatch('agent/init')
            }
        },
        async logout({commit, dispatch,getters}) {
            await fetch(getters.serverURL+"/oidc/logout",{
                method: 'GET',
                credentials: 'include'
            })
            commit('clearUser')
            await dispatch('agent/destroy')
        },
        loadUser({commit}) {
            commit('loadUser')
        },
        startUserSetup({commit}){
            commit('setUserSetupStatus', 'inprogress')
        },
        completeUserSetup({commit}, failure){
            commit('setUserSetupStatus', failure ? 'failed': 'success')
        }
    },
    getters: {
        getCurrentUser(state) {
            return state.username ? {username: state.username, metadata: state.metadata, setupStatus: state.setupStatus} : undefined
        },
    },
    modules: {
        agent: {
            namespaced: true,
            state: {
                instance: null,
                notifiers: null,
                agentName: null
            },
            mutations: {
                setInstance(state, {instance, user}) {
                    state.instance = instance
                    state.agentName = user
                },
                addNotifier(state, notifier) {
                    if (state.notifiers) {
                        state.notifiers.push(notifier)
                    } else {
                        state.notifiers = [notifier]
                    }
                },
                startNotifier(state, notifier) {
                    state.instance.startNotifier(notifier.callback, notifier.topics)
                },
                startAllNotifiers(state) {
                    if (!state.notifiers) {
                        return
                    }
                    state.notifiers.forEach(function (notifier) {
                        state.instance.startNotifier(notifier.callback, notifier.topics)
                    })
                }
            },
            actions: {
                async init({commit, rootState, state, rootGetters}) {
                    if (state.instance && state.agentName == rootState.user.username) {
                        return
                    }

                    if (!rootState.user.username) {
                        console.error('user should be logged in to initialize agent instance')
                        throw 'invalid user state'
                    }

                    let opts = {}
                    Object.assign(opts, rootGetters.getAgentOpts, {
                        'agent-default-label': rootState.user.username,
                        'db-namespace': rootState.user.username
                    })

                    let agent = await new Agent.Framework(opts)
                    await new WalletManager(agent).getWalletMetadata(rootState.user.username).then(
                        async resp => {
                            commit('setUserMetadata', JSON.stringify(resp))
                        }
                    )

                    let messenger = new Messenger(agent)

                    for (const {name, purpose, type} of msgServices) {
                        await messenger.register(name, purpose, type)
                    }

                    commit('setInstance', {instance: agent, user: rootState.user.username})
                    commit('startAllNotifiers')
                },
                async destroy({commit, state}) {
                    if (state.instance) {
                        await state.instance.destroy()
                    }
                    commit('setInstance', {})
                },
                addNotifier({commit, state}, notifier) {
                    commit('addNotifier', notifier)
                    if (state.instance) {
                        commit('startNotifier', notifier)
                    }
                }
            },
            getters: {
                getInstance(state) {
                    return state.instance
                },
                isInitialized(state) {
                    return state.instance != null
                }
            }
        }
    }

}
