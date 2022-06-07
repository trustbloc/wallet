/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import jp from 'jsonpath';
import { v4 as uuidv4 } from 'uuid';
import { PresentationExchange } from '../common/presentationExchange';
import { BlindedRouter, CredentialManager, DIDComm, DIDManager } from '@trustbloc/wallet-sdk';
import { getCredentialType, normalizeQuery } from '../';
import { toRaw, markRaw } from 'vue';

const manifestType = 'IssuerManifestCredential';
const didDocReqMsgType = 'https://trustbloc.dev/adapter/1.0/diddoc-req';
const didDocResMsgType = 'https://trustbloc.dev/adapter/1.0/diddoc-resp';

var blindedRoutingDisabled = {
  sharePeerDID: () => {},
};

/**
 * WalletGetByQuery provides CHAPI get vp features
 * @param agent instance & credential event
 * @class
 */
export class WalletGetByQuery {
  constructor(agent, protocolHandler, opts, user) {
    this.agent = agent;
    this.protocolHandler = protocolHandler;
    let { query } = this.protocolHandler.getEventData();

    let presExchQuery = normalizeQuery(jp.query(query, `$[?(@.type=="PresentationExchange")]`));
    this.invitation = markRaw(jp.query(query, '$[?(@.type=="DIDConnect")].invitation'));

    console.log('presExchQuery', JSON.stringify(presExchQuery, null, 2));
    /*
              TODO:
               * current assumption - expecting only one governance VC in request, may be support for multiple.
               * correlate governance VC with requesting party so that consent for trust gets shown only once.
               * verify governance VC proof.
               * verify requesting party in governance framework to make sure this party of behaving properly.
               * request party to get challenged to produce a VP that the governance credential agency has accredited them.
             */
    let govnVCs = jp.query(
      query,
      `$[?(@.type=="DIDConnect")].credentials[?(@.type.indexOf("GovernanceCredential") != -1)]`
    );
    this.govnVC = govnVCs.length > 0 ? govnVCs[0] : undefined;

    this.blindedRouter = blindedRoutingDisabled;
    if (opts.blindedRouting === true) {
      this.blindedRouter = new BlindedRouter(agent);
    }
    this.didManager = new DIDManager({ agent, user });
    this.didcomm = new DIDComm({ agent, user });
    this.credentialManager = new CredentialManager({ agent, user });
    this.presenationExchange = new PresentationExchange(presExchQuery[0].credentialQuery[0]);
  }

  requirementDetails() {
    return this.presenationExchange.requirementDetails();
  }

  async connectMediator() {
    await this.agent.mediator.reconnectAll();
  }

  async getPresentationSubmission(token) {
    let { contents } = await this.credentialManager.getAll(token);
    let vcs = Object.keys(contents).map((k) => contents[k]);

    return this.presenationExchange.createPresentationSubmission(vcs);
  }

  async createAndSendPresentation(user, presentationSubmission) {
    console.log('user:', user);
    console.log('presentationSubmission', presentationSubmission);
    console.log('here 80');
    let newPresentationSubmission;
    try {
      console.log('here 82');
      console.log('invitaion', this.invitation);
      if (this.invitation.length > 0) {
        console.log('here 85');
        console.log('user:', user);
        console.log('user.profile:', user.profile);

        newPresentationSubmission = await this._getAuthorizationCredentials(
          presentationSubmission,
          user.profile
        );
        console.log('here 92');
      }
      console.log('here 95');

      let { controller, proofType, verificationMethod } = user.preference;
      let { domain, challenge } = this.protocolHandler.getEventData();
      let { token } = user.profile;
      console.log('here 100');

      let { presentation } = await this.credentialManager.present(
        token,
        { presentation: newPresentationSubmission },
        {
          controller,
          proofType,
          verificationMethod,
          domain,
          challenge,
        }
      );
      console.log('here 113');
      console.log('presentation,', presentation);
      this.protocolHandler.present(presentation);
    } catch (e) {
      console.error(e);
    }
  }

  cancel() {
    this.protocolHandler.cancel();
  }

  sendNoCredentials() {
    this.protocolHandler.present({});
  }

