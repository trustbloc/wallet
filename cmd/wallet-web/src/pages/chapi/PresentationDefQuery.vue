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
        <div class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
            <div v-if="errors.length">
                <b>Failed with following error(s):</b>
                <md-field style="margin-top: -15px">
                    <ul>
                        <li v-for="error in errors" :key="error">{{ error }}</li>
                    </ul>
                </md-field>
            </div>

            <div>
                <h4 class="md-subheading">
                    This verifier would like you to share below information,
                </h4>
            </div>

            <div style="padding-bottom: 10px">
                <governance :govn-v-c="govnVC" :request-origin="requestOrigin" :issuer="false"/>
            </div>


            <md-card v-for="requirement in requirements" :key="requirement.name"
                     :value="requirement">

                <md-card-expand>
                    <div class="md-title" style="margin-left: 10px; margin-top: 10px">{{requirement.name}}</div>
                    <md-card-actions md-alignment="space-between">
                        <div class="md-subhead">{{requirement.purpose}}</div>
                        <md-card-expand-trigger>
                            <md-button class="md-icon-button">
                                <md-icon>keyboard_arrow_down</md-icon>
                            </md-button>
                        </md-card-expand-trigger>
                    </md-card-actions>

                    <md-card-expand-content>
                        <md-card-content>
                            {{requirement.rule}}
                            <ul>
                                <li v-for="descriptor in requirement.descriptors" :key="descriptor.name">
                                    <b>{{descriptor.name}} </b>{{descriptor.purpose}}
                                    <ul>
                                        <li v-for="constraint in descriptor.constraints" :key="constraint">
                                            {{ constraint}}
                                        </li>
                                    </ul>
                                </li>
                            </ul>
                        </md-card-content>
                    </md-card-expand-content>
                </md-card-expand>
            </md-card>


            <div v-if="!credentialWarning.length">
                <div v-if="credsFound">
                    <p id="result-header-1">
                        <md-icon style="color: #0E9A00; height: 40px;font-size: 20px !important;">done</md-icon>
                        Found {{ credsFound}} stored credential{{credsFound > 1 ? 's' : ''}} matching
                        above criteria in your wallet,
                    </p>

                    <md-list class="md-triple-line" style="margin-top: -10px">
                        <div v-for="(vc, key) in vcsFound" :key="key">
                            <md-list-item v-if="!isManifest(vc)">
                                <md-icon class="md-primary md-size-2x">perm_identity</md-icon>

                                <div class="md-list-item-text">
                                    <span>{{vc.name}}</span>
                                    <div class="md-subhead">{{vc.description}}</div>
                                </div>

                                <md-checkbox v-model="selectedVCs[key]" v-bind:id="'select-vc-' + key"></md-checkbox>
                            </md-list-item>
                        </div>
                    </md-list>
                </div>


                <div v-if="credsFound && issuersFound" style="margin: 30px"></div>

                <div v-if="issuersFound">
                    <p id="result-header-2">
                        <md-icon style="color: #0E9A00; height: 40px;font-size: 20px !important;">done</md-icon>
                        Found {{ issuersFound}} issuer{{issuersFound > 1 ? 's' : ''}} who can issue
                        credentials matching above criteria in your wallet,
                    </p>

                    <md-list class="md-triple-line" style="margin-top: -10px">
                        <div v-for="(vc, key) in vcsFound" :key="key">
                            <md-list-item v-if="isManifest(vc)">
                                <md-icon class="md-primary md-size-2x">security</md-icon>

                                <div class="md-list-item-text">
                                    <span>{{vc.name}}</span>
                                    <div class="md-subhead">{{vc.description}}</div>
                                </div>

                                <md-checkbox v-model="selectedVCs[key]" v-bind:id="'select-vc-' + key"></md-checkbox>
                            </md-list-item>
                        </div>
                    </md-list>
                </div>


                <div style="margin-left: 30%; margin-top: 5px">
                    <md-button v-on:click="createPresentation"
                               class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                               id="share-credentials" :disabled=isShareDisabled>Share
                    </md-button>
                    <md-button v-on:click="cancel" style="margin-left: 5px" class="md-cancel-text" id="cancelBtn">
                        Decline
                    </md-button>
                </div>
            </div>

            <div v-else>
                <div>
                    <md-empty-state
                            md-icon="devices_other"
                            md-label="No credentials found"
                            :md-description="credentialWarning">
                        <md-button class="md-primary md-raised" v-on:click="noCredential">Close</md-button>
                    </md-empty-state>
                </div>
            </div>

        </div>
    </div>
</template>
<script>

    import {filterCredentialsByType, getCredentialType, WalletGetByQuery} from "./wallet"
    import Governance from "./Governance.vue";
    import {mapGetters} from 'vuex'

    const warning = "No credentials found in your wallet for information asked above"
    const manifestCredType = 'IssuerManifestCredential'

    export default {
        components: {Governance},
        created: async function () {
            this.wallet = new WalletGetByQuery(this.getAgentInstance(), this.$parent.credentialEvent, this.getAgentOpts())
            await this.wallet.connect()

            this.requestOrigin = this.$parent.credentialEvent.credentialRequestOrigin

            try {
                this.presentation = await this.wallet.getPresentationSubmission()
            } catch (e) {
                this.credentialWarning = "Some unexpected error occurred, please try again later"
                this.loading = false
                return
            }

            this.requirements = this.wallet.requirementDetails()

            let vcsFound = []
            this.presentation.verifiableCredential.forEach(function (vc) {
                vcsFound.push(vc)
            })

            // TODO handling multiple governance VCs
            this.govnVC = this.wallet.govnVC.length > 0 ? this.wallet.govnVC[0] : undefined

            this.loading = false

            if (vcsFound.length == 0) {
                this.credentialWarning = warning
                return
            }

            this.vcsFound = vcsFound
            this.credsFound = filterCredentialsByType(vcsFound, [manifestCredType]).length
            this.issuersFound = filterCredentialsByType(vcsFound, [manifestCredType], true).length
        },
        data() {
            return {
                vcsFound: [],
                selectedVCs: [],
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
                searched: [],
                reason: "",
                requirements: []
            };
        },
        methods: {
            ...mapGetters(['getCurrentUser', 'getAgentOpts']),
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            createPresentation: async function () {
                this.loading = true
                await this.wallet.createAndSendPresentation(this.getCurrentUser().username, this.presentation, this.selectedVCs)
                this.loading = false
            },
            cancel: async function () {
                this.wallet.cancel()
            },
            noCredential: async function () {
                this.wallet.sendNoCredntials()
            },
            isManifest(vc) {
                return getCredentialType(vc.type) == manifestCredType
            }
        },
        computed: {
            isShareDisabled() {
                return this.selectedVCs.length == 0
            },
        },
    }
</script>
