/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


export default {
    state: {
        loginHandleState: false,
    },
    mutations: {
        setLoginHandleSetup(state, val) {
            state.loginHandleState = val
            localStorage.setItem('loginHandleState', val)
        },
    },
    actions: {
        popUpOpenedSetup({commit}) {
            commit('setLoginHandleSetup', 'true')
        },
        popUpClosedSetup({commit}) {
            window.top.close()
            commit('setLoginHandleSetup', 'false')
        },
    },
    getters: {
        getLoginHandleSetup(state) {
            return state.loginHandleState
        }
    },
}