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
        async queryConnections({commit, getters}) {
            let aries = getters['aries/getInstance']
            // retrieves all agent connections
            let res = await aries.didexchange.queryConnections()
            if (res.hasOwnProperty('results')) {
                // sets connections
                commit('updateConnections', res.results)
            }

            return res.results
        },
        async createInvitation(ctx, label) {
            let aries = ctx.getters['aries/getInstance']
            // creates invitation through the out-of-band protocol
            let res = await aries.outofband.createInvitation({label: label})

            return res.invitation
        },
        acceptExchangeRequest({dispatch, getters}, id) {
            let aries = getters['aries/getInstance']
            aries.didexchange.acceptExchangeRequest({id: id}).then(() => dispatch('queryConnections'))
        },
        async acceptInvitation({dispatch, getters}, payload) {
            let aries = getters['aries/getInstance']
            // accepts invitation thought out-of-band protocol
            let res = await aries.outofband.acceptInvitation(payload)

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
