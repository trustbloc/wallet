/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';

const routerCreateInvitationPath = `/didcomm/invitation`
const stateCompleteMessageType = 'https://trustbloc.dev/didexchange/1.0/state-complete'


export async function connectToMediator(agent, mediatorEndpoint){
    let invitation = await createInvitationFromRouter(mediatorEndpoint)
    let resp = await agent.mediatorclient.connect({
        myLabel: 'agent-default-label', invitation, stateCompleteMessageType
    })

    if (resp.connectionID) {
        console.log("router registered successfully!", resp.connectionID)
    } else {
        console.log("router was not registered!")
    }
}

export async function getMediatorConnections(agent, single) {
    let resp = await agent.mediator.getConnections()
    if (!resp.connections || resp.connections.length === 0) {
        return ""
    }

    if (single) {
        return resp.connections[Math.floor(Math.random() * resp.connections.length)]
    }

    return resp.connections.join(",");
}

export const createInvitationFromRouter = async (endpoint) => {
    const response = await axios.get(`${endpoint}${routerCreateInvitationPath}`)
    return response.data.invitation
}
