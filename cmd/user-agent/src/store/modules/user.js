/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from "../../pages/chapi/wallet";

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
        async login({commit}, username) {
            commit('setUser', username)

            await new WalletManager().getWalletMetadata(username).then(
                resp => {
                    commit('setUserMetadata', JSON.stringify(resp))
                }
            )

        },
        logout({commit}) {
            commit('clearUser')
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
}
