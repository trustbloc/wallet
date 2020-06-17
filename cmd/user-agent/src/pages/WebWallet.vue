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
                            <br><br>
                            <fieldset>
                                <legend style="font-size: x-large;font-weight: bold">Wallet Operations</legend>

                                <md-button v-on:click="sendVC" class="md-raised md-success">Store VC in Wallet
                                </md-button>
                                &nbsp;
                                <md-button v-on:click="sendVP" class="md-raised md-success">Store VP in Wallet
                                </md-button>
                                &nbsp;
                                <md-button v-on:click="getCredential" class="md-raised md-success">Get Credential from
                                    Wallet
                                </md-button>

                                <md-button v-on:click="getVP" class="md-raised md-success">Query Credentials from Wallet
                                </md-button>
                                &nbsp;
                                <md-button v-on:click="didAuth" class="md-raised md-success">Authenticate Wallet
                                </md-button> &nbsp;

                                <md-button v-on:click="didConnect" class="md-raised md-success">Connect to Wallet
                                </md-button> &nbsp;
                            </fieldset>

                            <br>
                            <fieldset>
                                <legend style="font-size: x-large;font-weight: bold">Other Miscellaneous Operations
                                </legend>
                                <md-button v-on:click="validateSchema" class="md-raised md-success">Validate
                                    Presentation
                                    Definition
                                </md-button> &nbsp;
                            </fieldset>

                        </div>
                    </md-card-content>
                </md-card>

            </div>
        </div>
    </div>
</template>

