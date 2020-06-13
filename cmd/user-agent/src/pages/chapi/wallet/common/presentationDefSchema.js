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
