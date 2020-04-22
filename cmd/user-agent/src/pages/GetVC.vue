/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="md-layout">
            <div
                    class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
                <md-card class="md-card-plain">
                    <md-card-header data-background-color="green">
                        <h4 class="title">Choose your credential</h4>
                    </md-card-header>
                    <md-card-content style="background-color: white;">
                        <md-field>
                        </md-field>
                        <label>
                            <md-icon>how_to_reg</md-icon>
                            Issuer</label><br>
                        <select v-model="selectedDID" id="selectDID" style="color: grey; width: 300px; height: 35px;">
                            <option v-for="did in savedDIDs" :key="did" :value="did.id">
                                {{did.name}}
                            </option>
                        </select><br><br>
                        <label>
                            <md-icon>fingerprint</md-icon>
                            Credential</label><br>
                        <select v-model="selectedVC" style="color: grey; width: 300px; height: 35px;">
                            <option v-for="vc in savedVCs" :key="vc" :value="vc.id">
                                {{vc.name}}
                            </option>
                        </select>
                        <md-field style="margin-top: -15px">
                        </md-field>
                        <md-button v-on:click="createPresentation" class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100">Send VC
                        </md-button>
                    </md-card-content>
                </md-card>

            </div>
        </div>
    </div>
</template>
<script>

    export default {
        beforeCreate: async function () {
            this.aries = await this.$arieslib
            // Load the Credentials in the drop down
            this.aries.verifiable.getCredentials()
                .then(resp => {
                        const data = resp.result
                        console.log("data from rsp", data)
                        if (data.length == 0) {
                            console.log("unable to get saved VCs")
                            return
                        }

                        this.savedVCs.length = 0
                        data.forEach((item) => {
                            this.savedVCs.push({id:item.id, name:item.name})
                        })

                        this.selectedVC = this.savedVCs[0].id
                    }
                ).catch(err => {
                    console.log('get credentials failed : errMsg=' + err)
                }
            )

            this.$polyfill.loadOnce()
            this.credentialEvent = await this.$webCredentialHandler.receiveCredentialEvent();

            // Load the DIDs in the drop down
           await this.aries.vdri.getDIDRecords()
                .then(resp => {
                        const data = resp.result
                        console.log("data from did resp", JSON.stringify(resp))
                        if (data.length == 0) {
                            console.log("unable to get saved DIDs")
                            return
                        }

                        this.savedDIDs.length = 0
                        data.forEach((item) => {
                            this.savedDIDs.push({id:item.id, name:item.name})
                        })

                        this.selectedDID = this.savedDIDs[0].id
                        console.log("What are the stored dids", this.selectedDID)
                    }
                ).catch(err => {
                    console.log('get DIDs failed : errMsg=' + err)
                }
            )
        },
        data() {
            return {
                savedVCs: [{id: "", name: "Select VC"}],
                selectedVC: "",
                savedDIDs: [{id: "", name: "Select DID"}],
                selectedDID: ""
            };
        },
        methods: {
            //TODO support multiple VCs + create presentation
            createPresentation: async function () {
                let data
                await this.aries.verifiable.getCredential({
                    id: this.selectedVC
                }).then(resp => {
                        data = JSON.stringify(JSON.parse(resp.verifiableCredential))
                    }
                ).catch(err => {
                    data = err
                    console.log('get vc failed : errMsg=' + err)
                })

                // Call Credential Handler callback
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: "Response",
                        data: data
                    });
                }))
            }
        }

    }
</script>
