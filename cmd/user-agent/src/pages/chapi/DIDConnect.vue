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


            <md-card class="md-card-plain">
                <md-card-header data-background-color="green">
                    <h4 class="title">Wallet Connection</h4>
                </md-card-header>

                <md-card-content v-if="!credentialWarning.length" style="background-color: white;">
                    <div v-if="errors.length">
                        <b>Failed with following error(s):</b>
                        <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                        </ul>
                    </div>

                    <h3 style="font-size: 20px"> Hi {{ this.walletUser.id }},</h3>
                    <h3 style="font-size: 20px"><span style="color: #0E9A00">{{ requestOrigin }}</span> would like to
                        connect to your wallet
                    </h3>

                    <md-button v-on:click="connect" style="margin-right: 5px; margin-left: 25%"
                               class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                               id="didconnect">Allow
                    </md-button>

                    <md-button v-on:click="cancel" class="md-cancel-text" id="cancelBtn" style="margin-left: 5px">
                        Deny
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

    import {DIDConn, WalletManager} from "./wallet"

    export default {
        beforeCreate: async function () {
            this.walletUser = await new WalletManager().getRegisteredUser()
            if (!this.walletUser) {
                //this can never happen, but still one extra layer of security
                this.credentialWarning = 'Wallet is not registered'
            }
            const aries = await this.$arieslib
            this.wallet = new DIDConn(aries, this.$parent.credentialEvent, this.walletUser)
            this.requestOrigin = this.$parent.credentialEvent.credentialRequestOrigin
            this.loading = false
        },
        data() {
            return {
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
            };
        },
        methods: {
            cancel: async function () {
                this.wallet.cancel()
            },
            connect: async function () {
                this.loading = true
                await this.wallet.connect()
                this.loading = false
            }
        }
    }
</script>