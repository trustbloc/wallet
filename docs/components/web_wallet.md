# User Agent Web Wallet

 User agent Web Wallet is based on [CHAPI](https://w3c-ccg.github.io/credential-handler-api/) developed using [credential handler polyfill](https://github.com/digitalbazaar/credential-handler-polyfill), [Vue](https://vuejs.org)

## Steps to start user agent web wallet

- run below command to build and start user agent

    ```bash
    make user-agent-start
    ```

## Web wallet demo

Once user agent is started using previous step,

- Login to web wallet from user web wallet [dashboard](https://user-ui-agent.example.com:8091/dashboard).

  Proceed with pre-filled username/password for the login. Once login is successful, you will get a prompt from your browser to allow wallet to manage credentials, choose 'Allow'. 
- Go to web wallet [demo page](https://user-ui-agent.example.com:8091/webwallet)

  Web wallet demo page can be used to perform all supported wallet operations. Click on sample requests based on operation you wish to perform and click on `STORE` or `GET` buttons. 
  
  Web wallet allows only 2 kind of requests for requesting and storing credentials,
   1. Get - A web app (a Relying Party or verifier) can request a credential using get requests.
   2. Store - A web app (for example, a credential issuer such as a university or institution) can ask to store a credential using store requests.
  
 
## Operations supported by user agent web wallet
User agent web wallet supports various wallet operations and also [DIDComm](https://github.com/hyperledger/aries-rfcs/blob/master/concepts/0005-didcomm/README.md). 
- **DID auth**: This operation can be used when a website that wants user to authenticate (prove they are authorized to act as a particular DID subject, aka "DIDAuth" over wallet)
- **Store Credential**: This operation can be used when an issuer website wants to give user a credential to store in their wallet
- **Get Credential**: This operation can be used when a website (a Relying Party or verifier) wants to request a credential from user.
- **DID Connect**: This operation can be used by a website(or agent) to perform [DID Exchange](https://github.com/hyperledger/aries-rfcs/blob/master/features/0023-did-exchange/README.md) with web wallet for all future DIDComm operations.  
- **Presentation Exchange**: This operation can be used by a website (a Relying Party or verifier) to request submission of proof from wallet in align with requester's [proof requirements](https://identity.foundation/presentation-exchange/). This request can also be made in conjuction with DID Connect operation, so that wallet can request further information from requester (a Relying party or verifier) over DIDComm channel.

## Login to wallet
After successful login, web wallet will load [polyfill library](https://github.com/digitalbazaar/credential-handler-polyfill) and a web credential handler will be installed which will serve all your credential handler request as per [Credential Handler API](https://w3c-ccg.github.io/credential-handler-api/) (aka CHAPI) standards.

In addition to that, user agent web wallet will perform some onboarding operations during first login of an user to support various DIDComm operations which includes,

 - Registering with router(also known as mediator). Since user agent web wallet is technically also an agent which runs in your browser. Since it has no inbound transport and cannot be online always, it has to register 
 with [router](https://github.com/hyperledger/aries-framework-go/blob/master/docs/didcomm_mediator.md) to route the messages to it by asking for permission. On successful grant, agent receives the endpoint and routing key details.
 Sample configuration of how to pass a mediator URL in user agent setup can be found [here](https://github.com/trustbloc/edge-agent/blob/master/test/bdd/fixtures/user-agent/docker-compose.yml#L71)
 - Creating TrustBloc DID. By default user agent web wallet uses trustbloc DIDs for generating verifiable presentation response over CHAPI.
 - Saving created DIDs and other metadata in storage for future operations.

## Storing credential(s) in wallet
A single or multiple credentials can be stored into wallet through `store()` call using a CHAPI request as given in sample below. 

In below sample, a presentation with single or multiple credentials is stored into wallet,

 ```
 const webCredential = new WebCredential('VerifiablePresentation', {
   '@context': 'https://www.w3.org/2018/credentials/v1',
   ...presentation
 });

const result = await navigator.credentials.store(webCredential);
if(!result) {
  console.log('store credential operation did not succeed');
}

```

In above example we are creating `WebCredential` instance of type `VerifiablePresentation` and credential data is being sent in presentation format.
Currently web wallet also supports sending credential to wallet as `VerifiableCredential` type. But using `VerifiablePresentation` format is always recommended.

## Getting credential from wallet
Credentials can be requested from web wallet through `get()` call using a CHAPI request as given in sample below,

```javascript
let credentialQuery = {
  "web": {
    "VerifiablePresentation": {
      "query": [
        {
          "type": "QueryByExample",
          "credentialQuery": {
            "reason": "Please present a credential for JaneDoe.",
            "example": {
              "@context": [
                "https://www.w3.org/2018/credentials/v1",
                "https://www.w3.org/2018/credentials/examples/v1"
              ],
              "type": [
                "UniversityDegreeCredential"
              ]
            }
          }
        }
      ],
      "challenge": "26e18e44-7c01-4e16-bbf9-1463e00df582",
      "domain": "example.com"
    }
  }
}

const webCredential = await navigator.credentials.get(credentialQuery);

```
Note: currently, web wallet lists all the saved credentials and applies filter based on type `query.credentialQuery.example.type`.

## DID Auth from wallet
A DID Auth request asking to prove ownership of a DID can be sent to wallet as given in sample below (wallet `get()` with type `"DIDAuth"` ),

```javascript
let didAuthQuery = {
  "web": {
    "VerifiablePresentation": {
      "query": {
        "type": "DIDAuth"
      },
      "challenge": "db926a16-791d-4a47-9d81-f9d5012bae0c",
      "domain": "example.com"
    }
  }
}

const webCredential = await navigator.credentials.get(didAuthQuery);

```

## Sending Presentation Exchange request to web wallet
A relying party or verifier can request submission of proof from wallet in align with requester's [proof requirements](https://identity.foundation/presentation-exchange/) as given in sample below (wallet `get()` with type `"PresentationDefinitionQuery"`),

Once presentation exchange request received, Web wallet will query saved credentials in wallet based on criteria defined in presentation exchange request. 

```javascript
let presentationExchangeQuery = {
                                  "web": {
                                    "VerifiablePresentation": {
                                      "query": [
                                        {
                                          "type": "PresentationDefinitionQuery",
                                          "presentationDefinitionQuery": {
                                            "submission_requirements": [
                                              {
                                                "name": "Education Qualification",
                                                "purpose": "We need to know if you are qualified for this job",
                                                "rule": "pick",
                                                "count": 1,
                                                "from": [
                                                  "E"
                                                ]
                                              },
                                              {
                                                "name": "Citizenship Information",
                                                "purpose": "You must be legally allowed to work in United States",
                                                "rule": "all",
                                                "from": [
                                                  "C"
                                                ]
                                              }
                                            ],
                                            "input_descriptors": [
                                              {
                                                "id": "citizenship_input_1",
                                                "group": [
                                                  "C"
                                                ],
                                                "schema": {
                                                  "uri": [
                                                    "https://w3id.org/citizenship/v1",
                                                    "https://w3id.org/citizenship/v2",
                                                    "https://w3id.org/citizenship/v3"
                                                  ],
                                                  "name": "US Permanent resident card"
                                                },
                                                "constraints": {
                                                  "fields": [
                                                    {
                                                      "path": [
                                                        "$.credentialSubject.lprCategory"
                                                      ],
                                                      "filter": {
                                                        "type": "string",
                                                        "pattern": "C09|C52|C57"
                                                      }
                                                    }
                                                  ]
                                                }
                                              },
                                              {
                                                "id": "degree_input_1",
                                                "group": [
                                                  "E"
                                                ],
                                                "schema": {
                                                  "uri": [
                                                    "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld"
                                                  ],
                                                  "name": "University degree certificate",
                                                  "purpose": "We need your education qualification details."
                                                },
                                                "constraints": {
                                                  "fields": [
                                                    {
                                                      "path": [
                                                        "$.credentialSubject.degree.type"
                                                      ],
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
                                                "group": [
                                                  "E"
                                                ],
                                                "schema": {
                                                  "uri": [
                                                    "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld"
                                                  ],
                                                  "name": "Diploma certificate",
                                                  "purpose": "We need your education qualification details."
                                                },
                                                "constraints": {
                                                  "fields": [
                                                    {
                                                      "path": [
                                                        "$.credentialSubject.degree.type"
                                                      ],
                                                      "purpose": "Should have valid diploma",
                                                      "filter": {
                                                        "type": "string",
                                                        "pattern": "Diploma"
                                                      }
                                                    },
                                                    {
                                                      "path": [
                                                        "$.credentialSubject.degree.coop"
                                                      ],
                                                      "purpose": "Should have co-op experience",
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
                                        }
                                      ],
                                      "challenge": "df840294-787e-442f-824a-1ccb5d6c1da9",
                                      "domain": "example.com"
                                    }
                                  }
                                }

const webCredential = await navigator.credentials.get(presentationExchangeQuery);

```

## DIDComm by web wallet
User agent web wallet supports DID communication for secured interactions with issuer and verifiers. 

In order to perform DID communication with web wallet, an issuer or a verifier has to connect to web wallet using **DID Connect**, a wallet `get()` call with query type `"DIDConnect"`.

### Issuer connecting to web wallet
In order to establish DIDComm connection with web wallet, an issuer should perform wallet `get()` operation by sending query with type `"DIDConnect"` and a DIDComm invitation as given in sample below.
In addition to `DIDComm invitation`, an issuer can also send credentials to store as given below.

```javascript
let didConnectQuery = {
    "web": {
        "VerifiablePresentation": {
            "query": {
                "type": "DIDConnect"
            },
            "invitation": {
                "@id": "0d72ebac-6210-46b4-a736-43a3bf90f976",
                "@type": "https://didcomm.org/oob-invitation/1.0/invitation",
                "label": "ewr",
                "protocols": [
                    "https://didcomm.org/didexchange/1.0"
                ],
                "service": [
                    {
                        "ID": "6941a000-b7c8-4535-be49-814a2289bc26",
                        "Priority": 0,
                        "Properties": null,
                        "RecipientKeys": [
                            "EYbrpVudXxPJnowxCycPhRz6NtdQdu9sZKYoHHyq2Fc4"
                        ],
                        "RoutingKeys": [
                            "CTkGyS7PjPgYDHhgv6WJwjgUU9ABbrT68x1X4HzwazM4"
                        ],
                        "ServiceEndpoint": "https://localhost:10091",
                        "Type": "did-communication"
                    }
                ]
            },
            // credentials issuer wants to send like manifests, governance credential or any other credentials 
            "credentials": [],
            "challenge": "6919ac64-9771-4343-a50f-318bca774d86",
            "domain": "example.com"
        }
    }
}

const webCredential = await navigator.credentials.get(didConnectQuery);
```


### Verifier or Relying party connecting to web wallet
A verifier or relying party can connect to wallet simply by sending `get()` query of type `"DIDConnect"`.

Also, while performing presentation exchange with web wallet, a verifier or relying party can also send `"DIDConnect"` request along with presentation exchange request as given in the sample below

```javascript
let presentationExchangeDIDCommQuery = {
                                           "web": {
                                               "VerifiablePresentation": {
                                                   "query": [
                                                       {
                                                           "type": "PresentationDefinitionQuery",
                                                           "presentationDefinitionQuery": {...}
                                                       },
                                                       {
                                                           "type": "DIDConnect",
                                                           "invitation": {
                                                               "@id": "0d72ebac-6210-46b4-a736-43a3bf90f976",
                                                               "@type": "https://didcomm.org/oob-invitation/1.0/invitation",
                                                               "label": "ewr",
                                                               "protocols": [
                                                                   "https://didcomm.org/didexchange/1.0"
                                                               ],
                                                               "service": [
                                                                   {
                                                                       "ID": "6941a000-b7c8-4535-be49-814a2289bc26",
                                                                       "Priority": 0,
                                                                       "Properties": null,
                                                                       "RecipientKeys": [
                                                                           "EYbrpVudXxPJnowxCycPhRz6NtdQdu9sZKYoHHyq2Fc4"
                                                                       ],
                                                                       "RoutingKeys": [
                                                                           "CTkGyS7PjPgYDHhgv6WJwjgUU9ABbrT68x1X4HzwazM4"
                                                                       ],
                                                                       "ServiceEndpoint": "https://localhost:10091",
                                                                       "Type": "did-communication"
                                                                   }
                                                               ]
                                                           }
                                                       }
                                                   ],
                                                   "challenge": "37c66a97-b2c9-4666-b3d5-66d01b02707b",
                                                   "domain": "example.com"
                                               }
                                           }
                                       }
                                       
     
const webCredential = await navigator.credentials.get(presentationExchangeDIDCommQuery);                                  
```

## Web wallet DIDComm flow explained

Credentials can be shared by web wallet which introduces an issuer who can issue requested credential(s) to a relying party or a verifier through DID communication with authorization of wallet user. 


- Step 1: An issuer connects to web wallet by sending DID communication invitation and manifest credential, a credential containing information about list credentials that can be issued by this issuer on demand.
- Step 2: A relying party or verifier sends presentation exchange request with DID communication invitation asking for credentials matching specific criteria.
- Step 3: Wallet queries its saved credentials based on criteria requested `Step 2.`

    - If wallet finds matching saved credentials, then it prepares presentation submission response based on credentials found.
    - If wallet unable to find any matching saved credentials, then it goes through the list of saved manifest credentials sent by all previously connected issuers. If wallet manages to find issuer who can issue requested credentials then it sends issue credential request to that issuer to issue authorization credential for given relyin party or verifier.
    Authorization credential shared by issuer will be sent as part of the response to relying party. Relying party will use this authorization credential to connect to actual issuer to get actual credential.
    - If wallet unable to find any issuer who can fulfill requested criteria then it updates presentation submission response accordingly (with no result for asked criteria).

Note: Web wallet asks user consent every time before sending back any CHAPI response. Web wallet never reuses peer DIDs it created for one issuer with any other verifiers or vice versa. 

Refer following pages for more details, 
- A Complete end to end DIDComm flow examples can be found [here](https://trustbloc.readthedocs.io/en/latest/adapters.html#flows).
- More details about issuer flow can be found [here](https://github.com/trustbloc/edge-adapter/blob/master/docs/issuer/README.md)
- More details about relying party flow can be found [here](https://github.com/trustbloc/edge-adapter/blob/master/docs/rp/integration/wallets.md)


## Web wallet DIDComm flow with blinded routing

DIDComm flow with web wallet can also be achieved using blinded routing where identity of issuer and verifiers will be hidden from each other. More details about blinded routing can be found [here](https://trustbloc.readthedocs.io/en/latest/blinded_routing.html) 

Blinded routing DIDComm flow is very similar to DIDComm flow explained in previous section except the following difference.

- After successfully connecting to an issuer or to a verifier,

    - web wallet requests a new peer DID from invitee (an issuer or verifier). 
    - web wallet creates a new peer DID from router for new peer DID it received from invitee.
    - web wallet shares newly created peer DID with invitee.
    
