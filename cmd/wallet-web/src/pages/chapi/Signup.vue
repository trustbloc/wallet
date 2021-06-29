<template>
  <div class="bg-gray-light min-w-screen min-h-screen flex items-center justify-center px-5 py-5
              bg-scroll xl:bg-onboarding sm:bg-onboarding md:bg-onboarding bg-no-repeat">
    <div class="bg-gradient-dark rounded-xl overflow-hidden text-neutrals-softWhite h-xl w-xl">
      <!--Trustbloc Intro div  -->
      <div class="grid grid-cols-2 bg-flare bg-no-repeat w-auto h-auto lg:divide-x divide-neutrals-medium divide-opacity-25">
        <div class="col-span-1 hidden md:block py-14 px-8">
          <ul class="list-none md:list-disc text-sm px-4 py-4">
            <div class="flex items-center">
              <img class="h-2 w-8 mr-2 mt-4 md-4 my-2 mx-4" src="@/assets/img/logo.png">
              <h2 class="font-semibold text-xl lg:text-3xl">TrustBloc</h2>
            </div>
            <ul class="w-full rounded-lg mt-2 mb-3 font-normal">
              <li class="pl-3 pr-4 py-3 text-sm">
                <div class="flex-1 flex items-center overflow-y-auto max-w-full">
                  <img class="flex w-8 h-full" src="@/assets/img/onboarding-icon-1.svg" />
                  <span class="text-sm px-6 text-neutrals-softWhite align-middle">Keep your digital identity safe</span>
                </div>
              </li>
              <li class="pl-3 pr-4 py-3 text-sm">
                <div class="flex-1 flex items-center overflow-y-auto max-w-full">
                  <img class="flex w-8 h-full" src="@/assets/img/onboarding-icon-2.svg" />
                  <span class="text-sm px-6 text-neutrals-softWhite align-middle">Store digital IDs, certifications, and
                    more—all in one secure wallet</span>
                </div>
              </li>
              <li class="pl-3 pr-4 py-3 text-sm">
                <div class="flex-1 flex items-center overflow-y-auto max-w-full">
                  <img class="flex w-8 h-full" src="@/assets/img/onboarding-icon-3.svg" />
                  <span class="text-sm px-6 text-neutrals-softWhite align-middle">Verify your identity in person or online</span>
                </div>
              </li>
            </ul>
          </ul>
        </div>
        <!--Trustbloc Sign-up provider div -->
        <div class="col-span-1 md:block py-14 px-4">
          <div class="text-center mb-10 items-center">
            <h1 class="font-semibold text-2xl text-neutrals-softWhite">Sign up. It's free
            </h1>
          </div>
            <md-card-content v-if="loading" style="margin: 10% 20% 10% 30%">
              <beat-loader :color="'white'" :size="20"></beat-loader>
              <transition name="fade" mode="out-in">
                <div class="text-white" style="padding-top: 10px;" :key="messageNum">{{messages[messageNum]}}</div>
              </transition>
            </md-card-content>
            <div class="flex mx-8">
              <div class="w-full px-8 mb-8">
                <button class="text-sm items-center flex py-2 px-6
                        font-bold bg-neutrals-softWhite text-neutrals-dark rounded-sm align-center"
                    v-for="(provider, index) in providers" :key="index"
                           v-on:click="beginOIDCLogin(provider.id)">
                  <img class="object-contain md:object-scale-down max-h-4 max-w-2 mr-2"
                       :src="provider.logoURL"/>
                  Continue with {{ provider.name }}
                </button>
              </div>
            </div>
        <div class="text-center mb-10">
          <p class="font-normal text-neutrals-softWhite">Already have an account?
            <a class="text-primary-blue underline" href=""> Sign in </a>
          </p>
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
