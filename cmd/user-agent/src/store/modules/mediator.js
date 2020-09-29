/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';
import {waitForEvent, POST_STATE} from '../../events'

const stateCompleted = 'completed'
const topicDidExchangeStates = 'didexchange_states'
const errRouterNotRegistered = 'router not registered'

export default {
    actions: {
        async loadMediatorState({commit, getters}) {
            let aries = getters['aries/getInstance']
            let res = await aries.mediator.getConnection().catch(err => {
                if (!err.message.includes(errRouterNotRegistered)) {
                    throw err
                }
            })

            let connectionID = '';
            if (res) {
                connectionID = res.connectionID
            }

            commit('updateMediatorConnID', connectionID)
        },
        async unregisteredMediator({getters, commit}) {
            if (!getters.isMediatorRegistered) {
                return
            }

            let aries = getters['aries/getInstance']
            await aries.mediator.unregister({id: getters.getMediatorConnID})
            commit('updateMediatorConnID', '')
        },
        async registeredMediator({dispatch, getters}, routerURL) {
            if (getters.isMediatorRegistered) return;

            let invitation = await axios.get(routerURL + '/didcomm/invitation')
            let conn = await dispatch('acceptInvitation', {
                my_label: getters.agentDefaultLabel,
                invitation: invitation.data.invitation,
            })

            let connID = conn['connection_id']

            let aries = getters['aries/getInstance']
            await waitForEvent(aries, {
                type: POST_STATE,
                stateID: stateCompleted,
                connectionID: connID,
                topic: topicDidExchangeStates,
            })

            await aries.mediator.register({connectionID: connID})

            await dispatch('loadMediatorState')
        },
    },
    mutations: {
        updateMediatorConnID(state, conID) {
            state.mediatorConnID = conID
        },
    },
    state: {
        mediatorConnID: '',
    },
    getters: {
        getMediatorConnID(state) {
            return state.mediatorConnID
        },
        isMediatorRegistered(state, getters) {
            return getters.getMediatorConnID !== ''
        },
    }
}
