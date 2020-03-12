/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <div> Create QR Code</div>
        <br>
        <div>
            <select id="selectVC">
                <option value="0">Select VC</option>
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
        document.getElementById('getVCBtn').addEventListener('click', async () => {
            // Get the VC ID from UI selection
            let e = document.getElementById("selectVC");
            let vcID = e.options[e.selectedIndex].value;
            if (vcID == 0) {
                alert("please select vc")
                return
            }

            // Get the VC data
            let vcData
            await window.$aries.verifiable.getCredential({
                id: vcID
            }).then(resp => {
                    vcData = JSON.stringify(JSON.parse(resp.verifiableCredential))
                }
            ).catch(err =>
                console.log('generateQRCode - get vc failed : errMsg=' + err)
            )

            // Generate QR code
            let QRCode = require('qrcode')
            QRCode.toDataURL(vcData, function (err, url) {
                let canvas = document.getElementById('qr-result')
                canvas.src = url
            })
        });
    }

    export default {
        beforeCreate: async function () {
            // Load the Credentials in the drop down
            let aries = await this.$arieslib
            await aries.verifiable.getCredentials()
                .then(resp => {
                        const data = resp.result
                        if (data && data.length !== 0) {
                            let dropdown = document.getElementById('selectVC');
                            let option;
                            for (let i = 0; i < data.length; i++) {
                                option = document.createElement('option');
                                option.text = data[i].name;
                                option.value = data[i].id;
                                dropdown.add(option);
                            }
                        } else {
                            console.log('no credentials exists')
                        }
                    }
                ).catch(err => {
                        console.log('get credentials failed : errMsg=' + err)
                    }
                )

            window.$webCredentialHandler = this.$webCredentialHandler
            window.$aries = aries

            generateQRCode()
        },
    }

</script>
