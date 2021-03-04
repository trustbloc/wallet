/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="md-layout">
            <div class="md-layout-item">
                <md-card>
                    <md-card-header style="background-color:#00bcd4">
                        <h3 class="title">
                            <md-icon>fingerprint</md-icon>
                            Wallet Operations
                        </h3>
                    </md-card-header>

                    <md-card-content>

                        <div id="#ivc">
                            <md-textarea v-model="interopData" id="vcDataTextArea" rows="30" style="width: 100%"/>
                            <br>
                            <div>
                                <b>Sample requests:</b>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('vp', 'store')">Store Presentation
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('getvp', 'get')">Request Presentation
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('bbs', 'get')">Selective Disclosure
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('pexq', 'get')">Presentation Exchange
                                    Query
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('pexq-didcomm', 'get')">Presentation
                                    Exchange Query With DIDComm
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('pexq-didcomm-govnvc', 'get')">Presentation
                                    Exchange Query With DIDComm & Governance VC
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('didauth', 'get')">DID Auth</md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('didconn', 'get')">DID Connect
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('didconn-manifest', 'get')">DID Connect
                                    with manifest
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('didconn-manifest-usrc', 'get')">DID Connect
                                    with manifest and user credential
                                </md-chip>
                                <md-chip class="request-sample" md-clickable v-on:click="prefillRequest('didconn-manifest-usrc-govvc', 'get')">DID Connect
                                    with manifest, user credential and governance VC
                                </md-chip>
                            </div>

                            <div>
                                <b>Wallet Operations:</b>
                                <br>
                                <md-button v-on:click="store" class="md-success" :disabled=disableStore>Store
                                </md-button>
                                &nbsp;
                                <md-button v-on:click="get" class="md-success" :disabled=disableGet>Get
                                </md-button>
                            </div>

                            <div v-if="responses.length" style="color: #0E9A00">
                                <p v-for="response in responses" :key="response">{{ response }}</p>
                            </div>
                            <div v-if="errors.length" style="color: #fb4934">
                                <b>Please correct the following error(s):</b>
                                <ul>
                                    <li v-for="error in errors" :key="error">{{ error }}</li>
                                </ul>
                            </div>

                        </div>
                    </md-card-content>
                </md-card>

            </div>
        </div>
    </div>
</template>

<script>
    import {WebCredential} from 'credential-handler-polyfill/WebCredential.js';
    import {getSample} from './webWalletSamples';


    export default {
        created: async function () {
            let opts = this.$store.getters.getAgentOpts
            if (!opts) {
                this.errors.push("Please login to your webwallet before running this demo")
                return
            }

            await this.$polyfill.loadOnce(opts.credentialMediatorURL)
        },
        data() {
            return {
                interopData: "",
                mode: "",
                errors: [],
                responses: []
            };
        },
        methods: {
            clearResults: async function () {
                this.errors.length = 0
                this.responses.length = 0
            },
            prefillRequest: function (id, mode) {
                this.interopData = JSON.stringify(getSample(id), null, 2)
                this.mode = mode
            },
            store: async function () {
                this.clearResults()
                if (this.interopData.length == 0) {
                    this.errors.push("Invalid presentation")
                    return
                }
                const webCredentialWrapper = new WebCredential('VerifiablePresentation', JSON.parse(this.interopData));
                const result = await navigator.credentials.store(webCredentialWrapper);
                console.log('Result received via store() request:', result);
                this.responses.push("Successfully stored verifiable presentation to wallet.")
            },
            get: async function () {
                this.clearResults()
                if (this.interopData.length == 0) {
                    this.errors.push("Invalid query")
                    return
                }
                const result = await navigator.credentials.get(JSON.parse(this.interopData));
                if (!result) {
                    this.errors.push("Failed to get result")
                    return
                }

                this.showResp(result.data)
                this.responses.push("Successfully got response from wallet.")
            },
            showResp: async function (data) {
                if ((typeof data) == "object") {
                    this.interopData = JSON.stringify(data, null, 2)
                } else {
                    this.responses.push("Warning: received unexpcted string data type")
                    this.interopData = data
                }
            },
        },
        computed: {
            disableStore() {
                return this.mode == 'get'
            },
            disableGet() {
                return this.mode == 'store'
            },
        }
    }
</script>
<style lang="css">
    .request-sample {
        margin: 2px;
    }

</style>
