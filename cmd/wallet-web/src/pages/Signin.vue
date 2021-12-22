<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="
      justify-between
      items-center
      px-6
      min-h-screen
      md:bg-onboarding
      flex flex-col
      bg-scroll bg-no-repeat bg-neutrals-softWhite bg-onboarding-sm
    "
  >
    <div class="flex flex-col flex-grow justify-center items-center">
      <!--Trustbloc Sign-up provider div -->
      <toast-notification
        v-if="systemError"
        :title="t('Signin.errorToast.title')"
        :description="t('Signin.errorToast.description')"
        type="error"
      ></toast-notification>
      <div
        class="
          overflow-hidden
          justify-start
          items-center
          px-6
          mx-6
          mt-20
          w-full
          sm:w-screen
          max-w-xl
          h-auto
          text-xl
          md:text-3xl
          bg-gradient-dark
          rounded-xl
          flex flex-col
        "
      >
        <Logo class="py-12" />
        <div class="items-center mb-10 md:mb-8 text-center">
          <span class="text-2xl md:text-4xl font-bold text-neutrals-white">
            {{ t('Signin.heading') }}
          </span>
        </div>
        <div
          class="
            grid grid-cols-1
            gap-5
            w-full
            sm:px-32
            h-64
            mb-12
            content-center
            justify-items-center
          "
        >
          <Spinner v-if="loading" />
          <button
            v-for="(provider, index) in providers"
            v-else
            :key="index"
            class="
              items-center
              w-full
              h-11
              text-sm
              font-bold
              text-neutrals-dark
              bg-neutrals-softWhite
              rounded-md
              flex flex-wrap
            "
            @click="beginOIDCLogin(provider.id)"
            @keyup.enter="beginOIDCLogin(provider.id)"
          >
            <img :id="provider.id" :src="provider.signInLogoUrl" />
          </button>
        </div>
        <div class="text-center mb-8">
          <p class="text-base font-normal text-neutrals-softWhite">
            {{ t('Signin.redirect') }}
            <router-link class="text-primary-blue whitespace-nowrap underline-blue" to="signup"
              >{{ t('Signin.signup') }}
            </router-link>
          </p>
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script>
import { CHAPIHandler, RegisterWallet } from '@/utils/mixins';
import { DeviceLogin } from '@trustbloc/wallet-sdk';
import Footer from '@/components/Footer/Footer.vue';
import Logo from '@/components/Logo/Logo.vue';
import Spinner from '@/components/Spinner/Spinner.vue';
import useBreakpoints from '@/plugins/breakpoints.js';
import { mapActions, mapGetters } from 'vuex';
import axios from 'axios';
import { useI18n } from 'vue-i18n';

export default {
  components: {
    Footer,
    Logo,
    Spinner,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      providers: [],
      statusMsg: '',
      loading: true,
      systemError: false,
      breakpoints: useBreakpoints(),
    };
  },
  computed: {
    isLoggedIn() {
      return this.isUserLoggedIn();
    },
    isSuspended() {
      return this.isLoginSuspended();
    },
  },
  watch: {
    isLoggedIn: {
      async handler() {
        // watch for use login state and proceed with load OIDC user step.
        if (this.isLoggedIn) {
          await this.refreshOpts();
          try {
            await this.loadOIDCUser();
          } catch (e) {
            this.systemError = true;
            this.loading = false;
          }
          if (this.getCurrentUser()) {
            await this.finishOIDCLogin();
            this.handleSuccess();
          }
        }
      },
    },
    isSuspended() {
      this.loading = false;
    },
  },
  created: async function () {
    await this.fetchProviders();
    //TODO: issue-601 Implement cookie logic with information from the backend.
    this.deviceLogin = new DeviceLogin(this.getAgentOpts()['edge-agent-server']);

    // user intended to destination
    const redirect = this.$route.params['redirect'];
    this.redirect = redirect
      ? {
          name: redirect,
          params: { locale: this.$store.getters.getLocale.base },
          query: this.$route.query,
        }
      : {
          name: 'vaults',
          params: { locale: this.$store.getters.getLocale.base },
          query: this.$route.query,
        };

    console.debug('redirecting to', this.redirect);

    // if intended target doesn't require CHAPI.
    this.disableCHAPI = this.$route.params.disableCHAPI;

    // load user.
    this.loadUser();

    // if session found, then no need to login.
    if (this.getCurrentUser()) {
      this.handleSuccess();
      return;
    }

    // show default view with signup options.
    this.loading = false;
  },
  methods: {
    ...mapActions({
      loadUser: 'loadUser',
      loadOIDCUser: 'loadOIDCUser',
      startUserSetup: 'startUserSetup',
      completeUserSetup: 'completeUserSetup',
      refreshUserPreference: 'refreshUserPreference',
      refreshOpts: 'initOpts',
      activateCHAPI: 'activateCHAPI',
    }),
    ...mapGetters([
      'getCurrentUser',
      'getAgentOpts',
      'serverURL',
      'hubAuthURL',
      'isUserLoggedIn',
      'isLoginSuspended',
    ]),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    beginOIDCLogin: function (providerID) {
      this.loading = true;
      this.popupwindow('/loginhandle?provider=' + providerID, '', 700, 770);
    },
    popupwindow(url, title, w, h) {
      var left = screen.width / 2 - w / 2;
      var top = screen.height / 2 - h / 2;
      return window.open(
        url,
        title,
        'menubar=yes,status=yes, replace=true, width=' +
          w +
          ', height=' +
          h +
          ', top=' +
          top +
          ', left=' +
          left
      );
    },
    // Fetching the providers from hub-auth
    fetchProviders: async function () {
      await axios.get(this.hubAuthURL() + '/oauth2/providers').then((response) => {
        this.providers = response.data.authProviders;
        // Sort the list of the providers based on the order property in the object.
        this.providers.sort((prov1, prov2) => prov1.order - prov2.order);
      });
    },
    finishOIDCLogin: async function () {
      let user = this.getCurrentUser();
      this.registerUser(user);

      if (!this.breakpoints.xs && !this.breakpoints.sm && !this.disableCHAPI) {
        // all credential handlers registration should happen here, ex: CHAPI etc
        let chapi = new CHAPIHandler(
          this.$polyfill,
          this.$webCredentialHandler,
          this.getAgentOpts().credentialMediatorURL
        );

        await chapi.install(this.getCurrentUser().username);
        this.activateCHAPI();
      }
    },
    async registerUser(user) {
      if (!user.preference) {
        this.startUserSetup();

        let user = this.getCurrentUser();
        // first time login, register this user
        await new RegisterWallet(this.getAgentInstance(), this.getAgentOpts()).register(
          {
            name: user.username,
            user: user.profile.user,
            token: user.profile.token,
          },
          this.completeUserSetup
        );
        this.refreshUserPreference();
      }
    },
    handleSuccess() {
      this.$router.push(this.redirect);
    },
  },
};
</script>
