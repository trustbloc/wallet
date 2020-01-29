/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <select id="selectVC">
        <option value="0">Select VC</option>
        </select>
        <br>
        <button id='getVCBtn'>Send VC</button>
    </div>
</template>

<script>
    async function handleWalletReceiveEvent() {
        const credentialEvent = await window.$webCredentialHandler.receiveCredentialEvent();
        window.console.log('Received event:', credentialEvent);
        document.getElementById('getVCBtn').addEventListener('click', () => {
            var e = document.getElementById("selectVC");
            var vcValue = e.options[e.selectedIndex].value;
            if (vcValue == 0) {
                alert("please select vc")
                return
            }
            window.getVC(credentialEvent,vcValue)
        });

    }
    export default {
        beforeCreate:function(){
            // Add an event listener
            document.addEventListener("afterLoadingWasm", function() {
                var select = document.getElementById("selectVC");
                window.populateVC(select)
            });
            window.$webCredentialHandler=this.$webCredentialHandler
            this.$polyfill.loadOnce().then(handleWalletReceiveEvent)
        },
    }





</script>

