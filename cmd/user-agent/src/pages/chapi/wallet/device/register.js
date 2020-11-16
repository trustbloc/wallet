/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';
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
        function bufferDecode(value) {
            return Uint8Array.from(atob(value), c => c.charCodeAt(0));
        }

        function bufferEncode(value) {
            return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
                .replace(/\+/g, "-")
                .replace(/\//g, "_")
                .replace(/=/g, "");
        }

        axios.get('device/register/begin', {
                headers: {
                    Cookie: "device_user=walletuser;"
                }
            },
            function (data) {
                return data
            },
            'json')
            .then((credentialCreationOptions) => {
                console.log(credentialCreationOptions)
                credentialCreationOptions.publicKey.challenge = bufferDecode(credentialCreationOptions.publicKey.challenge);
                credentialCreationOptions.publicKey.user.id = bufferDecode(credentialCreationOptions.publicKey.user.id);
                if (credentialCreationOptions.publicKey.excludeCredentials) {
                    for (var i = 0; i < credentialCreationOptions.publicKey.excludeCredentials.length; i++) {
                        credentialCreationOptions.publicKey.excludeCredentials[i].id = bufferDecode(credentialCreationOptions.publicKey.excludeCredentials[i].id);
                    }
                }

                return navigator.credentials.create({
                    publicKey: credentialCreationOptions.publicKey
                })
            })
            .then((credential) => {
                console.log(credential)
                let attestationObject = credential.response.attestationObject;
                let clientDataJSON = credential.response.clientDataJSON;
                let rawId = credential.rawId;

                axios.post(
                    '/register/finish',
                    JSON.stringify({
                        id: credential.id,
                        rawId: bufferEncode(rawId),
                        type: credential.type,
                        response: {
                            attestationObject: bufferEncode(attestationObject),
                            clientDataJSON: bufferEncode(clientDataJSON),
                        },
                    }),
                    {
                        headers: {
                            Cookie: "device_user=walletuser;"
                        }
                    },
                    function (data) {
                        return data
                    },
                    'json')
            })
            // eslint-disable-next-line no-unused-vars
            .then((success) => {
                alert("successfully registered")
                return
            })
            .catch((error) => {
                console.log(error)
                alert("failed to register")
            })
    }
}
