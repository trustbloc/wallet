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

                    <md-table v-model="searched" md-sort="name" md-card @md-selected="onSelect">
                        <md-table-empty-state
                                md-label="No credentials found"
                                :md-description="`No credentials found for this  type. Try a different type search or add a new credential.`">
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

    import {WalletGetByQuery, WalletManager} from "./wallet"

    export default {
        beforeCreate: async function () {
            const aries = await this.$arieslib
            this.wallet = new WalletGetByQuery(aries, this.$parent.credentialEvent)
            this.credentialEvent = this.$parent.credentialEvent
            this.requestOrigin = this.credentialEvent.credentialRequestOrigin
            this.registeredWalletUser = await new WalletManager().getRegisteredUser()
            if (!this.registeredWalletUser) {
                //this can never happen, but still one extra layer of security
                this.credentialWarning = 'Wallet is not registered'
            }

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
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
                searched: [],
                reason: "",
            };
        },
        methods: {
            searchOnTable() {
                this.searched = this.savedVCs
            },
            loadCredentials: async function () {
                this.savedVCs.length = 0

                try {
                    this.savedVCs = await this.wallet.getCredentialRecords(this.registeredWalletUser.did)

                    if (this.savedVCs.length == 0) {
                        this.credentialWarning = 'No Saved Credentials Found'
                        return
                    }

                } catch (err) {
                    this.errors.push('Failed to get credentials')
                    console.log('get credentials failed, error:', err)
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
                await this.wallet.createAndSendPresentation(this.registeredWalletUser, this.selectedVCs)
                this.loading = false
            }
        }
    }
</script>
