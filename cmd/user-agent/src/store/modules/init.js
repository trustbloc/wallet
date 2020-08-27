/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export default {
    actions: {
        async initStore({dispatch, commit, getters}, opts) {
            commit('updateAriesStartupOpts', opts.aries)
            commit('updateTrustblocStartupOpts', opts.trustbloc)

            await dispatch('loadMode')
            await dispatch('loadMediatorState')
            await dispatch('queryConnections')
            await dispatch('getCredentials')

            let mediatorURL = opts.trustbloc['walletMediatorURL']
            if (mediatorURL && !getters.isDevMode) {
                try {
                    await dispatch('registeredMediator', mediatorURL)
                } catch (error) {
                    console.error(error);
                }
            }
        },
    },
    mutations: {
        updateAriesStartupOpts(state, opts) {
            state.ariesStartupOpts = opts
        },
        updateTrustblocStartupOpts(state, opts) {
            state.trustblocStartupOpts = opts
        },
    },
    state: {
        ariesStartupOpts: {},
        trustblocStartupOpts: {},
    },
    getters: {
        getAriesStartupOpts(state) {
            return state.ariesStartupOpts
        },
        agentDefaultLabel(state) {
            return state.ariesStartupOpts['agent-default-label']
        },
        getTrustblocStartupOpts(state) {
            return state.trustblocStartupOpts
        },
    }
}