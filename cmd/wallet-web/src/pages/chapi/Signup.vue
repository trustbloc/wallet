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
          </div>
        </div>
      </div>
    </div>
    <md-card-content v-if="loading" style="margin: 10% 20% 10% 30%">
      <beat-loader :color="'white'" :size="20"></beat-loader>
      <transition name="fade" mode="out-in">
        <div class="text-white" style="padding-top: 10px;" :key="messageNum">{{messages[messageNum]}}</div>
      </transition>
    </md-card-content>
    <div class="container mx-auto h-full" v-else>
      <div class="w-full lg:w-4/12">
        <div class="relative w-full mb-6 shadow-lg rounded-lg border-0">
          <div class="flex-auto px-4 lg:px-10 py-10 pt-0">
            <div class="flex items-center justify-center m-4" >
              <md-button class="text-lg py-3 px-12 border module-border-wrap md-button
                   " v-for="(provider, index) in providers" :key="index"
                 v-on:click="beginOIDCLogin(provider.id)">
                <img class="object-scale-down h-8 w-16 text-black"
                     :src="provider.logoURL"/>
                {{ provider.name }}
              </md-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {DeviceLogin, RegisterWallet} from "./wallet"
import {mapActions, mapGetters} from 'vuex'
import {BeatLoader} from "@saeris/vue-spinners";
import axios from 'axios';

export default {
  created: async function () {
    await this.fetchProviders()
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
  components: {
    BeatLoader,
  },
  data() {
    return {
      providers:[],
      hubOauthProvider: this.hubURL(),
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
  methods: {
    ...mapActions({
      loadUser: 'loadUser',
      loadOIDCUser: 'loadOIDCUser',
      refreshUserMetadata: 'refreshUserMetadata',
      startUserSetup: 'startUserSetup',
      completeUserSetup: 'completeUserSetup'
    }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'serverURL', 'hubURL']),
    ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
    beginOIDCLogin: async function (providerID) {
      window.location.href = this.serverURL() + "/oidc/login?provider=" + providerID
    },

    // Fetching the providers from hub-auth
    fetchProviders: async function() {
      await axios.get(this.hubURL() +'/oauth2/providers')
          .then(response => {
            this.providers = response.data.authProviders
          })
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
  background:transparent !important;
  font-size: large;
  font-family: sans-serif;
  padding: 6px 36px;
}
.md-button:hover{
  background:transparent !important;
}
.md-button:active:after{
  background:transparent !important;
}
</style>
