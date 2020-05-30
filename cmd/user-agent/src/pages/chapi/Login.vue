/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div style="height: 213px">
        <div v-if="loading" style="padding-left: 30%">
            <div class="md-layout">
                <md-progress-spinner :md-diameter="150" class="md-accent" :md-stroke="10"
                                     md-mode="indeterminate"></md-progress-spinner>
            </div>
        </div>
        <div v-else>
            <div v-if="!statusMsg.length">
                <md-field md-clearable>
                    <label>Enter your username here</label>
                    <md-input v-model="username"></md-input>
                </md-field>

                <md-field md-clearable :md-toggle-password="false">
                    <label>Enter your password here</label>
                    <md-input v-model="password" type="password"></md-input>
                </md-field>

                <md-button v-on:click="login" class="md-raised md-success" id="storeVCBtn">
                    Login
                </md-button>
            </div>


            <div v-else style="text-align: center;">
                <div>
                    <md-icon class="md-size-5x" :style="statusStyle">{{ statusIcon}}</md-icon>
                </div>
                <h3 style="color: #0d47a1"> {{ statusMsg }}</h3>
            </div>

        </div>
    </div>
</template>

<script>

    import {DIDManager, RegisterWallet} from "./wallet"


    export default {
        beforeCreate: async function () {
            this.registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler)
            const aries = await this.$arieslib
            const opts = await this.$trustblocStartupOpts
            this.didManager = new DIDManager(aries, this.$trustblocAgent, opts)
        },
        data() {
            return {
                username: 'walletuser',
                password: 'mysecurepassword',
                statusMsg: '',
                loading: false,
                statusIcon: 'devices_other',
                statusStyle: 'color: #0E9A00;'
            };
        },
        methods: {
            login: async function () {
                this.loading = true
                this.statusMsg = "wallet registered successfully !!"

                try {
                    await this.createDID()
                    await this.registrar.register(this.username)
                } catch (e) {
                    console.error(e)
                    this.showFailure(e)
                }

                this.loading = false
            },
            createDID: async function () {
                // create DID
                let did = await this.didManager.createDID("Ed25519", "Ed25519Signature2018")

                // save DID
                await this.didManager.saveDID(this.username, did)

                // save DID metadata
                this.didManager.storeDIDMetadata(did.id, {
                    signatureType: "Ed25519Signature2018",
                    friendlyName: this.username
                })

                console.log(`created DID ${did.id} successfully for user ${this.username}`)
            },
            showFailure(e){
                this.statusMsg = e.toString()
                if(this.statusMsg.includes("did name already exists")) {
                    this.statusIcon = 'warning'
                    this.statusStyle = 'color: #E9DC51;'
                    this.statusMsg = "Wallet already registered for this username"
                } else {
                    this.statusIcon = 'error'
                    this.statusStyle = 'color: #cc2127;'
                }
            }
        }
    }

</script>

