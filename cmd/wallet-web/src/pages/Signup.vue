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
      md:bg-onboarding
      flex flex-col
      min-w-screen min-h-screen
      bg-scroll bg-no-repeat bg-neutrals-softWhite bg-onboarding-sm
    "
  >
    <div class="flex flex-col flex-grow justify-center items-center">
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
            bg-no-repeat bg-flare
            divide-x divide-neutrals-medium divide-opacity-25
          "
        >
          <div class="hidden md:block col-span-1 py-24 pr-16">
            <Logo class="mb-12" href="" />

            <div class="flex overflow-y-auto flex-1 items-center mb-8 max-w-full">
              <img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-1.svg" />
              <span class="pl-5 text-base text-neutrals-white align-middle">
                {{ i18n.leftContainer.span1 }}
              </span>
            </div>

            <div class="flex overflow-y-auto flex-1 items-center mb-8 max-w-full">
              <img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-2.svg" />
              <span class="pl-5 text-base text-neutrals-white align-middle">
                {{ i18n.leftContainer.span2 }}
              </span>
            </div>

            <div class="flex overflow-y-auto flex-1 items-center max-w-full">
              <img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-3.svg" />
              <span class="pl-5 text-base text-neutrals-white align-middle">
                {{ i18n.leftContainer.span3 }}
              </span>
            </div>
          </div>
          <!--Trustbloc Sign-up provider div -->
          <div class="md:block object-none object-center col-span-1">
            <div class="px-6 md:pt-16 md:pr-0 md:pb-12 md:pl-16">
              <Logo class="md:hidden justify-center my-2 mt-12" href="" />
              <div class="items-center mb-10 text-center">
                <h1 class="text-2xl md:text-4xl font-bold text-neutrals-white">
                  {{ i18n.heading }}
                </h1>
              </div>
              <div class="flex justify-center content-center py-24 w-full min-h-xl">
                <Spinner v-if="loading" />
                <button
                  v-for="(provider, index) in providers"
                  v-else
                  :key="index"
                  class="
                    items-center
                    py-2
                    px-4
                    mb-4
                    w-full
                    max-w-xs
                    h-11
                    text-sm
                    font-bold
                    text-neutrals-dark
                    bg-neutrals-softWhite
                    rounded-md
                    flex flex-wrap
                  "
                  @click="beginOIDCLogin(provider.id)"
                >
                  <img class="inline-block object-contain mr-2 max-h-6" :src="provider.logoURL" />
                  <span id="signUpText" class="flex flex-wrap">{{ provider.signUpText }}</span>
                </button>
              </div>
              <div class="py-10 md:pt-12 md:pb-0 text-center">
                <p class="text-base font-normal text-neutrals-white">
                  {{ i18n.redirect }}
                  <a class="text-primary-blue whitespace-nowrap underline-blue" href="/signin">{{
                    i18n.signin
                  }}</a>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <ContentFooter />
  </div>
</template>

<script>
import { CHAPIHandler, RegisterWallet } from './mixins';
import { DeviceLogin } from '@trustbloc/wallet-sdk';
import ContentFooter from '@/pages/layout/ContentFooter.vue';
import Logo from '@/components/Logo/Logo.vue';
import Spinner from '@/components/Spinner/Spinner.vue';

import { mapActions, mapGetters } from 'vuex';
import axios from 'axios';

export default {
  components: {
    ContentFooter,
    Logo,
    Spinner,
  },
  data() {
    return {
      providers: [],
      statusMsg: '',
      loading: true,
    };
  },
  computed: {
    i18n() {
      return this.$t('Signup');
    },
    isLoggedIn() {
      return this.isUserLoggedIn();
    },
  },
  watch: {
    async isLoggedIn() {
      // watch for use login state and proceed with load OIDC user step.
      if (this.isUserLoggedIn()) {
        await this.refreshOpts();
        await this.loadOIDCUser();

        if (this.getCurrentUser()) {
          await this.finishOIDCLogin();
          this.handleSuccess();
        }
      }
    },
  },
  created: async function () {
    await this.fetchProviders();
    //TODO: issue-601 Implement cookie logic with information from the backend.
    this.deviceLogin = new DeviceLogin(this.getAgentOpts()['edge-agent-server']);

    // user intended to destination
    let redirect = this.$route.params['redirect'];
    this.redirect = redirect ? { name: redirect } : `${__webpack_public_path__}dashboard`;

    console.debug('redirecting to', this.redirect);

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
    }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'serverURL', 'hubAuthURL', 'isUserLoggedIn']),
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
      });
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
      await chapi.install(this.getCurrentUser().username);
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
