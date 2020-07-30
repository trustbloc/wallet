# User Agent Web Wallet

 User agent Web Wallet is based on [CHAPI](https://w3c-ccg.github.io/credential-handler-api/) developed using [credential handler polyfill](https://github.com/digitalbazaar/credential-handler-polyfill), [Vue](https://vuejs.org)

## Steps to start user agent web wallet

- run below command to build and start user agent

    ```bash
    make user-agent-start
    ```

## Web wallet demo

Once user agent is started using previous step,

- Login to web wallet from user web wallet [dashboard](https://localhost:8091/dashboard). 

  Proceed with pre-filled username/password for the login. Once login is successful, you will get a prompt from your browser to allow wallet to manage credentials, choose 'Allow'. 
- Go to web wallet [demo page](https://localhost:8091/webwallet)

  Web wallet demo page can be used to perform all supported wallet operations. Click on sample requests based on operation you wish to perform and click on `STORE` or `GET` buttons. 
  
  Web wallet allows only 2 kind of requests for requesting and storing credentials,
   1. Get - A web app (a Relying Party or verifier) can request a credential using get requests.
   2. Store - A web app (for example, a credential issuer such as a university or institution) can ask to store a credential using store requests.
  
 
## Operations supported by user agent web wallet
User agent web wallet supports various wallet operations and also [DIDComm](https://github.com/hyperledger/aries-rfcs/blob/master/concepts/0005-didcomm/README.md). 
- **DID auth**: This operation can be used when a website that wants user to authenticate (prove they are authorized to act as a particular DID subject, aka "DIDAuth" over wallet)
- **Store Credential**: This operation can be used when an issuer website wants to give user a credential to store in their wallet
- **Get Credential**: This operation can be used when a website (a Relying Party or verifier) wants to request a credential from user.
- **DID Connect**: This operation can be used by a website(or agent) to perform [DID Exchange](https://github.com/hyperledger/aries-rfcs/blob/master/features/0023-did-exchange/README.md) with web wallet for all future DID Comm operations.  
- **Presentation Exchange**: This operation can be used by a website (a Relying Party or verifier) to request submission of proof from wallet in align with requester's [proof requirements](https://identity.foundation/presentation-exchange/). This request can also be made in conjuction with DID Connect operation, so that wallet can request further information from requester (a Relying party or verifier) over DID comm channel.

