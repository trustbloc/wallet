/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="gradient">
  <!-- This is an example component -->
    <div class="px-4 py-24">
      <div
          class="md:max-w-6xl md:mx-auto md:flex md:items-center md:justify-between">
        <div class="flex justify-between items-center">
          <img class="h-8 w-8 mr-2" src="@/assets/img/logo.png" alt="">
          <a class="font-semibold text-white text-2xl lg:text-4xl tracking-tight" href="javascript:history.back()">TrustBloc</a>
          <div
              class="inline-block cursor-pointer md:hidden">
            <div class="bg-gray-400 w-8 mb-2" style="height: 2px;"></div>
            <div class="bg-gray-400 w-8 mb-2" style="height: 2px;"></div>
            <div class="bg-gray-400 w-8" style="height: 2px;"></div>
          </div>
        </div>

        <div>
          <div class="hidden md:block  text-lg">
            <a href="#" class="inline-block py-1 md:py-4 text-gray-100 mr-6 hover:text-gray-500 font-inline">Features</a>
            <a href="#" class="inline-block py-1 md:py-4 text-gray-100 font-inline hover:text-gray-500 mr-6">Security</a>
            <a href="#" class="inline-block py-1 md:py-4 text-gray-100 font-inline hover:text-gray-500">Technologies</a>
          </div>
        </div>
        <div class="hidden md:block">
          <a href="#" class="inline-block py-1 md:py-4 text-gray-100 hover:text-gray-500 mr-6"></a>
          <a href="#" class="text-lg py-3 px-12 border module-border-wrap text-gray-100 hover:bg-gray-100 hover:text-black rounded-xl">Get Started</a>
        </div>
      </div>
    </div>

    <div class="md:overflow-hidden">
      <div class="px-4 py-20 md:py-4">
        <div class="md:max-w-6xl md:mx-auto">
          <div class="md:flex md:flex-wrap">
            <div class="md:w-1/2 text-center md:text-left md:pt-16">
              <h1 class="font-bold text-white text-2xl md:text-5xl leading-tight mb-4">
                Secure Verified Credential Storage
              </h1>
              <md-card-content v-if="loading" style="margin: 10% 20% 10% 30%">
                <beat-loader :color="'white'" :size="20"></beat-loader>
                <transition name="fade" mode="out-in">
                  <div class="text-white" style="padding-top: 10px;" :key="messageNum">{{messages[messageNum]}}</div>
                </transition>
              </md-card-content>
              <md-card-content v-else>
                <form>
                    <md-button v-on:click="beginOIDCLogin"
                               class="text-lg py-3 px-12 border module-border-wrap text-gray-100 hover:bg-gray-100 hover:text-black rounded-xl"
                               id="loginBtn">
                               Get Started
                    </md-button>
                    <!-- Todo remove this once the design is finalized -->
                    <md-button v-if="registered" v-on:click="loginDevice"
                               class="md-dense md-raised md-success login-button" id="loginDeviceBtn">
                      Sign-In Touch/Face ID
                    </md-button>
                </form>
              </md-card-content>
            </div>
            <div class="md:px-24 relative">
              <div class="hidden md:block">
                <img class="object-scale-down h-84 w-full" src="@/assets/img/home_wallet.png" alt="">
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- component -->
      <footer class="footer footer-gradient relative py-24  border-b-2 border-blue-700">
          <div class="container mx-auto px-6 border-gray-300 flex flex-col items-center">
            <div class="sm:w-2/3 py-16">
              <p class="text-sm text-white text-center font-semibold  mb-2">
                Copyright &copy; <a href="https://securekey.com/">SecureKey Technologies</a> and the TrustBloc Contributors.
              </p>
            </div>
          </div>
      </footer>
    </div>
  </div>
</template>

<script>
    import {DeviceLogin, RegisterWallet} from "./wallet"
    import {mapActions, mapGetters} from 'vuex'
    import {BeatLoader} from "@saeris/vue-spinners";

    export default {
        created: async function () {
            this.startLoading()
            //TODO: issue-601 Implement cookie logic with information from the backend.
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

            this.stopLoading()
            this.loading = false
        },
        data() {
            return {
                statusMsg: '',
                loading: true,
                registered: false,
                messageNum: 0,
                messages: [
                    "Getting Started..",
                    "This could take a minute.",
                    "Please do not refresh the page.",
                    "Please wait...",
                ]
            };
        },
        components: {
            BeatLoader,
        },
        methods: {
            ...mapActions({
                loadUser: 'loadUser',
                loadOIDCUser: 'loadOIDCUser',
                refreshUserMetadata: 'refreshUserMetadata',
                startUserSetup: 'startUserSetup',
                completeUserSetup: 'completeUserSetup'
            }),
            ...mapGetters(['getCurrentUser', 'getAgentOpts', 'serverURL']),
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            beginOIDCLogin: async function () {
                window.location.href = this.serverURL() + "/oidc/login"
            },

            sleep(ms) {
                return new Promise(resolve => setTimeout(resolve, ms));
            },
            finishOIDCLogin: async function () {
                let user = this.getCurrentUser()

                let registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, this.getAgentInstance(),
                    this.getAgentOpts())

                try {
                    if (!user.metadata) {
                        this.startUserSetup()

                        // first time login, register this user
                        registrar.register(user.username, this.completeUserSetup)
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
            startLoading() {
                this.intervalID = setInterval(() => {
                    this.messageNum++;
                    this.messageNum %= this.messages.length;
                }, 3000);
            },
            stopLoading() {
                clearInterval(this.intervalID);
            },

        }
    }

</script>
<style scoped>
    .login-button {
        text-transform: none !important; /*For Lower case use lowercase*/
        font-size: 16px !important;
        width: 100%;
    }
     .gradient {
       background: linear-gradient(to bottom ,#14061D, #261131,#0c0116,#13113F,#1A0C22);
     }
    .footer-gradient {
      background: linear-gradient(to bottom ,#000000, #0c0116,#14061D, #261131);
    }
    .module-border-wrap {
      border: 1px solid;
      border-image-slice: 1;
      border-image-source: linear-gradient(to left,  #3F5FD3, #743ad5, #d53a9d, #CD3A67);
    }
    a {
      color: white !important;
    }
    /*--Remove this once vue-material css is removed */
    .md-button{
      text-transform: none !important;
      background:linear-gradient(to bottom ,#14061D, #261131,#0c0116,#13113F,#1A0C22);
      font-size: large;
      font-family: sans-serif;
      padding: 6px 36px;
    }
</style>
