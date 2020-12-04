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
                            <md-tab id="tab-home" md-label="Digital Identity Dashboard" md-icon="contact_mail">
                                <div class="md-layout-item md-layout md-gutter">
                                    <div class="md-layout-item md-medium-size-75 md-xsmall-size-75 md-size-75">
                                        <md-card-content>
                                            <md-table v-model="searched" class="md-scrollbar" md-sort-order="asc">
                                                <md-table-toolbar>
                                                    <div class="md-toolbar-section-start">
                                                        <h1 class="md-title"><b>Digital Identity Table</b></h1>

                                                        <md-field md-clearable class="md-toolbar-section-end">
                                                            <md-input placeholder="Search by name..."
                                                                      v-model="searchQuery"/>
                                                        </md-field>
                                                    </div>


                                                </md-table-toolbar>
                                                <md-table-row>
                                                    <md-table-head><b>Name</b></md-table-head>
                                                    <md-table-head><b>Digital Identity</b></md-table-head>
                                                    <md-table-head><b>Signature Type</b></md-table-head>
                                                </md-table-row>
                                                <md-table-row v-for="(data) in resultQuery" :key="data.id">
                                                    <md-table-cell>{{data.friendlyName}}</md-table-cell>
                                                    <md-table-cell>{{data.id}}</md-table-cell>
                                                    <md-table-cell>{{data.signatureType}}</md-table-cell>
                                                </md-table-row>
                                            </md-table>
                                        </md-card-content>
                                        <md-field>
                                        </md-field>
                                    </div>
                                    <div class="md-layout-item">
                                        <md-card>
                                            <md-card-header data-background-color="green">
                                                <h4 class="title"><b>Create New Trustbloc Digital Identity</b></h4>
                                            </md-card-header>
                                            <md-card-content>
                                                <md-field>
                                                </md-field>
                                                <select id="selectKey" v-model="selectType"
                                                        style="color: grey; width: 200px; height: 35px;">
                                                    <option value="" disabled="disabled">Select Key Type</option>
                                                    <option value="Ed25519">Ed25519</option>
                                                    <option value="P256">P256</option>
                                                </select>
                                                <md-field style="margin-top: -15px">
                                                </md-field>
                                                <select id="signKey" v-model="signType"
                                                        style="color: grey; width: 200px; height: 35px;">
                                                    <option value="" disabled="disabled">Select Signature Type</option>
                                                    <option value="Ed25519Signature2018">Ed25519Signature2018</option>
                                                    <option value="JsonWebSignature2020">JsonWebSignature2020</option>
                                                </select>
                                                <md-field style="margin-top: -15px">
                                                </md-field>
                                                <div class="md-layout-item md-size-100">
                                                    <md-field maxlength="5">
                                                        <label class="md-helper-text">Type Digital Identity friendly name
                                                            here</label>
                                                        <md-input v-model="friendlyName" id="friendlyName"
                                                                  required></md-input>
                                                    </md-field>
                                                </div>

                                                <md-button
                                                        class="md-button md-info md-square  md-large-size-100 md-size-100"
                                                        id='createDIDBtn' v-on:click="createDID"><b>Create and Save
                                                    Digital Identity</b>
                                                </md-button>
                                                <div v-if="errors.length">
                                                    <b>Please correct the following error(s):</b>
                                                    <ul>
                                                        <li v-for="error in errors" :key="error">{{ error }}</li>
                                                    </ul>
                                                </div>
                                                <div v-if="loading"
                                                     style="margin-left: 40%;margin-top: 20%;height: 200px;">
                                                    <div class="md-layout">
                                                        <md-progress-spinner :md-diameter="100" class="md-primary"
                                                                             :md-stroke="10"
                                                                             md-mode="indeterminate"></md-progress-spinner>
                                                    </div>
                                                </div>
                                                <md-field>
                                                    <md-textarea v-model="didDocTextArea" readonly
                                                                 style="min-height:300px;">
                                                    </md-textarea>
                                                </md-field>
                                            </md-card-content>
                                        </md-card>
                                    </div>
                                </div>
                            </md-tab>
                            <md-tab id="tab-pages" md-label="Save Any Digital Identity" md-icon="contacts">
                                <md-card class="md-card-plain">
                                    <md-card-header data-background-color="green">
                                        <h4 class="title">Save Any Digital Identity</h4>
                                    </md-card-header>
                                    <md-card-content>
                                        <div class="md-layout-item md-size-100">
                                            <md-icon>line_style</md-icon>
                                            <label class="md-helper-text">Enter Digital Identity</label>
                                            <md-field maxlength="5">
                                                <md-input v-model="didID" id="did" required></md-input>
                                            </md-field>
                                        </div>
                                        <div class="md-layout-item md-size-100">
                                            <md-icon>vpn_key</md-icon>
                                            <label class="md-helper-text">Enter Private Key JWK</label>
                                            <md-field maxlength="5">
                                                <md-input v-model="privateKeyJwk" id="privateKeyJwk"
                                                          required></md-input>
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
                                            <md-icon>ballot</md-icon>
                                            <select id="selectSignKey" v-model="selectSignKey"
                                                    style="color: grey; width: 300px; height: 35px;">
                                                <option value="">Select Signature Type</option>
                                                <option value="Ed25519Signature2018">Ed25519Signature2018</option>
                                            </select>
                                            <md-field style="margin-top: -15px">
                                            </md-field>
                                        </div>
                                        <div class="md-layout-item md-size-100">
                                            <md-field maxlength="5">
                                                <label class="md-helper-text">Type Digital Identity friendly name here</label>
                                                <md-input v-model="anyDIDFriendlyName" id="anyDIDFriendlyName"
                                                          required></md-input>
                                            </md-field>
                                        </div>
                                        <md-button
                                                class="md-button md-success md-square md-theme-default md-large-size-100 md-size-100"
                                                id='saveDIDBtn' v-on:click="saveAnyDID">Resolve and Save Digital Identity
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

    import {DIDManager} from "./chapi/wallet";
    import {mapGetters} from 'vuex'

    export default {
        created: async function () {
            this.agent = this.getAgentInstance()

            this.didManager = new DIDManager(this.agent, this.getAgentOpts())
            this.searched = this.myData
            this.username = this.getCurrentUser().username

            await this.loadDIDMetadata()
            this.searched = this.myData
        },
        methods: {
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            ...mapGetters(['getCurrentUser', 'getAgentOpts']),
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
                if ((this.signType == "")) {
                    this.errors.push("signature type required")
                    return;
                }

                this.loading = true

                let did
                try {
                    did = await this.didManager.createTrustBlocDID(this.selectType, this.signType)
                } catch (e) {
                    this.loading = false;
                    this.didDocTextArea = `failed to create did: ${e.toString()}`
                    return;
                }

                this.didDocTextArea = JSON.stringify(did, undefined, 2);

                // saving did
                try {
                     await this.didManager.saveDID(this.username, this.friendlyName, this.signType, did)
                } catch (e) {
                    this.loading = false;
                    this.didDocTextArea = `failed to save did: ${e.toString()}`
                    console.error("failed to save did", e)
                    return;
                }

                // saving did metadata
                await this.didManager.storeDIDMetadata(did.id, {
                  signatureType: this.signType,
                  friendlyName: this.friendlyName,
                  id: did.id,
                })

                this.loading = false;

            },

            saveAnyDID: async function () {
                this.saveErrors.length = 0
                if (this.didID.length == 0) {
                    this.saveErrors.push("did id required.")
                    return
                }
                if (this.privateKeyJwk.length == 0) {
                    this.saveErrors.push("private key jwk required.")
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
                if ((this.selectSignKey == "")) {
                    this.saveErrors.push("signature type required")
                    return;
                }

                var resp
                try {
                    resp = await this.agent.vdr.resolveDID({
                        id: this.didID,
                    })
                } catch (err) {
                    this.saveErrors.push(err)
                    return
                }

                var obj = JSON.parse(this.privateKeyJwk);

                try {
                    await this.agent.kms.importKey(obj)
                } catch (err) {
                    this.saveErrors.push(err)
                    return
                }

                // saving did
                try {
                    await this.didManager.saveDID(this.username, this.anyDIDFriendlyName, this.selectSignKey, resp.did)
                } catch (e) {
                    this.anyDidDocTextArea = `failed to save did : ${e.toString()}`
                    console.error("failed to save the did", e)
                }

                await this.didManager.storeDIDMetadata(resp.did.id, {
                  signatureType: this.selectSignKey,
                  keyID: this.keyID,
                  friendlyName: this.anyDIDFriendlyName,
                  id: resp.did.id,
                })
                this.anyDidDocTextArea = JSON.stringify(resp.did, undefined, 2);
            },
            loadDIDMetadata: async function () {
               try {
                   this.myData = await this.didManager.getAllDIDMetadata()
               }catch (e) {
                    console.error("failed to load DIDs", e)
               }
            }

        },
        computed: {
            resultQuery() {
                if (this.searchQuery) {
                    return this.myData.filter((data) => {
                        return this.searchQuery.toLowerCase().split(' ').every(v => data.friendlyName.toLowerCase().includes(v))
                    })
                } else {
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
                signType: "",
                didID: "",
                privateKeyJwk: "",
                keyID: "",
                anyDIDFriendlyName: "",
                errors: [],
                saveErrors: [],
                myData: {},
                search: null,
                searchQuery: null,
                searched: [],
                loading: false,
            };
        }

    }
</script>
<style>
    select {
        -webkit-appearance: none;
        -moz-appearance: none;
        appearance: none;
        /* Some browsers will not display the caret when using calc, so we put the fallback first */
        background: url("http://cdn1.iconfinder.com/data/icons/cc_mono_icon_set/blacks/16x16/br_down.png") white no-repeat 98.5% !important; /* !important used for overriding all other customisations */
        background: url("http://cdn1.iconfinder.com/data/icons/cc_mono_icon_set/blacks/16x16/br_down.png") white no-repeat calc(100% - 10px) !important; /* Better placement regardless of input width */
    }
</style>
