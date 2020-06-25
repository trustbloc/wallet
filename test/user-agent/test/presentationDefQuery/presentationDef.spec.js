/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {expect} from 'chai'
import {PresentationExchange} from '../../../../cmd/user-agent/src/pages/chapi/wallet'
import {degreeCertificare, prCardv1, prCardv2, samplePresentationDefQuery, pdCardManifestVC} from './testdata.js'


describe('presentation definition query schema validation', () => {
    it('presentation definition successful schema query', async () => {
        let defQ = new PresentationExchange(samplePresentationDefQuery)
        expect(defQ).to.not.be.null
    })

    it('submission_requirements[*] schema validation ', async () => {
        let sample = Object.assign({}, samplePresentationDefQuery)
        sample["submission_requirements"] = [{
            "rule": {
                "type": "pick",
                "count": 1,
                "from": ["A"]
            }
        }]
        let defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": {
                "type": "all",
                "from": ["A"]
            }
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": {
                "type": "pick",
                "from": ["A"]
            }
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property '.count'")
        }
        expect(defQ).to.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": {
                "type": "test",
                "from": ["A"]
            }
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be equal to one of the allowed values")
        }
        expect(defQ).to.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
        }]
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property 'rule'")
        }
        expect(defQ).to.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "sample-rule"
        }]
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be object")
        }
        expect(defQ).to.be.null
    })

    it('input_descriptors[*] schema validation ', async () => {
        let sample = Object.assign({}, samplePresentationDefQuery)
        delete sample.submission_requirements;

        let defQ = null

        delete sample.input_descriptors;
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property 'input_descriptors'")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = null
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be array")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = []
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should NOT have fewer than 1 items")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = [{}]
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property 'id'")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = [{
            "id": "banking_input_1"
        }]
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property 'schema'")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = [{
            "id": "banking_input_1",
            "schema": {
                "uri": "https://bank-standards.com/customer.json",
                "name": "Bank Account Information",
                "purpose": "We need your bank and account information."
            }
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["input_descriptors"] = [{
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
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["input_descriptors"] = [{
            "id": "employment_input",
            "group": ["B"],
            "schema": {
                "uri": "https://business-standards.org/schemas/employment-history.json",
                "name": "Employment History",
                "purpose": "We need your bank and account information."
            },
            "constraints": {}
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["input_descriptors"] = [{
            "id": "employment_input",
            "group": ["B"],
            "schema": {
                "uri": "https://business-standards.org/schemas/employment-history.json",
                "name": "Employment History",
                "purpose": "We need your bank and account information."
            },
            "constraints": {"fields": []}
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should NOT have fewer than 1 items")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = [{
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
                        "filter": {
                            "type": "date",
                            "minimum": "1999-5-16"
                        }
                    }
                ]
            }
        }]
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property 'path'")
        }
        expect(defQ).to.be.null


        sample["input_descriptors"] = [{
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
                        "filter": "sample-filter"
                    }
                ]
            }
        }]
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be object")
        }
        expect(defQ).to.be.null
    })
})

describe('presentation definition submission requirements', () => {
    it('all submission requirements should be available in input descriptors', async () => {
        let sample = Object.assign({}, samplePresentationDefQuery)
        sample["submission_requirements"] = [{
            "rule": {
                "type": "pick",
                "count": 1,
                "from": ["A", "B", "C"]
            }
        }]

        let defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["submission_requirements"] = [{
            "rule": {
                "type": "pick",
                "count": 1,
                "from": ["A", "B", "C", "X"]
            }
        }]

        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("Couldn't find matching group in descriptors for 'submission_requirements'")
        }
        expect(defQ).to.be.null
    })
})


