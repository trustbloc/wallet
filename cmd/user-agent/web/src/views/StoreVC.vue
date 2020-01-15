/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <textarea id="vcDataTextArea" readonly rows="30" cols="200"/>
        <br>
        <button id='storeVCBtn'>Store VC</button>
    </div>
</template>

<script>
    import * as polyfill from "credential-handler-polyfill";
    import * as webCredentialHandler from "web-credential-handler";

    async function handleWalletReceiveEvent() {
        const credentialEvent = await webCredentialHandler.receiveCredentialEvent();
        const vcData = credentialEvent.credential.data
        window.console.log('Received event:', vcData);
        document.getElementById('vcDataTextArea').value=vcData
        document.getElementById('storeVCBtn').addEventListener('click', () => {
            window.storeVC(credentialEvent)
        });
    }
    polyfill.loadOnce().then(handleWalletReceiveEvent)
    export default {
        polyfill
    }
</script>

