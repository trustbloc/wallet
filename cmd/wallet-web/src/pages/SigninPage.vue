<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useStore } from 'vuex';
import axios from 'axios';
import { useI18n } from 'vue-i18n';
import { CHAPIHandler, RegisterWallet } from '@/mixins';
import { DeviceLogin } from '@trustbloc/wallet-sdk';
import useBreakpoints from '@/plugins/breakpoints.js';
import FooterComponent from '@/components/Footer/FooterComponent.vue';
import ToastNotificationComponent from '@/components/ToastNotification/ToastNotificationComponent.vue';
import LogoComponent from '@/components/Logo/LogoComponent.vue';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

// Local Variables
const loading = ref(true);
const providers = ref([]);
const systemError = ref(false);
const providerPopup = ref({ closed: false });
const disableCHAPI = ref(false);
const deviceLogin = ref();
const redirect = ref('');

// Hooks
const router = useRouter();
const route = useRoute();
const store = useStore();
const { t, locale } = useI18n();
const breakpoints = useBreakpoints();
const polyfill = inject('polyfill');
const webCredentialHandler = inject('webCredentialHandler');

// Store Getters
const currentUser = computed(() => store.getters['getCurrentUser']);
const agentOpts = computed(() => store.getters['getAgentOpts']);
const agentInstance = computed(() => store.getters['agent/getInstance']);
const hubAuthURL = computed(() => store.getters['hubAuthURL']);
const isUserLoggedIn = computed(() => store.getters['isUserLoggedIn']);
const isLoginSuspended = computed(() => store.getters['isLoginSuspended']);

// Store Actions
const loadUser = () => store.dispatch('loadUser');
const loadOIDCUser = () => store.dispatch('loadOIDCUser');
const startUserSetup = () => store.dispatch('startUserSetup');
const completeUserSetup = () => store.dispatch('completeUserSetup');
const refreshUserPreference = () => store.dispatch('refreshUserPreference');
const refreshOpts = () => store.dispatch('initOpts');
const activateCHAPI = () => store.dispatch('activateCHAPI');

// Watchers
watch(
  () => isUserLoggedIn.value,
  async (isUserLoggedIn) => {
    if (isUserLoggedIn) {
      await refreshOpts();
      try {
        await loadOIDCUser();
      } catch (e) {
        systemError.value = true;
        loading.value = false;
      }
      if (currentUser.value) {
        await finishOIDCLogin();
        handleSuccess();
      }
    }
  }
);

watch(
  () => isLoginSuspended.value,
  () => {
    loading.value = false;
  }
);

// Methods
function openProviderPopup(url, title, w, h) {
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
}

function initiateOIDCLogin(providerID) {
  loading.value = true;
  providerPopup.value = openProviderPopup('/loginhandle?providerID=' + providerID, '', 700, 770);
}

async function finishOIDCLogin() {
  await registerUser();
  if (!breakpoints.xs && !breakpoints.sm && !disableCHAPI.value) {
    // all credential handlers registration should happen here, ex: CHAPI etc
    const chapi = new CHAPIHandler(
      polyfill,
      webCredentialHandler,
      agentOpts.value.credentialMediatorURL
    );

    await chapi.install(currentUser.value.username);
    activateCHAPI();
  }
}
async function registerUser() {
  if (!currentUser.value.preference) {
    startUserSetup();
    // first time login, register this user
    await new RegisterWallet(agentInstance.value, agentOpts.value).register(
      {
        name: currentUser.value.username,
        user: currentUser.value.profile.user,
        token: currentUser.value.profile.token,
      },
      completeUserSetup
    );
    refreshUserPreference();
  }
}

function handleSuccess() {
  router.push(redirect.value);
}

onMounted(async () => {
  try {
    const rawProviders = await axios.get(hubAuthURL.value + '/oauth2/providers');
    providers.value = rawProviders.data.authProviders.sort(
      (prov1, prov2) => prov1.order - prov2.order
    );
    loading.value = false;
  } catch (e) {
    systemError.value = true;
    console.error('failed to fetch providers', e);
  }
  // TODO: issue-601 Implement cookie logic with information from the backend.
  deviceLogin.value = new DeviceLogin(agentOpts.value['edge-agent-server']);

  // user intended to destination
  redirect.value = route.params['redirect']
    ? {
        name: route.params['redirect'],
        params: { locale: store.getters.getLocale.base },
        query: route.query,
      }
    : {
        name: 'vaults',
        params: { locale: store.getters.getLocale.base },
        query: route.query,
      };

  console.debug('redirecting to', redirect.value);

  // if intended target doesn't require CHAPI.
  disableCHAPI.value = route.params.disableCHAPI;

  // load user.
  loadUser();

  // if session found, then no need to login.
  if (currentUser.value) {
    handleSuccess();
    return;
  }

  // show default view with signup options.
  loading.value = false;
});
</script>

<template>
  <div
    class="flex flex-col justify-between items-center px-6 min-h-screen bg-scroll bg-neutrals-softWhite bg-no-repeat bg-onboarding-sm md:bg-onboarding"
  >
    <div class="flex flex-col flex-grow justify-center items-center">
      <ToastNotificationComponent
        v-if="systemError"
        :title="t('Signin.errorToast.title')"
        :description="t('Signin.errorToast.description')"
        type="error"
      />
      <div
        class="flex overflow-hidden flex-col justify-start items-center px-6 mx-6 w-full sm:w-screen max-w-xl h-auto text-xl md:text-3xl bg-gradient-dark rounded-xl"
      >
        <LogoComponent class="py-12" />
        <div class="items-center mb-10 md:mb-8 text-center">
          <span class="text-2xl md:text-4xl font-bold text-neutrals-white">
            {{ t('Signin.heading') }}
          </span>
        </div>
        <div
          class="grid grid-cols-1 gap-5 justify-items-center content-center sm:px-32 mb-12 w-full h-64"
        >
          <SpinnerIcon v-if="loading" />
          <button
            v-for="(provider, index) in providers"
            v-else
            :id="provider.id"
            :key="index"
            class="flex flex-wrap items-center w-full h-11 text-sm font-bold text-neutrals-dark bg-neutrals-softWhite rounded-md"
            @click="initiateOIDCLogin(provider.id)"
            @keyup.enter="initiateOIDCLogin(provider.id)"
          >
            <img :src="provider.signInIconUrl[locale]" />
          </button>
        </div>
        <div class="mb-12 text-center">
          <p class="text-base font-normal text-neutrals-softWhite">
            {{ t('Signin.redirect') }}
            <router-link
              class="text-primary-blue whitespace-nowrap underline-blue"
              :to="{ name: 'signup' }"
              >{{ t('Signin.signup') }}
            </router-link>
          </p>
        </div>
      </div>
    </div>
    <FooterComponent />
  </div>
</template>