describe('generate presentation submission  with no submission requirements', () => {
    it('generate presentation submission with only schema rules', async () => {
        let allCreds = [prCardv1, prCardv2, degreeCertificare]

        // matching with one schema
        let presDef = {
            input_descriptors: [
                {
                    "id": "citizenship_input_1",
                    "group": ["C"],
                    "schema": {
                        "uri": "https://w3id.org/citizenship/v1",
                        "name": "US Permanent resident card"
                    },
                }
            ]
        }

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal([{
            id: 'citizenship_input_1',
            path: '$.verifiableCredential.[0]'
        }])
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv1])


        // matching with one of multiple schemas
        presDef = {
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
                }
            ]
        }

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal([
            {id: 'citizenship_input_1', path: '$.verifiableCredential.[0]'},
            {id: 'citizenship_input_1', path: '$.verifiableCredential.[1]'}
        ])
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv1, prCardv2])

        // matching with one of multiple input descriptors
        presDef = {
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
                },
                {
                    "id": "university_degree_input_1",
                    "group": ["C"],
                    "schema": {
                        "uri": [
                            "https://www.example.com/2020/udc-example/v1"
                        ],
                        "name": "University degree certificate"
                    },
                }
            ]
        }

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal([
            {id: 'citizenship_input_1', path: '$.verifiableCredential.[0]'},
            {id: 'citizenship_input_1', path: '$.verifiableCredential.[1]'},
            {id: 'university_degree_input_1', path: '$.verifiableCredential.[2]'}
        ])
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv1, prCardv2, degreeCertificare])


        // no matches
        presDef = {
            input_descriptors: [
                {
                    "id": "driving_license_input_1",
                    "group": ["C"],
                    "schema": {
                        "uri": [
                            "https://www.example.com/driving-license/v1",
                            "https://www.example.com/driving-license/v2",
                        ],
                        "name": "Driving License"
                    },
                }
            ]
        }

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.be.empty
        expect(presSubmission.verifiableCredential).to.be.empty
    })


    it('generate presentation submission with using constraints', async () => {

        // pr card 3 with different commuter classification
        var prCardv3 = JSON.parse(JSON.stringify(prCardv2))
        prCardv3.credentialSubject.commuterClassification = "C2"

        let allCreds = [prCardv1, prCardv2, prCardv3, degreeCertificare]

        let presDef = {
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
                }
            ]
        }

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal([{
            id: 'citizenship_input_1',
            path: '$.verifiableCredential.[0]'
        }])
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv2])
    })

    it('generate presentation submission using multiple descriptors with constraints', async () => {

        var mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        degreeCertificare.issuer.id = "did:web:jake.university"
        degreeCertificare.credentialSubject.degree.type = "MastersDegree"

        var secondDegree = JSON.parse(JSON.stringify(degreeCertificare))
        secondDegree.issuer.id = "did:web:fake.university"
        secondDegree.credentialSubject.degree.type = "BachelorDegree"

        let allCreds = [prCardv1, prCardv2, degreeCertificare, mastersDegree, secondDegree]

        let presDef = {
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
                                "purpose": "Should be masters or bachelors degree",
                                "filter": {
                                    "type": "string",
                                    "pattern": "BachelorDegree|MastersDegree"
                                }
                            }
                        ]
                    }
                }
            ]
        }

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "citizenship_input_1", "path": "$.verifiableCredential.[0]"},
                {"id": "degree_input_1", "path": "$.verifiableCredential.[1]"},
                {"id": "degree_input_1", "path": "$.verifiableCredential.[2]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv2, degreeCertificare, mastersDegree])
    })


})


