/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';
import {waitForEvent, POST_STATE} from '../../events'

const stateCompleted = 'completed'
const topicDidExchangeStates = 'didexchange_states'

export default {
    actions: {
        async loadMediatorState({commit, getters}) {
            let agent = getters['agent/getInstance']
            let res = await agent.mediator.getConnections()

            commit('updateMediatorConnections', res.connections)
        },
        async unregisteredMediator({dispatch, getters}, connID) {
            if (!getters.getMediatorConnections.includes(connID)) {
                return
            }

            let agent = getters['agent/getInstance']
            await agent.mediator.unregister({connectionID: connID})
            await dispatch('loadMediatorState')
        },
        async registeredMediator({dispatch, getters}, routerURL) {
            let invitation = await axios.get(routerURL + '/didcomm/invitation')
            let agent = getters['agent/getInstance']
            // accepts invitation thought out-of-band protocol
            let conn = await agent.outofband.acceptInvitation({
                my_label: getters.agentDefaultLabel,
                invitation: invitation.data.invitation,
            })

            let connID = conn['connection_id']

            await waitForEvent(agent, {
                type: POST_STATE,
                stateID: stateCompleted,
                connectionID: connID,
                topic: topicDidExchangeStates,
            })

            await agent.mediator.register({connectionID: connID})

            await dispatch('queryConnections')
            await dispatch('loadMediatorState')
        },
    },
    mutations: {
        updateMediatorConnections(state, connections) {
            state.mediatorConnections = connections
        },
    },
    state: {
        mediatorConnections: [],
    },
    getters: {
        getMediatorConnections(state) {
            return state.mediatorConnections
        },
        isMediatorRegistered(state, getters) {
            return getters.getMediatorConnections && getters.getMediatorConnections.length > 0
        },
    }
}
