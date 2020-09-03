/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//TODO integrate with wallet manager in addition to local storage
export default {
    actions: {
        setUser({commit}, username) {
            commit('setUsername', username)
        },
        resetUser({commit}) {
            commit('removeUsername')
        },
        initUserStore({commit}) {
            commit('initiateUserStore')
        },
    },
    mutations: {
        setUsername(state, val) {
            state.username = val
            localStorage.setItem('username', val)
        },
        removeUsername(state) {
            state.username = null
            localStorage.removeItem('username')
        },
        initiateUserStore(state) {
            state.username = localStorage.getItem('username');
        }
    },
    state: {
        username: null,
    },
    getters: {
        getUser(state) {
            return state.username
        },
    },
}
