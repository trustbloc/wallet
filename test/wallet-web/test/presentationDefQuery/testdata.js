/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export const samplePresentationDefQuery = {
    "submission_requirements": [
        {
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "pick",
            "count": 1,
            "from": ["A"]
        },
        {
            "name": "Employment Information",
            "purpose": "We need to know that you are currently employed.",
            "rule": "all",
            "from": ["B"]
        },
        {
            "name": "Citizenship Information",
            "rule": "pick",
            "count": 1,
            "from": ["C"]
        }
    ],
    "input_descriptors": [
        {
            "id": "banking_input_1",
            "group": ["A"],
            "schema": {
                "uri": "https://bank-standards.com/customer.json",
                "name": "Bank Account Information",
                "purpose": "We need your bank and account information."
            },
            "constraints": {
                "fields": [
                    {
                        "path": ["$.issuer", "$.vc.issuer", "$.iss"],
                        "purpose": "The credential must be from one of the specified issuers",
                        "filter": {
                            "type": "string",
                            "pattern": "did:example:123|did:example:456"
                        }
                    },
                    {
                        "path": ["$.credentialSubject.account[*].account_number", "$.vc.credentialSubject.account[*].account_number"],
                        "purpose": "We need your bank account number for processing purposes",
                        "filter": {
                            "type": "string",
                            "minLength": 10,
                            "maxLength": 12
                        }
                    },
                    {
                        "path": ["$.credentialSubject.account[*].routing_number", "$.vc.credentialSubject.account[*].routing_number"],
                        "purpose": "You must have an account with a German, US, or Japanese bank account",
                        "filter": {
                            "type": "string",
                            "pattern": "^DE|^US|^JP"
                        }
                    }
                ]
            }
        },
        {
            "id": "banking_input_2",
            "group": ["A"],
            "schema": {
                "uri": [
                    "https://bank-schemas.org/1.0.0/accounts.json",
                    "https://bank-schemas.org/2.0.0/accounts.json"
                ],
                "name": "Bank Account Information",
                "purpose": "We need your bank and account information."
            },
            "constraints": {
                "fields": [
                    {
                        "path": ["$.issuer", "$.vc.issuer", "$.iss"],
                        "purpose": "The credential must be from one of the specified issuers",
                        "filter": {
                            "type": "string",
                            "pattern": "did:example:123|did:example:456"
                        }
                    },
                    {
                        "path": ["$.credentialSubject.account[*].id", "$.vc.credentialSubject.account[*].id"],
                        "purpose": "We need your bank account number for processing purposes",
                        "filter": {
                            "type": "string",
                            "minLength": 10,
                            "maxLength": 12
                        }
                    },
                    {
                        "path": ["$.credentialSubject.account[*].route", "$.vc.credentialSubject.account[*].route"],
                        "purpose": "You must have an account with a German, US, or Japanese bank account",
                        "filter": {
                            "type": "string",
                            "pattern": "^DE|^US|^JP"
                        }
                    }
                ]
            }
        },
        {
            "id": "employment_input",
            "group": ["B"],
            "schema": {
                "uri": "https://business-standards.org/schemas/employment-history.json",
                "name": "Employment History",
                "purpose": "We need your bank and account information."
            },
            "constraints": {
                "fields": [
                    {
                        "path": ["$.jobs[*].active"],
                        "filter": {
                            "type": "boolean",
                            "pattern": "true"
                        }
                    }
                ]
            }
        },
        {
            "id": "citizenship_input_1",
            "group": ["C"],
            "schema": {
                "uri": "https://eu.com/claims/DriversLicense.json",
                "name": "EU Driver's License"
            },
            "constraints": {
                "fields": [
                    {
                        "path": ["$.issuer", "$.vc.issuer", "$.iss"],
                        "purpose": "The credential must be from one of the specified issuers",
                        "filter": {
                            "type": "string",
                            "pattern": "did:example:gov1|did:example:gov2"
                        }
                    },
                    {
                        "path": ["$.dob"],
                        "filter": {
                            "type": "date",
                            "minimum": "1999-5-16"
                        }
                    }
                ]
            }
        },
        {
            "id": "citizenship_input_2",
            "group": ["C"],
            "schema": {
                "uri": "hub://did:foo:123/Collections/schema.us.gov/passport.json",
                "name": "US Passport"
            },
            "constraints": {
                "issuers": ["did:foo:gov3"],
                "fields": [
                    {
                        "path": ["$.birth_date"],
                        "filter": {
                            "type": "date",
                            "minimum": "1999-5-16"
                        }
                    }
                ]
            }
        }
    ]
}

