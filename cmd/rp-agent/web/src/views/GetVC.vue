/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <button id='getVCBtn'>Get VC</button>
    </div>
</template>

<script>
    async function onClickGet() {
      const credentialQuery = JSON.parse('{"web": {"VerifiablePresentation": {}}}');
      const result = await navigator.credentials.get(credentialQuery);
        if(!result) {
           alert("null result");
            return;
        }
      alert(result.data);
    }
    export default {
        beforeCreate:function(){
            this.$polyfill.loadOnce()
        },
        mounted:function () {
            document.getElementById('getVCBtn').addEventListener('click', onClickGet);
        }
    }
</script>

