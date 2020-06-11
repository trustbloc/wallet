/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export const samplePresentationDefQuery = {
    "submission_requirements": [
        {
            "name": "Banking Information",
            "purpose": "We need to know if you have an established banking history.",
            "rule": {
                "type": "pick",
                "count": 1,
                "from": ["A"]
            }
        },
        {
            "name": "Employment Information",
            "purpose": "We need to know that you are currently employed.",
            "rule": {
                "type": "all",
                "from": ["B"]
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
