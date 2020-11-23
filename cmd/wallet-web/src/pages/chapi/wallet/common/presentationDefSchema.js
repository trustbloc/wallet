/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


/**
 * presentationDefSchema is Presentation Definition Schema based on RFC - https://identity.foundation/presentation-exchange/#presentation-definition
 *
 */
export const presentationDefSchema = {
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
                            "type": "string",
                            "enum": ["all", "pick"]
                        },
                        "count": {
                            "type": "integer",
                            "minimum": 1,
                        },
                        "min": {
                            "type": "integer",
                            "minimum": 0,
                        },
                        "max": {
                            "type": "integer",
                            "minimum": 1,
                            "exclusiveMinimum": {"$data": "1/min"}
                        },
                        "from": {
                            "type": ["array"],
                            "items":
                                {
                                    "type": "string",
                                }
                        },
                        "from_nested": {
                            "$ref": "#/properties/submission_requirements"
                        }
                    },
                    "required": ["rule"],
                    "if": {
                        "properties": {
                            "rule": {"const": "pick"}
                        },
                        "required": ["rule"]
                    },
                    "then": {
                        "anyOf": [
                            {
                                "required": [
                                    "count"
                                ]
                            },
                            {
                                "required": [
                                    "min"
                                ]
                            },
                            {
                                "required": [
                                    "max"
                                ]
                            }
                        ]
                    },
                    "oneOf": [
                        {
                            "required": [
                                "from"
                            ]
                        },
                        {
                            "required": [
                                "from_nested"
                            ]
                        }
                    ]
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
