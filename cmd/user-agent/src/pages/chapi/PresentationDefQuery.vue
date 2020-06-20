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
                    <md-icon style="color: #0E9A00; height: 40px;font-size: 30px !important;">verified_user</md-icon>
                    <span style="margin-left: 7px; font-weight: 700">{{ requestOrigin }} </span> would like you to share
                    below information,
                </h4>
            </div>

            <md-card style="margin-top: -5px" v-for="requirement in requirements" :key="requirement.name"
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
                <div>
                    <h4 class="md-subheading" id="result-header">
                        <md-icon style="color: #0E9A00; height: 40px;font-size: 20px !important;">done</md-icon>
                        Found {{ vcsFound.length}} credentials matching above criteria,
                    </h4>
                </div>

                <md-list class="md-triple-line" style="margin-top: -10px">
                    <md-list-item v-for="(vc, key) in vcsFound" :key="key">
                        <md-icon style="font-size: 50px !important;">security</md-icon>

                        <div class="md-list-item-text">
                            <span>{{vc.name}}</span>
                            <div class="md-subhead">{{vc.description}}</div>
                        </div>

                        <md-checkbox v-model="selectedVCs[key]" v-bind:id="'select-vc-' + key"></md-checkbox>
                    </md-list-item>

                </md-list>


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
                        <md-button class="md-primary md-raised" v-on:click="noCredential" >Close</md-button>
                    </md-empty-state>
                </div>
            </div>

        </div>
    </div>
</template>
<script>

    import {WalletGetByQuery, WalletManager} from "./wallet"

    const warning = "No credentials found in your wallet for above information asked"

    export default {
        beforeCreate: async function () {
            const aries = await this.$arieslib
            this.wallet = new WalletGetByQuery(aries, this.$parent.credentialEvent)

            this.requestOrigin = this.$parent.credentialEvent.credentialRequestOrigin
            this.registeredWalletUser = await new WalletManager().getRegisteredUser()
            if (!this.registeredWalletUser) {
                //this can never happen, but still one extra layer of security
                this.credentialWarning = 'Wallet is not registered'
            }

            this.presentation = await this.wallet.getPresentationSubmission()
            this.requirements = this.wallet.requirementDetails()

            let vcsFound = []
            this.presentation.verifiableCredential.forEach(function (vc) {
                vcsFound.push(vc)
            })
            this.vcsFound = vcsFound
            this.credentialWarning = vcsFound.length == 0 ? warning : ""

            this.loading = false
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
            createPresentation: async function () {
                this.loading = true
                await this.wallet.createAndSendPresentation(this.registeredWalletUser, this.presentation, this.selectedVCs)
                this.loading = false
            },
            cancel: async function () {
                this.wallet.cancel()
            },
            noCredential:async function () {
                this.wallet.sendNoCredntials()
            }
        },
        computed: {
            isShareDisabled() {
                return this.selectedVCs.length == 0
            },
        },
    }
</script>