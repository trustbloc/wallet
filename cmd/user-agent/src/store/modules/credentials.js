/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


export default {
    actions: {
        onIssueCredentialState({dispatch}, notice) {
            if (notice.payload.Type !== "post_state") {
                return
            }

            dispatch('getCredentials')
        },
        async getCredentials({commit, getters}) {
            // retrieves all agent credentials
            let res = await window.$aries.verifiable.getCredentials()
            if (!res.hasOwnProperty('result')) {
                return
            }

            res.result.forEach(function (v) {
                getters.completedConnections.forEach(function (conn) {
                    if (conn.MyDID !== v.my_did || conn.TheirDID !== v.their_did) {
                        return
                    }

                    v.label = conn.TheirLabel
                    if (!v.label) {
                        v.label = conn.ConnectionID
                    }
                })
            })

            // sets connections
            commit('updateCredentials', res.result.filter(v => v.label))

            return res.result
        },
    },
    mutations: {
        updateCredentials(state, credentials) {
            state.credentials = credentials
        },
    },
    state: {
        credentials: [],
    },
    getters: {
        allCredentials(state) {
            return state.credentials
        },
        allCredentialsCount(state, {allCredentials}) {
            return allCredentials.length
        },
    },
}