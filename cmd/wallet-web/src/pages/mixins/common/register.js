/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {
  CollectionManager,
  connectToMediator,
  DIDManager,
  WalletUser,
} from '@trustbloc/wallet-sdk';

const proofType = 'Ed25519Signature2018';
const waitForStateComplete = true;

/**
 * RegisterWallet registers webcredential handler and manages wallet metadata in underlying db
 * @param agent instance, wallet options.
 * @class
 */
export class RegisterWallet {
  constructor(agent, opts) {
    this.agent = agent;
    this.didManager = new DIDManager(agent, opts);
    this.collectionManager = new CollectionManager(agent, opts);
    this.mediatorEndpoint = opts.walletMediatorURL;
  }

  // wallet user registration and setup process
  async register(profile, callback) {
    // register mediator, create and save DID
    let failure;
    try {
      await Promise.all([
        this._connectToMediator(),
        this._assignDID(profile),
        this._createDefaultVault(profile),
      ]);
    } catch (e) {
      console.error('failed to setup wallet user', e);
      failure = e;
    }

    await this.agent.store.flush();

    if (callback) {
      callback(failure);
    }
  }

  async _connectToMediator() {
    console.time('time tracking: register mediator time');
    if (this.mediatorEndpoint) {
      try {
        await connectToMediator(this.agent, this.mediatorEndpoint, { waitForStateComplete });

        console.debug(`registered with mediator successfully for user`);
      } catch (e) {
        // mediator registration isn't mandatory, so handle errors
        console.warn(
          'unable to connect to mediator, registered wallet may not support DID Exchange, cause:',
          e.toString()
        );
        throw e;
      }
    } else {
      console.warn(
        'unable to find to mediator wallet URL, registered wallet may not support DID Exchange'
      );
    }
    console.timeEnd('time tracking: register mediator time');
  }

  async _assignDID(profile) {
    console.time('time tracking: create did time');
    let { user, name, token } = profile;
    let didManager = new DIDManager({ agent: this.agent, user });
    let docres = await didManager.createOrbDID(token, {
      purposes: ['authentication', 'assertionMethod'],
    });

    let controller = docres.didDocument.id;

    if (docres.didDocumentMetadata.equivalentId && docres.didDocumentMetadata.equivalentId.length) {
      controller = docres.didDocumentMetadata.equivalentId[0];
      console.info('use equivalentId ' + controller);
    }

    let walletUser = new WalletUser({ agent: this.agent, user });
    await walletUser.savePreferences(token, { name, controller, proofType });
    console.timeEnd('time tracking: create did time');
  }

  async _createDefaultVault(profile) {
    let { user, token } = profile;
    let name = 'Default Vault';
    let collectionManager = new CollectionManager({ agent: this.agent, user });
    let vaultID = await collectionManager.create(token, { name });
    if (vaultID) {
      console.info('Default vault is created');
    }
  }
}
