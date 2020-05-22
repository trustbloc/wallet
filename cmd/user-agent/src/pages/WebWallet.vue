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
                        <div v-if="errors.length">
                            <b>Please correct the following error(s):</b>
                            <ul>
                                <li v-for="error in errors" :key="error">{{ error }}</li>
                            </ul>
                        </div>
                        <div v-if="responses.length">
                            <p v-for="response in responses" :key="response">{{ response }}</p>
                        </div>
                        <div id="#ivc">
                            <md-textarea v-model="interopData" id="vcDataTextArea" cols="150" rows="30"/>
                            <br>
                            <md-button v-on:click="sendVC" class="md-raised md-success">Store VC in Wallet</md-button>
                            &nbsp;
                            <md-button v-on:click="sendVP" class="md-raised md-success">Store VP in Wallet</md-button>
                            &nbsp;
                            <md-button v-on:click="getVP" class="md-raised md-success">Get VP from Wallet</md-button>
                            &nbsp;
                            <md-button v-on:click="didAuth" class="md-raised md-success">Authenticate</md-button> &nbsp;
                        </div>
                    </md-card-content>
                </md-card>

            </div>
        </div>
    </div>
</template>

<script>
    import {WebCredential} from 'credential-handler-polyfill/WebCredential.js';

    const sampleVC = {
        "@context": ["https://www.w3.org/2018/credentials/v1", "https://www.w3.org/2018/credentials/examples/v1"],
        "credentialStatus": {
            "id": "http://issuer.vc.rest.example.com:8070/status/1",
            "type": "CredentialStatusList2017"
        },
        "credentialSubject": {
            "degree": {"degree": "MIT", "type": "BachelorDegree"},
            "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
            "name": "Jayden Doe",
            "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
        },
        "id": "https://example.com/credentials/b33e4c8c-3cfc-4b7c-96d9-8d875b9a661a",
        "issuanceDate": "2020-03-16T22:37:26.544Z",
        "issuer": {
            "id": "did:trustbloc:testnet.trustbloc.local:EiABBmUZ7Jjp-mlxWJInqp3Ak2v82QQtCdIUS5KSTNGq9Q==",
            "name": "myprofile_ud1"
        },
        "proof": {
            "created": "2020-04-17T04:17:48Z",
            "proofPurpose": "assertionMethod",
            "proofValue": "CAQJKqd0MELydkNdPh7TIwgKhcMt_ypQd8AUdTbFUU4VVQVpPhEZLjg1U-1lBJyluRejsNbHZCJDRptPkBuqAQ",
            "type": "Ed25519Signature2018",
            "verificationMethod": "did:trustbloc:testnet.trustbloc.local:EiABBmUZ7Jjp-mlxWJInqp3Ak2v82QQtCdIUS5KSTNGq9Q==#key-1"
        },
        "type": ["VerifiableCredential", "UniversityDegreeCredential"]
    }

    export default {
        beforeCreate: async function () {
            this.aries = await this.$arieslib
            this.$polyfill.loadOnce()
        },
        data() {
            return {
                interopData: JSON.stringify(sampleVC, 0, 2),
                errors: [],
                responses: []
            };
        },
        methods: {
            clearResults: async function () {
                this.errors.length = 0
                this.responses.length = 0
            },
            sendVC: async function () {
                this.clearResults()
                if (this.interopData.length == 0) {
                    this.errors.push("Invalid VC")
                    return
                }
                const webCredentialWrapper = new WebCredential('VerifiableCredential', JSON.parse(this.interopData));
                const result = await navigator.credentials.store(webCredentialWrapper);
                console.log('Result received via store() request:', result);
                this.responses.push("Successfully stored verifiable credential to wallet.")
            },
            sendVP: async function () {
                this.clearResults()
                if (this.interopData.length == 0) {
                    this.errors.push("Invalid VC")
                    return
                }
                const webCredentialWrapper = new WebCredential('VerifiablePresentation', JSON.parse(this.interopData));
                const result = await navigator.credentials.store(webCredentialWrapper);
                console.log('Result received via store() request:', result);
                this.responses.push("Successfully stored verifiable presentation to wallet.")
            },
            getVC: async function () {
                this.clearResults()
                const credentialQuery = {
                    "web": {
                        "VerifiableCredential": {
                            query: [
                                {
                                    type: "QueryByExample",
                                    credentialQuery: {
                                        reason: "Please present a credential for JaneDoe.",
                                        example: {
                                            "@context": ["https://www.w3.org/2018/credentials/v1", "https://www.w3.org/2018/credentials/examples/v1"],
                                            type: ["UniversityDegreeCredential"]
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
                const result = await navigator.credentials.get(credentialQuery);
                if (!result) {
                    this.errors.push("Failed to get result")
                    return
                }
                this.showResp(result.data)
                this.responses.push("Successfully got verifiable credential from wallet.")
            },
            getVP: async function () {
                this.clearResults()
                const credentialQuery = {
                    web: {
                        VerifiablePresentation: {
                            query: [
                                {
                                    type: "QueryByExample",
                                    credentialQuery: {
                                        reason: "Please present a credential for JaneDoe.",
                                        example: {
                                            "@context": ["https://www.w3.org/2018/credentials/v1", "https://www.w3.org/2018/credentials/examples/v1"],
                                            type: ["UniversityDegreeCredential"]
                                        }
                                    }
                                }
                            ],
                            challenge: "",
                            domain: ""
                        }
                    }
                }
                const result = await navigator.credentials.get(credentialQuery);
                if (!result) {
                    this.errors.push("Failed to get result")
                    return
                }
                this.showResp(result.data)
                this.responses.push("Successfully got verifiable presentation from wallet.")
            },
            didAuth: async function () {
                this.clearResults()
                const credentialQuery = JSON.parse('{"web": {"VerifiablePresentation": {"query":{"type":"DIDAuth"},' +
                    '"challenge":"54f3da1a-d1af-4c25-b1a6-90315dda62fc","domain":"issuer.interop.example.com"}}}');
                const result = await navigator.credentials.get(credentialQuery);
                if (!result) {
                    this.errors.push("Failed to get result")
                    return
                }

                this.showResp(result.data)
                this.responses.push("Successfully got DID authorization from wallet.")
            },
            showResp: async function (data) {
                if ((typeof data) == "object") {
                    this.interopData = JSON.stringify(data, null, 2)
                } else {
                    this.responses.push("Warning: received unexpcted string data type")
                    this.interopData = data
                }
            }
        }
    }
</script>