  async _getAuthorizationCredentials(presentationSubmission, profile) {
    let { token } = profile;
    console.log('token,', token);
    console.log('this.invitation[0]', this.invitation[0]);
    console.log('here 139');
    let rpConn;
    try {
      rpConn = await this.didcomm.connect(token, this.invitation[0], {
        userAnyRouterConnection: true,
        waitForCompletion: true,
      });
    } catch (e) {
      console.log('error 141', e);
    }

    console.log('here 144');

    // share peer DID with RP for blinded routing
    await this.blindedRouter.sharePeerDID(rpConn.result);
    console.log('here 148');

    // request new peer DID from RP
    let didDocRes = await this.agent.messaging.send({
      connection_ID: rpConn.result.ConnectionID,
      message_body: {
        '@id': uuidv4(),
        '@type': didDocReqMsgType,
        sent_time: new Date().toJSON(),
      },
      await_reply: { messageType: didDocResMsgType, timeout: 20000000000 }, //TODO (#531): Reduce timeout once EDV storage speed is improved. Note: this value used to be 2 seconds (now it's 20).
    });
    console.log('here 160');

    // response could be byte array from RP
    let rpDIDDoc = Array.isArray(didDocRes.response.message.data.didDoc)
      ? JSON.parse(String.fromCharCode.apply(String, didDocRes.response.message.data.didDoc))
      : didDocRes.response.message.data.didDoc;

    let peerDID = (await this.didManager.createPeerDID(token)).didDocument;
    let agent = this.agent;
    let credManager = this.credentialManager;
    let acceptCredPool = new Map();

    await Promise.all(
      presentationSubmission.verifiableCredential.map(async (vc, index) => {
        if (getCredentialType(vc.type) != manifestType) {
          return;
        }

        let connectionID = await credManager.getManifestConnection(token, vc.id);
        let { result } = await agent.didexchange.queryConnectionByID({ id: connectionID });

        // TODO `request_credential.requests~attach.data.json.subjectDID` to be removed once adapters are updated
        let resp = await agent.issuecredential.sendRequest({
          my_did: result.MyDID,
          their_did: result.TheirDID,
          request_credential: {
            'requests~attach': [
              {
                lastmod_time: new Date(),
                data: {
                  json: {
                    subjectDID: peerDID.id,
                    subjectDIDDoc: {
                      id: peerDID.id,
                      doc: peerDID,
                    },
                    requestingPartyDIDDoc: {
                      id: rpDIDDoc.id,
                      doc: rpDIDDoc,
                    },
                  },
                },
              },
            ],
          },
        });

        console.log('sent request credential message', resp.piid);

        acceptCredPool.set(resp.piid, { index });
      })
    );

    console.log(`${acceptCredPool.size} issue credential requests sent`);

    await waitForCredentials(agent, acceptCredPool);

    acceptCredPool.forEach(function (value) {
      presentationSubmission.verifiableCredential[value.index] = value.credential;
    });

    return presentationSubmission;
  }
}

// retainOnlySelected retain only selected VCs and their respective descriptors
function retainOnlySelected(presentationSubmission, selectedIndexes) {
  let descriptors = [];
  let vcs = [];

  let vcCount = 0;
  selectedIndexes.forEach(function (selected, index) {
    presentationSubmission.verifiableCredential.forEach(function (vc, vcIndex) {
      if (selected && index == vcIndex) {
        vcs.push(vc);

        let vcDescrs = jp.query(
          presentationSubmission,
          `$.presentation_submission.descriptor_map[?(@.path=="$.verifiableCredential[${vcIndex}]")].id`
        );
        vcDescrs.forEach(function (id) {
          descriptors.push({ id, path: `$.verifiableCredential[${vcCount}]` });
        });

        vcCount++;
      }
    });
  });

  presentationSubmission.verifiableCredential = vcs;
  presentationSubmission.presentation_submission.descriptor_map = descriptors;

  return presentationSubmission;
}

function waitForCredentials(agent, pool) {
  let processed = 0;
  // eslint-disable-next-line no-async-promise-executor
  return new Promise(async (resolve, reject) => {
    setTimeout(() => reject(new Error('timeout waiting for incoming credentials')), 600000); // TODO (#531): Reduce timeout once EDV storage speed is improved.

    for (;;) {
      // eslint-disable-next-line no-await-in-loop
      let resp = await agent.issuecredential.actions();
      for (let action of resp.actions) {
        if (pool.has(action.PIID)) {
          let credID = uuidv4();
          agent.issuecredential.acceptCredential({
            piid: action.PIID,
            names: [credID],
          });

          // eslint-disable-next-line no-await-in-loop
          let metadata = await getCredentialMetadata(agent, credID);
          // eslint-disable-next-line no-await-in-loop
          let credential = await agent.verifiable.getCredential(metadata);

          pool.get(action.PIID).credential = JSON.parse(credential.verifiableCredential);

          processed++;
        }
      }

      if (processed < pool.size) {
        console.log(`received ${processed} out of ${pool.size} credentials, retrying`);
        // eslint-disable-next-line no-await-in-loop
        await sleep(1000);
        continue;
      }

      break;
    }

    console.log(`received all ${processed} credentials`);
    resolve();
  });
}

async function getCredentialMetadata(agent, name) {
  const retries = 30; // TODO (#531): Reduce number of retries once EDV storage speed is improved.
  for (let i = 0; i < retries; i++) {
    try {
      // eslint-disable-next-line no-await-in-loop
      return await agent.verifiable.getCredentialByName({
        name: name,
      });
    } catch (e) {
      console.log(`credential '${name}' not found, retrying`);
    }

    // eslint-disable-next-line no-await-in-loop
    await sleep(500);
  }

  throw new Error('no credential found');
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
