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
                <md-card-content>
                <md-tabs class="md-success" md-alignment="left">
                    <md-tab id="tab-home" md-label="Create Trustbloc DID" md-icon="contact_mail">
                        <md-card class="md-card-plain">
                            <md-card-header data-background-color="blue">
                                <h4 class="title">Create trustbloc DID</h4>
                                <p class="category"> Create button will create and save the did in DID store</p>
                            </md-card-header>
                            <md-card-content>
                                <md-field>
                                </md-field>
                                <select id="selectKey" v-model="selectType" style="color: grey; width: 200px; height: 35px;">
                                    <option value="" disabled="disabled">Select Key Type</option>
                                    <option value="Ed25519">Ed25519</option>
                                    <option value="P256">P256</option>
                                </select>
                                <md-field style="margin-top: -15px">
                                </md-field>
                                <select id="signKey" v-model="signType" style="color: grey; width: 200px; height: 35px;">
                                    <option value="" disabled="disabled">Select Signature Type</option>
                                    <option value="Ed25519Signature2018">Ed25519Signature2018</option>
                                    <option value="JsonWebSignature2020">JsonWebSignature2020</option>
                                </select>
                                <md-field style="margin-top: -15px">
                                </md-field>
                                <div class="md-layout-item md-size-100">
                                    <md-field maxlength="5">
                                        <label class="md-helper-text">Type DID friendly name here</label>
                                        <md-input v-model="friendlyName" id="friendlyName" required></md-input>
                                    </md-field>
                                </div>
                                <md-button class="md-button md-success md-square md-theme-default md-large-size-100 md-size-100"
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
                    </md-tab>
                    <md-tab id="tab-pages" md-label="Save Any DID" md-icon="contacts">
                        <md-card class="md-card-plain">
                            <md-card-header data-background-color="blue">
                                <h4 class="title">Save Any DID</h4>
                                <p class="category"> Create button will resolve and save the did in DID store</p>
                            </md-card-header>
                            <md-card-content>
                                <md-field>
                                </md-field>
                                <div class="md-layout-item md-size-100">
                                    <md-icon>line_style</md-icon><label class="md-helper-text">Enter DID</label>
                                    <md-field maxlength="5">
                                        <md-input v-model="didID" id="did" required></md-input>
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-size-100">
                                    <md-icon>vpn_key</md-icon> <label class="md-helper-text">Enter Private Key</label>
                                    <md-field maxlength="5">
                                        <md-input v-model="privateKey" id="privateKey" required></md-input>
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-size-100">
                                    <md-icon>memory</md-icon> <select id="selectSignKey" v-model="signType" style="color: grey; width: 300px; height: 35px;">
                                    <option value="" disabled="disabled">Select Signature Type</option>
                                    <option value="Ed25519Signature2018">Ed25519Signature2018</option>
                                    <option value="JsonWebSignature2020">JsonWebSignature2020</option>
                                    </select>
                                    <md-field style="margin-top: -15px">
                                    </md-field>
                                </div>
                                <md-button class="md-button md-success md-square md-theme-default md-large-size-100 md-size-100"
                                           id='saveDIDBtn' v-on:click="saveAnyDID">Resolve and Save DID
                                </md-button>
                                <div v-if="errors.length">
                                    <b>Please correct the following error(s):</b>
                                    <ul>
                                        <li v-for="error in errors" :key="error">{{ error }}</li>
                                    </ul>
                                </div>
                                <md-field>
                                    <md-textarea v-model="anyDidDocTextArea" readonly style="min-height:360px;">
                                    </md-textarea>
                                </md-field>
                            </md-card-content>
                        </md-card>
                    </md-tab>

                </md-tabs>
                </md-card-content>
                </md-card>
            </div>
        </div>
    </div>
</template>

<script>
    import axios from 'axios'
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
                if (this.selectType == "Ed25519") {
                    generateKeyType = "ED25519"
                }
                if (this.selectType == "P256") {
                    generateKeyType = "ECDSAP256IEEE1363"
                }

                const keyset = await window.$aries.kms.createKeySet({keyType: generateKeyType})
                const recoveryKeyset = await window.$aries.kms.createKeySet({keyType: generateKeyType})

                const createDIDRequest = {
                    "publicKeys": [{
                        "id": keyset.keyID,
                        "type": "JwsVerificationKey2020",
                        "value": keyset.publicKey,
                        "encoding": "Jwk",
                        "keyType": this.selectType,
                        "usage": ["ops", "general"]
                    }, {
                        "id": recoveryKeyset.keyID,
                        "type": "JwsVerificationKey2020",
                        "value": recoveryKeyset.publicKey,
                        "encoding": "Jwk",
                        "keyType": this.selectType,
                        "recovery": true
                    }
                    ]
                };

                const t = await new window.$trustblocAgent.Framework(JSON.parse(window.$trustblocStartupOpts))
                await t.didclient.createDID(createDIDRequest).then(
                    resp => {
                        // TODO generate public key from generic wasm
                        // TODO pass public key to createDID
                        this.didDocTextArea = JSON.stringify(resp.DID, undefined, 2);

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
            },
            saveAnyDID: async function () {
                this.errors.length = 0
                if (this.didID.length == 0) {
                    this.errors.push("did id required.")
                    return
                }
                if (this.privateKey.length == 0) {
                    this.errors.push("private key required.")
                    return
                }
                if ((this.selectType == "")) {
                    this.errors.push("key type required")
                    return;
                }
                // todo logic will be revisited soon
                for (let i = 0; i < 10; i++) {
                    let success = false
                    await axios.get(`${this.resolver}/` + this.didID).then((response) => {
                        console.log(response.data)
                        success = true
                        this.anyDidDocTextArea = JSON.stringify(response.data, undefined, 2);
                    }).catch(function (error) {
                        console.log('will try to resolve again : attempt=', i + 1 + error);
                    });
                    if (success) {
                        let objectStoreName = "didKeyStore"
                        let db
                        let did = this.didID
                        let pk = this.privateKey
                        let resolverLink = this.resolver

                        let request = indexedDB.open('store-keys', 1);

                        console.log('did', JSON.stringify(this.didID));

                        request.onsuccess = function (event) {
                            db = event.target.result;
                            let keys = {
                                key: did,
                                privateKey: pk,
                                resolverLink: resolverLink,
                                timestamp: Date.now()
                            };
                            console.log('data', keys);
                            db.transaction([objectStoreName], "readwrite").objectStore(objectStoreName).add(keys);

                        };
                        request.onupgradeneeded = function (event) {
                            db = event.target.result;
                            let store = db.createObjectStore(objectStoreName, { keyPath: 'key' });
                            console.log('[onsuccess]', store);
                        };
                        break;
                    }
                }

            }
        },
        data() {
            return {
                didDocTextArea: "",
                anyDidDocTextArea: "",
                friendlyName: "",
                selectType:"",
                signType:"",
                didID:"",
                privateKey:"",
                errors: [],
            };
        }
    }
</script>

