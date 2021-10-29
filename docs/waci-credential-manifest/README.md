# WACI Credential Manifest

## Wallet and Credential Interactions for Issuance

**Specification Status:** Draft

**Authors:**
[Rolson Quadras](https://github.com/rolsonquadras) (SecureKey)

---

## Abstract

There are interactions between a wallet and an issuer that require passing
information between the two. This specification provides an initial protocol
definition for these interactions.

This document describes an interoperability profile which incorporates elements
from a number of existing specifications and protocols, without assuming or
requiring an implementer to understand all of them. It inherits its overall
structure from
[the current pre-draft of WACI Offer/Claim](https://identity.foundation/wallet-and-credential-interactions/#offerclaim),
but makes use of elements from the
[DIDComm v2.0](https://github.com/decentralized-identity/didcomm-messaging)
messaging protocol, along with
[Issue Credential v3](./issue_credential/README.md)
message formats and
[DIF Credential Manifest](https://identity.foundation/credential-manifest/)
data objects. This version of the specification also restricts itself to
[Verifiable Credentials](https://www.w3.org/TR/vc-data-model/).

It is anticipated that future versions of this specification will add
support for a much broader range of messaging and data-sharing formats than
those used in v0.1.

## Flow Diagram

![WACI Issuance flow diagram](./images/waci-issuance-flow.svg)

## Interoperability Profile

### DIDComm

#### Step 1 : Generate Out-Of-Band (OOB) message

The user scans a QR code or Clicks on a redirect URL from an Issuer website to save the credential in the wallet. The Out-Of-Band (OOB) message structure is [similar to WACI-PEx OOB](https://identity.foundation/waci-presentation-exchange/#step-1-generate-qr-code) with different `goal_code`.

```json=
{
   "type":"https://didcomm.org/out-of-band/2.0/invitation",
   "id":"f137e0db-db7b-4776-9530-83c808a34a42",
   "from":"did:example:issuer",
   "body":{
      "goal_code":"streamlined-vc",
      "accept":[
         "didcomm/v2"
      ]
   }
}
```

#### Step 2 : Issue Credential - Propose Credential

The wallet (user agent) initiates Issuer interaction by sending [Issue Credential - Propose Credential](https://github.com/hyperledger/aries-rfcs/tree/main/features/0453-issue-credential-v2#propose-credential) message with `pthid` same as `id` from OOB message. This provides context to the Issuer and helps map to the original session.

```json=
{
   "type":"https://didcomm.org/issue-credential/3.0/propose-credential",
   "id":"7f62f655-9cac-4728-854a-775ba6944593",
   "pthid":"f137e0db-db7b-4776-9530-83c808a34a42",
   "from":"did:example:holder",
   "to":[
      "did:example:issuer"
   ]
}
```

#### Step 3 : Issue Credential - Offer Credential (Credential Manifest)

The Issuer will send the [Issue Credential - Offer Credential](https://github.com/hyperledger/aries-rfcs/tree/main/features/0453-issue-credential-v2#offer-credential) message to Holder. The message attachment will have [Credential Manifest message](https://identity.foundation/credential-manifest/#credential-manifest-2) from [Credential Manifest Spec](https://identity.foundation/credential-manifest/), which will contain challenge/domain along with optional presentation definition (in case issuer needs any other credential before issuing the new credential) and output descriptor to show credential preview data to the user.

In this case, the issuer wants a Permanent Resident Card (PRC) inorder to issue a Drivers License (DL).

```json=
{
   "type":"https://didcomm.org/issue-credential/3.0/offer-credential",
   "id":"07c44208-06a9-4f8a-a5ce-8ce953270d4b",
   "pthid":"f137e0db-db7b-4776-9530-83c808a34a42",
   "from":"did:example:issuer",
   "to":[
      "did:example:holder"
   ],
   "body":{

   },
   "attachments":[
      {
         "id":"e00e11d4-906d-4c88-ba72-7c66c7113a78",
         "media_type":"application/json",
         "format":"dif/credential-manifest/manifest@v1.0",
         "data":{
            "json":{
               "options":{
                  "challenge":"508adef4-b8e0-4edf-a53d-a260371c1423",
                  "domain":"9rf25a28rs96"
               },
               "id":"WA-DL-CLASS-A",
               "version":"0.1.0",
               "issuer":{
                  "id":"did:example:123?linked-domains=3",
                  "name":"Washington State Government",
                  "styles":{

                  }
               },
               "presentation_definition":{
                  "id":"32f54163-7166-48f1-93d8-ff217bdb0654",
                  "input_descriptors":[
                     {
                        "id":"prc",
                        "name":"Permanent Resident Card",
                        "purpose":"We need PRC to verify your status.",
                        "schema":[
                           {
                              "uri":"https://w3id.org/citizenship#PermanentResidentCard"
                           }
                        ]
                     }
                  ]
               },
               "output_descriptors":[
                  {
                     "id":"driver_license_output",
                     "schema":"https://schema.org/EducationalOccupationalCredential",
                     "display":{
                        "title":{
                           "path":[
                              "$.name",
                              "$.vc.name"
                           ],
                           "fallback":"Washington State Driver License"
                        },
                        "subtitle":{
                           "path":[
                              "$.class",
                              "$.vc.class"
                           ],
                           "fallback":"Class A, Commercial"
                        },
                        "description":{
                           "text":"License to operate a vehicle with a gross combined weight rating (GCWR) of 26,001 or more pounds, as long as the GVWR of the vehicle(s) being towed is over 10,000 pounds."
                        },
                        "properties":[
                           {
                              "path":[
                                 "$.donor",
                                 "$.vc.donor"
                              ],
                              "fallback":"Unknown",
                              "label":"Organ Donor"
                           }
                        ]
                     },
                     "styles":{
                        "thumbnail":{
                           "uri":"https://dol.wa.com/logo.png",
                           "alt":"Washington State Seal"
                        },
                        "hero":{
                           "uri":"https://dol.wa.com/happy-people-driving.png",
                           "alt":"Happy people driving"
                        },
                        "background":{
                           "color":"#ff0000"
                        },
                        "text":{
                           "color":"#d4d400"
                        }
                     }
                  }
               ]
            }
         }
      }
   ]
}
```

#### Step 4 : Issue Credential - Request Credential (Credential Application)

The User sends the [Credential Application message](https://identity.foundation/credential-manifest/#credential-application) as an attachment in [Issue Credential - Request Credential](https://github.com/hyperledger/aries-rfcs/tree/main/features/0453-issue-credential-v2#request-credential).

```json=
{
   "type":"https://didcomm.org/issue-credential/3.0/request-credential",
   "id":"c6686159-ef49-45b2-938f-51818da14723",
   "pthid":"f137e0db-db7b-4776-9530-83c808a34a42",
   "from":"did:example:holder",
   "to":[
      "did:example:issuer"
   ],
   "body":{

   },
   "attachments":[
      {
         "id":"e00e11d4-906d-4c88-ba72-7c66c7113a78",
         "media_type":"application/json",
         "format":"dif/credential-manifest/application@v1.0",
         "data":{
            "json":{
               "@context":[
                  "https://www.w3.org/2018/credentials/v1",
                  "https://identity.foundation/credential-manifest/application/v1"
               ],
               "type":[
                  "VerifiablePresentation",
                  "CredentialApplication"
               ],
               "credential_application":{
                  "id":"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
                  "manifest_id":"WA-DL-CLASS-A",
                  "format":{
                     "ldp_vc":{
                        "proof_type":[
                           "JsonWebSignature2020",
                           "EcdsaSecp256k1Signature2019"
                        ]
                     }
                  }
               },
               "presentation_submission":{
                  "id":"a30e3b91-fb77-4d22-95fa-871689c322e2",
                  "definition_id":"32f54163-7166-48f1-93d8-ff217bdb0654",
                  "descriptor_map":[
                     {
                        "id":"prc",
                        "path":"$.verifiableCredential[0]"
                     }
                  ]
               },
               "verifiableCredential":[
                  {
                     "@context":[
                        "https://www.w3.org/2018/credentials/v1",
                        "https://w3id.org/citizenship#PermanentResidentCard"
                     ],
                     "id":"ccc017c2-6cb1-4e40-bcf4-fca6cf95fe7b",
                     "type":[
                        "VerifiableCredential",
                        "PermanentResidentCard"
                     ],
                     "issuer":"did:foo:123",
                     "issuanceDate":"2020-01-01T19:73:24Z",
                     "credentialSubject":{
                        "id":"did:example:ebfeb1f712ebc6f1c276e12ec21",
                        "givenName":"john"
                     },
                     "proof":{
                        "type":"EcdsaSecp256k1VerificationKey2019",
                        "created":"2017-06-18T21:19:10Z",
                        "proofPurpose":"assertionMethod",
                        "verificationMethod":"https://example.edu/issuers/keys/1",
                        "jws":"..."
                     }
                  }
               ],
               "proof":{
                  "type":"RsaSignature2018",
                  "created":"2018-09-14T21:19:10Z",
                  "proofPurpose":"authentication",
                  "verificationMethod":"did:example:ebfeb1f712ebc6f1c276e12ec21#keys-1",
                  "challenge":"1f44d55f-f161-4938-a659-f8026467f126",
                  "domain":"4jt78h47fh47",
                  "jws":"..."
               }
            }
         }
      }
   ]
}
```

#### Step 5 : Issue Credential - Issue Credential (Credential Fulfilment)

The Issuer sends the [Credential Fulfilment message](https://identity.foundation/credential-manifest/#credential-fulfillment) as an attachment in [Issue Credential - Issue Credential](https://github.com/hyperledger/aries-rfcs/tree/main/features/0453-issue-credential-v2#issue-credential) which will contain the Verifiable Credentials.

```json=
{
   "type":"https://didcomm.org/issue-credential/3.0/issue-credential",
   "id":"7a476bd8-cc3f-4d80-b784-caeb2ff265da",
   "pthid":"f137e0db-db7b-4776-9530-83c808a34a42",
   "from":"did:example:issuer",
   "to":[
      "did:example:holder"
   ],
   "body":{

   },
   "attachments":[
      {
         "id":"e00e11d4-906d-4c88-ba72-7c66c7113a78",
         "media_type":"application/json",
         "format":"dif/credential-manifest/application@v1.0",
         "data":{
            "json":{
               "@context":[
                  "https://www.w3.org/2018/credentials/v1",
                  "https://identity.foundation/credential-manifest/fulfillment/v1"
               ],
               "type":[
                  "VerifiablePresentation",
                  "CredentialFulfillment"
               ],
               "credential_fulfillment":{
                  "id":"a30e3b91-fb77-4d22-95fa-871689c322e2",
                  "manifest_id":"32f54163-7166-48f1-93d8-ff217bdb0653",
                  "descriptor_map":[
                     {
                        "id":"citizenship_output_1",
                        "format":"ldp_vc",
                        "path":"$.verifiableCredential[1]"
                     }
                  ]
               },
               "verifiableCredential":[
                  {
                     "@context":"https://www.w3.org/2018/credentials/v1",
                     "id":"https://eu.com/claims/DriversLicense",
                     "type":[
                        "EUDriversLicense"
                     ],
                     "issuer":"did:foo:123",
                     "issuanceDate":"2010-01-01T19:73:24Z",
                     "credentialSubject":{
                        "id":"did:example:ebfeb1f712ebc6f1c276e12ec21",
                        "license":{
                           "number":"34DGE352",
                           "dob":"07/13/80"
                        }
                     },
                     "proof":{
                        "type":"RsaSignature2018",
                        "created":"2017-06-18T21:19:10Z",
                        "proofPurpose":"assertionMethod",
                        "verificationMethod":"https://example.edu/issuers/keys/1",
                        "jws":"..."
                     }
                  }
               ],
               "proof":{
                  "type":"RsaSignature2018",
                  "created":"2018-09-14T21:19:10Z",
                  "proofPurpose":"authentication",
                  "verificationMethod":"did:example:ebfeb1f712ebc6f1c276e12ec21#keys-1",
                  "challenge":"1f44d55f-f161-4938-a659-f8026467f126",
                  "domain":"4jt78h47fh47",
                  "jws":"..."
               }
            }
         }
      }
   ]
}
```

#### Step 6 : Issue Credential - Ack

https://github.com/hyperledger/aries-rfcs/tree/main/features/0453-issue-credential-v2#issue-credential

```json=
{
   "type":"https://didcomm.org/issue-credential/3.0/ack",
   "id":"d1fb78ad-c452-4c52-a7a0-b68b3e82cdd3",
   "pthid":"f137e0db-db7b-4776-9530-83c808a34a42",
   "from":"did:example:holder",
   "to":[
      "did:example:issuer"
   ],
   "body":{

   }
}
```
