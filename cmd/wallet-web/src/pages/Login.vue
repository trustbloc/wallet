/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <div class="gradient">
    <!-- This is an example component -->
    <div class="py-24 px-4">
      <div class="md:flex md:justify-between md:items-center md:mx-auto md:max-w-6xl">
        <div class="flex justify-between items-center">
          <Logo class="" href="javascript:history.back()" />
          <div class="inline-block md:hidden cursor-pointer">
            <div class="mb-2 w-8 bg-gray-400" style="height: 2px"></div>
            <div class="mb-2 w-8 bg-gray-400" style="height: 2px"></div>
            <div class="w-8 bg-gray-400" style="height: 2px"></div>
          </div>
        </div>
        <div>
          <div class="hidden md:block text-lg">
            <a
              href="#"
              class="inline-block py-1 md:py-4 mr-6 text-gray-100 hover:text-gray-500 font-inline"
              >Features</a
            >
            <a
              href="#"
              class="inline-block py-1 md:py-4 mr-6 text-gray-100 font-inline hover:text-gray-500"
              >Security</a
            >
            <a
              href="#"
              class="inline-block py-1 md:py-4 text-gray-100 font-inline hover:text-gray-500"
              >Technologies</a
            >
          </div>
        </div>
        <div class="hidden md:block">
          <a href="#" class="inline-block py-1 md:py-4 mr-6 text-gray-100 hover:text-gray-500"></a>
          <a
            href="#"
            class="
              py-3
              px-12
              text-lg
              border
              module-border-wrap
              text-gray-100
              hover:shadow-lg hover:text-black
            "
            >Get Started</a
          >
        </div>
      </div>
    </div>

    <div class="md:overflow-hidden">
      <div class="py-20 md:py-4 px-4">
        <div class="md:mx-auto md:max-w-6xl">
          <div class="md:flex md:flex-wrap md:justify-center">
            <div class="md:pt-16 md:w-1/2 text-center md:text-left">
              <h1 class="mb-4 text-2xl md:text-5xl font-bold leading-tight text-white">
                Secure Verified Credential Storage
              </h1>
              <md-card-content v-if="loading" style="margin: 10% 20% 10% 30%">
                <beat-loader :color="'white'" :size="20"></beat-loader>
                <transition name="fade" mode="out-in">
                  <div :key="messageNum" class="text-white" style="padding-top: 10px">
                    {{ messages[messageNum] }}
                  </div>
                </transition>
              </md-card-content>
              <md-card-content v-else>
                <form>
                  <md-button
                    id="loginBtn"
                    class="py-3 px-12 text-lg border module-border-wrap md-button"
                    @click="beginOIDCLogin"
                  >
                    Get Started
                  </md-button>
                  <!-- Todo remove this once the design is finalized -->
                  <md-button
                    v-if="registered"
                    id="loginDeviceBtn"
                    class="md-dense md-raised md-success login-button"
                    @click="loginDevice"
                  >
                    Sign-In Touch/Face ID
                  </md-button>
                </form>
              </md-card-content>
            </div>
          </div>
        </div>
      </div>
      <!-- component -->
      <footer class="relative py-24 border-b-2 footer footer-gradient border-blue-700">
        <div class="container flex flex-col items-center px-6 mx-auto border-gray-300">
          <div class="py-16 sm:w-2/3">
            <p class="mb-2 text-sm text-center text-white font-light">
              Copyright &copy; <a href="https://securekey.com/">SecureKey Technologies</a> and the
              TrustBloc Contributors.
            </p>
          </div>
        </div>
      </footer>
    </div>
  </div>
</template>

<script>
import { CHAPIHandler, RegisterWallet } from './mixins';
import { DeviceLogin } from '@trustbloc/wallet-sdk';
import { mapActions, mapGetters } from 'vuex';
import { BeatLoader } from '@saeris/vue-spinners';
import Logo from '@/components/Logo/Logo.vue';

export default {
  components: {
    BeatLoader,
    Logo,
  },
  data() {
    return {
      statusMsg: '',
      loading: true,
      registered: false,
      messageNum: 0,
      messages: [
        'Getting Started..',
        'This could take a minute.',
        'Please do not refresh the page.',
        'Please wait...',
      ],
    };
  },
  created: async function () {
    this.startLoading();

    //TODO: issue-601 Implement cookie logic with information from the backend.
    this.deviceLogin = new DeviceLogin(this.getAgentOpts()['edge-agent-server']);

    let redirect = this.$route.params['redirect'];
    this.redirect = redirect ? { name: redirect } : `${__webpack_public_path__}`;
    this.loadUser();

    // you find session, no need to login
    if (this.getCurrentUser()) {
      console.log('redirecting to ', this.$route.params['redirect']);
      this.handleSuccess();
      return;
    }

    // call server to get user info and process login
    await this.loadOIDCUser();

    // register and let user into wallet
    if (this.getCurrentUser()) {
      await this.finishOIDCLogin();
      this.handleSuccess();
      return;
    }

    if (this.$cookies.isKey('registerSuccess')) {
      this.registered = true;
    }

    this.stopLoading();
    this.loading = false;
  },
  methods: {
    ...mapActions({
      loadUser: 'loadUser',
      loadOIDCUser: 'loadOIDCUser',
      startUserSetup: 'startUserSetup',
      completeUserSetup: 'completeUserSetup',
      refreshUserPreference: 'refreshUserPreference',
    }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'serverURL', 'getStaticAssetsUrl']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    beginOIDCLogin: function () {
      window.location.href = this.serverURL() + '/oidc/login';
    },

    sleep(ms) {
      return new Promise((resolve) => setTimeout(resolve, ms));
    },
    finishOIDCLogin: async function () {
      let user = this.getCurrentUser();
      this.registerUser(user);

      // all credential handlers registration should happen here, ex: CHAPI, WACI etc
      let chapi = new CHAPIHandler(
        this.$polyfill,
        this.$webCredentialHandler,
        this.getAgentOpts().credentialMediatorURL
      );
      await chapi.install(user.username);
    },
    async registerUser(user) {
      if (!user.preference) {
        this.startUserSetup();

        let user = this.getCurrentUser();
        let registrar = new RegisterWallet(this.getAgentInstance(), this.getAgentOpts());

        // first time login, register this user
        await registrar.register(
          {
            name: user.username,
            user: user.profile.user,
            token: user.profile.token,
          },
          this.completeUserSetup
        );
      }

      this.refreshUserPreference();
    },
    handleSuccess() {
      this.$router.push(this.redirect);
    },
    handleFailure(e) {
      console.error('login failure: ', e);
      this.statusMsg = e.toString();
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
  },
};
</script>
<style scoped>
.login-button {
  text-transform: none !important; /*For Lower case use lowercase*/
  font-size: 16px !important;
  width: 100%;
}

.gradient {
  background: linear-gradient(to bottom, #14061d, #261131, #0c0116, #13113f, #1a0c22);
}

.footer-gradient {
  background: linear-gradient(to bottom, #000000, #0c0116, #14061d, #261131);
}

.module-border-wrap {
  border: 1px solid;
  border-image-slice: 1;
  border-image-source: linear-gradient(to left, #3f5fd3, #743ad5, #d53a9d, #cd3a67);
}

a {
  color: white !important;
}

/*--Remove this once vue-material css is removed */
.md-button {
  text-transform: none !important;
  background: transparent !important;
  font-size: large;
  font-family: sans-serif;
  padding: 6px 36px;
}

.md-button:hover {
  background: transparent !important;
}

.md-button:active:after {
  background: transparent !important;
}
</style>
