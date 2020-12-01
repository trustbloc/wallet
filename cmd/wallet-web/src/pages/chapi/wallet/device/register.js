/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';

const client = axios.create({
    withCredentials: true
})
/**
 * DeviceRegister provides device registration features
 * @param agent instance
 * @class
 */
export class DeviceRegister {
    constructor(agent) {
        this.agent = agent
    }

    async register() {
     //   let alertMessage;
        const serverURL = this.agent['edge-agent-server'];
        var registerSuccess = 'none';
        await client.get(serverURL+'/device/register/begin',null,
            function (data) {
                return data
            },
            'json')
            .then((credentialCreationOptions) => {
                console.log(credentialCreationOptions);
                credentialCreationOptions.data.publicKey.challenge = bufferDecode(credentialCreationOptions.data.publicKey.challenge);
                credentialCreationOptions.data.publicKey.user.id = bufferDecode(credentialCreationOptions.data.publicKey.user.id);
                if (credentialCreationOptions.data.publicKey.excludeCredentials) {
                    for (var i = 0; i < credentialCreationOptions.data.publicKey.excludeCredentials.length; i++) {
                        credentialCreationOptions.data.publicKey.excludeCredentials[i].id = bufferDecode(credentialCreationOptions.data.publicKey.excludeCredentials[i].id);
                    }
                }

                return navigator.credentials.create({
                    publicKey: credentialCreationOptions.data.publicKey
                })
            })
            .then((credential) => {
                console.log(credential);
                let attestationObject = credential.response.attestationObject;
                let clientDataJSON = credential.response.clientDataJSON;
                let rawId = credential.rawId;

                client.post(
                    serverURL+'/device/register/finish',
                    JSON.stringify({
                        id: credential.id,
                        rawId: bufferEncode(rawId),
                        type: credential.type,
                        response: {
                            attestationObject: bufferEncode(attestationObject),
                            clientDataJSON: bufferEncode(clientDataJSON),
                        },
                    }),
                    function (data) {
                        return data
                    },
                    'json')
            })
            // eslint-disable-next-line no-unused-vars
            .then((success) => {
                registerSuccess = 'success';
                return
            })
            .catch((error) => {
                console.log(error);
                registerSuccess = 'failure';
            })
        return registerSuccess;
    }
}

function bufferDecode(value) {
    return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}


function bufferEncode(value) {
    return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=/g, "");
}
