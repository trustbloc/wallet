/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {Messenger, WalletManager} from "../../pages/chapi/wallet";
import * as Agent from "@trustbloc/agent-sdk";

// TODO message type domain needs to be finalized
const msgServices = [
    {name: 'request-peer-did', type: 'https://didcomm.org/peerdidrequest/1.0/message'},
    {name: 'create-conn-resp', type: 'https://trustbloc.dev/blinded-routing/1.0/create-conn-resp'},
    {name: 'diddoc-resp', type: 'https://trustbloc.dev/blinded-routing/1.0/diddoc-resp'},
    {name: 'register-route-res', type: 'https://trustbloc.dev/blinded-routing/1.0/register-route-resp'},
    {name: 'diddoc-res', type: 'https://trustbloc.dev/adapter/1.0/diddoc-resp'},
]

export default {
    state: {
        username: null,
        metadata: null,
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
        clearUser(state) {
            state.username = null
            state.metadata = null

            localStorage.removeItem('user')
            localStorage.removeItem('metadata')
        },
        loadUser(state) {
            state.username = localStorage.getItem('user');
            state.metadata = localStorage.getItem('metadata');
        }
    },
    actions: {
        async login({commit, dispatch}, username) {
            commit('setUser', username)

            await new WalletManager().getWalletMetadata(username).then(
                async resp => {
                    commit('setUserMetadata', JSON.stringify(resp))
                    await dispatch('agent/init')
                }
            )

        },
        async refreshUserMetadata({commit, state}) {
            if (!state.username) {
                throw 'invalid operation, user not logged in'
            }

            await new WalletManager().getWalletMetadata(state.username).then(
                async resp => {
                    commit('setUserMetadata', JSON.stringify(resp))
                }
            )
        },
        async loadOIDCUser({commit, dispatch}) {
            let userInfo = await fetch("/oidc/userinfo")

            if (userInfo.ok) {
                let profile = await userInfo.json()
                console.log("received user data: " + JSON.stringify(profile, null, 2))

                commit('setUser', profile.sub)

                await new WalletManager().getWalletMetadata(profile.sub).then(
                    async resp => {
                        commit('setUserMetadata', JSON.stringify(resp))
                        await dispatch('agent/init')
                    }
                )
            }
        },
        async logout({commit, dispatch}) {
            await fetch("/oidc/logout")
            commit('clearUser')
            await dispatch('agent/destroy')
        },
        loadUser({commit}) {
            commit('loadUser')
        }
    },
    getters: {
        getCurrentUser(state) {
            return state.username ? {username: state.username, metadata: state.metadata} : undefined
        }
    },
    modules: {
        agent: {
            namespaced: true,
            state: {
                instance: null,
                notifiers: null,
                agentName: null,
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
