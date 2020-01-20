/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <textarea id="vcDataTextArea" readonly rows="30" cols="200"/>
        <br>
        <input id="friendlyName" placeholder="friendly name">
        <br>
        <button id='storeVCBtn'>Store VC</button>
    </div>
</template>

<script>
    async function handleWalletReceiveEvent() {
        const credentialEvent = await window.$webCredentialHandler.receiveCredentialEvent();
        const vcData = credentialEvent.credential.data
        window.console.log('Received vc data:', vcData);
        document.getElementById('vcDataTextArea').value=vcData
        document.getElementById('storeVCBtn').addEventListener('click', () => {
            const friendlyName =  document.getElementById('friendlyName').value
            if (!friendlyName) {
                alert("please enter friendly name")
                return
            }
            window.storeVC(credentialEvent,friendlyName)
        });
    }
    export default {
        beforeCreate:function(){
                window.$webCredentialHandler=this.$webCredentialHandler
            this.$polyfill.loadOnce().then(handleWalletReceiveEvent)
        },
    }
</script>

