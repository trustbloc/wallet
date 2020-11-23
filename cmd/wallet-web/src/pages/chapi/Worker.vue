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
            let opts = this.$store.getters.getAgentOpts
            try {
                await this.$polyfill.loadOnce(opts ? opts.credentialMediatorURL : undefined);
            } catch (e) {
                console.error('Error in loadOnce:', e);
            }

            console.log('Worker Polyfill loaded..!!');

            return this.$webCredentialHandler.activateHandler({
                mediatorOrigin : opts.credentialMediatorURL,
                async get(event) {
                    console.log('Received get() event:', event.event);
                    return {type: 'redirect', url: `${__webpack_public_path__}GetFromWallet`};
                },
                async store(event) {
                    console.log('Received store() event:', event.event);
                    return {type: 'redirect', url: `${__webpack_public_path__}StoreInWallet`};
                }
            })
        }
    }
</script>
