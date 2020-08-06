/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


export default {
    actions: {
        onDidExchangeState({dispatch}, notice) {
            if (notice.payload.Type !== "post_state") {
                return
            }

            dispatch('queryConnections')
        },
        async queryConnections({commit}) {
            // retrieves all agent connections
            let res = await window.$aries.didexchange.queryConnections()
            if (res.hasOwnProperty('results')) {
                // sets connections
                commit('updateConnections', res.results)
            }

            return res.results
        },
        async createInvitation(ctx, label) {
            // creates invitation through the out-of-band protocol
            let res = await window.$aries.outofband.createInvitation({label: label})

            return res.invitation
        },
        acceptExchangeRequest({dispatch}, id) {
            window.$aries.didexchange.acceptExchangeRequest({id: id}).then(() => dispatch('queryConnections'))
        },
        async acceptInvitation({dispatch}, payload) {
            // accepts invitation thought out-of-band protocol
            let res = await window.$aries.outofband.acceptInvitation(payload)

            dispatch('queryConnections')

            return res;
        },
    },
    mutations: {
        updateConnections(state, connections) {
            state.connections = connections
        },
    },
    state: {
        connections: [],
    },
    getters: {
        allConnections(state) {
            return state.connections
        },
        pendingConnections(state) {
            return state.connections.filter(conn => conn.State === 'requested' && conn.Namespace === 'their')
        },
        pendingConnectionsCount(state, {pendingConnections}) {
            return pendingConnections.length
        },
        completedConnections(state) {
            return state.connections.filter(conn => conn.State === 'completed')
        },
        completedConnectionsCount(state, {completedConnections}) {
            return completedConnections.length
        },
    },
}