describe('generate presentation submission with submission requirements', () => {
    it('generate presentation submission using multiple submission requirements - scenario 1', async () => {

        var mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        degreeCertificare.issuer.id = "did:web:jake.university"
        degreeCertificare.credentialSubject.degree.type = "MastersDegree"

        var secondDegree = JSON.parse(JSON.stringify(degreeCertificare))
        secondDegree.issuer.id = "did:web:fake.university"
        secondDegree.credentialSubject.degree.type = "BachelorDegree"

        var diploma = JSON.parse(JSON.stringify(degreeCertificare))
        diploma.issuer.id = "did:web:trustbloc.university"
        diploma.credentialSubject.degree.type = "PostGraduationDiploma"

        let allCreds = [prCardv1, prCardv2, degreeCertificare, mastersDegree, secondDegree, diploma]

        let presDef = {
            submission_requirements: [
                {
                    "name": "Degree Information",
                    "purpose": "We need to know if you are qualified for this job",
                    "rule": {
                        "type": "pick",
                        "count": 1,
                        "from": ["D"]
                    }
                },
                {
                    "name": "Citizenship Information",
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
                                "purpose": "Should be masters or bachelors degree",
                                "filter": {
                                    "type": "string",
                                    "pattern": "PostGraduationDiploma"
                                }
                            }
                        ]
                    }
                }
            ]
        }

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_1", "path": "$.verifiableCredential.[0]"},
                {"id": "degree_input_1", "path": "$.verifiableCredential.[1]"},
                {"id": "degree_input_2", "path": "$.verifiableCredential.[2]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential.[3]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([degreeCertificare, mastersDegree, diploma, prCardv2])
    })

    it('generate presentation submission using multiple submission requirements - scenario 2', async () => {

        var mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        degreeCertificare.issuer.id = "did:web:jake.university"
        degreeCertificare.credentialSubject.degree.type = "MastersDegree"

        var secondDegree = JSON.parse(JSON.stringify(degreeCertificare))
        secondDegree.issuer.id = "did:web:faber.university"
        secondDegree.credentialSubject.degree.type = "BachelorDegree"
        secondDegree.credentialSubject.degree.coop = "Y"

        var diploma = JSON.parse(JSON.stringify(degreeCertificare))
        diploma.issuer.id = "did:web:trustbloc.university"
        diploma.credentialSubject.degree.type = "PostGraduationDiploma"

        let allCreds = [prCardv1, prCardv2, degreeCertificare, mastersDegree, secondDegree, diploma]

        let presDef = {
            submission_requirements: [
                {
                    "name": "Degree Information",
                    "purpose": "We need to know if you are qualified for this job",
                    "rule": {
                        "type": "all",
                        "from": ["D"]
                    }
                },
                {
                    "name": "Citizenship Information",
                    "rule": {
                        "type": "pick",
                        "count": 1,
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
                                "path": ["$.credentialSubject.degree.coop"],
                                "purpose": "Should include co-op",
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

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_1", "path": "$.verifiableCredential.[0]"},
                {"id": "degree_input_2", "path": "$.verifiableCredential.[1]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential.[2]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([secondDegree, secondDegree, prCardv2])
    })

    it('generate presentation submission using multiple submission requirements and manifest credentials', async () => {

        var mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        degreeCertificare.issuer.id = "did:web:jake.university"
        degreeCertificare.credentialSubject.degree.type = "MastersDegree"

        var secondDegree = JSON.parse(JSON.stringify(degreeCertificare))
        secondDegree.issuer.id = "did:web:faber.university"
        secondDegree.credentialSubject.degree.type = "BachelorDegree"
        secondDegree.credentialSubject.degree.coop = "Y"

        var diploma = JSON.parse(JSON.stringify(degreeCertificare))
        diploma.issuer.id = "did:web:trustbloc.university"
        diploma.credentialSubject.degree.type = "PostGraduationDiploma"

        let allCreds = [prCardv1, prCardv2, pdCardManifestVC, degreeCertificare, mastersDegree, secondDegree, diploma]

        let presDef = {
            submission_requirements: [
                {
                    "name": "Degree Information",
                    "purpose": "We need to know if you are qualified for this job",
                    "rule": {
                        "type": "all",
                        "from": ["D"]
                    }
                },
                {
                    "name": "Citizenship Information",
                    "rule": {
                        "type": "pick",
                        "count": 1,
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
                            "https://w3id.org/citizenship/v3",
                            "https://w3id.org/citizenship/v4"
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
                                "path": ["$.credentialSubject.degree.coop"],
                                "purpose": "Should include co-op",
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

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_1", "path": "$.verifiableCredential.[0]"},
                {"id": "degree_input_2", "path": "$.verifiableCredential.[1]"}
            ]
        )
        expect(presSubmission.presentation_location.descriptor_map).to.deep.equal(
            [
                {"id": "citizenship_input_1", "path": "$.verifiableCredential.[2]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([secondDegree, secondDegree, pdCardManifestVC])
    })
})


describe('generate requirement details with from presentation definition', () => {
    it('get requirement details from well described definition', async () => {
        let defQ = new PresentationExchange(samplePresentationDefQuery)
        expect(defQ).to.not.be.null

        let reqDetails = defQ.requirementDetails()
        expect(reqDetails).to.have.lengthOf(3);

        let bankingInfo = reqDetails[0]
        let employmentInfo = reqDetails[1]
        let citizenshipInfo = reqDetails[2]

        expect(bankingInfo).to.deep.equal(
            {
                "name": "Banking Information",
                "purpose": "We need to know if you have an established banking history.",
                "rule": "at least 1 of each condition should be met",
                "descriptors": [
                    {
                        "name": "Bank Account Information",
                        "purpose": "We need your bank and account information.",
                        "constraints": [
                            "The credential must be from one of the specified issuers",
                            "We need your bank account number for processing purposes",
                            "You must have an account with a German, US, or Japanese bank account"
                        ]
                    },
                    {
                        "name": "Bank Account Information",
                        "purpose": "We need your bank and account information.",
                        "constraints": [
                            "The credential must be from one of the specified issuers",
                            "We need your bank account number for processing purposes",
                            "You must have an account with a German, US, or Japanese bank account"
                        ]
                    }
                ]
            }
        )

        expect(employmentInfo).to.deep.equal(
            {
                "name": "Employment Information",
                "purpose": "We need to know that you are currently employed.",
                "rule": "all conditions should be met",
                "descriptors": [
                    {
                        "name": "Employment History",
                        "purpose": "We need your bank and account information.",
                        "constraints": []
                    }
                ]
            }
        )


        expect(citizenshipInfo).to.deep.equal(
            {
                "name": "Citizenship Information",
                "purpose": "We need below information from your wallet",
                "rule": "at least 1 of each condition should be met",
                "descriptors": [
                    {
                        "name": "EU Driver's License",
                        "purpose": "",
                        "constraints": [
                            "The credential must be from one of the specified issuers"
                        ]
                    },
                    {
                        "name": "US Passport",
                        "purpose": "",
                        "constraints": []
                    }
                ]
            }
        )
    })

    it('get requirement details from definition without submission requirements', async () => {
        let query = JSON.parse(JSON.stringify(samplePresentationDefQuery))
        delete query.submission_requirements

        let defQ = new PresentationExchange(query)
        expect(defQ).to.not.be.null

        let reqDetails = defQ.requirementDetails()
        expect(reqDetails).to.deep.equal([
            {
                "name": "Requested information",
                "purpose": "We need below information from your wallet",
                "rule": "all conditions should be met",
                "descriptors": [
                    {
                        "name": "Bank Account Information",
                        "purpose": "We need your bank and account information.",
                        "constraints": [
                            "The credential must be from one of the specified issuers",
                            "We need your bank account number for processing purposes",
                            "You must have an account with a German, US, or Japanese bank account"
                        ]
                    },
                    {
                        "name": "Bank Account Information",
                        "purpose": "We need your bank and account information.",
                        "constraints": [
                            "The credential must be from one of the specified issuers",
                            "We need your bank account number for processing purposes",
                            "You must have an account with a German, US, or Japanese bank account"
                        ]
                    },
                    {
                        "name": "Employment History",
                        "purpose": "We need your bank and account information.",
                        "constraints": []
                    },
                    {
                        "name": "EU Driver's License",
                        "purpose": "",
                        "constraints": [
                            "The credential must be from one of the specified issuers"
                        ]
                    },
                    {
                        "name": "US Passport",
                        "purpose": "",
                        "constraints": []
                    }
                ]
            }
        ])
    })

    it('get requirement details from definition without name and purpose', async () => {
        let query = JSON.parse(JSON.stringify(samplePresentationDefQuery))

        query.submission_requirements = [
            {
                "rule": {
                    "type": "all",
                    "from": ["B"]
                }
            },
            {
                "rule": {
                    "type": "pick",
                    "count": 1,
                    "from": ["C"]
                }
            }
        ]

        query.input_descriptors = [
            {
                "id": "employment_input",
                "group": ["B"],
                "schema": {
                    "uri": "https://business-standards.org/schemas/employment-history.json"
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
                    "uri": "https://eu.com/claims/DriversLicense.json"
                },
                "constraints": {
                    "fields": [
                        {
                            "path": ["$.issuer", "$.vc.issuer", "$.iss"],
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
                    "uri": "hub://did:foo:123/Collections/schema.us.gov/passport.json"
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

        let defQ = new PresentationExchange(query)
        expect(defQ).to.not.be.null

        let reqDetails = defQ.requirementDetails()
        expect(reqDetails).to.deep.equal([
            {
                "name": "Requested information #1",
                "purpose": "We need below information from your wallet",
                "rule": "all conditions should be met",
                "descriptors": [
                    {
                        "name": "Condition details are not provided in request",
                        "purpose": "",
                        "constraints": []
                    }
                ]
            },
            {
                "name": "Requested information #2",
                "purpose": "We need below information from your wallet",
                "rule": "at least 1 of each condition should be met",
                "descriptors": [
                    {
                        "name": "Condition details are not provided in request",
                        "purpose": "",
                        "constraints": []
                    },
                    {
                        "name": "Condition details are not provided in request",
                        "purpose": "",
                        "constraints": []
                    }
                ]
            }
        ])
    })
})




