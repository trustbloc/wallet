/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<script>
    async function installHandler() {
        try {
            await window.$polyfill.loadOnce();
        } catch(e) {
            window.console.error('Error in loadOnce:', e);
        }

        window.console.log('Polyfill loaded.');

        const registration = await  window.$webCredentialHandler
            .installHandler({url: '/Worker'})

        await registration.credentialManager.hints.set(
            'test', {
                name: 'TestUser',
                enabledTypes: ['VerifiablePresentation', 'VerifiableCredential', 'AlumniCredential']
            });
        alert("Registration is Completed")
    }
    export default {
        beforeCreate:function(){
            window.$polyfill=this.$polyfill,
            window.$webCredentialHandler=this.$webCredentialHandler
            installHandler();
        },
    }
</script>

