/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div v-if="loading" style="margin-left: 40%;margin-top: 10%;height: 200px;">
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
                <md-card-content style="background-color: white;">
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
                            <md-option v-for="issuer in issuers" :key="issuer" :value="issuer.id">{{issuer.name}}
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
            </md-card>

        </div>

        <div v-else class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
            <h4> {{ requestOrigin }} has requested a credential from you </h4>
            <md-card class="md-card-plain">
                <md-card-header data-background-color="green">
                    <h4 class="title">Share your credential</h4>
                </md-card-header>
                <md-card-content style="background-color: white;">
                    <div v-if="errors.length">
                        <b>Failed with following error(s):</b>
                        <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                        </ul>
                    </div>


                    <md-field>
                        <label>
                            <md-icon>how_to_reg</md-icon>
                            Select Identity: </label>
                        <md-select v-model="selectedIssuer">
                            <md-option v-for="issuer in issuers" :key="issuer" :value="issuer.id">{{issuer.name}}
                            </md-option>
                        </md-select>
                    </md-field>


                    <md-field>
                        <label>
                            <md-icon>fingerprint</md-icon>
                            Credentials: </label>
                        <md-select v-model="selectedVCs" name="movies" id="movies" multiple>
                            <md-option v-for="vc in savedVCs" :key="vc" :value="vc.id">{{vc.name}}</md-option>
                        </md-select>
                    </md-field>


                    <md-button v-on:click="createPresentation"
                               class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                               id="share">Share
                    </md-button>
                </md-card-content>
            </md-card>

        </div>
    </div>
</template>
<script>

    const getDomainAndChallenge = (vp) => {
        let {challenge, domain, query} = vp;

        if (query.challenge) {
            challenge = query.challenge;
        }

        if (query.domain) {
            domain = query.domain;
        }

        return {domain, challenge};
    }

    export default {
        beforeCreate: async function () {
            this.$polyfill.loadOnce()
            this.credentialEvent = await this.$webCredentialHandler.receiveCredentialEvent();
            this.aries = await this.$arieslib
            this.requestOrigin = this.credentialEvent.credentialRequestOrigin
            const vp = this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation
            // TODO What if there are multiple queries?
            const query = Array.isArray(vp.query) ? vp.query[0] : vp.query;

            if (query && query.type === 'DIDAuth') {
                this.authView = true
                const {domain, challenge} = getDomainAndChallenge(vp);
                this.domain = (domain) ? domain : event.credentialRequestOrigin.split('//').pop()
                this.challenge = challenge
                await this.loadIssuers()
                this.loading = false
                return
            }

            this.aries.verifiable.getCredentials()
                .then(resp => {
                        const data = resp.result
                        if (!data || data.length == 0) {
                            this.errors.push('No saved credentials')
                            this.loading = false
                            return
                        }

                        this.savedVCs.length = 0
                        data.forEach((item, id) => {
                            this.savedVCs.push({id: id, name: item.name, key: item.id})
                        })

                        this.loadIssuers()
                        this.loading = false
                    }
                ).catch(err => {
                    this.errors.push('Failed to get credentials')
                    console.log('get credentials failed, error:', err)
                    this.loading = false
                }
            )
        },
        data() {
            return {
                authView: false,
                savedVCs: [{id: 0, name: "Select VC"}],
                selectedVCs: [],
                issuers: [{id: 0, name: "Select Identity"}],
                selectedIssuer: 0,
                errors: [],
                requestOrigin: "",
                loading: true
            };
        },
        methods: {
            createPresentation: async function () {
                let data = await this.getSelectedCredentials()
                if (data.vcs) {
                    await this.aries.verifiable.generatePresentation({
                        verifiableCredential: data.vcs,
                        did: this.issuers[this.selectedIssuer].key,
                        signatureType:"JsonWebSignature2020",
                        // TODO can be an option in view
                        skipVerify: true,
                        // TODO - domain & challenge
                    }).then(resp => {
                            data = JSON.stringify(resp.verifiablePresentation)
                        }
                    ).catch(err => {
                        data = err
                        console.log('failed to create presentation, errMsg:', err)
                    })
                }

                if (data.retry) {
                    this.errors.push(data.retry)
                    return
                }

                // Call Credential Handler callback
                this.credentialEvent.respondWith(new Promise(function (resolve) {
                    return resolve({
                        dataType: "Response",
                        data: data
                    });
                }))
            },
            getSelectedCredentials: async function () {
                if (this.selectedVCs.length == 0) {
                    return {retry: "Please select at least one credential"}
                }

                try {
                    let vcs = []
                    for (let selectedVC of this.selectedVCs) {
                        const resp = await this.aries.verifiable.getCredential({
                            id: this.savedVCs[selectedVC].key
                        })
                        vcs.push(JSON.parse(resp.verifiableCredential))
                    }
                    return {vcs: vcs}
                } catch (e) {
                    return e
                }
            },
            loadIssuers: async function () {
                await this.aries.vdri.getDIDRecords().then(
                    resp => {
                        const data = resp.result
                        if (!data || data.length == 0) {
                            this.errors.push("No issuers found to select, please create an issuer")
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
                let data
                await this.aries.verifiable.generatePresentation({
                    presentation: {
                        "@context": "https://www.w3.org/2018/credentials/v1",
                        "type": "VerifiablePresentation",
                        "holder": `${this.issuers[this.selectedIssuer].key}`,
                    },
                    signatureType:"JsonWebSignature2020",
                    domain: this.domain,
                    challenge: this.challenge,
                    did: this.issuers[this.selectedIssuer].key
                }).then(resp => {
                        data = JSON.stringify(resp.verifiablePresentation)
                        //TODO bug in aries to show '"verifiableCredential": null,' in empty presentations
                        if (data.hasOwnProperty('verifiableCredential')) {
                            delete data.verifiableCredential
                        }
                    }
                ).catch(err => {
                    data = err
                    console.log('failed to create presentation, errMsg:', err)
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
<style lang="scss">
    .md-menu-content {
        width: auto;
        max-width: 100% !important;
    }

    .md-menu-content .md-ripple > span {
        position: relative !important;
    }
</style>