export const prCardv1 = {
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
    "description": "Permanent Resident Card",
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
}

export const prCardv2 = {
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://w3id.org/citizenship/v2"
    ],
    "id": "http://example.gov/credentials/ff98f978-588f-4eb0-b17b-60c18e1dac2c",
    "type": [
        "VerifiableCredential",
        "PermanentResidentCard"
    ],
    "name": "Permanent Resident Card",
    "description": "Permanent Resident Card",
    "issuer": "did:web:example.two",
    "issuanceDate": "2019-12-03T12:19:52Z",
    "expirationDate": "2029-12-03T12:19:52Z",
    "credentialSubject": {
        "id": "did:example:b34ca6cd37bbf23",
        "type": [
            "PermanentResident",
            "Person"
        ],
        "givenName": "JOE",
        "familyName": "SMITH",
        "gender": "Male",
        "image": "data:image/png;base64,iVBORw0KGgo...kJggg==",
        "residentSince": "2015-01-01",
        "lprCategory": "C09",
        "lprNumber": "999-999-999",
        "commuterClassification": "C1",
        "birthCountry": "Bahamas",
        "birthDate": "1967-07-20"
    },
    "proof": {
        "type": "Ed25519Signature2018",
        "created": "2019-12-11T03:50:55Z",
        "verificationMethod": "did:web:example#z6MksHh7qHWvybLg5QTPPdG2DgEjjduBDArV9EF9mRiRzMBN",
        "proofPurpose": "assertionMethod",
        "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..SeUoIpwN_1Zrwc9zcl5NuvI88eJh6mWcxUMROHLrRg9Ubrz1YBhprPjcIZVE9JikK2DOO75pwC06fEwmu4GUAw"
    }
}

export const degreeCertificare = {
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://www.example.com/2020/udc-example/v1"
    ],
    "type": [
        "VerifiableCredential",
        "UniversityDegreeCredential"
    ],
    "id": "http://example.gov/credentials/ff98f978-588f-4eb0-b17b-60c18e1dac2c",
    "issuanceDate": "2020-03-16T22:37:26.544Z",
    "issuer": {
        "id": "did:web:faber.university",
        "name": "University"
    },
    "credentialSubject": {
        "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
        "degree": {
            "type": "BachelorDegree",
            "degree": "MIT"
        },
        "name": "Jayden Doe",
        "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
    },
    "proof": {
        "type": "Ed25519Signature2018",
        "created": "2019-12-11T03:50:55Z",
        "verificationMethod": "did:web:example#z6MksHh7qHWvybLg5QTPPdG2DgEjjduBDArV9EF9mRiRzMBN",
        "proofPurpose": "assertionMethod",
        "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..hckTjvJ8umMsgAdsgPfAYGKYh-4IqqvuGOX_kbNMwpkwyslFj_XKl06wgJDDMLkmnvHHEk74FDBUL_F_0mdeAA"
    }
}

export const pdCardManifestVC = {
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://trustbloc.github.io/context/vc/issuer-manifest-credential-v1.jsonld"
    ],
    "type": [
        "VerifiableCredential",
        "IssuerManifestCredential"
    ],
    "name": "Example Issuer Manifest Credential",
    "description": "List of verifiable credentials provided by example issuer",
    "id": "http://example.gov/credentials/ff98f978-588f-4eb0-b17b-60c18e1dac2c",
    "issuanceDate": "2020-03-16T22:37:26.544Z",
    "issuer": "did:factom:5d0dd58757119dd437c70d92b44fbf86627ee275f0f2146c3d99e441da342d9f",
    "credentialSubject": {
        "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
        "contexts": [
            "https://w3id.org/citizenship/v3",
            "https://w3id.org/citizenship/v4"
        ]
    }
}

export const driversLicenseVC = {
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://trustbloc.github.io/context/vc/examples/mdl-v1.jsonld"
    ],
    "type": [
        "VerifiableCredential",
        "mDL"
    ],
    "id": "http://example.gov/credentials/ff98f978-588f-4eb0-b17b-60c18e1dac2c",
    "issuanceDate": "2020-03-16T22:37:26.544Z",
    "issuer": {
        "id": "did:gov:transport_ministry"
    },
    "credentialSubject": {
        "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
        "document_number": "123-456-789",
        "given_name": "Jayden",
        "family_name": "Doe"
    },
    "proof": {
        "type": "Ed25519Signature2018",
        "created": "2019-12-11T03:50:55Z",
        "verificationMethod": "did:gov:transport_ministry#z6MksHh7qHWvybLg5QTPPdG2DgEjjduBDArV9EF9mRiRzMBN",
        "proofPurpose": "assertionMethod",
        "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..hckTjvJ8umMsgAdsgPfAYGKYh-4IqqvuGOX_kbNMwpkwyslFj_XKl06wgJDDMLkmnvHHEk74FDBUL_F_0mdeAA"
    }
}

