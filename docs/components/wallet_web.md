# Wallet Web

Wallet Web is based on [CHAPI](https://w3c-ccg.github.io/credential-handler-api/) developed using [credential handler polyfill](https://github.com/digitalbazaar/credential-handler-polyfill), [Vue](https://vuejs.org), [Aries JS Worker](https://github.com/hyperledger/aries-framework-go/blob/main/cmd/aries-js-worker/README.md).

This wallet uses [Wallet SDK](https://github.com/trustbloc/agent-sdk/tree/main/cmd/wallet-js-sdk) built on top of [Aries Verifiable Credential wallet](https://github.com/hyperledger/aries-framework-go/blob/main/docs/vc_wallet.md) based on [Universal Wallet Specifications](https://w3c-ccg.github.io/universal-wallet-interop-spec/) implementation.

## Key components

Here are the major components used by Wallet Web.

- **Storage** - used for storing wallet contents. Supports many storage providers, but Encrypted Data Vaults(EDV) is highly recommended. Refer [TrustBloc EDV](https://github.com/trustbloc/edv) for TrustBloc Encrypted Data Vault implementations.
  - Supports EDV and indexedDB for Aries Web Assembly based user agents.
  - Supports EDV, couchDB, levelDB, mongoDB and many more storage providers in case of Aries REST based user agents.
- **Key Management System (KMS)** - used for managing keys for authorization, cryptographic operations, encrypted data vaults etc. Supports both local KMS and WebKMS, but WebKMS is highly recommended. Refer [TrustBloc KMS](https://github.com/trustbloc/kms) for TrustBloc WebKMS implementations.

- **DID Management (VDR)** - For creating, resolving DIDs. By Default Web Wallet supports Orb DIDs. Wallet Web also provides features to import any DIDs with their keys.
  Refer this [link](https://trustbloc.github.io/did-method-orb/) to learn more about Orb DIDs.

- **Authentication** - Any OIDC implementation can be integrated into Wallet Web for user login and authentication.
  Refer [TrustBloc Hub Auth](https://github.com/trustbloc/hub-auth) for TrustBloc end user authentication services.

Refer `User Agent (Wallet)` in [Architectural Diagram](https://github.com/trustbloc/sandbox/blob/main/docs/components/adapter_components.md) to learn more about relationship of Wallet Web with various TrustBloc components.

## Running Local Wallet Web Demo

Once Wallet Web is started using previous step,

- Sign up to Wallet Web (https://wallet.trustbloc.local:8091/).

  Proceed with pre-filled username/password for the login. Once login is successful, you will get a prompt from your browser to allow wallet to manage credentials, choose 'Allow'.

- Go to Wallet Operations Demo Page [demo page](https://demo-adapter.trustbloc.local:8094/web-wallet)

  Wallet Web demo page can be used to perform all supported wallet operations. Click on sample requests based on operation you wish to perform and click on `STORE` or `GET` buttons.

  Wallet Web allows only 2 kind of requests for requesting and storing credentials,

  1.  Get - A web app (a Relying Party or verifier) can request a credential using get requests.
  2.  Store - A web app (for example, a credential issuer such as a university or institution) can ask to store a credential using store requests.

## Operations Supported by Wallet Web

Wallet Web supports various wallet operations and also [DIDComm](https://github.com/hyperledger/aries-rfcs/blob/master/concepts/0005-didcomm/README.md).

- **DID auth**: This operation can be used when a website that wants user to authenticate (prove they are authorized to act as a particular DID subject, aka "DIDAuth" over wallet)
- **Store Credential**: This operation can be used when an issuer website wants to give user a credential to store in their wallet
- **Get Credential**: This operation can be used when a website (a Relying Party or verifier) wants to request a credential from user.
  The Wallet Web supports various query types including the ones in [verifiable presentation request specifications](https://w3c-ccg.github.io/vp-request-spec/).
- **DID Connect**: This operation can be used by a website(or agent) to perform [DID Exchange](https://github.com/hyperledger/aries-rfcs/blob/master/features/0023-did-exchange/README.md) with Wallet Web for all future DIDComm operations.
- **Presentation Exchange**: This operation can be used by a website (a Relying Party or verifier) to request submission of proof from Wallet Web in align with requester's [proof requirements](https://identity.foundation/presentation-exchange/). This request can also be made in conjuction with DID Connect operation, so that wallet can request further information from requester (a Relying party or verifier) over DIDComm channel.

## Wallet Web Signup

After successful signup, Wallet Web will load [polyfill library](https://github.com/digitalbazaar/credential-handler-polyfill) and a web credential handler will be installed which will serve all your credential handler request as per [Credential Handler API](https://w3c-ccg.github.io/credential-handler-api/) (aka CHAPI) standards.

In addition to that, Wallet Web will perform some onboarding operations during the signup of the user to support various DIDComm operations which includes,

- Registering with router(also known as mediator). Since Wallet Web is technically also an agent which runs in your browser. Since it has no inbound transport and cannot be online always, it has to register
  with [router](https://github.com/hyperledger/aries-framework-go/blob/master/docs/didcomm_mediator.md) to route the messages to it by asking for permission. On successful grant, agent receives the endpoint and routing key details.
  Sample configuration of how to pass a mediator URL in user agent setup can be found [here](https://github.com/trustbloc/wallet/blob/master/test/fixtures/wallet-web/docker-compose-web.yml#L24)
- Creating TrustBloc DID. By default Wallet Web uses trustbloc DIDs for generating verifiable presentation response over CHAPI.
- Saving created DIDs and other metadata in storage for future operations.

## Storing Credential(s) in Wallet Web

A single or multiple credentials can be stored into wallet through `store()` call using a CHAPI request as given in sample below.

In below sample, a presentation with single or multiple credentials is stored into wallet,

```js
const webCredential = new WebCredential('VerifiablePresentation', {
  '@context': 'https://www.w3.org/2018/credentials/v1',
  ...presentation,
});

const result = await navigator.credentials.store(webCredential);
if (!result) {
  console.log('store credential operation did not succeed');
}
```

In above example we are creating `WebCredential` instance of type `VerifiablePresentation` and credential data is being sent in presentation format.
Currently Wallet Web also supports sending credential to wallet as `VerifiableCredential` type. But using `VerifiablePresentation` format is always recommended.

## Getting Credential from Wallet Web

Credentials can be requested from Wallet Web through `get()` call using a `CHAPI` request as given in sample below,

```js
let credentialQuery = {
  web: {
    VerifiablePresentation: {
      query: [
        {
          type: 'QueryByExample',
          credentialQuery: {
            reason: 'Please present a credential for JaneDoe.',
            example: {
              '@context': [
                'https://www.w3.org/2018/credentials/v1',
                'https://www.w3.org/2018/credentials/examples/v1',
              ],
              type: ['UniversityDegreeCredential'],
            },
          },
        },
      ],
      challenge: '26e18e44-7c01-4e16-bbf9-1463e00df582',
      domain: 'example.com',
    },
  },
};

const webCredential = await navigator.credentials.get(credentialQuery);
```

Wallet Web supports multiple queries where different query types can be mixed together in a single request.

## DID Auth from Wallet Web

A DID Auth request asking to prove ownership of a DID can be sent to wallet as given in sample below (wallet `get()` with type `"DIDAuth"` ),

```js
let didAuthQuery = {
  web: {
    VerifiablePresentation: {
      query: {
        type: 'DIDAuth',
      },
      challenge: 'db926a16-791d-4a47-9d81-f9d5012bae0c',
      domain: 'example.com',
    },
  },
};

const webCredential = await navigator.credentials.get(didAuthQuery);
```

## Sending Presentation Exchange request to Wallet Web

A relying party or verifier can request submission of proof from Wallet Web in align with requester's [proof requirements](https://identity.foundation/presentation-exchange/) as given in sample below (wallet `get()` with type `"PresentationExchange"`),

Once presentation exchange request received, Wallet Web will query saved credentials in the storage based on criteria defined in presentation exchange request.

```js
let presentationExchangeQuery = {
  web: {
    VerifiablePresentation: {
      query: [
        {
          type: 'PresentationExchange',
          credentialQuery: {
            submission_requirements: [
              {
                name: 'Education Qualification',
                purpose: 'We need to know if you are qualified for this job',
                rule: 'pick',
                count: 1,
                from: ['E'],
              },
              {
                name: 'Citizenship Information',
                purpose: 'You must be legally allowed to work in United States',
                rule: 'all',
                from: ['C'],
              },
            ],
            input_descriptors: [
              {
                id: 'citizenship_input_1',
                group: ['C'],
                schema: {
                  uri: [
                    'https://w3id.org/citizenship/v1',
                    'https://w3id.org/citizenship/v2',
                    'https://w3id.org/citizenship/v3',
                  ],
                  name: 'US Permanent resident card',
                },
                constraints: {
                  fields: [
                    {
                      path: ['$.credentialSubject.lprCategory'],
                      filter: {
                        type: 'string',
                        pattern: 'C09|C52|C57',
                      },
                    },
                  ],
                },
              },
              {
                id: 'degree_input_1',
                group: ['E'],
                schema: {
                  uri: ['https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld'],
                  name: 'University degree certificate',
                  purpose: 'We need your education qualification details.',
                },
                constraints: {
                  fields: [
                    {
                      path: ['$.credentialSubject.degree.type'],
                      purpose: 'Should be masters or bachelors degree',
                      filter: {
                        type: 'string',
                        pattern: 'BachelorDegree|MastersDegree',
                      },
                    },
                  ],
                },
              },
              {
                id: 'degree_input_2',
                group: ['E'],
                schema: {
                  uri: ['https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld'],
                  name: 'Diploma certificate',
                  purpose: 'We need your education qualification details.',
                },
                constraints: {
                  fields: [
                    {
                      path: ['$.credentialSubject.degree.type'],
                      purpose: 'Should have valid diploma',
                      filter: {
                        type: 'string',
                        pattern: 'Diploma',
                      },
                    },
                    {
                      path: ['$.credentialSubject.degree.coop'],
                      purpose: 'Should have co-op experience',
                      filter: {
                        type: 'string',
                        pattern: 'Y',
                      },
                    },
                  ],
                },
              },
            ],
          },
        },
      ],
      challenge: 'df840294-787e-442f-824a-1ccb5d6c1da9',
      domain: 'example.com',
    },
  },
};

const webCredential = await navigator.credentials.get(presentationExchangeQuery);
```

Note: `PresentaionExchange` can be mixed with other credential query types like QueryByExample, QueryByFrame etc. But it is not recommended since it might produce multiple presentations as a response.

## DIDComm by Wallet Web

Wallet Web supports DID communication for secured interactions with issuer and verifiers.

In order to perform DID communication with Wallet Web, an issuer or a verifier has to connect to Wallet Web `get()` operation with query type `"DIDConnect"`.

### Issuer Connecting to Wallet Web

In order to establish DIDComm connection with Wallet Web, an issuer should perform wallet `get()` operation by sending query with type `"DIDConnect"` and a DIDComm invitation as given in sample below.
In addition to `DIDComm invitation`, an issuer can also send credentials to store as given below.

```js
let didConnectQuery = {
  web: {
    VerifiablePresentation: {
      query: {
        type: 'DIDConnect',
      },
      invitation: {
        '@id': '2629b7f4-f8f9-43fc-8964-65740e73d0ab',
        '@type': 'https://didcomm.org/out-of-band/1.0/invitation',
        label: 'issuer',
        services: [
          {
            id: '2c757b3f-2f57-44bc-b9d2-0c2301571f41',
            recipientKeys: ['did:key:z6MkkYU4VfCzss5JuHQiHiXS9GKVHVrs5GFrA4RRTakNu7o2'],
            serviceEndpoint: 'https://adapter-issuer-didcomm.stg.trustbloc.dev',
            type: 'did-communication',
          },
        ],
        handshake_protocols: ['https://didcomm.org/didexchange/1.0'],
      },
      // credentials issuer wants to send like manifests, governance credential or any other credentials
      credentials: [],
      challenge: '6919ac64-9771-4343-a50f-318bca774d86',
      domain: 'example.com',
    },
  },
};

const webCredential = await navigator.credentials.get(didConnectQuery);
```

### Verifier or Relying Party Connecting to Wallet Web

A verifier or relying party can connect to Wallet Web simply by sending `get()` query of type `"DIDConnect"`.

Also, while performing presentation exchange with Wallet Web, a verifier or relying party can also send `"DIDConnect"` request along with presentation exchange request as given in the sample below

```js
let presentationExchangeDIDCommQuery = {
                                           "web": {
                                               "VerifiablePresentation": {
                                                   "query": [
                                                       {
                                                           "type": "PresentationExchange",
                                                           "credentialQuery": {...}
                                                       },
                                                       {
                                                            "type": "DIDConnect",
                                                            "invitation": {
                                                                    "@id": "2629b7f4-f8f9-43fc-8964-65740e73d0ab",
                                                                    "@type": "https://didcomm.org/out-of-band/1.0/invitation",
                                                                    "label": "issuer",
                                                                     "services": [{
                                                                        "id": "2c757b3f-2f57-44bc-b9d2-0c2301571f41",
                                                                        "recipientKeys": ["did:key:z6MkkYU4VfCzss5JuHQiHiXS9GKVHVrs5GFrA4RRTakNu7o2"],
                                                                        "serviceEndpoint": "https://adapter-rp-didcomm.stg.trustbloc.dev",
                                                                        "type": "did-communication"
                                                                    }],
                                                                 "handshake_protocols": ["https://didcomm.org/didexchange/1.0"]
                                                            },
                                                       }
                                                   ],
                                                   "challenge": "37c66a97-b2c9-4666-b3d5-66d01b02707b",
                                                   "domain": "example.com"
                                               }
                                           }
                                       }


const webCredential = await navigator.credentials.get(presentationExchangeDIDCommQuery);
```

## Wallet Web DIDComm Flow Explained

Credentials can be shared by Wallet Web which introduces an issuer who can issue requested credential(s) to a relying party or a verifier through DID communication with authorization of the Wallet Web's user.

- Step 1: An issuer connects to Wallet Web by sending DID communication invitation and manifest credential (a credential containing information about the list of credentials that can be issued by this issuer on demand).
- Step 2: A relying party or verifier sends presentation exchange request with DID communication invitation requesting credentials matching specific criteria.
- Step 3: Wallet Web queries saved credentials based on the criteria requested in `Step 2.`

  - If Wallet Web finds any match from the saved credentials, then it prepares presentation submission response based on credentials found.
  - If Wallet Web is unable to find any match from the saved credentials, then it goes through the list of saved manifest credentials sent by all previously connected issuers. If Wallet Web manages to find issuer who can issue requested credentials then it sends issue credential request to that issuer to issue authorization credential for given relying party or verifier.
    Authorization credential shared by issuer will be sent as part of the response to relying party. Relying party will use this authorization credential to connect to actual issuer to get actual credential.
  - If Wallet Web is unable to find any issuer who can fulfill requested criteria, then it updates presentation submission response with the message `no result for asked criteria`.

Note: Wallet Web asks user consent every time before sending back any CHAPI response. Wallet Web never reuses peer DIDs it created for one issuer with any other verifiers or vice versa.

Refer following pages for more details,

- A Complete end to end DIDComm flow examples can be found [here](https://trustbloc.readthedocs.io/en/latest/adapters.html#flows).
- More details about issuer flow can be found [here](https://github.com/trustbloc/edge-adapter/blob/master/docs/issuer/README.md)
- More details about relying party flow can be found [here](https://github.com/trustbloc/edge-adapter/blob/master/docs/rp/integration/wallets.md)

## Wallet Web DIDComm Flow with Blinded Routing

DIDComm flow with Wallet Web can also be achieved using blinded routing where identity of issuer and verifiers will be hidden from each other. More details about blinded routing can be found [here](https://trustbloc.readthedocs.io/en/latest/blinded_routing.html)

Blinded routing DIDComm flow is very similar to DIDComm flow explained in previous section except the following difference.

- After successfully connecting to an issuer or to a verifier,

  - Wallet Web requests a new peer DID from invitee (an issuer or verifier).
  - Wallet Web creates a new peer DID from router for new peer DID it received from invitee.
  - Wallet Web shares newly created peer DID with invitee.

## Customizing you Wallet Web

By default Wallet Web is configured to [orb DID](https://trustbloc.github.io/did-method-orb/) and [Ed25519Signature2018](https://w3c-ccg.github.io/lds-ed25519-2018/) signature suite for all its signing operations.
Wallet Web creates and assigns a new orb DID as a controller during user signup (i.e setup process).
But your Wallet Web's profile can also be customized to use any other DID methods and signature types.

To customize signing DIDs and signature types, go to Wallet Web's settings, you will see you current digital identity preferences showing your current DID (aka controller), verification method and signature type.
Once customized, Wallet Web will use these controller, verification method & signature type to add digital proofs to presentations presented from this wallet (i.e. signing verifiable presentation).

- **Customizing Signature type**: you can customize your signature suite by changing `Update Signature Type` section in digital identity preferences page.
- **Customizing Signing key**: you can customize your signing key (aka verification method) by selecting verification method of your choice from `Key ID` section.
- **Customizing controller**: you can customize your controller by choosing a DID of your choice in `Update Identity` section.
  You have to create a new orb DID or import any other DID in order to get them listed in `Update Identity` section. Remember, wallet always creates one orb DID for user during registration.
  - **Creating a new orb DID**: Navigate to `CREATE ORB DIGITAL IDENTITY` tab and create your new orb DID. You can even customize your orb DID by choosing `Key Type, Signature Suite` & `Key Purpose` of your choice while creating it.
    Once a new orb DID is created successfully it will show up in your list of DIDs in `DIGITAL IDENITTY PREFERNCE` page.
  - **Import a DID**: You can also import any DID along with its private key by navigating to `IMPORT ANY DIGITAL IDENTITY` tab. Wallet supports importing private keys in both `base58` & `jwk` formats.
    Any DID which can be resolved by wallet can be imported. Refer [this](https://github.com/trustbloc/wallet/blob/main/test/fixtures/wallet-web/config.json) for resolver configuration of demo web wallet instance.
    Once imported successfully it will show up in your list of DIDs in `DIGITAL IDENITTY PREFERNCE` page.
