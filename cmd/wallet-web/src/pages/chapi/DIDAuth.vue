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
                        <md-select v-model="selectedIssuer" id="select-did">
                            <md-option v-for="{id, name} in issuers" :key="id" :value="id" :id="name">
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
    </div>

</template>
<script>

    import {DIDAuth} from "./wallet"
    import {mapGetters} from 'vuex'

    export default {
        created: async function () {
            this.wallet = new DIDAuth(this.getAgentInstance(), this.$parent.credentialEvent)
            this.requestOrigin = this.$parent.credentialEvent.credentialRequestOrigin

            await this.loadIssuers()
            this.loading = false
        },
        data() {
            return {
                issuers: [{id: 0, name: "Select Identity"}],
                selectedIssuer: 0,
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
            };
        },
        methods: {
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            loadIssuers: async function () {
                try {
                    this.issuers = await this.wallet.getDIDRecords()

                    if (this.issuers.length == 0) {
                        this.credentialWarning = 'Issuers not found, please create an issuer'
                        return
                    }

                } catch (err) {
                    this.errors.push(err)
                }
            },
            cancel: async function () {
                this.wallet.cancel()
            },
            authorize: async function () {
                this.loading = true
                await this.wallet.authorize(this.issuers[this.selectedIssuer].key)
                this.loading = false
            }
        }
    }
</script>
