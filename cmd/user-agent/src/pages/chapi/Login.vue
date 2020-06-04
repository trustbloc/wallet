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

                <md-button v-on:click="login" class="md-raised md-success" id="loginBtn">
                    Login
                </md-button>
            </div>


            <div v-else style="text-align: center;">
                <div>
                    <md-icon class="md-size-5x" :style="statusStyle">{{ statusIcon}}</md-icon>
                </div>
                <div>
                    <label name="statusMsg" style="font-size: 20px"> {{ statusMsg }}</label>
                </div>
                    <md-button  v-if="registered" v-on:click="logout" class="md-raised md-cancel-text" id="logoutBtn">
                        Logout
                    </md-button>
            </div>

        </div>
    </div>
</template>

<script>

    import {RegisterWallet} from "./wallet"

    export default {
        beforeCreate: async function () {
            const aries = await this.$arieslib
            const opts = await this.$trustblocStartupOpts
            this.registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, aries, this.$trustblocAgent, opts)

            const wuser = await this.registrar.getRegisteredUser()
            if (wuser) {
                this.registered = true
                this.statusMsg = `Wallet already registered for ${wuser.id}`
            }
        },
        data() {
            return {
                username: 'walletuser',
                password: 'mysecurepassword',
                statusMsg: '',
                loading: false,
                statusIcon: 'devices_other',
                statusStyle: 'color: #0E9A00;',
                registered: false,
            };
        },
        methods: {
            login: async function () {
                this.loading = true
                this.statusMsg = "wallet registered successfully !!"

                // TODO OIDC login intergation, for now all logins will succeed
                try {
                    await this.registrar.register(this.username)
                } catch (e) {
                    console.error(e)
                    this.handleFailure(e)
                    return
                }
                this.handleSuccess()
            },
            logout: async function () {
                this.loading = true
                // TODO OIDC logout intergation
                await this.registrar.unregister()
                this.resetView()
            },
            handleFailure(e) {
                this.statusMsg = e.toString()
                this.statusIcon = 'error'
                this.statusStyle = 'color: #cc2127;'
                this.registered = false
                this.loading = false
            },
            handleSuccess() {
                this.registered = true
                this.statusMsg = `wallet registered successfully for ${this.username}`
                this.loading = false
            },
            resetView() {
                this.statusMsg = ''
                this.registered = false
                this.loading = false
            }
        }
    }

</script>

