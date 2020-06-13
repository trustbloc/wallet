/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Ajv from 'ajv';
import jp from 'jsonpath';
import {presentationDefSchema} from './presentationDefSchema';


const presentationSubmissionTemplate = `{
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://identity.foundation/presentation-exchange/submission/v1"
    ],
    "type": ["VerifiablePresentation", "PresentationSubmission"],
    "presentation_submission": {
        "descriptor_map": []
    },
    "verifiableCredential": []
}`

/**
 * PresentationDefinition represents Presentation Definitions objects
 * to articulate what proofs an entity requires to make a decision about an interaction with a Subject.
 *
 * RFC: https://identity.foundation/presentation-exchange
 *
 * Note:
 *  when 2 credentials match to same input descriptor id, showing both in presentation submission descriptor map
 *
 * @param presentation definition requirement object
 * @class
 */
export class PresentationExchange {

    constructor(requirement) {
        validateSchema(requirement)

        this.requirementObjs = requirement["submission_requirements"]
        this.descriptors = requirement["input_descriptors"]
        this.applyRules = () => this.requirementObjs && this.requirementObjs.length > 0

        if (this.applyRules()) (
            this._filterDescriptors()
        )
    }

    _filterDescriptors() {
        // validate groups defined in 'submission_requirements'
        var requiredRules = jp.query(this.requirementObjs, '$..rule.from[*]');
        var availableRules = jp.query(this.descriptors, '$..group[*]');
        if (!requiredRules.every(v => availableRules.includes(v))) {
            throw [{message: "Couldn't find matching group in descriptors for 'submission_requirements'"}]
        }

        // retain descriptors only needed by required rules
        let descriptors = this.descriptors.filter(descriptor => descriptor.group.some(v => requiredRules.includes(v)))

        let descrsByGroup = new Map()
        requiredRules.forEach(function (rule) {
            descrsByGroup[rule] = descriptors.filter(descriptor => descriptor.group.includes(rule))
        })

        this.descriptorsByGroup = descrsByGroup
        this.descriptors = descriptors
    }

    createPresentationSubmission(credentials) {
        let results = []

        if (this.applyRules()) {
            results = evaluateByRules(credentials, this.descriptorsByGroup, this.requirementObjs)
        } else {
            results = evaluateAll(credentials, this.descriptors)
        }

        return prepareSubmission(results)
    }
}

// Schema validator
var ajv = new Ajv();
var validate = ajv.compile(presentationDefSchema);

function validateSchema(data) {
    let valid = validate(data);
    if (!valid) {
        throw validate.errors
    }
}

// match matches given credential with descriptor
function match(credential, descriptor) {
    // match schema
    let schemas = Array.isArray(descriptor.schema.uri) ? descriptor.schema.uri : [descriptor.schema.uri]
    let contexts = Array.isArray(credential["@context"]) ? credential["@context"] : [credential["@context"]]
    let schemaMatched = contexts.some(v => schemas.includes(v))

    if (!schemaMatched) {
        // schema not matched, skip this credential
        return false
    }

    if (!descriptor.constraints || !descriptor.constraints.fields || descriptor.constraints.fields.length == 0) {
        // if no constraints declared, credential matched !!
        return true
    }

    // found constraints, apply filter using constraints,
    let filterMatched
    for (let f in descriptor.constraints.fields) {
        let field = descriptor.constraints.fields[f]
        let valueFound
        // look for matching value
        for (let p in field.path) {
            valueFound = jp.query(credential, field.path[p])
            if (!valueFound || valueFound.length > 0) {
                break
            }
        }

        // no matching path found in given credential
        if (valueFound && valueFound.length == 0) {
            filterMatched = false
            break
        }

        // if filter present, then apply filter
        if (field.filter) {
            valueFound = valueFound.filter(v => ajv.validate(field.filter, v))
            if (valueFound.length == 0) {
                // only if the result is valid, proceed iterating the rest of the fields entries
                filterMatched = false
                break
            }
        }

        filterMatched = true
    }

    return filterMatched
}

// prepareSubmission creates presentation submission for all matched credentials
function prepareSubmission(results) {
    let presentationSubmission = JSON.parse(presentationSubmissionTemplate)

    results.forEach(function (result, index) {
        //TODO add VC only once if it matches 2 conditions
        presentationSubmission.verifiableCredential.push(result.credential)
        presentationSubmission.presentation_submission.descriptor_map.push(
            {
                id: result.id,
                path: `$.verifiableCredential.[${index}]`
            }
        )
    })

    return presentationSubmission
}

// evaluateAll evaluates credentials based on all input descriptors
function evaluateAll(credentials, descriptors) {
    let results = []

    credentials.forEach(function (credential) {
        descriptors.forEach(function (descriptor) {
            if (match(credential, descriptor)) {
                results.push({credential, id: descriptor.id})
            }
        })
    })

    return results
}

// evaluateByRules evaluates credentials based on submission rules
function evaluateByRules(credentials, descrsByGroup, submissions) {
    let results = []

    submissions.forEach(function (submission) {

        submission.rule.from.forEach(function (rule) {
            let descriptors = descrsByGroup[rule]
            let mustPass = submission.rule.type == "all" ? descriptors.length : submission.rule.count

            credentials.forEach(function (credential) {
                let matches = descriptors.filter(d => match(credential, d))
                if (matches.length >= mustPass) {
                    matches.forEach(function (match) {
                        results.push({credential, id: match.id})
                    })
                }
            })

        })

    })

    return results
}






