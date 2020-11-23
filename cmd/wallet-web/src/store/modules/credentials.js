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
            let agent = getters['agent/getInstance']
            // retrieves all agent credentials
            let res = await agent.verifiable.getCredentials()
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
            commit('updateCredentials', res.result)

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
        associatedCredentials(state, {allCredentials}) {
            return allCredentials.filter(v => v.label)
        },
        associatedCredentialsCount(state, {associatedCredentials}) {
            return associatedCredentials.length
        },
    },
}
