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
    <div v-else class="md-layout w-screen flex flex-col items-center">
        <div class="md-layout-item max-w-screen-md">
            <div>
                <h5 class="md-subheading">
                    This verifier would like you to share below information,
                </h5>
                <div class="md-alignment-center-center">
                    <p class="md-body-1">

                        <span class="md-subheading" style="color: #025C8F;">
                            <md-icon style="color: #00cc66;" class="md-size-1x">verified_user</md-icon>
                            {{ requestOrigin }} </span> <br>
                        <span class="md-caption" v-if="reason.length">{{ reason }} </span>
                    </p>
                </div>
            </div>

            <div v-if="errors.length">
                <b>Failed with following error(s):</b>
                <md-field style="margin-top: -15px">
                    <ul>
                        <li v-for="error in errors" :key="error">{{ error }}</li>
                    </ul>
                </md-field>
            </div>

            <div>
                <md-card v-for="(record, key) in allRecords" :key="key">
                    <md-card-header>
                        <md-card-header-text>
                            <div class="md-title">{{record.credential.name ? record.credential.name :
                                getCredentialType(record.credential.type)}}
                            </div>
                            <div class="md-subhead">{{record.credential.description}}</div>
                        </md-card-header-text>

                        <md-card-media>
                            <md-icon class="md-size-3x">{{getRandomIcon(key)}}</md-icon>
                        </md-card-media>
                    </md-card-header>

                    <md-card-content>
                        The verifier would like to access below information from your <span class="md-subheading">{{record.credential.name ? record.credential.name  : getCredentialType(record.credential.type)}}</span>
                        <md-list>
                            <div v-for="(subj, skey) in record.output.credentialSubject" :key="skey">
                                <p class="md-caption" v-if="displayContent(skey)">
                                    {{skey}} : {{subj}}
                                </p>
                            </div>
                        </md-list>
                    </md-card-content>

                    <md-card-actions v-if="isDecisionMade(key)">
                        <md-button v-if="isAllowed(key)" class="md-accent" style="background-color: #fb4934 !important;" v-on:click="deny(key)">Deny</md-button>
                        <md-button v-if="!isAllowed(key)" class="md-accent" style="background-color: #fb4934 !important;"><md-icon>cancel</md-icon>Denied</md-button>

                        <md-button v-if="!isAllowed(key)" class="md-accent" style="background-color: #0E9A00 !important;" v-on:click="allow(key)">Allow</md-button>
                        <md-button v-if="isAllowed(key)" class="md-accent" style="background-color: #0E9A00 !important;"><md-icon>check</md-icon>Allowed</md-button>
                    </md-card-actions>
                    <md-card-actions v-else>
                        <md-button class="md-accent" style="background-color: #fb4934 !important;" v-on:click="deny(key)">Deny</md-button>
                        <md-button class="md-accent" style="background-color: #0E9A00 !important;" v-on:click="allow(key)">Allow</md-button>
                    </md-card-actions>
                </md-card>

                <md-card-content class="flex flex-row justify-between" style="padding-left: 0; padding-right: 0;">
                    <md-button v-on:click="share"
                               class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                               id="share-credentials" :disabled=isShareDisabled>Share
                    </md-button>
                    <md-button v-on:click="cancel" class="md-cancel-text" id="cancelBtn">
                        Decline
                    </md-button>
                </md-card-content>
            </div>

        </div>
    </div>
</template>
<script>

    import {SelectiveDisclosure} from "./wallet"
    import {mapGetters} from 'vuex'

    const nonDisplayContent = ['id', 'type']

    export default {
        created: async function () {
            this.wallet = new SelectiveDisclosure(this.getAgentInstance(), this.$parent.credentialEvent)
            this.reason = this.wallet.credentialQuery.reason
            this.requestOrigin = this.$parent.credentialEvent.credentialRequestOrigin

            await this.loadCredentials()
            this.loading = false
        },
        data() {
            return {
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
                reason: "",
                allIcons: ['account_box', 'contacts', 'person', 'person_outline', 'card_membership', 'portrait', 'bento'],
                selectedFrames: [],
                decisionMade: []
            };
        },
        methods: {
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            ...mapGetters(['getCurrentUser']),
            loadCredentials: async function () {
                try {
                    this.allRecords = await this.wallet.queryByFrame()
                } catch (err) {
                    this.errors.push('failed to get saved credentials')
                    console.log('get credentials failed, error:', err)
                    return
                }

                if (this.allRecords.length == 0) {
                    this.errors.push('No credentials found matching request criteria')
                }
            },
            async share() {
                this.loading = true
                let filtered = this.allRecords.filter((k, v) => this.selectedFrames.includes(v))
                this.wallet.generatePresentation(this.getCurrentUser().username, filtered)
            },
            cancel() {
                this.wallet.cancel()
            },
            getRandomIcon(key) {
                return this.allIcons[key % this.allIcons.length]
            },
            getCredentialType(types) {
                return types.filter(type => type != "VerifiableCredential")[0]
            },
            displayContent(k) {
                return !nonDisplayContent.includes(k)
            },
            allow(key) {
                this.selectedFrames.push(key)
                this.decisionMade.push(key)
            },
            deny(key) {
                this.selectedFrames = this.selectedFrames.filter(k => k != key)
                this.decisionMade.push(key)
            },
            isAllowed(key) {
                return this.selectedFrames.includes(key)
            },
            isDecisionMade(key) {
                return this.decisionMade.includes(key)
            },
        },
        computed: {
            isShareDisabled() {
                return this.selectedFrames.length == 0;
            },
        },
    }
</script>
