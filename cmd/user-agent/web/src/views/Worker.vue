/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<script>
    async function activateHandler() {
        try {
            await window.$polyfill.loadOnce();
        } catch(e) {
            window.console.error('Error in loadOnce:', e);
        }

        return window.$webCredentialHandler.activateHandler({
            async store(event) {
                window.console.log('Received store() event:', event.event);
                return {type: 'redirect', url: '/StoreVC'};
            }
        })
    }
    export default {
        beforeCreate:function(){
            window.$polyfill=this.$polyfill,
            window.$webCredentialHandler=this.$webCredentialHandler
            activateHandler();
        },
    }
</script>

