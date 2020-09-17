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
    </div>
</template>
<script>

    import {searchByTypeAndHolder, WalletGet} from "./wallet"
    import {mapGetters} from 'vuex'

    export default {
        created: async function () {
            this.wallet = new WalletGet(this.getAriesInstance(), this.$parent.credentialEvent)
            this.search = this.wallet.search
            this.reason = this.wallet.reason
            this.credentialEvent = this.$parent.credentialEvent
            this.requestOrigin = this.credentialEvent.credentialRequestOrigin


            await this.loadIssuers()
            await this.loadCredentials()
            this.loading = false

            // perform search while loading
            this.searched = []
            this.searchOnTable()
        },
        data() {
            return {
                savedVCs: [{id: 0, name: "Select VC"}],
                selectedVCs: [],
                issuers: [{id: 0, name: "Select Identity"}],
                selectedIssuer: 0,
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
                search: "",
                searched: [],
                reason: "",
            };
        },
        methods: {
            ...mapGetters('aries', {getAriesInstance: 'getInstance'}),
            searchOnTable() {
                this.searched = searchByTypeAndHolder(this.savedVCs, this.search, this.issuers[this.selectedIssuer].key)
            },
            loadCredentials: async function () {
                this.savedVCs.length = 0

                try {
                    this.savedVCs = await this.wallet.getCredentialRecords()

                    if (this.savedVCs.length == 0) {
                        this.credentialWarning = 'No Saved Credentials Found'
                        return
                    }

                } catch (err) {
                    this.errors.push('Failed to get credentials')
                    console.log('get credentials failed, error:', err)
                }
            },
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
            onSelect(items) {
                this.selectedVCs = items
            },
            createPresentation: async function () {
                if (this.selectedVCs.length == 0) {
                    this.errors.push("Please select at least one credential")
                    return
                }

                this.loading = true
                await this.wallet.createAndSendPresentation(this.issuers[this.selectedIssuer].key, this.selectedVCs)
                this.loading = false
            }
        }
    }
</script>