<script>
    import {WebCredential} from 'credential-handler-polyfill/WebCredential.js';
    import {PresentationExchange} from "./chapi/wallet"

    const sampleVC = {
        "@context": [
            "https://www.w3.org/2018/credentials/v1"
        ],
        type: "VerifiablePresentation",
        verifiableCredential: [
            {
                "@context": [
                    "https://www.w3.org/2018/credentials/v1",
                    "https://w3id.org/citizenship/v1"
                ],
                "id": "http://example.gov/credentials/ff98f978-588f-4eb0-b17b-60c18e1dac2c",
                "type": [
                    "VerifiableCredential",
                    "PermanentResidentCard"
                ],
                "name": "Permanent Resident Card",
                "description": "Permanent Resident Card of Mr.John Smith",
                "issuer": {
                    "id": "did:web:example.world",
                    "name": "Border Service, NY"
                },
                "issuanceDate": "2019-12-03T12:19:52Z",
                "expirationDate": "2029-12-03T12:19:52Z",
                "credentialSubject": {
                    "id": "did:example:b56ca6cd37bbf23",
                    "type": [
                        "PermanentResident",
                        "Person"
                    ],
                    "givenName": "JOHN",
                    "familyName": "SMITH",
                    "gender": "Male",
                    "image": "data:image/png;base64,iVBORw0KGgo...kJggg==",
                    "residentSince": "2015-01-01",
                    "lprCategory": "C09",
                    "lprNumber": "999-999-999",
                    "commuterClassification": "C1",
                    "birthCountry": "Bahamas",
                    "birthDate": "1958-07-17"
                },
                "proof": {
                    "type": "Ed25519Signature2018",
                    "created": "2019-12-11T03:50:55Z",
                    "verificationMethod": "did:web:example#z6MksHh7qHWvybLg5QTPPdG2DgEjjduBDArV9EF9mRiRzMBN",
                    "proofPurpose": "assertionMethod",
                    "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..SeUoIpwN_1Zrwc9zcl5NuvI88eJh6mWcxUMROHLrRg9Ubrz1YBhprPjcIZVE9JikK2DOO75pwC06fEwmu4GUAw"
                }
            },
            {
                "@context": [
                    "https://www.w3.org/2018/credentials/v1",
                    "https://www.w3.org/2018/credentials/examples/v1",
                    "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld"
                ],
                "id": "http://example.gov/credentials/3732",
                "type": ["VerifiableCredential", "UniversityDegreeCredential"],
                "name": "Bachelor Degree",
                "description": "Bachelor of Science and Arts of Mr.John Smith",
                "issuer": "did:key:z6MkjRagNiMu91DduvCvgEsqLZDVzrJzFrwahc4tXLt9DoHd",
                "issuanceDate": "2020-03-16T22:37:26.544Z",
                "credentialSubject": {
                    "id": "did:key:z6MkjRagNiMu91DduvCvgEsqLZDVzrJzFrwahc4tXLt9DoHd",
                    "degree": {"type": "BachelorDegree", "name": "Bachelor of Science and Arts"}
                },
                "proof": {
                    "type": "Ed25519Signature2018",
                    "created": "2020-03-16T22:37:26Z",
                    "verificationMethod": "did:key:z6MkjRagNiMu91DduvCvgEsqLZDVzrJzFrwahc4tXLt9DoHd#z6MkjRagNiMu91DduvCvgEsqLZDVzrJzFrwahc4tXLt9DoHd",
                    "proofPurpose": "assertionMethod",
                    "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..7gJwYBvJuXYrFa_hpuWxknm3R5Czas_NDL-Bh7LnURA1PwjH0uBqMy4W4pgYeat3xYa12gZBkmIR0VmgY3qQCw"
                }
            }
        ]
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
            getCredential: async function () {
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
                this.responses.push("Successfully got credential from wallet.")
            },
            getVP: async function () {
                await this.clearResults()
                const credentialQuery = {
                    web: {
                        VerifiablePresentation: {
                            query: [
                                {
                                    type: "PresentationDefinitionQuery",
                                    presentationDefinitionQuery:
                                        {
                                            "submission_requirements": [
                                                {
                                                    "name": "Education Qualification",
                                                    "purpose": "We need to know if you are qualified for this job",
                                                    "rule": {
                                                        "type": "pick",
                                                        "count": 1,
                                                        "from": ["E"]
                                                    }
                                                },
                                                {
                                                    "name": "Citizenship Information",
                                                    "purpose": "You must be legally allowed to work in United States",
                                                    "rule": {
                                                        "type": "all",
                                                        "from": ["C"]
                                                    }
                                                }
                                            ],

                                            input_descriptors: [
                                                {
                                                    "id": "citizenship_input_1",
                                                    "group": ["C"],
                                                    "schema": {
                                                        "uri": [
                                                            "https://w3id.org/citizenship/v1",
                                                            "https://w3id.org/citizenship/v2"
                                                        ],
                                                        "name": "US Permanent resident card"
                                                    },
                                                    "constraints": {
                                                        "fields": [
                                                            {
                                                                "path": ["$.credentialSubject.lprCategory"],
                                                                "filter": {
                                                                    "type": "string",
                                                                    "pattern": "C09|C52|C57"
                                                                }
                                                            }
                                                        ]
                                                    }
                                                },
                                                {
                                                    "id": "degree_input_1",
                                                    "group": ["E"],
                                                    "schema": {
                                                        "uri": [
                                                            "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld"
                                                        ],
                                                        "name": "University degree certificate",
                                                        "purpose": "We need your education qualification details."
                                                    },
                                                    "constraints": {
                                                        "fields": [
                                                            {
                                                                "path": ["$.credentialSubject.degree.type"],
                                                                "purpose": "Should be masters or bachelors degree",
                                                                "filter": {
                                                                    "type": "string",
                                                                    "pattern": "BachelorDegree|MastersDegree"
                                                                }
                                                            }
                                                        ]
                                                    }
                                                },
                                                {
                                                    "id": "degree_input_2",
                                                    "group": ["E"],
                                                    "schema": {
                                                        "uri": [
                                                            "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld"
                                                        ],
                                                        "name": "Diploma certificate",
                                                        "purpose": "We need your education qualification details."
                                                    },
                                                    "constraints": {
                                                        "fields": [
                                                            {
                                                                "path": ["$.credentialSubject.degree.type"],
                                                                "purpose": "Should have valid diploma",
                                                                "filter": {
                                                                    "type": "string",
                                                                    "pattern": "Diploma"
                                                                }
                                                            },
                                                            {
                                                                "path": ["$.credentialSubject.degree.coop"],
                                                                "purpose": "Should have co-op experience",
                                                                "filter": {
                                                                    "type": "string",
                                                                    "pattern": "Y"
                                                                }
                                                            }
                                                        ]
                                                    }
                                                }
                                            ]
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
            didConnect: async function () {
                this.clearResults()
                if (this.interopData.length == 0) {
                    this.errors.push("Invalid invitation")
                    return
                }


                let invitation = JSON.parse(this.interopData)
                if (invitation['@type'] != 'https://didcomm.org/didexchange/1.0/invitation') {
                    this.errors.push("Invalid invitation, expecting did comm invitation")
                    return
                }


                const connectionRequest = {
                    web: {
                        VerifiablePresentation: {
                            query: {type: "DIDConnect"},
                            invitation: JSON.parse(this.interopData),
                            challenge: "54f3da1a-d1af-4c25-b1a6-90315dda62fc",
                            domain: "issuer.interop.example.com"
                        }
                    }
                };
                const result = await navigator.credentials.get(connectionRequest);
                if (!result) {
                    this.errors.push("Failed to get result")
                    return
                }

                this.showResp(result.data)
                this.responses.push("Successfully got DID connection response from wallet.")
            },
            showResp: async function (data) {
                if ((typeof data) == "object") {
                    this.interopData = JSON.stringify(data, null, 2)
                } else {
                    this.responses.push("Warning: received unexpcted string data type")
                    this.interopData = data
                }
            },
            validateSchema: async function () {
                this.clearResults()
                if (this.interopData.length == 0) {
                    this.errors.push("Invalid presentation def")
                    return
                }


                let presDef = JSON.parse(this.interopData)

                try {
                    new PresentationExchange(presDef)
                } catch (e) {
                    if (Array.isArray(e)) {
                        const err = this.errors
                        e.forEach(function (p) {
                            err.push(p)
                        })

                        return
                    }
                    this.errors.push(e.toString())
                    return
                }


                this.responses.push("Successfully validated sample JSON")
            },
        }
    }
</script>

