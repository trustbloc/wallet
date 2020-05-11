/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div v-if="loading" style="margin-left: 40%;margin-top: 20%;height: 200px;">
        <div class="md-layout">
            <md-progress-spinner :md-diameter="100" class="md-accent" :md-stroke="10"
                                 md-mode="indeterminate"></md-progress-spinner>
        </div>
    </div>
    <div v-else class="md-layout">

        <div v-if="authView" class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
            <h4> {{ requestOrigin }} would like you to authenticate </h4>
            <md-card class="md-card-plain">
                <md-card-header data-background-color="green">
                    <h4 class="title">DID Authorization</h4>
                </md-card-header>
                <md-card-content v-if="!credentialWarning.length" style="background-color: white;">
                    <div v-if="errors.length">
                        <b>Failed with following error(s):</b>
                        <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                        </ul>
                    </div>


                    <md-field>
                        <label>
                            <md-icon>how_to_reg</md-icon>
                            Select a Subject DID: </label>
                        <md-select v-model="selectedIssuer">
                            <md-option v-for="{id, name} in issuers" :key="id" :value="id">
                                {{name}}
                            </md-option>
                        </md-select>
                    </md-field>

                    <md-button v-on:click="cancel" class="md-cancel-text" id="cancelBtn" style="margin-right: 5px">
                        Cancel
                    </md-button>

                    <md-button v-on:click="authorize" style="margin-left: 5px"
                               class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                               id="authenticate">Authenticate
                    </md-button>
                </md-card-content>

                <md-card-content v-else style="background-color: white;">
                    <md-empty-state md-size=250
                                    class="md-accent"
                                    md-rounded
                                    md-icon="link_off"
                                    :md-label="credentialWarning">
                    </md-empty-state>
                </md-card-content>
            </md-card>

        </div>

        <!-- Generate Presentation View-->
        <div v-if="isVP && !authView " class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
            <h4> {{ requestOrigin }} has requested a credential from you </h4>
            <h5 v-if="reason.length"><b>Reason:</b> {{ reason }} </h5>
            <md-card class="md-card-plain">

                <div v-if="errors.length">
                    <b>Failed with following error(s):</b>
                    <md-field style="margin-top: -15px">
                        <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                        </ul>
                    </md-field>
                </div>

                <md-card-header data-background-color="green">
                    <h4 class="title">Share your credential</h4>
                </md-card-header>


                <md-card-content id="getvc-content" v-if="!credentialWarning.length" style="background-color: white;">

                    <md-field style="margin-bottom: -2%;">
                        <label>
                            <md-icon>how_to_reg</md-icon>
                            Select Identity: </label>
                        <md-select v-model="selectedIssuer" @md-selected="searchOnTable">
                            <md-option v-for="{id, name} in issuers" :key="id" :value="id">
                                {{name}}
                            </md-option>
                        </md-select>
                    </md-field>

                    <md-table v-model="searched" md-sort="name" md-card @md-selected="onSelect">
                        <md-table-toolbar>
                            <md-field md-clearable class="md-toolbar-section-end">
                                <md-input placeholder="Search by Credential Type..." v-model="search"
                                          @input="searchOnTable"/>
                            </md-field>
                        </md-table-toolbar>

                        <md-table-empty-state
                                md-label="No credentials found"
                                :md-description="`No credentials found for this '${search}' type. Try a different type search or add a new credential.`">
                        </md-table-empty-state>

                        <md-table-row slot="md-table-row" slot-scope="{ item }" md-selectable="multiple" md-auto-select>
                            <md-table-cell md-label="Name" md-sort-by="name">{{ item.name }}</md-table-cell>
                            <md-table-cell id="cred-type" md-label="Credential Type" md-sort-by="type">{{ item.type }}
                            </md-table-cell>
                        </md-table-row>

                    </md-table>

                    <md-button v-on:click="createPresentation"
                               class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                               id="share-credentials">Share
                    </md-button>
                </md-card-content>

                <md-card-content v-else style="background-color: white;">
                    <md-empty-state md-size=250
                                    class="md-accent"
                                    md-rounded
                                    md-icon="link_off"
                                    :md-label="credentialWarning">
                    </md-empty-state>
                </md-card-content>

            </md-card>
        </div>

        <!-- Generate Credential View-->
        <div v-if="isVC" class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
            <h4> {{ requestOrigin }} has requested a credential from you </h4>
            <md-card class="md-card-plain">
                <md-card-header data-background-color="green">
                    <h4 class="title">Share your credential</h4>
                </md-card-header>
                <md-card-content v-if="!credentialWarning.length" style="background-color: white;">
                    <div v-if="errors.length">
                        <b>Failed with following error(s):</b>
                        <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                        </ul>
                    </div>

                    <md-field>
                        <label>
                            <md-icon>fingerprint</md-icon>
                            Credentials: </label>
                        <md-select v-model="selectedVC" name="creds" id="creds">
                            <md-option v-for="vc in savedVCs" :key="vc" :value="vc.id">{{vc.name}}</md-option>
                        </md-select>
                    </md-field>


                    <md-button v-on:click="createCredential"
                               class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                               id="share-credential">Share
                    </md-button>
                </md-card-content>

                <md-card-content v-else style="background-color: white;">
                    <md-empty-state md-size=250
                                    class="md-accent"
                                    md-rounded
                                    md-icon="link_off"
                                    :md-label="credentialWarning">
                    </md-empty-state>
                </md-card-content>

            </md-card>
        </div>

    </div>
