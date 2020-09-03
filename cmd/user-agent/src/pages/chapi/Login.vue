/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div style="margin: 0 40%">
        <md-card md-with-hover>
            <md-card-header data-background-color="green">
                <div class="md-title">Login to your wallet</div>
            </md-card-header>

            <md-card-content v-if="loading" style="margin: 10% 22% 10% 30%">
                <md-progress-spinner :md-diameter="100" class="md-accent" :md-stroke="10"
                                     md-mode="indeterminate"></md-progress-spinner>
            </md-card-content>

            <md-card-content v-else>
                <form>
                    <md-field md-clearable>
                        <label>Enter your username here</label>
                        <md-input v-model="username"></md-input>
                    </md-field>

                    <md-field md-clearable :md-toggle-password="false">
                        <label>Enter your password here</label>
                        <md-input v-model="password" type="password"></md-input>
                    </md-field>

                    <div style="font-weight: 400; color: #d73a49; margin: 2px" v-if="statusMsg.length">
                        <md-icon style="color: red">error</md-icon>
                        {{ statusMsg }}
                    </div>

                    <md-button v-on:click="login" class="md-raised md-success" id="loginBtn">
                        Login
                    </md-button>
                </form>

            </md-card-content>

        </md-card>
    </div>
</template>

<script>
    import {RegisterWallet} from "./wallet"

    export default {
        beforeCreate: async function () {
            this.$store.dispatch('initUserStore')
            if (this.$store.getters.getUser) {
                this.$router.push("/dashboard");
                return
            }

            this.registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, await this.$arieslib,
                this.$trustblocAgent, await this.$trustblocStartupOpts)
        },
        data() {
            return {
                username: 'walletuser',
                password: 'mysecurepassword',
                statusMsg: '',
                loading: false,
            };
        },
        methods: {
            login: async function () {
                this.loading = true
                // TODO OIDC login intergation, for now all logins will succeed
                try {
                    await this.registrar.register(this.username)
                } catch (e) {
                    console.error(e)
                    this.statusMsg = e.toString()
                    this.loading = false
                    return
                }
                this.handleSuccess()
                this.loading = false
            },
            handleSuccess() {
                this.$store.dispatch('setUser', this.username)
                this.$router.push("/dashboard");
            }
        }
    }

</script>

