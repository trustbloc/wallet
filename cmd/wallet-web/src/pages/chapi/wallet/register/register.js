/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from './walletManager'
import {AgentMediator} from "../didcomm/mediator";
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
        super()

        this.polyfill = polyfill
        this.wcredHandler = wcredHandler
        this.didManager = new DIDManager(agent, opts)
        this.mediator = new AgentMediator(agent)
        this.mediatorEndpoint = opts.walletMediatorURL
        this.credentialMediator = opts.credentialMediatorURL
    }

    // wallet user registration and setup process
    async register(user) {
        // register mediator
        let invitation
        if (this.mediatorEndpoint) {
            //TODO read router endpoint from config
            try {
                let connected = await this.mediator.isAlreadyConnected()
                if (!connected) {
                    await this.mediator.connect(this.mediatorEndpoint)
                }

                invitation = await this.mediator.createInvitation()

                console.debug(`registered with mediator successfully for user ${user}`)
            } catch (e) {
                // mediator registration isn't mandatory, so handle errors
                console.warn("unable to connect to mediator, registered wallet may not support DID Exchange, cause:", e.toString())
            }
        } else {
            console.warn("unable to find to mediator wallet URL, registered wallet may not support DID Exchange")
        }

        // create DID
        let did = await this.didManager.createTrustBlocDID(keyType, signType)

        // save wallet metadata
        // TODO wallet metadata to be saved after saveDID [ #332]
        await this.storeWalletMetadata(user, {
            signatureType: signType,
            did: did.id,
            invitation: invitation
        })

        // save DID
        await this.didManager.saveDID(user, `${user}_${uuid()}`, signType, did)

        console.debug(`created DID ${did.id} successfully for user ${user}`)
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
}
