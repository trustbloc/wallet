/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from "../../pages/chapi/wallet";
import * as Aries from "@trustbloc-cicd/aries-framework-go";

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
                    await dispatch('aries/init')
                }
            )

        },
        async logout({commit, dispatch}) {
            commit('clearUser')
            await dispatch('aries/destroy')
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
        aries: {
            namespaced: true,
            state: {
                instance: null,
                opts: null,
                agentName: null,
            },
            mutations: {
                setInstance(state, {instance, user}) {
                    state.instance = instance
                    state.agentName = user
                },
                setOpts(state, opts) {
                    state.opts = opts
                },
            },
            actions: {
                async init({commit, rootState, state}) {
                    if (!state.opts) {
                        console.error('aries opts should be set before initializing aries')
                        throw 'invalid aries opts'
                    }

                    if (!rootState.user.username) {
                        console.error('user should be logged in to initialize aries instance')
                        throw 'invalid user state'
                    }

                    if (state.instance && state.agentName == rootState.user.username) {
                        return
                    }

                    let opts = {}
                    Object.assign(opts, state.opts, {
                        'agent-default-label': rootState.user.username,
                        'db-namespace': rootState.user.username
                    })

                    let aries = await new Aries.Framework(opts)

                    commit('setInstance', {instance: aries, user: rootState.user.username})
                },
                async destroy({commit, state}) {
                    if (state.instance) {
                        await state.instance.destroy()
                    }
                    commit('setInstance', {})
                },
                setOpts({commit}, opts) {
                    commit('setOpts', opts)
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
