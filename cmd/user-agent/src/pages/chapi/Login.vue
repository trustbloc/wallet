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
    import {mapActions, mapGetters} from 'vuex'

    export default {
        beforeCreate: async function () {
            let redirect = this.$route.params['redirect']
            this.redirect = redirect ? {name: redirect} : '/'

            this.$store.dispatch('loadUser')
            if (this.$store.getters.getCurrentUser) {
                this.$router.push(this.redirect);
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
            ...mapActions({loginUser: 'login'}),
            ...mapGetters(['getCurrentUser']),
            // this function to be called after successful OIDC login
            login: async function () {
                this.loading = true

                await this.loginUser(this.username)
                let user = this.getCurrentUser()

                try {
                    if (!user || !user.metadata) {
                        // first time login, register this user
                        await this.registrar.register(this.username)
                    }

                    await this.registrar.installHandlers(this.username)
                } catch (e) {
                    this.handleFailure(e)
                    this.loading = false
                    return
                }

                this.handleSuccess()
                this.loading = false
            },
            handleSuccess() {
                this.$router.push(this.redirect);
            },
            handleFailure(e) {
                console.error(e)
                this.statusMsg = e.toString()
            }
        }
    }

</script>

