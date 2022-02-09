/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { CHAPIEventHandler } from '../';
import { BlindedRouter, CredentialManager, DIDComm } from '@trustbloc/wallet-sdk';
import jp from 'jsonpath';

const manifestCredType = 'IssuerManifestCredential';
const governanceCredType = 'GovernanceCredential';

var blindedRoutingDisabled = {
  sharePeerDID: () => {},
};

/**
 * DIDConn provides CHAPI did connection/exchange features
 * @param agent instance & credential event
 * @class
 */
export class DIDConn {
  constructor(agent, profile, startupOpts, protocolHandler) {
    this.didcomm = new DIDComm({ agent, user: profile.user });
    this.blindedRouter = startupOpts.blindedRouting
      ? new BlindedRouter(agent)
      : blindedRoutingDisabled;
    this.protocolHandler = protocolHandler;
    this.credentialManager = new CredentialManager({ agent, user: profile.user });
    this.profile = profile;

    let { domain, challenge, invitation, credentials = [] } = this.protocolHandler.getEventData();

    this.domain = domain;
    this.challenge = challenge;
    this.invitation = invitation;

    /*
              TODO:
               * current assumption - expecting only one governance VC in request, may be support for multiple.
               * correlate governance VC with requesting party so that consent for trust gets shown only once.
               * verify governance VC proof.
               * verify requesting party in governance framework to make sure this party of behaving properly.
               * request party to get challenged to produce a VP that the governance credential agency has accredited them.
             */
    // govn vc
    let govnVCs = jp.query(credentials, `$[?(@.type.indexOf('${governanceCredType}') != -1)]`);
    this.govnVC = govnVCs.length > 0 ? govnVCs[0] : undefined;

    // manifest VC
    let manifest = jp.query(credentials, `$[?(@.type.indexOf('${manifestCredType}') != -1)]`);
    this.manifestVC = manifest.length > 0 ? manifest[0] : undefined;

    // user credentials
    this.userCredentials = jp.query(
      credentials,
      `$[?(@.type.indexOf('${governanceCredType}') == -1 && @.type.indexOf('${manifestCredType}') == -1)]`
    );
  }

  async connect(preference) {
    let connection = await this.didcomm.connect(this.profile.token, this.invitation, {
      userAnyRouterConnection: true,
    });

    await this.blindedRouter.sharePeerDID(connection.result);

    // save credentials + manifestVCs
    let saveQueue = [];
    if (this.manifestVC) {
      saveQueue.push(
        this.credentialManager.saveManifestVC(
          this.profile.token,
          this.manifestVC,
          connection.result.ConnectionID
        )
      );
    }
    if (this.userCredentials.length > 0) {
      saveQueue.push(
        this.credentialManager.save(this.profile.token, { credentials: this.userCredentials })
      );
    }
    await Promise.all(saveQueue);

    await this._createConnectionResponse(connection.result, preference);
  }

  async _createConnectionResponse(connection, preference) {
    let { controller, proofType, verificationMethod } = preference;

    let credential = {
      '@context': [
        'https://www.w3.org/2018/credentials/v1',
        'https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld',
      ],
      issuer: controller,
      issuanceDate: new Date(),
      type: ['VerifiableCredential', 'DIDConnection'],
      credentialSubject: {
        id: connection.ConnectionID,
        threadID: connection.ThreadID,
        inviteeDID: connection.MyDID,
        inviterDID: connection.TheirDID,
        inviterLabel: connection.TheirLabel,
        connectionState: connection.State,
      },
    };

    let issued = await this.credentialManager.issue(this.profile.token, credential, {
      controller,
      proofType,
      verificationMethod,
    });

    let { presentation } = await this.credentialManager.present(
      this.profile.token,
      { rawCredentials: [issued.credential] },
      {
        controller,
        proofType,
        verificationMethod,
        domain: this.domain,
        challenge: this.challenge,
      }
    );

    this.protocolHandler.present(presentation);
  }

  cancel() {
    this.protocolHandler.cancel();
  }
}
