/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <textarea id="didDocTextArea" readonly rows="30" cols="200"/>
        <br>
        <button id='createDIDBtn' v-on:click="createDID">Create DID</button>
    </div>
</template>

<script>
    export default {
        beforeCreate: async function () {
            window.$aries = await this.$arieslib
            window.$trustblocAgent = this.$trustblocAgent
            window.$trustblocStartupOpts = await this.$trustblocStartupOpts

        },
        methods: {
            createDID: async function () {
                // TODO remove create DID page after this issue is done https://github.com/trustbloc/edge-agent/issues/59
                const keyset= await window.$aries.kms.createKeySet()
                // create did request
                const createDIDRequest = {
                    "publicKey":{
                        "id":"#key-1",
                        "type":"Ed25519VerificationKey2018",
                        "value":keyset.signaturePublicKey
                    },
                };

                const t= await new window.$trustblocAgent.Framework(JSON.parse(window.$trustblocStartupOpts))
                await t.didclient.createDID(createDIDRequest).then(
                    resp => {
                    // TODO generate public key from generic wasm
                    // TODO pass public key to createDID
                    document.getElementById('didDocTextArea').value = JSON.stringify(resp)
                    })
                    .catch(err => {
                        document.getElementById('didDocTextArea').value = err
                    })
                await t.destroy()
            }
        }
    }
</script>

