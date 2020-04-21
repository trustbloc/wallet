/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="md-layout" style="margin-top: -5%;">
            <div class="md-layout-item">
                <md-textarea v-model="vcdata" id="vcDataTextArea" style="visibility:hidden"/>
                <form>
                    <md-card>
                        <md-card-header style="background-color:#00bcd4">
                            <h3 class="title">
                                <md-icon>fingerprint</md-icon>
                                 Credential
                            </h3>
                        </md-card-header>

                        <md-card-content>
                            <div class="md-layout">

                                <div class="md-layout-item md-small-size-100 md-size-33">
                                    <md-field>
                                        <label>
                                            <md-icon>face</md-icon>
                                            Subject</label>
                                        <md-input v-model="subject" disabled id="subject"></md-input>
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-small-size-100 md-size-33">
                                    <md-field>
                                        <label>
                                            <md-icon>how_to_reg</md-icon>
                                            Issuer</label>
                                        <md-input v-model="issuer" disabled id="issuer"></md-input>
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-small-size-100 md-size-33">
                                    <md-field>
                                        <label>
                                            <md-icon>today</md-icon>
                                            Issuance Date</label>
                                        <md-input v-model="issuance" disabled id="issueDate"></md-input>
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-size-100">
                                    <md-field maxlength="5">
                                        <label class="md-helper-text">Type friendly name here</label>
                                        <md-input v-model="friendlyName" id="friendlyName" required></md-input>
                                    </md-field>
                                </div>

                                <div v-if="errors.length">
                                    <b>Please correct the following error(s):</b>
                                    <ul>
                                        <li v-for="error in errors" :key="error">{{ error }}</li>
                                    </ul>
                                </div>

                                <div class="md-layout-item md-size-100 text-center">
                                    <md-button v-on:click="cancel" class="md-cancel-text" id="cancelBtn">Cancel
                                    </md-button>
                                    <md-button v-on:click="store" class="md-raised md-success" id="storeVCBtn"
                                               :disabled=isDisabled>Confirm
                                    </md-button>
                                </div>
                            </div>
                        </md-card-content>
                    </md-card>
                </form>

            </div>
        </div>
    </div>
</template>

<script>
    export default {
        beforeCreate: async function () {
            window.$webCredentialHandler = this.$webCredentialHandler
            this.$polyfill.loadOnce()

            const credentialEvent = await window.$webCredentialHandler.receiveCredentialEvent();
            this.vcdata = credentialEvent.credential.data
            this.credentialEvent = credentialEvent

            window.console.log("Received VC:", this.vcdata);

            const vc = JSON.parse(this.vcdata)

            // set issuance date
            if (!vc.issuanceDate) {
                this.issuance = new Date(vc.issuanceDate);
            } else {
                this.issuance = new Date()
            }

            // set issuer
            if (vc.issuer && vc.issuer.id) {
                this.issuer = vc.issuer.id
            } else if (vc.issuer) {
                this.issuer = vc.issuer
            }

            // set type as subject
            if (vc.type && Array.isArray(vc.type)) {
                vc.type.forEach((item) => {
                    if (item != 'VerifiableCredential') {
                        this.subject = item
                    }
                });
            }


            window.$aries = await this.$arieslib
            // enable send vc button loaded
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
                subject: "",
                issuer: "",
                issuance: "",
                friendlyName: "",
                vcdata: "",
                errors: [],
            };
        },
        methods: {
            store: async function () {
                this.errors.length = 0
                if (this.friendlyName.length == 0) {
                   this.errors.push("friendly name required.")
                    return
                }

                // Save the VC
                let status = 'success'
                await window.$aries.verifiable.saveCredential({
                    name: this.friendlyName,
                    verifiableCredential: this.vcdata
                }).then(() => {
                        console.log('vc save success')
                    }
                ).catch(err => {
                    status = err.toString()
                    console.log('vc save failed : errMsg=' + err)
                })

                // Call Credential Handler callback
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: "Response",
                        data: status
                    });
                }))
            },
            cancel: async function () {
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: "Response",
                        data: 'cancelled'
                    });
                }))
            }
        }
    }
</script>

