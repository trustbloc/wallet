/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
<template>

    <div v-if="loading" class="w-screen" style="margin-left: 40%;margin-top: 20%;height: 200px;">
        <div class="md-layout">
            <md-progress-spinner :md-diameter="100" class="md-accent" :md-stroke="10"
                                 md-mode="indeterminate"></md-progress-spinner>
        </div>
    </div>

    <div v-else class="md-layout w-screen flex justify-center">
        <div class="md-layout-item max-w-screen-md">

            <md-card class="md-card-plain">
                <md-card-header>
                    <h4 class="title">Authenticate Your Wallet</h4>
                </md-card-header>

                <md-card-content style="background-color: white; ">

                    <div v-if="errors.length">
                        <b>Failed with following error(s):</b>
                        <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                        </ul>
                    </div>

                    <md-card-content class="viewport">
                        This issuer would like to you to authenticate.
                        <governance :govn-v-c="govnVC" :request-origin="requestOrigin"/>
                    </md-card-content>

                    <md-divider></md-divider>

                    <md-card-content class="md-layout md-alignment-center-center">
                        <md-button v-on:click="authorize"
                                   style="margin-right: 5%"
                                   class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                                   id="didauth">Authenticate
                        </md-button>
                        <md-button v-on:click="cancel" class="md-cancel-text" id="cancelBtn">
                            Cancel
                        </md-button>
                    </md-card-content>

                </md-card-content>
            </md-card>

        </div>
    </div>

</template>
<script>

    import {CHAPIEventHandler} from "./wallet"
    import {CredentialManager} from "@trustbloc/wallet-sdk"
    import Governance from "./Governance.vue";
    import {mapGetters} from 'vuex'

    export default {
        components: {Governance},
        created: async function () {
            this.chapiHandler = new CHAPIEventHandler(this.$parent.credentialEvent)
            let {query} = this.chapiHandler.getEventData()

            let {user, token} = this.getCurrentUser().profile
            this.credentialManager = new CredentialManager({agent: this.getAgentInstance(), user})

            try {
                let {results} = await this.credentialManager.query(token, Array.isArray(query) ? query : [query])
                this.presentation = results[0]
            } catch (e) {
                console.error('failed to prepare DIDAuth response:', e)
                this.errors.push('failed to handle request, try again later.')
                this.loading = false
            }

            this.requestOrigin = this.chapiHandler.getRequestor()
            this.loading = false
        },
        data() {
            return {
                issuers: [{id: 0, name: "Select Identity"}],
                errors: [],
                requestOrigin: "",
                loading: true,
                govnVC: null,
            };
        },
        methods: {
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            ...mapGetters(['getCurrentUser']),
            cancel: async function () {
                this.chapiHandler.cancel()
            },
            authorize: async function () {
                this.loading = true

                let {profile, preference} = this.getCurrentUser()
                let {controller, proofType, verificationMethod} = preference
                let {domain, challenge} = this.chapiHandler.getEventData()

                let {presentation} = await this.credentialManager.present(profile.token, {presentation: this.presentation}, {
                    controller, proofType, domain, challenge, verificationMethod
                })

                this.chapiHandler.present(presentation)

                this.loading = false
            }
        }
    }
</script>
