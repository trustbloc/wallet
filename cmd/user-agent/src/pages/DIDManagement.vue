/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="md-layout">
            <div class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
                <md-card class="md-card-plain">
                <md-card-content>
                <md-tabs class="md-success" md-alignment="left">
                    <md-tab id="tab-home" md-label="DID Dashboard" md-icon="contact_mail">
                        <div class="md-layout-item md-layout md-gutter">
                            <div class="md-layout-item md-medium-size-75 md-xsmall-size-75 md-size-75">
                            <md-card-content>
                                <md-table v-model="searched" class="md-scrollbar" md-sort-order="asc">
                                    <md-table-toolbar>
                                        <div class="md-toolbar-section-start">
                                            <h1 class="md-title"><b>DID Table</b></h1>

                                            <md-field md-clearable class="md-toolbar-section-end">
                                                <md-input placeholder="Search by name..." v-model="searchQuery"  />
                                            </md-field>
                                        </div>


                                    </md-table-toolbar>
                                    <md-table-row>
                                        <md-table-head><b>Name</b></md-table-head>
                                        <md-table-head><b>DID</b></md-table-head>
                                        <md-table-head><b>Signature Type</b></md-table-head>
                                        <md-table-head><b>PrivateKey Type</b></md-table-head>
                                    </md-table-row>
                                    <md-table-row v-for="(data) in resultQuery" :key="data.id">
                                        <md-table-cell>{{data.friendlyName}}</md-table-cell>
                                        <md-table-cell>{{data.id}}</md-table-cell>
                                        <md-table-cell>{{data.signatureType}}</md-table-cell>
                                        <md-table-cell>{{data.privateKeyType}}</md-table-cell>
                                    </md-table-row>
                                </md-table>
                            </md-card-content>
                            <md-field>
                            </md-field>
                            </div>
                            <div class="md-layout-item">
                            <md-card>
                            <md-card-header data-background-color="green">
                                <h4 class="title"><b>Create New Trustbloc DID</b></h4>
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

                                <md-button class="md-button md-info md-square  md-large-size-100 md-size-100"
                                           id='createDIDBtn' v-on:click="createDID"><b>Create and Save DID</b>
                                </md-button>
                                <div v-if="errors.length" >
                                    <b>Please correct the following error(s):</b>
                                    <ul>
                                        <li v-for="error in errors" :key="error">{{ error }}</li>
                                    </ul>
                                </div>
                                <div v-if="loading" style="margin-left: 40%;margin-top: 20%;height: 200px;">
                                    <div class="md-layout">
                                        <md-progress-spinner :md-diameter="100" class="md-primary" :md-stroke="10"
                                                             md-mode="indeterminate"></md-progress-spinner>
                                    </div>
                                </div>
                                <md-field>
                                    <md-textarea v-model="didDocTextArea"  readonly style="min-height:300px;">
                                    </md-textarea>
                                </md-field>
                            </md-card-content>
                          </md-card>
                            </div>
                        </div>
                    </md-tab>
                    <md-tab id="tab-pages" md-label="Save Any DID" md-icon="contacts">
                        <md-card class="md-card-plain">
                            <md-card-header data-background-color="green">
                                <h4 class="title">Save Any DID</h4>
                            </md-card-header>
                            <md-card-content>
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
                                    <md-icon>aspect_ratio
                                        <md-tooltip md-direction="top">Enter key ID for above private key
                                        </md-tooltip>
                                    </md-icon>
                                    <label class="md-helper-text">Enter matching Key ID</label>
                                    <md-field maxlength="5">
                                        <md-input v-model="keyID" id="keyID" required></md-input>
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-size-100">
                                    <md-icon>memory</md-icon>
                                    <select id="privateKeyType" v-model="privateKeyType" style="color: grey; width: 200px; height: 35px;">
                                    <option value="">Select Key Type</option>
                                    <option value="Ed25519">Ed25519</option>
                                </select>
                                    <md-field style="margin-top: -15px">
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-size-100">
                                    <md-icon>ballot</md-icon> <select id="selectSignKey" v-model="selectSignKey" style="color: grey; width: 300px; height: 35px;">
                                    <option value="">Select Signature Type</option>
                                    <option value="Ed25519Signature2018">Ed25519Signature2018</option>
                                    </select>
                                    <md-field style="margin-top: -15px">
                                    </md-field>
                                </div>
                                <div class="md-layout-item md-size-100">
                                    <md-field maxlength="5">
                                        <label class="md-helper-text">Type DID friendly name here</label>
                                        <md-input v-model="anyDIDFriendlyName" id="anyDIDFriendlyName" required></md-input>
                                    </md-field>
                                </div>
                                <md-button class="md-button md-success md-square md-theme-default md-large-size-100 md-size-100"
                                           id='saveDIDBtn' v-on:click="saveAnyDID">Resolve and Save DID
                                </md-button>
                                <div v-if="saveErrors.length">
                                    <b>Please correct the following error(s):</b>
                                    <ul>
                                        <li v-for="error in saveErrors" :key="error">{{ error }}</li>
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
    export default {
        beforeCreate: async function () {
            window.$aries = await this.$arieslib
            window.$trustblocAgent = this.$trustblocAgent
            window.$trustblocStartupOpts = await this.$trustblocStartupOpts
            this.searched = this.myData
            await this.loadDIDMetadata()
        },
        methods: {
            createDID: async function () {
                var m = new Map([["Ed25519Signature2018","Ed25519VerificationKey2018"], ["JsonWebSignature2020" ,"JwsVerificationKey2020" ]]);

                this.errors.length = 0
                if (this.friendlyName.length == 0) {
                    this.errors.push("friendly name required.")
                    return
                }
                if ((this.selectType == "")) {
                    this.errors.push("key type required")
                    return;
                }
                if ((this.signType == "")) {
                    this.errors.push("signature type required")
                    return;
                }

                this.loading = true
                let generateKeyType
                if (this.selectType == "Ed25519") {
                    generateKeyType = "ED25519"
                }
                if (this.selectType == "P256") {
                    generateKeyType = "ECDSAP256IEEEP1363"
                }

                const keyset = await window.$aries.kms.createKeySet({keyType: generateKeyType})
                const recoveryKeyset = await window.$aries.kms.createKeySet({keyType: generateKeyType})
                const opsKeyset = await window.$aries.kms.createKeySet({keyType: generateKeyType})


                const createDIDRequest = {
                    "publicKeys": [{
                        "id": opsKeyset.keyID,
                        "type": m.get("JsonWebSignature2020"),
                        "value": opsKeyset.publicKey,
                        "encoding": "Jwk",
                        "keyType": this.selectType,
                        "usage": ["ops"]
                    }, {
                        "id": keyset.keyID,
                        "type": m.get(this.signType),
                        "value": keyset.publicKey,
                        "encoding": "Jwk",
                        "keyType": this.selectType,
                        "usage": ["general", "auth"]
                    }, {
                        "id": recoveryKeyset.keyID,
                        "type": m.get(this.signType),
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
                        this.loading = false;
                        this.didDocTextArea = err
                    })

                await t.destroy()

                const did = JSON.parse(this.didDocTextArea)
                // saving did in the did store
                await window.$aries.vdri.saveDID({
                        name: this.friendlyName,
                        did: did
                    }
                ).then(
                    console.log("successfully saved the did")
                ).catch(err => {
                    this.loading = false;
                    this.didDocTextArea = 'failed to save the did : ' + err
                        console.log('failed to save the did : errMsg=' + err)
                    }
                )

                this.storeDIDMetadata(did.id,"","", this.signType, "", this.friendlyName)
                this.loading = false;

            },
            saveAnyDID: async function () {
                this.saveErrors.length = 0
                if (this.didID.length == 0) {
                    this.saveErrors.push("did id required.")
                    return
                }
                if (this.privateKey.length == 0) {
                    this.saveErrors.push("private key required.")
                    return
                }
                if (this.keyID.length == 0) {
                    this.errors.push("key ID (verification method) matching private key is required.")
                    return
                }
                if (this.anyDIDFriendlyName.length == 0) {
                    this.saveErrors.push("friendly name required.")
                    return
                }
                if ((this.privateKeyType == "")) {
                    this.saveErrors.push("key type required")
                    return;
                }
                if ((this.selectSignKey == "")) {
                    this.saveErrors.push("signature type required")
                    return;
                }

                var resp
                try {
                    resp = await window.$aries.vdri.resolveDID({
                        id: this.didID,
                    })
                } catch (err) {
                    this.saveErrors.push(err)
                    return
                }

                this.anyDidDocTextArea = JSON.stringify(resp.did, undefined, 2);


                // saving did in the did store
                await window.$aries.vdri.saveDID({
                        name: this.anyDIDFriendlyName,
                        did: resp.did
                    }
                ).then(
                    console.log("successfully saved the did")
                ).catch(err => {
                        this.anyDidDocTextArea = 'failed to save the did : ' + err
                        console.log('failed to save the did : errMsg=' + err)
                    }
                )

                this.storeDIDMetadata(resp.did.id,this.privateKey,this.privateKeyType,this.selectSignKey, this.keyID, this.anyDIDFriendlyName)
            },
            storeDIDMetadata: function (did,privateKey,privateKeyType,signatureType,keyID,friendlyName) {
                var openDB = indexedDB.open("did-metadata", 1);

                openDB.onupgradeneeded = function () {
                    var db = {}
                    db.result = openDB.result;
                    db.store = db.result.createObjectStore("metadata", {keyPath: "id"});
                };

                openDB.onsuccess = function() {
                    var db = {};
                    db.result = openDB.result;
                    db.tx = db.result.transaction("metadata", "readwrite");
                    db.store = db.tx.objectStore("metadata");
                    db.store.put({
                                id: did,
                                friendlyName: friendlyName,
                                privateKey: privateKey,
                                privateKeyType:privateKeyType,
                                signatureType: signatureType,
                                keyID: keyID});
                    console.log("stored did metadata to db")
                }

            },
            loadDIDMetadata: async function () {
               let openDB = indexedDB.open("did-metadata", 1);
                openDB.onupgradeneeded = function() {
                    var db = {}
                    db.result = openDB.result;
                    db.store = db.result.createObjectStore("metadata", {keyPath: "id"});
                };
              const callback = (data) => this.myData = data
                openDB.onsuccess = function() {
                    let db = {};
                    db.result = openDB.result;
                    db.tx = db.result.transaction("metadata", "readonly");
                    db.store = db.tx.objectStore("metadata");
                    let request = db.store.getAll();
                    request.onsuccess = function() {
                      callback(request.result)
                    };

                }
            }

        },
        computed: {
            resultQuery(){
                if(this.searchQuery){
                    return this.myData.filter((data)=>{
                        return this.searchQuery.toLowerCase().split(' ').every(v => data.friendlyName.toLowerCase().includes(v))
                    })
                }else{
                    return this.myData;
                }
            }
        },
        data() {
            return {
                didDocTextArea: "",
                anyDidDocTextArea: "",
                friendlyName: "",
                selectType: "",
                selectSignKey: "",
                privateKeyType: "",
                signType: "",
                didID: "",
                privateKey: "",
                keyID: "",
                anyDIDFriendlyName: "",
                errors: [],
                saveErrors: [],
                myData:{},
                search: null,
                searchQuery:null,
                searched: [],
                loading: false,
            };
        },
        created(){
            this.searched = this.myData
        }

    }
</script>
<style>
    select {
        -webkit-appearance: none;
        -moz-appearance: none;
        appearance: none;
    }
</style>
