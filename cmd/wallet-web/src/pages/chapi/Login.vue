/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="viewport">
        <md-card class="md-layout-item md-size-50 md-small-size-100">
            <md-card-header data-background-color="green">
                <div class="md-title text-center">Web Wallet Login</div>
            </md-card-header>
            <md-card-content>
                <img src="@/assets/img/wallet.jpg">
            </md-card-content>
            <div class="md-body-1 center-header">Log in or register with the same sign-in information you use for other online services (for example, online banking).</div>
            <md-card-content v-if="loading" style="margin: 10% 22% 10% 30%">
                <md-progress-spinner :md-diameter="100" class="md-accent" :md-stroke="10"
                                     md-mode="indeterminate"></md-progress-spinner>
            </md-card-content>

            <md-card-content v-else>
                <form>
                    <md-card-content>
                        <md-button v-on:click="beginOIDCLogin" class="md-dense md-raised md-success" id="loginBtn">
                            Sign-In Partner Login/Register
                        </md-button>
                            <md-button  v-if="registered" v-on:click="loginDevice" class="md-dense md-raised md-success" id="loginDeviceBtn">
                                Sign-In Touch/Face ID
                            </md-button>
                    </md-card-content>
                </form>
            </md-card-content>

        </md-card>
    </div>
</template>

<script>
    import {DeviceLogin, RegisterWallet} from "./wallet"
    import {mapActions, mapGetters} from 'vuex'
    import Vue from 'vue';

    export default {
        created: async function () {
            //TODO: issue-601 Implement cookie logic with information from the backend.
           this.registerSuccess = Vue.$cookies.get("device");
            this.deviceLogin = new DeviceLogin(this.getAgentOpts());
            let redirect = this.$route.params['redirect']
            this.redirect = redirect ? {name: redirect} : `${__webpack_public_path__}`
            this.loadUser()
            if (this.getCurrentUser()) {
              await this.getAgentInstance().store.flush()
              this.handleSuccess()
                return
            }

            await this.loadOIDCUser()

            try {
              await this.refreshUserMetadata()
            } catch (e) {
              console.warn("first time login: ignore", e)
            }

            if (this.getCurrentUser()) {
              await this.finishOIDCLogin()
              await this.getAgentInstance().store.flush()
              this.handleSuccess()
              return
            }
            if (this.$cookies.isKey('registerSuccess')) {
                this.registered = true;
            }

            this.loading = false
        },
        data() {
            return {
                statusMsg: '',
                loading: true,
                registered: false,
            };
        },
        methods: {
            ...mapActions({loadUser: 'loadUser', loadOIDCUser: 'loadOIDCUser', refreshUserMetadata: 'refreshUserMetadata'}),
            ...mapGetters(['getCurrentUser', 'getAgentOpts', 'serverURL']),
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            beginOIDCLogin: async function() {
                window.location.href = this.serverURL() + "/oidc/login"
            },

           sleep(ms) {
               return new Promise(resolve => setTimeout(resolve, ms));
             },
          finishOIDCLogin: async function() {
                let user = this.getCurrentUser()

                let registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, this.getAgentInstance(),
                    this.getAgentOpts())

                try {
                    if (!user.metadata) {
                        // first time login, register this user
                        await registrar.register(user.username)
                    }

                  await this.getAgentInstance().store.flush()

                  await registrar.installHandlers(user.username)
                } catch (e) {
                    this.handleFailure(e)
                }
            },
            handleSuccess() {
              this.$router.push(this.redirect);
            },
            handleFailure(e) {
                console.error("login failure: ", e)
                this.statusMsg = e.toString()
            },
            loginDevice: async function () {
                await this.deviceLogin.login();
            },

        }
    }

</script>
<style lang="css">
    .md-button {
          text-transform: none !important; /*For Lower case use lowercase*/
          font-size: 16px !important;
           width: 100%;
           max-width: 100%;
           position:absolute;

    }
    .viewport {
        width: 40%;
        max-width: 100%;
        display: inline-flex;
        justify-content: center;
        align-content: center;
        top:25%;
        left: 30%;
        position:absolute;
        margin-left:auto;
        margin-right:auto
    }
</style>
