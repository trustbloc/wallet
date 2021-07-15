/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {getMediatorConnections} from "@trustbloc/wallet-sdk"

export default {
    actions: {
        onDidExchangeState({dispatch}, notice) {
            if (notice.payload.Type !== "post_state") {
                return
            }

            dispatch('queryConnections')
        },
        async queryConnections({commit, getters}) {
            let agent = getters['agent/getInstance']
            // retrieves all agent connections
            let res = await agent.didexchange.queryConnections()
            if (res.hasOwnProperty('results')) {
                // sets connections
                commit('updateConnections', res.results)
            }

            return res.results
        },
        async createInvitation(ctx, label) {
            let agent = ctx.getters['agent/getInstance']
            // creates invitation through the out-of-band protocol
            let res = await agent.outofband.createInvitation({
                label: label,
                router_connection_id: await getMediatorConnections(agent, {single: true})
            })

            return res.invitation
        },
        async acceptExchangeRequest({dispatch, getters}, id) {
            let agent = getters['agent/getInstance']
            agent.didexchange.acceptExchangeRequest({
                id: id,
                router_connections: await getMediatorConnections(agent, {single: true}),
            }).then(() => dispatch('queryConnections'))
        },
        async acceptInvitation({dispatch, getters}, payload) {
            let agent = getters['agent/getInstance']
            // accepts invitation thought out-of-band protocol
            let res = await agent.outofband.acceptInvitation({
                ...payload,
                router_connections: await getMediatorConnections(agent, {single: true})
            })

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
