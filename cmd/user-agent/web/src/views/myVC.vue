/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <div> Create QR Code </div>
        <br>
          <div>
            <select id="selectVC">
                <option value="0" >Select VC</option>
            </select>
            <br>
            <br>
            <button id='getVCBtn'>Get QR Code</button>
            <br>
            <br>
            <img src="" id="qr-result"/>
         </div>
        </div>
</template>
<script>
    async function generateQRCode() {
        window.onload = function () {
            document.getElementById('getVCBtn').addEventListener('click', () => {
            var e = document.getElementById("selectVC");
            var vcValue = e.options[e.selectedIndex].value;
            if (vcValue == 0) {
                alert("please select vc")
                return
            }
             window.createQRCode(vcValue);
            });
        }
    }

    export default{
        beforeCreate:function(){
            // Add an event listener
            document.addEventListener("afterLoadingWasm", function() {
                var select = document.getElementById("selectVC");
                window.populateVC(select)
            });
            window.$webCredentialHandler=this.$webCredentialHandler
            generateQRCode()
        },
    }

</script>
