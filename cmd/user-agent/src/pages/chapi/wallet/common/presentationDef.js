/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Ajv from 'ajv';


/**
 * PresentationDefinition represents Presentation Definitions objects
 * to articulate what proofs an entity requires to make a decision about an interaction with a Subject.
 *
 * RFC: https://identity.foundation/presentation-exchange/#presentation-definition
 *
 * @param presentation definition requirement object
 * @class
 */
export class PresentationDefinition {

    constructor(requirement) {
        this._validateSchema(requirement)
        this.requirementObjs = requirement["submission_requirements"]
        this.descriptors = requirement["input_descriptors"]
    }

    _validateSchema(data) {
        let valid = validate(data);
        if (!valid) {
            throw validate.errors
        }
    }

    // TODO value parsing & input evaluation
}

// presentationDefSchema based on RFC  RFC: https://identity.foundation/presentation-exchange/#presentation-definition
const presentationDefSchema = {
    "required": [
        "input_descriptors"
    ],
    "properties": {
        "submission_requirements": {
            "type": ["array"],
            "items":
                {
                    "type": "object",
                    "properties": {
                        "name": {
                            "type": "string"
                        },
                        "purpose": {
                            "type": "string"
                        },
                        "rule": {
                            "type": "object",
                            "properties": {
                                "type": {
                                    "type": "string",
                                    "enum": ["all", "pick"]
                                },
                                "count": {
                                    "type": "integer"
                                },
                                "from": {
                                    "type": ["array"],
                                    "items":
                                        {
                                            "type": "string",
                                        }
                                }
                            },
                            "required": ["type", "from"],
                            "if": {
                                "properties": {
                                    "type": {"const": "pick"}
                                },
                                "required": ["type"]
                            },
                            "then": {"required": ["count"]}
                        }
                    },
                    "required": ["rule"]
                }

        },
        "input_descriptors": {
            "type": ["array"],
            "minItems": 1,
            "items":
                {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "string"
                        },
                        "group": {
                            "type": "array",
                            "items":
                                {
                                    "type": "string",
                                }
                        },
                        "schema": {
                            "type": "object",
                            "properties": {
                                "uri": {
                                    "type": ["array", "string"],
                                    "items":
                                        {
                                            "type": "string",
                                        }
                                },
                                "name": {
                                    "type": "string"
                                },
                                "purpose": {
                                    "type": "string"
                                },
                            },
                            "required": ["uri"]
                        },
                        "constraints": {
                            "type": "object",
                            "properties": {
                                "fields": {
                                    "type": ["array"],
                                    "minItems": 1,
                                    "items":
                                        {
                                            "type": "object",
                                            "properties": {
                                                "path": {
                                                    "type": ["array"],
                                                    "items":
                                                        {
                                                            "type": "string",
                                                        }
                                                },
                                                "purpose": {
                                                    "type": "string"
                                                },
                                                "filter": {
                                                    "type": "object"
                                                },
                                            },
                                            "required": ["path"]
                                        }
                                },
                            }
                        }
                    },
                    "required": ["id", "schema"]
                }
        }
    }
}

// Schema validator
var ajv = new Ajv();
var validate = ajv.compile(presentationDefSchema);