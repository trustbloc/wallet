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

        document.getElementById('getVCBtn').addEventListener('click', async () => {
            // Get the VC ID from UI selection
            var e = document.getElementById("selectVC");
            var vcID = e.options[e.selectedIndex].value;
            if (vcID == 0) {
                alert("please select vc")
                return
            }

            // Get the VC data
            let data
            await window.$aries.verifiable.getCredential({
                id: vcID
            }).then(resp => {
                    data = JSON.stringify(JSON.parse(resp.verifiableCredential))
                }
            ).catch(err => {
                data = err
                console.log('get vc failed : errMsg=' + err)
            })

            // Call Credential Handler callback
            credentialEvent.respondWith(new Promise(function (resolve) {
                return resolve({
                    dataType: "Response",
                    data: data
                });
            }))
        });
    }

    export default {
        beforeCreate: async function () {
            // Load the Credentials in the drop down
            let aries = await this.$arieslib
            aries.verifiable.getCredentials()
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
            this.$polyfill.loadOnce().then(handleWalletReceiveEvent)
            window.$aries = aries
        },
    }
</script>

