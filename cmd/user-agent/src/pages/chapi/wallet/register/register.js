/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from './walletManager'
import {AgentMediator} from "../didcomm/connections";
import {DIDManager} from "../didmgmt/didManager";


var uuid = require('uuid/v4')

const keyType = "Ed25519"
const signType = "Ed25519Signature2018"
var allowedTypes = ['VerifiablePresentation', 'VerifiableCredential']

/**
 * RegisterWallet registers webcredential handler and manages wallet metadat in underlying db
 * @param polyfill, web credential handler, aries agent instance and DID manager
 * @class
 */
export class RegisterWallet extends WalletManager {
    constructor(polyfill, wcredHandler, aries, trustblocAgent, opts) {
        super()

        this.polyfill = polyfill
        this.wcredHandler = wcredHandler
        this.didManager = new DIDManager(aries, trustblocAgent, opts)
        this.mediator = new AgentMediator(aries)
        this.mediatorEndpoint = opts.walletMediatorURL
        this.credentialMediator = opts.credentialMediatorURL
    }

    async register(user) {
        // create DID
        let did = await this.didManager.createDID(keyType, signType)

        // save DID
        await this.didManager.saveDID(`${user}_${uuid()}`, signType, did)

        console.log(`created DID ${did.id} successfully for user ${this.username}`)

        if (!this.skipPolyfill) {
            try {
                await this.polyfill.loadOnce(this.credentialMediator);
            } catch (e) {
                console.error('Error in loadOnce:', e);
                throw "failed to register wallet, please try again later"
            }

            const registration = await this.wcredHandler.installHandler({url: `/worker`})

            await registration.credentialManager.hints.set(
                'edge', {
                    name: user,
                    enabledTypes: allowedTypes
                });
        }

        // register mediator
        let invitation
        if (this.mediatorEndpoint) {
            try {
                let connected = await this.mediator.isAlreadyConnected()
                if (!connected) {
                    //TODO read router endpoint from config
                    this.disconnectMediator = await this.mediator.connect(this.mediatorEndpoint)
                }
                invitation = await this.mediator.createInvitation()
            } catch (e) {
                // mediator registration isn't mandatory, so handle errors
                console.warn("unable to connect to mediator, registered wallet may not support DID Exchange, cause:", e.toString())
            }
        } else {
            console.warn("unable to find to mediator wallet URL, registered wallet may not support DID Exchange")
        }

        // clear any existing date
        await this.clear()

        // save wallet metadata
        await this.storeWalletMetadata(user, {
            signatureType: signType,
            did: did.id,
            invitation: invitation
        })
    }

    async unregister() {
        await this.clear()

        if (this.disconnectMediator) {
            this.disconnectMediator()
        }

        try {
            await this.polyfill.loadOnce(this.credentialMediator);
        } catch (e) {
            console.error('Error in loadOnce:', e);
            return
        }

        await this.wcredHandler.uninstallHandler({url: '/worker.html'})
    }
}
