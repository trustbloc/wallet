/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
<template>
    <div class="content">
    <md-card>
        <md-card-header data-background-color="blue">
            <h4 class="title">If you see registration pop-up, you have successfully registered</h4>
        </md-card-header>
    </md-card>
    </div>
</template>
<script>
    import Swal from 'sweetalert2'
    async function installHandler() {
        try {
            await window.$polyfill.loadOnce();
        } catch(e) {
            window.console.error('Error in loadOnce:', e);
        }

        window.console.log('Polyfill loaded.');

        const registration = await  window.$webCredentialHandler
            .installHandler({url: '/worker.html'})

        await registration.credentialManager.hints.set(
            'edge', {
                name: 'EdgeUser',
                enabledTypes: ['VerifiablePresentation', 'VerifiableCredential']
            });

        Swal.fire({
            title: 'Registration is completed',
            showClass: {
                popup: 'animated fadeInDown faster'
            },
            hideClass: {
                popup: 'animated fadeOutUp faster'
            }
        })
    }

    export default {
        beforeCreate:function(){
            window.$polyfill=this.$polyfill,
            window.$webCredentialHandler=this.$webCredentialHandler
            installHandler();
        },
    }
</script>

<style lang="scss" scoped>
    .md-toolbar + .md-toolbar {
        margin-top: 16px;
    }
</style>
