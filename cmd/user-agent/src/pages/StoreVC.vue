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
    const toLower = text => {
        return text.toString().toLowerCase()
    }

    const getCredentialType = (types) => {
        const result = types.filter(type => toLower(type) != "verifiablecredential")
        if (result.length > 0) {
            return result[0]
        }
        return ""
    }

    export default {
        beforeCreate: async function () {
            window.$webCredentialHandler = this.$webCredentialHandler
            this.$polyfill.loadOnce()

            const credentialEvent = await window.$webCredentialHandler.receiveCredentialEvent();
            console.log("Credential event received :", credentialEvent.credential)

            if (credentialEvent.credential.dataType == 'VerifiableCredential') {
                this.isVC = true
            } else if (credentialEvent.credential.dataType == 'VerifiablePresentation') {
                this.isVP = true
            } else {
                this.errors.push("unknown credential data type")
                return
            }

            this.credData = credentialEvent.credential.data
            this.credentialEvent = credentialEvent
            this.isVC ? this.populateVCData() : this.populateVPData()
            this.aries = await this.$arieslib
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
            populateVCData: function () {
                // set issuance date
                if (!this.credData.issuanceDate) {
                    this.issuance = new Date(this.credData.issuanceDate);
                } else {
                    this.issuance = new Date()
                }

                // set issuer
                if (this.credData.issuer && this.credData.issuer.id) {
                    this.issuer = this.credData.issuer.id
                } else if (this.credData.issuer) {
                    this.issuer = this.credData.issuer
                }

                // set type as subject
                if (this.credData.type && Array.isArray(this.credData.type)) {
                    this.subject = getCredentialType(this.credData.type)
                }
            },
            populateVPData: function () {
                this.issuance = new Date()
                this.subject = ''
                this.issuer = ''

                if (!this.credData.verifiableCredential) {
                    return
                }

                const allCreds = Array.isArray(this.credData.verifiableCredential) ? this.credData.verifiableCredential
                    : [this.credData.verifiableCredential];


                // set issuance date
                if (allCreds[0].issuanceDate) {
                    this.issuance = allCreds[0].issuanceDate
                }


                let tempIssuers = new Set()
                let tempSubject = new Set()

                allCreds.forEach((cred) => {
                    if (cred.issuer && cred.issuer.id) {
                        tempIssuers.add(cred.issuer.id)
                    } else if (cred.issuer) {
                        tempIssuers.add(cred.issuer)
                    }

                    if (cred.type && Array.isArray(cred.type)) {
                        cred.type.forEach((item) => {
                            if (item != 'VerifiableCredential') {
                                tempSubject.add(item)
                            }
                        });
                    }
                })

                let issuer = '', subject = ''
                tempIssuers.forEach(function(value) {
                    issuer = `${value}, ${issuer}`
                })

                tempSubject.forEach(function(value) {
                    subject = `${value}, ${subject}`
                })

                this.issuer = issuer
                this.subject = subject
            },
            store: async function () {
                this.errors.length = 0
                if (this.friendlyName.length == 0) {
                    this.errors.push("friendly name required.")
                    return
                }

                // Save the VC/VP
                let status = 'success'
                if (this.isVC) {
                    await this.aries.verifiable.saveCredential({
                        name: this.friendlyName,
                        verifiableCredential: JSON.stringify(this.credData)
                    }).then(() => {
                            console.log('vc save success')
                        }
                    ).catch(err => {
                        status = err.toString()
                        console.log('vc save failed : errMsg=' + err)
                    })
                } else {
                    let index = 0
                    for (let credItem of this.credData.verifiableCredential) {
                        await this.aries.verifiable.saveCredential({
                            name: `${this.friendlyName}_${getCredentialType(credItem.type)}_${++index}`,
                            verifiableCredential: JSON.stringify(credItem)
                        }).then(() => {
                                console.log('vc save success:')
                            }
                        ).catch(err => {
                            status = err.toString()
                            console.log('vc save failed : errMsg=' + err)
                        })
                    }
                }


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