</template>
<script>

    const vcType = "VerifiableCredential"
    const vpType = "VerifiablePresentation"

    const getDomainAndChallenge = (vp) => {
        let {challenge, domain, query} = vp;

        if (query && query.challenge) {
            challenge = query.challenge;
        }

        if (query && query.domain) {
            domain = query.domain;
        }

        return {domain, challenge};
    }


    const toLower = text => {
        return text.toString().toLowerCase()
    }

    const isKeyType = (type) => {
        return toLower(type) == toLower(vcType) || toLower(type) == toLower(vpType)
    }

    const searchByTypeAndHolder = (items, term, key) => {
        if (key) {
            items = items.filter(item => item.holder == key)
        }

        if (term) {
            return items.filter(item => toLower(item.type).includes(toLower(term)))
        }

        return items
    }

    const getCredentialType = (types) => {
        const result = types.filter(type => !isKeyType(type))
        if (result.length > 0) {
            return result[0]
        }
        return ""
    }


    export default {
        beforeCreate: async function () {
            this.$polyfill.loadOnce()
            this.credentialEvent = await this.$webCredentialHandler.receiveCredentialEvent();
            this.aries = await this.$arieslib
            this.requestOrigin = this.credentialEvent.credentialRequestOrigin


            if (this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation) {
                this.isVP = true
                this.verifiable = this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation
            } else if (this.credentialEvent.credentialRequestOptions.web.VerifiableCredential) {
                this.isVC = true
                this.verifiable = this.credentialEvent.credentialRequestOptions.web.VerifiableCredential
            } else {
                console.log("invalid web credential type, expected 'VerifiablePresentation' or 'VerifiableCredential'")
                return
            }

            console.log("incoming web credential", this.credentialEvent.credentialRequestOptions.web)

            this.query = {}
            if (this.verifiable.query) {
                // TODO What if there are multiple queries?
                this.query = Array.isArray(this.verifiable.query) ? this.verifiable.query[0] : this.verifiable.query;
            }

            this.setupDomainAndChallenge()

            if (this.query.type === 'DIDAuth') {
                this.authView = true
                await this.loadIssuers()
                this.loading = false
                return
            }

            if (this.isVP) {
                await this.loadIssuers()
            }
            await this.loadCredentials()

            this.loading = false

            if (this.query.credentialQuery && this.query.credentialQuery.reason) {
                this.reason = this.query.credentialQuery.reason
            }

            this.search = ""
            if (this.query.credentialQuery && this.query.credentialQuery.example && this.query.credentialQuery.example.type) {
                let t = this.query.credentialQuery.example.type
                let key = Array.isArray(t) ? t[0] : t
                if (!isKeyType(key)) {
                    this.search = key
                }
            }

            // perform search while loading
            this.searched = []
            this.searchOnTable()
        },
        data() {
            return {
                authView: false,
                savedVCs: [{id: 0, name: "Select VC"}],
                selectedVCs: [],
                selectedVC: 0,
                issuers: [{id: 0, name: "Select Identity"}],
                selectedIssuer: 0,
                errors: [],
                requestOrigin: "",
                loading: true,
                isVP: false,
                isVC: false,
                credentialWarning: "",
                search: "",
                searched: [],
                reason: "",
            };
        },
        methods: {
            getDIDMetadata: function (id) {
                return new Promise(function (resolve) {
                    var openDB = indexedDB.open("did-metadata", 1);

                    openDB.onupgradeneeded = function () {
                        var db = {}
                        db.result = openDB.result;
                        db.store = db.result.createObjectStore("metadata", {keyPath: "id"});
                    };

                    openDB.onsuccess = function () {
                        var db = {};
                        db.result = openDB.result;
                        db.tx = db.result.transaction("metadata", "readonly");
                        db.store = db.tx.objectStore("metadata");
                        let getData = db.store.get(id);
                        getData.onsuccess = function () {
                            resolve(getData.result);
                        };

                        db.tx.oncomplete = function () {
                            db.result.close();
                        };
                        console.log("got did metadata from db")
                    }
                });
            },
            searchOnTable() {
                this.searched = searchByTypeAndHolder(this.savedVCs, this.search, this.issuers[this.selectedIssuer].key)
            },
            setupDomainAndChallenge: function () {
                const {domain, challenge} = getDomainAndChallenge(this.verifiable);
                this.domain = domain
                this.challenge = challenge

                if (!this.domain && event.credentialRequestOrigin) {
                    this.domain = event.credentialRequestOrigin.split('//').pop()
                }
            },
            loadCredentials: async function () {
                await this.aries.verifiable.getCredentials()
                    .then(resp => {
                            this.savedVCs.length = 0
                            const data = resp.result
                            if (!data || data.length == 0) {
                                this.credentialWarning = 'No Saved Credentials Found'
                                return
                            }

                            data.forEach((item, id) => {
                                this.savedVCs.push({
                                    id: id,
                                    name: item.name,
                                    key: item.id,
                                    type: getCredentialType(item.type),
                                    holder: item.subjectId,
                                })
                            })
                        }
                    ).catch(err => {
                            this.errors.push('Failed to get credentials')
                            console.log('get credentials failed, error:', err)
                        }
                    )
            },
            loadIssuers: async function () {
                await this.aries.vdri.getDIDRecords().then(
                    resp => {
                        const data = resp.result
                        if (!data || data.length == 0) {
                            this.credentialWarning = 'Issuers not found, please create an issuer'
                            return
                        }

                        this.issuers = []
                        this.selectedIssuer = 0

                        data.forEach((item, id) => {
                            this.issuers.push({id: id, name: item.name, key: item.id})
                        })
                    })
                    .catch(err => {
                        this.errors.push(err)
                    })
            },
            onSelect(items) {
                this.selectedVCs = items
            },
            getSelectedCredentials: async function () {
                if (this.selectedVCs.length == 0) {
                    return {retry: "Please select at least one credential"}
                }

                try {
                    let vcs = []
                    for (let selectedVC of this.selectedVCs) {
                        const resp = await this.aries.verifiable.getCredential({
                            id: selectedVC.key
                        })
                        vcs.push(JSON.parse(resp.verifiableCredential))
                    }
                    return {vcs: vcs}
                } catch (e) {
                    return e
                }
            },
            getSelectedCredential: async function () {
                try {
                    const resp = await this.aries.verifiable.getCredential({
                        id: this.savedVCs[this.selectedVC].key
                    })

                    return JSON.parse(resp.verifiableCredential)
                } catch (e) {
                    return e
                }
            },
            cancel: async function () {
                // Call Credential Handler callback
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: "Response",
                        data: "DID Auth cancelled"
                    });
                }))
            },
            authorize: async function () {
                this.loading = true
                let didMetadata = await this.getDIDMetadata(this.issuers[this.selectedIssuer].key)

                let data
                await this.aries.verifiable.generatePresentation({
                    presentation: {
                        "@context": "https://www.w3.org/2018/credentials/v1",
                        "type": "VerifiablePresentation"
                    },
                    domain: this.domain,
                    challenge: this.challenge,
                    did: this.issuers[this.selectedIssuer].key,
                    signatureType: didMetadata.signatureType,
                    verificationMethod: didMetadata.keyID,
                }).then(resp => {
                        if (!resp.verifiablePresentation) {
                            data = "failed to create did auth presentation"
                            return
                        }

                        data = resp.verifiablePresentation
                        //TODO bug in aries to show '"verifiableCredential": null,' in empty presentations
                        if (data.hasOwnProperty('verifiableCredential')) {
                            delete data.verifiableCredential
                        }
                    }
                ).catch(err => {
                    data = err
                    console.log('failed to create presentation, errMsg:', err)
                })

                this.loading = false
                console.log("Response presentation:", data)

                // Call Credential Handler callback
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: vpType,
                        data: data
                    });
                }))
            },
            createPresentation: async function () {
                let creds = await this.getSelectedCredentials()
                if (creds.retry) {
                    this.errors.push(creds.retry)
                    return
                }

                this.loading = true

                let data
                if (creds.vcs) {
                    let didMetadata = await this.getDIDMetadata(this.issuers[this.selectedIssuer].key)

                    await this.aries.verifiable.generatePresentation({
                        verifiableCredential: creds.vcs,
                        did: this.issuers[this.selectedIssuer].key,
                        domain: this.domain,
                        challenge: this.challenge,
                        // TODO skipVerify can be an option in view
                        skipVerify: true,
                        signatureType: didMetadata.signatureType,
                        verificationMethod: didMetadata.keyID
                    }).then(resp => {
                            data = resp.verifiablePresentation
                        }
                    ).catch(err => {
                        data = err
                        console.log('failed to create presentation, errMsg:', err)
                    })
                }

                this.loading = false
                console.log("Response presentation:", data)

                // Call Credential Handler callback
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: vpType,
                        data: data
                    });
                }))
            },
            createCredential: async function () {
                this.loading = true

                let cred = await this.getSelectedCredential()
                console.log("Response credential:", cred)

                this.loading = false

                // Call Credential Handler callback
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: vcType,
                        data: cred
                    });
                }))
            }
        }

    }
</script>
<style lang="scss">
    .md-menu-content {
        width: auto;
        max-width: 100% !important;
    }

    .md-menu-content .md-ripple > span {
        position: relative !important;
    }

    .md-table-head {
        background-color: #00bcd4 !important;
        color: white !important;
        font-size: 10px !important;
        font-weight: 500 !important;
    }

    .md-table-head-label .md-sortable {
        padding-left: 10px !important;
    }

    .md-empty-state-icon {
        width: 60px !important;
        min-width: 60px !important;
        height: 50px !important;
        font-size: 50px !important;
        margin: 0;
    }

    .md-empty-state-label {
        font-size: 18px !important;
        font-weight: 500 !important;
        line-height: 20px !important;
    }

</style>
