/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>

    </div>
</template>

<script>
    export default {
        beforeCreate: async function () {
            let opts = await this.$trustblocStartupOpts
            try {
                await this.$polyfill.loadOnce(opts.credentialMediatorURL);
            } catch (e) {
                console.error('Error in loadOnce:', e);
            }

            console.log('Worker Polyfill loaded..!!');

            return this.$webCredentialHandler.activateHandler({
                mediatorOrigin : opts.credentialMediatorURL,
                async get(event) {
                    console.log('Received get() event:', event.event);
                    return {type: 'redirect', url: '/GetFromWallet'};
                },
                async store(event) {
                    console.log('Received store() event:', event.event);
                    return {type: 'redirect', url: '/StoreInWallet'};
                }
            })
        }
    }
</script>
