/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Ajv from 'ajv';
import jp from 'jsonpath';
import {presentationDefSchema} from './presentationDefSchema';
import {getCredentialType} from './util'


const presentationSubmissionTemplate = `{
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://trustbloc.github.io/context/vp/presentation-exchange-submission-v1.jsonld"
    ],
    "type": ["VerifiablePresentation", "PresentationSubmission"],
    "presentation_submission": {
        "descriptor_map": []
    },
    "presentation_location": {
        "descriptor_map": []
    },
    "verifiableCredential": []
}`

const defSubmissionRuleName = "Requested information"
const defSubmissionRulePurpose = "We need below information from your wallet"
const defSubmissionRule = "all conditions should be met"

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
        this.applyRules = this.requirementObjs && this.requirementObjs.length > 0
        if (this.applyRules) (
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

    createPresentationSubmission(credentials, manifests) {
        let results = []

        if (this.applyRules) {
            results = evaluateByRules(credentials, manifests, this.descriptorsByGroup, this.requirementObjs)
        } else {
            results = evaluateAll(credentials, manifests, this.descriptors)
        }

        return prepareSubmission(results)
    }

    requirementDetails() {
        let result = []

        if (this.applyRules) {
            let descrsByGroup = this.descriptorsByGroup
            this.requirementObjs.forEach(function (obj, index) {
                let r = {}
                let {name, purpose, rule} = obj

                r.name = name ? name : `${defSubmissionRuleName} #${index + 1}`
                r.purpose = purpose ? purpose : defSubmissionRulePurpose
                r.rule = rule.type == "all" ? "all conditions should be met" : `at least ${rule.count} of each condition should be met`
                r.descriptors = []

                rule.from.forEach(function (grp) {
                    descrsByGroup[grp].forEach(function (d) {
                        r.descriptors.push(getNameAndPurpose(d))
                    })
                })

                result.push(r)
            })
        } else {
            let r = {
                name: defSubmissionRuleName,
                purpose: defSubmissionRulePurpose,
                rule: defSubmissionRule,
                descriptors: []
            }

            this.descriptors.forEach(function (descriptor) {
                r.descriptors.push(getNameAndPurpose(descriptor))
            })

            result.push(r)
        }

        return result
    }
}


function getNameAndPurpose(descriptor) {
    let name = jp.query(descriptor, "$.schema.name")
    let purpose = jp.query(descriptor, "$.schema.purpose")
    let constraints = jp.query(descriptor, "$.constraints.fields[*].purpose")

    return {
        name: name.length > 0 ? name[0] : "Condition details are not provided in request",
        purpose: purpose.length > 0 ? purpose[0] : "",
        constraints
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
    if (getCredentialType(credential.type) == 'IssuerManifestCredential') {
        return false
    }

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

// matchManifest matches if descriptor schema exists in manifest credential contexts list
// TODO: manifests to have credential previews so that complete constraint checks can be run
function matchManifest(manifest, descriptor) {

    let schemas = Array.isArray(descriptor.schema.uri) ? descriptor.schema.uri : [descriptor.schema.uri]

    return manifest.credentialSubject.contexts.some(v => schemas.includes(v))
}

// prepareSubmission creates presentation submission for all matched credentials
function prepareSubmission(results) {
    let presentationSubmission = JSON.parse(presentationSubmissionTemplate)

    results.forEach(function (result, index) {
        //TODO add VC only once if it matches 2 conditions
        presentationSubmission.verifiableCredential.push(result.credential)

        if (result.manifest){
            presentationSubmission.presentation_location.descriptor_map.push(
                {
                    id: result.id,
                    path: `$.verifiableCredential.[${index}]`
                }
            )
        } else {
            presentationSubmission.presentation_submission.descriptor_map.push(
                {
                    id: result.id,
                    path: `$.verifiableCredential.[${index}]`
                }
            )
        }

    })

    return presentationSubmission
}

// evaluateAll evaluates credentials based on all input descriptors
function evaluateAll(credentials, manifests, descriptors) {
    let result = []

    descriptors.forEach(function (descriptor) {
        let matched = false

        credentials.forEach(function (credential) {
            if (match(credential, descriptor)) {
                matched = true
                result.push({credential, id: descriptor.id})
            }
        })

        // none of the credential matched, check for manifest credential matches
        if (!matched && manifests) {
            manifests.forEach(function (credential) {
                if (matchManifest(credential, descriptor)) {
                    result.push({credential, id: descriptor.id, manifest: true})
                }
            })
        }
    })

    return result
}

// evaluateByRules evaluates credentials based on submission rules
function evaluateByRules(credentials, manifests, descrsByGroup, submissions) {
    let result = []

    submissions.forEach(function (submission) {

        submission.rule.from.forEach(function (rule) {
            let descriptors = descrsByGroup[rule]
            let mustPass = submission.rule.type == "all" ? descriptors.length : submission.rule.count
            let matched = false

            credentials.forEach(function (credential) {
                let matches = descriptors.filter(d => match(credential, d))
                if (matches.length >= mustPass) {
                    matched = true
                    matches.forEach(function (match) {
                        result.push({credential, id: match.id})
                    })
                }
            })

            // none of the credential matched, check for manifest credential matches
            if (!matched && manifests) {
                manifests.forEach(function (credential) {
                    let matches = descriptors.filter(d => matchManifest(credential, d))
                    if (matches.length >= mustPass) {
                        matches.forEach(function (match) {
                            result.push({credential, id: match.id, manifest: true})
                        })
                    }
                })
            }
        })

    })

    return result
}






