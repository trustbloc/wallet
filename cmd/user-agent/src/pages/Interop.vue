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
                            <md-textarea v-model="vcdata" id="vcDataTextArea" cols="100" rows="20"/>
                            <br>
                            <md-button v-on:click="sendVC" class="md-raised md-success">Store in Wallet</md-button>
                            <md-button v-on:click="getVC" class="md-raised md-success">Get from Wallet</md-button>
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
            window.$aries = await this.$arieslib
            this.$polyfill.loadOnce()
        },
        data() {
            return {
                vcdata: JSON.stringify(sampleVC),
                errors: [],
                responses: []
            };
        },
        methods: {
            clearResults: async function() {
                this.errors.length = 0
                this.responses.length = 0
            },
            sendVC: async function () {
                this.clearResults()
                if (this.vcdata.length == 0) {
                    this.errors.push("Invalid VC")
                    return
                }
                const credentialType = 'VerifiableCredential';
                const webCredentialWrapper = new WebCredential(credentialType, this.vcdata);
                const result = await navigator.credentials.store(webCredentialWrapper);
                console.log('Result of receiving via store() request:', result);
                this.responses.push("Successfully stored verifiable credential to wallet.")
            },
            getVC: async function () {
                this.clearResults()
                const credentialQuery = JSON.parse('{"web": {"VerifiablePresentation": {}}}');
                const result = await navigator.credentials.get(credentialQuery);
                if (!result) {
                    this.errors.push("Failed to get result")
                    return
                }
                this.vcdata = result.data
                this.responses.push("Successfully got verifiable presentation from wallet.")
            }
        }
    }
</script>