export const driversLicenseEvidenceManifestVC = {
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://trustbloc.github.io/context/vc/issuer-manifest-credential-v1.jsonld"
    ],
    "type": [
        "VerifiableCredential",
        "IssuerManifestCredential"
    ],
    "name": "Example Issuer Manifest Credential",
    "description": "List of verifiable credentials provided by example issuer",
    "id": "http://example.gov/credentials/ff98f978-588f-4eb0-b17b-60c18e1dac2c",
    "issuanceDate": "2020-03-16T22:37:26.544Z",
    "issuer": "did:factom:5d0dd58757119dd437c70d92b44fbf86627ee275f0f2146c3d99e441da342d9f",
    "credentialSubject": {
        "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
        "contexts": ["https://trustbloc.github.io/context/vc/examples/driver-license-evidence-v1.jsonld"]
    }
}

export const samplePresentationDefQuery1 = {
    submission_requirements: [
        {
            "name": "Degree Information",
            "purpose": "We need to know if you are qualified for this job",
            "rule": "pick",
            "count": 1,
            "from": ["D"]
        },
        {
            "name": "Citizenship Information",
            "rule": "all",
            "from": ["C"]
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
                        "path": ["$.issuer.id", "$.vc.issuer.id", "$.issuer", "$.vc.issuer"],
                        "purpose": "The credential must be from one of the specified issuers",
                        "filter": {
                            "type": "string",
                            "pattern": "did:web:example.two|did:web:example.three"
                        }
                    },
                    {
                        "path": ["$.credentialSubject.commuterClassification"],
                        "filter": {
                            "type": "string",
                            "pattern": "C1"
                        }
                    }
                ]
            }
        },
        {
            "id": "degree_input_1",
            "group": ["D"],
            "schema": {
                "uri": [
                    "https://www.example.com/2020/udc-example/v1"
                ],
                "name": "University degree certificate"
            },
            "constraints": {
                "fields": [
                    {
                        "path": ["$.issuer.id", "$.vc.issuer.id", "$.issuer", "$.vc.issuer"],
                        "purpose": "The credential must be from one of the specified issuers",
                        "filter": {
                            "type": "string",
                            "pattern": "did:web:faber.university|did:web:jake.university"
                        }
                    },
                    {
                        "path": ["$.credentialSubject.degree.type"],
                        "purpose": "Should be masters of engineering",
                        "filter": {
                            "type": "string",
                            "pattern": "MastersOfEngineering"
                        }
                    }
                ]
            }
        },
        {
            "id": "degree_input_2",
            "group": ["D"],
            "schema": {
                "uri": [
                    "https://www.example.com/2020/udc-example/v1"
                ],
                "name": "University degree certificate"
            },
            "constraints": {
                "fields": [
                    {
                        "path": ["$.issuer.id", "$.vc.issuer.id", "$.issuer", "$.vc.issuer"],
                        "purpose": "The credential must be from one of the specified issuers",
                        "filter": {
                            "type": "string",
                            "pattern": "did:web:trustbloc.university|did:web:edge.university"
                        }
                    },
                    {
                        "path": ["$.credentialSubject.degree.type"],
                        "purpose": "Should be masters of science",
                        "filter": {
                            "type": "string",
                            "pattern": "MastersOfScience"
                        }
                    }
                ]
            }
        },
        {
            "id": "degree_input_3",
            "group": ["D"],
            "schema": {
                "uri": [
                    "https://www.example.com/2020/udc-example/v1"
                ],
                "name": "University degree certificate"
            },
            "constraints": {
                "fields": [
                    {
                        "path": ["$.issuer.id", "$.vc.issuer.id", "$.issuer", "$.vc.issuer"],
                        "purpose": "The credential must be from one of the specified issuers",
                        "filter": {
                            "type": "string",
                            "pattern": "did:web:faber.university|did:web:jake.university"
                        }
                    },
                    {
                        "path": ["$.credentialSubject.degree.type"],
                        "purpose": "Should be masters Degree",
                        "filter": {
                            "type": "string",
                            "pattern": "MastersDegree"
                        }
                    }
                ]
            }
        }
    ]
}
