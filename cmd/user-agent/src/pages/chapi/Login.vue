/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <md-card style="width: 20%; margin: 15% 0% 0% 40%" md-with-hover>
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
        created: async function () {
            let redirect = this.$route.params['redirect']
            this.redirect = redirect ? {name: redirect} : `${__webpack_public_path__}`

            this.loadUser()
            if (this.getCurrentUser()) {
                this.handleSuccess()
                return
            }
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
            ...mapActions({loginUser: 'login', loadUser: 'loadUser'}),
            ...mapGetters(['getCurrentUser', 'getAgentOpts']),
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            // this function to be called after successful OIDC login
            login: async function () {
                this.loading = true

                await this.loginUser(this.username)
                let user = this.getCurrentUser()

                let registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, this.getAgentInstance(),
                    this.getAgentOpts())

                try {
                    if (!user || !user.metadata) {
                        // first time login, register this user
                        await registrar.register(this.username)
                    }

                    await registrar.installHandlers(this.username)
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
