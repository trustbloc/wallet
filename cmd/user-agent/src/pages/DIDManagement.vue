/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="md-layout">
            <div
                    class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100"
            >
                <md-card class="md-card-plain">
                    <md-card-header data-background-color="green">
                        <h4 class="title">Create your DID</h4>
                        <p class="category"> Create button will create and save the did in DID store</p>
                    </md-card-header>
                    <md-card-content style="background-color: white;">
                        <md-field>
                        </md-field>
                        <select id="selectVC" v-model="selectType" style="color: grey; width: 200px; height: 35px;">
                            <option value="" disabled="disabled">Select Key Type</option>
                            <option value="Ed25519">Ed25519</option>
                            <option value="P256">P256</option>
                        </select>
                        <md-field style="margin-top: -15px">
                        </md-field>
                        <div class="md-layout-item md-size-100">
                            <md-field maxlength="5">
                                <label class="md-helper-text">Type DID friendly name here</label>
                                <md-input v-model="friendlyName" id="friendlyName" required></md-input>
                            </md-field>
                        </div>
                        <md-button class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                                    id='createDIDBtn' v-on:click="createDID">Create and Save DID
                        </md-button>
                        <div v-if="errors.length">
                            <b>Please correct the following error(s):</b>
                            <ul>
                                <li v-for="error in errors" :key="error">{{ error }}</li>
                            </ul>
                        </div>
                        <md-field>
                            <md-textarea v-model="didDocTextArea" readonly style="min-height:360px;">
                            </md-textarea>
                        </md-field>
                    </md-card-content>
                </md-card>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        beforeCreate: async function () {
            window.$aries = await this.$arieslib
            window.$trustblocAgent = this.$trustblocAgent
            window.$trustblocStartupOpts = await this.$trustblocStartupOpts

        },
        methods: {
            createDID: async function () {
                this.errors.length = 0
                if (this.friendlyName.length == 0) {
                    this.errors.push("friendly name required.")
                    return
                }
                if ((this.selectType == "")) {
                    this.errors.push("key type required")
                    return;
                }
                let generateKeyType
                if (this.selectType == "Ed25519"){
                    generateKeyType="ED25519"
                }
                if (this.selectType == "P256"){
                    generateKeyType="ECDSAP256IEEE1363"
                }

                const keyset= await window.$aries.kms.createKeySet({keyType: generateKeyType})
                const recoveryKeyset= await window.$aries.kms.createKeySet({keyType: generateKeyType})

                const createDIDRequest = {
                    "publicKeys":[{
                        "id":keyset.keyID,
                        "type":"JwsVerificationKey2020",
                        "value":keyset.publicKey,
                        "encoding":"Jwk",
                        "keyType":this.selectType,
                        "usage":["ops","general"]
                    }, {
                        "id":recoveryKeyset.keyID,
                        "type":"JwsVerificationKey2020",
                        "value":recoveryKeyset.publicKey,
                        "encoding":"Jwk",
                        "keyType":this.selectType,
                        "recovery":true
                    }
                    ]
                };

                const t= await new window.$trustblocAgent.Framework(JSON.parse(window.$trustblocStartupOpts))
                await t.didclient.createDID(createDIDRequest).then(
                    resp => {
                        // TODO generate public key from generic wasm
                        // TODO pass public key to createDID
                        this.didDocTextArea = JSON.stringify(resp.DID,undefined, 2);

                        // TODO create local index db store to store private keys only when user provided
                        // if (selectType == "P256"){
                        //     let storeName = "privateKeys"
                        //     let db
                        //
                        //     let request = indexedDB.open('store-keys', 1);
                        //
                        //     request.onsuccess = function(event) {
                        //         db = event.target.result;
                        //         let keys = {
                        //             key: friendlyName,
                        //             privateKey: resp.PrivateKey,
                        //             type: selectType,
                        //             timestamp: Date.now()
                        //           };
                        //         db.transaction([storeName], "readwrite").objectStore("privateKeys").add(keys);
                        //
                        //     };
                        //     request.onupgradeneeded = function(event) {
                        //         db = event.target.result;
                        //         let store = db.createObjectStore(storeName, {keyPath: 'key'});
                        //         console.log('[onsuccess]', store);
                        //     };
                        // }

                    })
                    .catch(err => {
                        this.didDocTextArea = err
                    })

                await t.destroy()

                // saving did in the did store
                window.$aries.vdri.saveDID({
                        name: this.friendlyName,
                        did: JSON.parse(this.didDocTextArea)
                    }
                ).then(
                    console.log("successfully saved the did")

                ).catch(err => {
                         this.didDocTextArea = 'failed to save the did : ' + err
                          console.log('failed to save the did : errMsg=' + err)
                    }
                )
            }
        },
        data() {
            return {
                didDocTextArea: "",
                friendlyName: "",
                selectType:"",
                errors: [],
            };
        }
    }
</script>
