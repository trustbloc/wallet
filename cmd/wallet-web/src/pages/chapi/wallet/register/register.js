/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from './walletManager'
import {connectToMediator} from "../didcomm/mediator";
import {DIDManager} from "../didmgmt/didManager";


var uuid = require('uuid/v4')

const keyType = "Ed25519"
const signType = "Ed25519Signature2018"
var allowedTypes = ['VerifiablePresentation', 'VerifiableCredential']

/**
 * RegisterWallet registers webcredential handler and manages wallet metadata in underlying db
 * @param polyfill, web credential handler, agent instance.
 * @class
 */
export class RegisterWallet extends WalletManager {
    constructor(polyfill, wcredHandler, agent, opts) {
        super(agent)

        this.agent = agent
        this.polyfill = polyfill
        this.wcredHandler = wcredHandler
        this.didManager = new DIDManager(agent, opts)
        this.mediatorEndpoint = opts.walletMediatorURL
        this.credentialMediator = opts.credentialMediatorURL
    }

    // wallet user registration and setup process
    async register(user) {
        // register mediator
        this._connectToMediator()

        // create and assign DID to this user & save
        this._assignDID(user)
    }

    // install credential handler polyfill handlers
    async installHandlers(user) {
        if (!this.skipPolyfill) {
            try {
                await this.polyfill.loadOnce(this.credentialMediator);
            } catch (e) {
                console.error('Error in loadOnce:', e);
                throw "failed to register wallet, please try again later"
            }

            const registration = await this.wcredHandler.installHandler({url: `${__webpack_public_path__}worker`})

            await registration.credentialManager.hints.set(
                'edge', {
                    name: user,
                    enabledTypes: allowedTypes
                });
        }
    }

    // uninstall credential handler polyfill handlers
    async uninstallHandlers() {
        try {
            await this.polyfill.loadOnce(this.credentialMediator);
        } catch (e) {
            console.error('Error in loadOnce:', e);
            return
        }

        await this.wcredHandler.uninstallHandler({url: `${__webpack_public_path__}worker`})
    }

    async _connectToMediator() {
        console.time('time tracking: register mediator time');
        if (this.mediatorEndpoint) {
            try {
                let resp = await this.agent.mediator.getConnections()
                if (!resp.connections || resp.connections.length == 0) {
                    await connectToMediator(this.agent, this.mediatorEndpoint)
                }

                console.debug(`registered with mediator successfully for user`)
            } catch (e) {
                // mediator registration isn't mandatory, so handle errors
                console.warn("unable to connect to mediator, registered wallet may not support DID Exchange, cause:", e.toString())
            }
        } else {
            console.warn("unable to find to mediator wallet URL, registered wallet may not support DID Exchange")
        }
        console.timeEnd('time tracking: register mediator time');
    }
    
    async _assignDID(user){
        console.time('time tracking: create and save DID');

        console.time('time tracking: create trustbloc did time');
        let did = await this.didManager.createTrustBlocDID(keyType, signType)
        console.timeEnd('time tracking: create trustbloc did time');

        console.time('time tracking: store wallet metadata time');
        await this.storeWalletMetadata(user, {
            signatureType: signType,
            did: did.id
        })
        console.timeEnd('time tracking: store wallet metadata time');

        // save DID
        console.time('time tracking: save did time');
        await this.didManager.saveDID(user, `${user}_${uuid()}`, signType, did)
        console.timeEnd('time tracking: save did time');

        console.debug(`created DID ${did.id} successfully for user ${user}`)

        console.timeEnd('time tracking: create and save DID');
    }
}
