/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {expect} from 'chai'
import {PresentationExchange} from '../../../../cmd/wallet-web/src/pages/chapi/wallet'
import {driversLicenseVC, driversLicenseEvidenceManifestVC, degreeCertificare, prCardv1, prCardv2, pdCardManifestVC, samplePresentationDefQuery, samplePresentationDefQuery1} from './testdata.js'


describe('presentation definition query schema validation', () => {
    it('presentation definition successful schema query', async () => {
        let defQ = new PresentationExchange(samplePresentationDefQuery)
        expect(defQ).to.not.be.null
    })

    it('submission_requirements[*] schema validations ', async () => {
        // pick rule
        let sample = Object.assign({}, samplePresentationDefQuery)
        sample["submission_requirements"] = [{
            "rule": "pick",
            "count": 1,
            "from": ["A"]
        }]
        let defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        // all rule
        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "all",
            "from": ["A"]
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        // pick rule count should be greater than zero
        sample["submission_requirements"] = [{
            "rule": "pick",
            "count": 0,
            "from": ["A"]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be >= 1")
        }
        expect(defQ).to.be.null

        // count is required for pick rule
        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "pick",
            "from": ["A"]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property '.count'")
        }
        expect(defQ).to.be.null

        // submission rule should be either 'all' or 'pick'
        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "test",
            "from": ["A"]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be equal to one of the allowed values")
        }
        expect(defQ).to.be.null

        // submission rule is required
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

        // submission rule 'from or from_nested' is required
        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "all"
        }]
        try {
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property '.from'")
        }
        expect(defQ).to.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "all",
            "from_nested": [
                {
                    "name": "Banking Information",
                    "purpose": "We need to know if you have an established banking history.",
                    "rule": "all",
                    "from": ["A"]
                }
            ]
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        // submission rule 'from and from_nested' both shouldn't be present
        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "all",
            "from": ["A"],
            "from_nested": [
                {
                    "name": "Banking Information",
                    "purpose": "We need to know if you have an established banking history.",
                    "rule": "all",
                    "from": ["A"]
                }
            ]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should match exactly one schema in oneOf")
        }
        expect(defQ).to.be.null

        // submission rule 'from_nested' should be 'submission rule' type
        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": "all",
            "from": ["A"],
            "from_nested": [
                {
                    "name": "Banking Information",
                    "purpose": "We need to know if you have an established banking history.",
                }
            ]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property 'rule'")
        }
        expect(defQ).to.be.null


        // one of 'count, min, max' is required when rule is 'pick'
        sample["submission_requirements"] = [{
            "rule": "pick",
            "min": 0,
            "from": ["A"]
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["submission_requirements"] = [{
            "rule": "pick",
            "max": 1,
            "from": ["A"]
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["submission_requirements"] = [{
            "rule": "pick",
            "count": 3,
            "from": ["A"]
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        // one of 'min >= 0' when rule is 'pick'
        sample["submission_requirements"] = [{
            "rule": "pick",
            "min": -1,
            "from": ["A"]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be >= 0")
        }
        expect(defQ).to.be.null

        // one of 'max > 0' when rule is 'pick'
        sample["submission_requirements"] = [{
            "rule": "pick",
            "max": 0,
            "from": ["A"]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be >= 1")
        }
        expect(defQ).to.be.null

        // one of 'max > min' when rule is 'pick' and 'min' is present
        sample["submission_requirements"] = [{
            "rule": "pick",
            "min": 3,
            "max": 3,
            "from": ["A"]
        }]
        try {
            defQ = null
            defQ = new PresentationExchange(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be > 1")
        }
        expect(defQ).to.be.null

        sample["submission_requirements"] = [{
            "rule": "pick",
            "min": 3,
            "max": 7,
            "from": ["A"]
        }]
        defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null
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
            "rule": "pick",
            "count": 1,
            "from": ["A", "B", "C"]
        }]

        let defQ = new PresentationExchange(sample)
        expect(defQ).to.not.be.null

        sample["submission_requirements"] = [{
            "rule": "pick",
            "count": 1,
            "from": ["A", "B", "C", "X"]
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
            path: '$.verifiableCredential[0]'
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
            {id: 'citizenship_input_1', path: '$.verifiableCredential[0]'},
            {id: 'citizenship_input_1', path: '$.verifiableCredential[1]'}
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
            {id: 'citizenship_input_1', path: '$.verifiableCredential[0]'},
            {id: 'citizenship_input_1', path: '$.verifiableCredential[1]'},
            {id: 'university_degree_input_1', path: '$.verifiableCredential[2]'}
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
            path: '$.verifiableCredential[0]'
        }])
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv2])
    })

    it('generate presentation submission using multiple descriptors with constraints', async () => {

        var mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        mastersDegree.issuer.id = "did:web:jake.university"
        mastersDegree.credentialSubject.degree.type = "MastersDegree"

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
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[0]"},
                {"id": "degree_input_1", "path": "$.verifiableCredential[1]"},
                {"id": "degree_input_1", "path": "$.verifiableCredential[2]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv2, degreeCertificare, mastersDegree])
    })
})


describe('generate presentation submission with submission requirements', () => {
    it('generate presentation submission using multiple submission requirements - scenario 1', async () => {

        var mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        mastersDegree.issuer.id = "did:web:jake.university"
        mastersDegree.credentialSubject.degree.type = "MastersDegree"

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
                {"id": "degree_input_1", "path": "$.verifiableCredential[0]"},
                {"id": "degree_input_1", "path": "$.verifiableCredential[1]"},
                {"id": "degree_input_2", "path": "$.verifiableCredential[2]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[3]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([degreeCertificare, mastersDegree, diploma, prCardv2])
    })

    it('generate presentation submission using multiple submission requirements - scenario 2', async () => {

        var mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        mastersDegree.issuer.id = "did:web:jake.university"
        mastersDegree.credentialSubject.degree.type = "MastersDegree"

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
                    "rule": "all",
                    "from": ["D"]
                },
                {
                    "name": "Citizenship Information",
                    "rule": "pick",
                    "count": 1,
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
                {"id": "degree_input_1", "path": "$.verifiableCredential[0]"},
                {"id": "degree_input_2", "path": "$.verifiableCredential[1]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[2]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([secondDegree, secondDegree, prCardv2])
    })

    it('generate presentation submission using multiple submission requirements and manifest credentials', async () => {

        let mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        mastersDegree.issuer.id = "did:web:jake.university"
        mastersDegree.credentialSubject.degree.type = "MastersDegree"

        let secondDegree = JSON.parse(JSON.stringify(degreeCertificare))
        secondDegree.issuer.id = "did:web:faber.university"
        secondDegree.credentialSubject.degree.type = "BachelorDegree"
        secondDegree.credentialSubject.degree.coop = "Y"

        let diploma = JSON.parse(JSON.stringify(degreeCertificare))
        diploma.issuer.id = "did:web:trustbloc.university"
        diploma.credentialSubject.degree.type = "PostGraduationDiploma"

        let allCreds = [prCardv1, prCardv2, degreeCertificare, mastersDegree, secondDegree, diploma]

        let presDef = {
            submission_requirements: [
                {
                    "name": "Degree Information",
                    "purpose": "We need to know if you are qualified for this job",
                    "rule": "all",
                    "from": ["D"]
                },
                {
                    "name": "Citizenship Information",
                    "rule": "pick",
                    "count": 1,
                    "from": ["C"]
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

        let presSubmission = defQ.createPresentationSubmission(allCreds, [pdCardManifestVC])
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_1", "path": "$.verifiableCredential[0]"},
                {"id": "degree_input_2", "path": "$.verifiableCredential[1]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[2]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([secondDegree, secondDegree, pdCardManifestVC])
    })

    it('generate presentation submission from no submission requirements, using manifest credentials, and the schema uri is referenced from a field constraint', async () => {

        let allCreds = [driversLicenseVC]

        let presDef = {
            input_descriptors: [
                {
                    "id": "driver_license_1",
                    "schema": {
                        "uri": ["https://trustbloc.github.io/context/vc/examples/mdl-v1.jsonld"],
                        "name": "Your driver's license."
                    }
                },
                {
                    "id": "driver_license_evidence_1",
                    "schema": {
                        "uri": ["https://trustbloc.github.io/context/vc/authorization-credential-v1.jsonld"],
                        "name": "Supporting evidence of your driver's license."
                    },
                    "constraints": {
                        "fields": [
                            {
                                "path": ["$.credentialSubject.scope[*].schema.uri"],
                                "filter": {
                                    "const": "https://trustbloc.github.io/context/vc/examples/driver-license-evidence-v1.jsonld"
                                }
                            }
                        ]
                    }
                }
            ]
        }

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds, [driversLicenseEvidenceManifestVC])
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.verifiableCredential).to.deep.equal([driversLicenseVC, driversLicenseEvidenceManifestVC])
    })

    it('generate presentation submission using multiple submission requirements - all "pick" rule scenarios', async () => {

        let mastersDegree = JSON.parse(JSON.stringify(degreeCertificare))
        mastersDegree.id = "sampleID"
        mastersDegree.issuer.id = "did:web:jake.university"
        mastersDegree.credentialSubject.degree.type = "MastersDegree"

        let collegeDiploma = JSON.parse(JSON.stringify(degreeCertificare))
        collegeDiploma.issuer.id = "did:web:fake.university"
        collegeDiploma.credentialSubject.degree.type = "CollegeDiploma"

        let postGradDiploma = JSON.parse(JSON.stringify(degreeCertificare))
        postGradDiploma.issuer.id = "did:web:trustbloc.university"
        postGradDiploma.credentialSubject.degree.type = "PostGraduationDiploma"

        let allCreds = [prCardv1, prCardv2, degreeCertificare, mastersDegree, collegeDiploma, postGradDiploma]
        let presDef = JSON.parse(JSON.stringify(samplePresentationDefQuery1))

        let defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        let presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_3", "path": "$.verifiableCredential[0]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[1]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([mastersDegree, prCardv2])

        // set count to 2
        presDef["submission_requirements"] = [
            {
                "name": "Degree Information",
                "purpose": "We need to know if you are qualified for this job",
                "rule": "pick",
                "count": 2,
                "from": ["D"]
            },
            {
                "name": "Citizenship Information",
                "rule": "all",
                "from": ["C"]
            }
        ]

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[0]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv2])

        // set min to 1
        presDef["submission_requirements"] = [
            {
                "name": "Degree Information",
                "purpose": "We need to know if you are qualified for this job",
                "rule": "pick",
                "min": 1,
                "from": ["D"]
            },
            {
                "name": "Citizenship Information",
                "rule": "all",
                "from": ["C"]
            }
        ]

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_3", "path": "$.verifiableCredential[0]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[1]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([mastersDegree, prCardv2])


        // set min to 2
        presDef["submission_requirements"] = [
            {
                "name": "Degree Information",
                "purpose": "We need to know if you are qualified for this job",
                "rule": "pick",
                "min": 2,
                "from": ["D"]
            },
            {
                "name": "Citizenship Information",
                "rule": "all",
                "from": ["C"]
            }
        ]

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[0]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv2])

        // set max to 1
        presDef["submission_requirements"] = [
            {
                "name": "Degree Information",
                "purpose": "We need to know if you are qualified for this job",
                "rule": "pick",
                "max": 1,
                "from": ["D"]
            },
            {
                "name": "Citizenship Information",
                "rule": "all",
                "from": ["C"]
            }
        ]

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_3", "path": "$.verifiableCredential[0]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[1]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([mastersDegree, prCardv2])

        // set max to 2
        presDef["submission_requirements"] = [
            {
                "name": "Degree Information",
                "purpose": "We need to know if you are qualified for this job",
                "rule": "pick",
                "max": 1,
                "from": ["D"]
            },
            {
                "name": "Citizenship Information",
                "rule": "all",
                "from": ["C"]
            }
        ]

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "degree_input_3", "path": "$.verifiableCredential[0]"},
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[1]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([mastersDegree, prCardv2])

        // set max to 1, but there are more than 1 match
        presDef["submission_requirements"] = [
            {
                "name": "Degree Information",
                "purpose": "We need to know if you are qualified for this job",
                "rule": "pick",
                "max": 1,
                "from": ["D"]
            },
            {
                "name": "Citizenship Information",
                "rule": "all",
                "from": ["C"]
            }
        ]

        presDef['input_descriptors'].push({
            "id": "degree_input_4",
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
                        "path": ["$.credentialSubject.degree.degree"],
                        "purpose": "Any degree from MIT",
                        "filter": {
                            "type": "string",
                            "pattern": "MIT"
                        }
                    }
                ]
            }
        })
        presDef['input_descriptors'].push({
            "id": "degree_input_5",
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
                        "purpose": "Any degree from MIT",
                        "filter": {
                            "type": "string",
                            "pattern": "BachelorDegree"
                        }
                    }
                ]
            }
        })

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {"id": "citizenship_input_1", "path": "$.verifiableCredential[0]"}
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([prCardv2])


        // set max to 2 & min to 1
        presDef["submission_requirements"] = [
            {
                "name": "Degree Information",
                "purpose": "We need to know if you are qualified for this job",
                "rule": "pick",
                "min": 1,
                "max": 2,
                "from": ["D"]
            },
            {
                "name": "Citizenship Information",
                "rule": "all",
                "from": ["C"]
            }
        ]

        defQ = new PresentationExchange(presDef)
        expect(defQ).to.not.be.null

        presSubmission = defQ.createPresentationSubmission(allCreds)
        expect(presSubmission).to.not.be.null
        expect(presSubmission.type).to.deep.equal(["VerifiablePresentation", "PresentationSubmission"])
        expect(presSubmission.presentation_submission).to.not.be.empty
        expect(presSubmission.presentation_submission.descriptor_map).to.deep.equal(
            [
                {
                    "id": "degree_input_4",
                    "path": "$.verifiableCredential[0]"
                },
                {
                    "id": "degree_input_5",
                    "path": "$.verifiableCredential[1]"
                },
                {
                    "id": "degree_input_3",
                    "path": "$.verifiableCredential[2]"
                },
                {
                    "id": "degree_input_4",
                    "path": "$.verifiableCredential[3]"
                },
                {
                    "id": "citizenship_input_1",
                    "path": "$.verifiableCredential[4]"
                }
            ]
        )
        expect(presSubmission.verifiableCredential).to.deep.equal([degreeCertificare, degreeCertificare, mastersDegree, mastersDegree, prCardv2])
    })
})


describe('generate requirement details from presentation definition', () => {
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
                "rule": "at least 1 condition(s) should be met",
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
                "rule": "at least 1 condition(s) should be met",
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
                    "rule": "all",
                    "from": ["B"]
            },
            {
                    "rule": "pick",
                    "count": 1,
                    "from": ["C"]
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
                "rule": "at least 1 condition(s) should be met",
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




