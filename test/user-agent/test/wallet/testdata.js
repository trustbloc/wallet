/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export const studentCardToStore = {
    "@context": [
        "https://www.w3.org/2018/credentials/v1"
    ],
    type: "VerifiablePresentation",
    verifiableCredential: [
        {
            "@context": [
                "https://www.w3.org/2018/credentials/v1",
                "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld",
                "https://trustbloc.github.io/context/vc/examples-v1.jsonld"
            ],
            "credentialStatus": {
                "id": "https://issuer.sandbox.trustbloc.dev/status/2",
                "type": "CredentialStatusList2017"
            },
            "credentialSubject": {
                "email": "foo@bar.com",
                "id": "did:trustbloc:testnet.trustbloc.local:EiBLPQGlDPox0hMsG_s3xltIRJs-lls3InzSgBDI-FxZIQ",
                "name": "Foo",
                "semester": "3",
                "studentid": "1234568",
                "type": "StudentCard",
                "university": "Faber College"
            },
            "description": "Student Card for Mr.Foo",
            "id": "http://example.com/edff99e0-87c5-4aa2-bbfe-fb7dbda35c36",
            "issuanceDate": "2020-05-27T20:36:05.301078505Z",
            "issuer": {
                "id": "did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw",
                "name": "trustbloc-ed25519signature2018-ed25519"
            },
            "name": "Student Card",
            "proof": {
                "created": "2020-05-27T20:36:17.702094979Z",
                "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..mzp-n5w2PhAKFkLqLY1YiWRAhS7lpzBktRfWjju-voy1Zvz2ubQQKba3rQixly75gfTw0xqMIeipzii62XDPAA",
                "proofPurpose": "assertionMethod",
                "type": "Ed25519Signature2018",
                "verificationMethod": "did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw#hQSD86WYa-HgSRboXUEF"
            },
            "type": [
                "VerifiableCredential",
                "StudentCard"
            ]
        }
    ]
}

export const studentCardAndDegreeToStore = {
    "@context": [
        "https://www.w3.org/2018/credentials/v1"
    ],
    "holder": "did:trustbloc:testnet.trustbloc.local:EiBLPQGlDPox0hMsG_s3xltIRJs-lls3InzSgBDI-FxZIQ",
    "type": "VerifiablePresentation",
    "verifiableCredential": [
        {
            "@context": [
                "https://www.w3.org/2018/credentials/v1",
                "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld",
                "https://trustbloc.github.io/context/vc/examples-v1.jsonld"
            ],
            "credentialStatus": {
                "id": "https://issuer.sandbox.trustbloc.dev/status/2",
                "type": "CredentialStatusList2017"
            },
            "credentialSubject": {
                "email": "foo@bar.com",
                "id": "did:trustbloc:testnet.trustbloc.local:EiBLPQGlDPox0hMsG_s3xltIRJs-lls3InzSgBDI-FxZIQ",
                "name": "Foo",
                "semester": "3",
                "studentid": "1234568",
                "type": "StudentCard",
                "university": "Faber College"
            },
            "description": "Student Card for Mr.Foo",
            "id": "http://example.com/f9b4d23f-5728-495c-b843-3c955e7d6de5",
            "issuanceDate": "2020-05-28T21:16:57.780923246Z",
            "issuer": {
                "id": "did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw",
                "name": "trustbloc-ed25519signature2018-ed25519"
            },
            "name": "Student Card",
            "proof": {
                "created": "2020-05-28T21:17:03.11910197Z",
                "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..1bmYv80QTQ90-CrUew_TuvwmNb1DkvkPc_DKhbPb8ULSy0pc9-xiQvbixwDj18msX9nHrtBP0FNquXdUzmOZCg",
                "proofPurpose": "assertionMethod",
                "type": "Ed25519Signature2018",
                "verificationMethod": "did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw#hQSD86WYa-HgSRboXUEF"
            },
            "type": [
                "VerifiableCredential",
                "StudentCard"
            ]
        },
        {
            "@context": [
                "https://www.w3.org/2018/credentials/v1",
                "https://www.w3.org/2018/credentials/examples/v1",
                "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld",
                "https://trustbloc.github.io/context/vc/examples-v1.jsonld"
            ],
            "credentialStatus": {
                "id": "https://issuer.sandbox.trustbloc.dev/status/2",
                "type": "CredentialStatusList2017"
            },
            "credentialSubject": {
                "degree": {
                    "degree": "Bachelor of Science and Arts",
                    "type": "BachelorDegree"
                },
                "id": "did:trustbloc:testnet.trustbloc.local:EiBLPQGlDPox0hMsG_s3xltIRJs-lls3InzSgBDI-FxZIQ",
                "name": "Jayden Doe"
            },
            "description": "University Degree Credential for Mr.Jayden Doe",
            "id": "http://example.com/a807e4e9-c534-4adb-a7cc-b5e6d9bb00e2",
            "issuanceDate": "2020-05-28T21:15:57.090081868Z",
            "issuer": {
                "id": "did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw",
                "name": "trustbloc-ed25519signature2018-ed25519"
            },
            "name": "University Degree Credential",
            "proof": {
                "created": "2020-05-28T21:16:04.427607084Z",
                "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..4Le2zJoXE7PBI1ySXkUWlUqx0a1LQa7lULKhGnoCoXlFWRy0S9jI31N5zJS-TA6udodi8WZhgeF5jv1Gf8LFCQ",
                "proofPurpose": "assertionMethod",
                "type": "Ed25519Signature2018",
                "verificationMethod": "did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw#hQSD86WYa-HgSRboXUEF"
            },
            "type": [
                "VerifiableCredential",
                "UniversityDegreeCredential"
            ]
        }
    ]
}

export const prcAndUdcVP = {
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

export const presentationDefQuery1 =  {
    submission_requirements: [
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
