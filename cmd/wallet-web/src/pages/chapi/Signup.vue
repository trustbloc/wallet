<template>
  <!-- Todo  Move this style to theme setup -->
  <div class="bg-gray-light min-w-screen min-h-screen flex items-center justify-center px-5 py-5 bg-scroll bg-onboarding bg-no-repeat">
    <!-- This is an example component -->
    <div class="bg-gradient-dark rounded-xl w-full overflow-hidden text-neutrals-softWhite" style="max-width:896px; max-height:467px;">
      <div class="md:flex w-full bg-flare bg-no-repeat">
        <div class="hidden md:block w-1/2 py-14 px-10">
          <div class="flex items-center">
            <img class="h-2 w-8 mr-2 mt-4 md-4 my-2 mx-4" src="@/assets/img/logo.png">
            <h2 class="font-semibold text-xl lg:text-3xl">TrustBloc</h2>
          </div>
          <ul class="list-none md:list-disc text-sm px-8 py-4">
            <li class="py-2">Keep your digital identity safe</li>
            <li class="py-2">Store digital IDs, certifications, and more—all in one secure wallet</li>
            <li class="py-2">Verify your identity in person or online</li>
          </ul>
        </div>
        <div class="w-full md:w-1/2 py-10 px-5 md:px-10">
          <div class="text-center mb-10">
            <h1 class="font-semibold text-2xl">Sign up. It's free</h1>
          </div>
          <div>
            <md-card-content v-if="loading" style="margin: 10% 20% 10% 30%">
              <beat-loader :color="'white'" :size="20"></beat-loader>
              <transition name="fade" mode="out-in">
                <div class="text-white" style="padding-top: 10px;" :key="messageNum">{{messages[messageNum]}}</div>
              </transition>
            </md-card-content>
            <div class="flex -mx-3">
              <div class="w-full px-3 mb-5">
                <md-button class="text-lg py-3 px-12 border bg-white md-button
                   " v-for="(provider, index) in providers" :key="index"
                           v-on:click="beginOIDCLogin(provider.id)">
                  <img class="object-scale-down h-8 w-16 text-black"
                       :src="provider.logoURL"/>
                  {{ provider.name }}
                </md-button>
              </div>
            </div>
          </div>
          <div class="text-center mb-10">
            <p>Already have an account? Sign in</p>
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
.footer-gradient {
  background: linear-gradient(to bottom ,#000000, #0c0116,#14061D, #261131);
}
.module-border-wrap {
  border: 1px solid;
  border-image-slice: 1;
  border-image-source: linear-gradient(to left,  #3F5FD3, #743ad5, #d53a9d, #CD3A67);
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
