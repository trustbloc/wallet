/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from './walletManager'

var uuid = require('uuid/v4')

const keyType = "Ed25519"
const signType = "Ed25519Signature2018"
var allowedTypes = ['VerifiablePresentation', 'VerifiableCredential']

/**
 * RegisterWallet registers webcredential handler and manages wallet metadat in underlying db
 * @param polyfill and web credential handler
 * @class
 */
export class RegisterWallet extends WalletManager {
    constructor(polyfill, wcredHandler, didManager) {
        super()

        this.polyfill = polyfill
        this.wcredHandler = wcredHandler
        this.didManager = didManager
    }

    async register(user) {
        // create DID
        let did = await this.didManager.createDID(keyType, signType)

        // save DID
        await this.didManager.saveDID(`${user}_${uuid()}`, did)

        console.log(`created DID ${did.id} successfully for user ${this.username}`)

        try {
            await this.polyfill.loadOnce();
        } catch (e) {
            console.error('Error in loadOnce:', e);
            throw "failed to register wallet, please try again later"
        }

        const registration = await this.wcredHandler.installHandler({url: '/worker.html'})

        await registration.credentialManager.hints.set(
            'edge', {
                name: user,
                enabledTypes: allowedTypes
            });

        // clear any existing date
        await this.clear()

        // save wallet metadata
        await this.storeWalletMetadata(user, {
            signatureType: signType,
            did: did.id
        })
    }

    async unregister() {
        await this.clear()

        try {
            await this.polyfill.loadOnce();
        } catch (e) {
            console.error('Error in loadOnce:', e);
            return
        }

        await this.wcredHandler.uninstallHandler({url: '/worker.html'})
    }
}