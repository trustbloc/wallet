/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="content">
            <div class="md-layout">
                <div>
                    <md-label v-if="showOfflineWarning" style="color: #1B5E20; font-size: 16px; margin: 10px">
                        <md-icon>warning</md-icon>
                        <b>Warning:</b> Failed to connect to server. Your wallet can not participate in secured communication.
                    </md-label>
                    <md-card md-with-hover v-if="verifiableCredentials.length">
                        <md-card-header data-background-color="green">
                            <h4 class="title">
                                <md-icon>content_paste</md-icon>
                                Your Stored Credentials
                            </h4>
                        </md-card-header>
                        <md-card-content>
                            <simple-table v-for="vc in verifiableCredentials"
                                          :key=vc.id
                                          :name=credDisplayName(vc)
                                          :description=vc.description
                                          :data=vc
                                          :headerIcon=icon>
                            </simple-table>
                        </md-card-content>
                    </md-card>
                    <md-empty-state v-else
                                    md-icon="devices_other"
                                    :md-label=error
                                    :md-description=errorDescription>
                    </md-empty-state>
                </div>

            </div>
        </div>
    </div>

</template>

<script>
    import {SimpleTable} from "@/components";
    import {filterCredentialsByType, getCredentialType} from "@/pages/chapi/wallet";
    import {mapActions, mapGetters} from 'vuex'

    const manifestCredType = "IssuerManifestCredential"
    const governanceCredType = "GovernanceCredential"

    export default {
        components: {
            SimpleTable,
        },
        created: async function () {
            // Load the Credentials
            await this.getCredentials()
            await this.fetchAllCredentials()
            await this.refreshUserMetadata()

            this.username = this.getCurrentUser().username
            this.showOfflineWarning = this.getTrustblocOpts().walletMediatorURL && !JSON.parse(this.getCurrentUser().metadata).invitation
        },
        methods: {
            ...mapGetters('aries', {getAriesInstance: 'getInstance'}),
            ...mapGetters(['getCurrentUser', 'allCredentials', 'getTrustblocOpts']),
            ...mapActions(['getCredentials', 'refreshUserMetadata']),
            fetchAllCredentials: async function () {
                this.verifiableCredentials = []
                try {
                    for (let c of filterCredentialsByType(this.allCredentials(), [manifestCredType, governanceCredType])) {
                        let resp = await this.getAriesInstance().verifiable.getCredential({
                            id: c.id
                        })
                        this.verifiableCredentials.push(JSON.parse(resp.verifiableCredential))
                    }
                } catch (e) {
                    console.error('failed to get all stored credentials', e)
                    this.error = 'Failed to get your stored credetntials'
                    this.errorDescription = 'Unable to get stored credentials from your wallet, please try again later.'
                }

            },
            credDisplayName: function (vc) {
                return vc.name ? vc.name : getCredentialType(vc.type)
            }
        },
        data() {
            return {
                verifiableCredentials: [],
                username: '',
                aries: null,
                icon: 'perm_identity',
                showOfflineWarning: false,
                error: 'No stored credentials',
                errorDescription: 'Your wallet is empty, there aren\'t any stored credentials to show.',
            }
        }
    }
</script>

<style>
    .title {
        text-transform: capitalize;
    }

    .md-content {
        overflow: auto;
        padding: 1px;
        font-size: 6px;
        line-height: 16px;
    }

</style>


