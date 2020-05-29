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