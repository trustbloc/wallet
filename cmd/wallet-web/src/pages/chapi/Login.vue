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
                  <md-card-content>
                    <md-button v-on:click="beginOIDCLogin" class="md-raised md-success" id="loginBtn">
                      Login with Credentials
                    </md-button>
                      <md-button v-on:click="loginDevice" class="md-raised md-success" id="loginDeviceBtn">
                          Login with TouchID
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

    export default {
        created: async function () {
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
            if (this.getCurrentUser()) {
              await this.finishOIDCLogin()
              await this.getAgentInstance().store.flush()
              this.handleSuccess()
              return
            }

            this.loading = false
        },
        data() {
            return {
                statusMsg: '',
                loading: true,
            };
        },
        methods: {
            ...mapActions({loadUser: 'loadUser', loadOIDCUser: 'loadOIDCUser'}),
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
