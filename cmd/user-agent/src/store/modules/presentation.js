/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


export default {
    actions: {
        onPresentProofState({dispatch}, notice) {
            if (notice.payload.Type !== "post_state") {
                return
            }

            dispatch('getPresentations')
        },
        async getPresentations({commit, getters}) {
            // retrieves all agent credentials
            let res = await window.$aries.verifiable.getPresentations()
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
            commit('updatePresentations', res.result)

            return res.result
        },
    },
    mutations: {
        updatePresentations(state, presentations) {
            state.presentations = presentations
        },
    },
    state: {
        presentations: [],
    },
    getters: {
        allPresentations(state) {
            return state.presentations
        },
        associatedPresentations(state, {allPresentations}) {
            return allPresentations.filter(v => v.label)
        },
        associatedPresentationsCount(state, {associatedPresentations}) {
            return associatedPresentations.length
        },
    },
}