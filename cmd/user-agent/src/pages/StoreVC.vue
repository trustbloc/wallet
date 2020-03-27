/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="md-layout" style="margin-top: -5%;">
    <div class="md-layout-item">
            <md-textarea v-model="textarea" id="vcDataTextArea" style="visibility:hidden"/>
        <form>
            <md-card>
                <md-card-header style="background-color:#00bcd4">
                    <h3 class="title"><md-icon>fingerprint</md-icon>Credential</h3>
                </md-card-header>

                <md-card-content>
                    <div class="md-layout">
                        <div class="md-layout-item md-small-size-100 md-size-33">
                            <md-field>
                                <label><md-icon>face</md-icon>Subject</label>
                                <md-input v-model="subject" disabled id="subject"></md-input>
                            </md-field>
                        </div>
                        <div class="md-layout-item md-small-size-100 md-size-33">
                            <md-field>
                                <label><md-icon>how_to_reg</md-icon>Issuer</label>
                                <md-input v-model="issuer" disabled id="issuer"></md-input>
                            </md-field>
                        </div>
                        <div class="md-layout-item md-small-size-100 md-size-33">
                            <md-field>
                                <label><md-icon>today</md-icon>Issuance Date</label>
                                <md-input v-model="issuance" disabled id="issueDate"></md-input>
                            </md-field>
                        </div>
                        <div class="md-layout-item md-size-100">
                            <md-field maxlength="5">
                                <label class="md-helper-text">Type friendly name here</label>
                                <md-input v-model="friendlyName" id="friendlyName" required></md-input>
                            </md-field>
                        </div>
                        <div class="md-layout-item md-size-100 text-center">
                                <md-button class="md-cancel-text"  id="cancelBtn">Cancel</md-button>
                                <md-button class="md-raised md-success"  id="storeVCBtn" :disabled=isDisabled>Confirm</md-button>
                        </div>
                    </div>
                </md-card-content>
            </md-card>
        </form>

    </div>
        </div>
    </div>
</template>

<!--TODO : Fix the hard coded data with the actual data issue-90 -->
<script>
    async function handleWalletReceiveEvent() {
        const credentialEvent = await window.$webCredentialHandler.receiveCredentialEvent();
        window.console.log('Credential event:', credentialEvent);

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

        document.getElementById('cancelBtn').addEventListener('click', async () => {
            credentialEvent.respondWith(new Promise(function (resolve) {
                return resolve({
                    dataType: "Response",
                    data: 'cancelled'
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
                subject: "Foo",
                issuer: "did:example:ebfeb1f712ebc6f1c276e12ec21",
                issuance: "2010-01-01T19:23:24Z",
            };
        },
    }
</script>

