/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {expect} from 'chai'
import {PresentationDefinition} from '../../../../cmd/user-agent/src/pages/chapi/wallet'
import {samplePresentationDefQuery} from './testdata.js'


describe('presentation definition query schema validation', () => {
    it('presentation definition successful schema query', async () => {
        let defQ = new PresentationDefinition(samplePresentationDefQuery)
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
        let defQ = new PresentationDefinition(sample)
        expect(defQ).to.not.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": {
                "type": "all",
                "from": ["A"]
            }
        }]
        defQ = new PresentationDefinition(sample)
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
            defQ = new PresentationDefinition(sample)
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
            defQ = new PresentationDefinition(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be equal to one of the allowed values")
        }
        expect(defQ).to.be.null

        sample["submission_requirements"] = [{
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
        }]
        try {
            defQ = new PresentationDefinition(sample)
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
            defQ = new PresentationDefinition(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be object")
        }
        expect(defQ).to.be.null
    })

    it('input_descriptors[*] schema validation ', async () => {
        let sample = Object.assign({}, samplePresentationDefQuery)
        sample["input_descriptors"] = null
        let defQ = null

        try {
            defQ = new PresentationDefinition(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be array")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = []
        try {
            defQ = new PresentationDefinition(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should NOT have fewer than 1 items")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = [{}]
        try {
            defQ = new PresentationDefinition(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should have required property 'id'")
        }
        expect(defQ).to.be.null

        sample["input_descriptors"] = [{
            "id": "banking_input_1"
        }]
        try {
            defQ = new PresentationDefinition(sample)
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
        defQ = new PresentationDefinition(sample)
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
        defQ = new PresentationDefinition(sample)
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
        defQ = new PresentationDefinition(sample)
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
            defQ = new PresentationDefinition(sample)
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
            "constraints": {"fields":[
                    {
                        "filter": {
                            "type": "date",
                            "minimum": "1999-5-16"
                        }
                    }
                ]}
        }]
        try {
            defQ = new PresentationDefinition(sample)
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
            "constraints": {"fields":[
                    {
                        "path": ["$.jobs[*].active"],
                        "filter": "sample-filter"
                    }
                ]}
        }]
        try {
            defQ = new PresentationDefinition(sample)
        } catch (e) {
            expect(e[0].message).to.have.string("should be object")
        }
        expect(defQ).to.be.null
    })


})
