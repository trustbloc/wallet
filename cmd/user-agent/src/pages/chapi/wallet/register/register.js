/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

/**
 * RegisterWallet registers webcredential handler
 * @param polyfill and web credential handler
 * @class
 */
export class RegisterWallet {
    constructor(polyfill, wcredHandler) {
        this.polyfill = polyfill
        this.wcredHandler = wcredHandler
    }

    async register(user) {
        try {
            await this.polyfill.loadOnce();
        } catch(e) {
           console.error('Error in loadOnce:', e);
           throw "failed to register wallet, please try again later"
        }

        const registration = await  this.wcredHandler.installHandler({url: '/worker.html'})

        await registration.credentialManager.hints.set(
            'edge', {
                name: user,
                enabledTypes: ['VerifiablePresentation', 'VerifiableCredential']
            });
    }
}