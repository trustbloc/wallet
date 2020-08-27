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
                <md-card-header>
                    <h4 class="title">Connect your wallet</h4>
                </md-card-header>

                <md-card-content v-if="!credentialWarning.length" style="background-color: white; ">
                    <div v-if="errors.length">
                        <b>Failed with following error(s):</b>
                        <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                        </ul>
                    </div>

                    <md-card-content class="viewport">
                        <p class="viewport"> Hello {{ this.walletUser.id }},</p>
                            <span style="color: #0E9A00">{{ requestOrigin }}</span> would like to
                        connect to your wallet for secured communication.
                    </md-card-content>

                    <md-card-content v-if="userCredentials.length" class="viewport">
                        Here are the credentials being sent to your wallet,

                        <md-list class="md-double-line">
                            <md-list-item v-for="credential in userCredentials" :key="credential">
                                <md-icon class="md-primary md-size-2x" >perm_identity</md-icon>

                                <div class="md-list-item-text">
                                    <span>{{credential.name ? credential.name : 'Credential name not provided'}}</span>
                                    <span>{{credential.description}}</span>
                                </div>


                            </md-list-item>
                        </md-list>
                    </md-card-content>

                    <md-divider></md-divider>

                    <md-card-content class="center-span">
                        <md-button v-on:click="connect"
                                   class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                                   id="didconnect">{{buttonLabel}}
                        </md-button>
                        <md-button v-on:click="cancel" style="margin-left: 5%" class="md-cancel-text" id="cancelBtn">
                            Cancel
                        </md-button>
                    </md-card-content>

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

            this.wallet = new DIDConn(await this.$arieslib,
                await new this.$trustblocAgent.Framework(await this.$trustblocStartupOpts), this.$trustblocStartupOpts,
                this.$parent.credentialEvent, this.walletUser)
            this.requestOrigin = this.$parent.credentialEvent.credentialRequestOrigin
            this.userCredentials = this.wallet.getUserCredentials()
            this.buttonLabel =  this.userCredentials.length > 0 ? 'Store & Connect' : 'Connect'

            this.loading = false
        },
        data() {
            return {
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
                userCredentials: [],
                buttonLabel : "Connect",
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
