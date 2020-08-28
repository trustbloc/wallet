/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="md-layout" style="margin-top: -5%;">
            <div class="md-layout-item">
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
                                <div class="md-layout md-gutter">
                                <div class="md-layout-item md-layout md-gutter">
                                    <div class="md-layout-item md-size-80">
                                        <md-button v-on:click="cancel" class="md-cancel-text" id="cancelBtn">Cancel
                                        </md-button>
                                    </div>
                                    <div class="md-layout-item md-size-20">
                                        <md-button v-on:click="store" class="md-raised md-success" id="storeVCBtn"
                                                                           :disabled=isDisabled>Confirm
                                         </md-button>
                                    </div>
                                  </div>
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
    import {isCredentialType, isVCType, getCredentialMetadata, WalletStore} from "./wallet"

    export default {
        beforeCreate: async function () {

            const credentialEvent = await this.$webCredentialHandler.receiveCredentialEvent();
            console.log("Credential event received :", credentialEvent.credential)

            if (!isCredentialType(credentialEvent.credential.dataType)) {
                this.errors.push("unknown credential data type", credentialEvent.credential.dataType)
                return
            }

            this.isVC = isVCType(credentialEvent.credential.dataType)
            this.credData = credentialEvent.credential.data
            this.dataType = credentialEvent.credential.dataType

            this.wallet = new WalletStore(await this.$arieslib,
                await new this.$trustblocAgent.Framework(await this.$trustblocStartupOpts), this.$trustblocStartupOpts,
                credentialEvent)

            // prefill form
            this.prefillForm()

            // enable send vc button once loaded
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
                credData: {},
                errors: [],
            };
        },
        methods: {
            prefillForm: function() {
                const {issuance, issuer, subject} = getCredentialMetadata(this.credData, this.dataType)
                this.issuance = issuance
                this.issuer = issuer
                this.subject = subject
                this.friendlyName = subject.concat(' ', issuance)
            },
            store: async function () {
                this.errors.length = 0
                if (this.friendlyName.length == 0) {
                    this.errors.push("friendly name required.")
                    return
                }

                await this.wallet.saveCredential(this.friendlyName, this.credData, this.isVC)
            },
            cancel: async function () {
                this.wallet.cancel()
            }
        }
    }
</script>

