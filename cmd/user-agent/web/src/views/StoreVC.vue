/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <textarea id="vcDataTextArea" readonly rows="30" cols="200"/>
        <br>
        <input id="friendlyName" placeholder="friendly name">
        <br>
        <button id='storeVCBtn' :disabled=isDisabled>Store VC</button>
    </div>
</template>

<script>
    async function handleWalletReceiveEvent() {
        const credentialEvent = await window.$webCredentialHandler.receiveCredentialEvent();
        const vcData = credentialEvent.credential.data

        window.console.log('Received vc data:', vcData);
        document.getElementById('vcDataTextArea').value = vcData

        document.getElementById('storeVCBtn').addEventListener('click', async () => {
            // Get the friendly name
            const friendlyName = document.getElementById('friendlyName').value
            if (!friendlyName) {
                alert("please enter friendly name")
                return
            }

            // Save the VC
            let status = 'success'
            await window.$aries.verifiable.saveCredential({
                name: friendlyName,
                verifiableCredential: vcData.toString()
            }).then(() => {
                    console.log('vc save success')
                }
            ).catch(err => {
                status = err.toString()
                console.log('vc save failed : errMsg=' + err)
            })

            // Call Credential Handler callback
            credentialEvent.respondWith(new Promise(function (resolve) {
                return resolve({
                    dataType: "Response",
                    data: status
                });
            }))
        });
    }

    export default {
        beforeCreate: async function () {
            window.$webCredentialHandler = this.$webCredentialHandler
            this.$polyfill.loadOnce().then(handleWalletReceiveEvent)
            window.$aries = await this.$arieslib

            // enable send vc button once aries is loaded
            this.sendButton = false
        },
        computed: {
            isDisabled() {
                return this.sendButton
            },
        },
        data() {
            return {
                sendButton: true,
            };
        },
    }
</script>

