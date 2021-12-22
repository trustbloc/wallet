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
      <toast-notification
        v-if="systemError"
        :title="t('Signup.errorToast.title')"
        :description="t('Signup.errorToast.description')"
        type="error"
      ></toast-notification>
      <div
        class="
          overflow-hidden
          mt-20
          md:max-w-4xl
          h-auto
          text-xl
          md:text-3xl
          bg-gradient-dark
          rounded-xl
        "
      >
        <!--Trustbloc Intro div  -->
        <div
          class="
            md:grid-cols-2 md:px-20
            w-full
            h-full
            grid grid-cols-1
            bg-no-repeat bg-onboarding-flare-lg
            divide-x divide-neutrals-medium divide-opacity-25
          "
        >
          <div class="hidden md:block col-span-1 py-24 pr-16">
            <Logo class="mb-12" />

            <div class="flex overflow-y-auto flex-1 items-center mb-8 max-w-full">
              <img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-1.svg" />
              <span class="pl-5 text-base text-neutrals-white align-middle">
                {{ t('Signup.leftContainer.span1') }}
              </span>
            </div>

            <div class="flex overflow-y-auto flex-1 items-center mb-8 max-w-full">
              <img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-2.svg" />
              <span class="pl-5 text-base text-neutrals-white align-middle">
                {{ t('Signup.leftContainer.span2') }}
              </span>
            </div>

            <div class="flex overflow-y-auto flex-1 items-center max-w-full">
              <img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-3.svg" />
              <span class="pl-5 text-base text-neutrals-white align-middle">
                {{ t('Signup.leftContainer.span3') }}
              </span>
            </div>
          </div>
          <!--Trustbloc Sign-up provider div -->
          <div class="md:block object-none object-center col-span-1">
            <div class="px-6 md:pt-16 md:pr-0 md:pb-12 md:pl-16">
              <Logo class="md:hidden justify-center my-2 mt-12" />
              <div class="items-center pb-6 text-center">
                <h1 class="text-2xl md:text-4xl font-bold text-neutrals-white">
                  {{ t('Signup.heading') }}
                </h1>
              </div>
              <div
                class="grid grid-cols-1 gap-5 w-full h-64 mb-8 content-center justify-items-center"
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
                  <img :id="provider.id" :src="provider.signUpLogoUrl" />
                </button>
              </div>
              <div class="text-center mb-8">
                <p class="text-base font-normal text-neutrals-white">
                  {{ t('Signup.redirect') }}
                  <router-link
                    class="text-primary-blue whitespace-nowrap underline-blue"
                    to="signin"
                    >{{ t('Signup.signin') }}</router-link
                  >
                </p>
              </div>
            </div>
          </div>
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
        // all credential handlers registration should happen here, ex: CHAPI, WACI etc
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
