<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div v-if="loading" class="chapi-container">
    <SpinnerIcon v-if="loading" />
  </div>

  <div v-else class="flex w-screen justify-center">
    <div>
      <HeaderComponent />

      <div class="chapi-container h-auto rounded-b border border-neutrals-black bg-gray-light">
        <div v-if="errors.length">
          <b>Failed with following error(s):</b>
          <ul>
            <li v-for="error in errors" :key="error" class="text-primary-valencia">{{ error }}</li>
          </ul>
        </div>
        <div class="p-5 text-sm text-neutrals-dark">
          {{ issuerName }} wants to connect to your wallet.
        </div>

        <hr class="mx-5 border border-neutrals-thistle" />
        <div class="py-6 px-5">
          <div class="z-10 flex w-full flex-row items-center justify-start p-5">
            <div class="h-12 w-12 flex-none border-opacity-10">
              <!-- todo issue-1055 Read meta data from external urls -->
              <img src="@/assets/img/generic-issuer-icon.svg" />
            </div>
            <div class="flex flex-col">
              <span
                class="flex-1 text-ellipsis pl-4 text-left text-sm font-bold text-neutrals-dark"
              >
                {{ issuerName }}
              </span>
              <div class="flex flex-row items-center justify-center pl-4">
                <img src="@/assets/img/small-lock-icon.svg" />
                <span class="flex-1 text-ellipsis pl-1 text-left text-xs text-neutrals-medium">
                  {{ requestOrigin }}
                </span>
              </div>
            </div>
          </div>
          <!-- Share icon -->
          <div class="flex h-14 w-14 flex-col border-opacity-10 pl-6">
            <img src="@/assets/img/share.svg" />
          </div>

          <!-- Wallet to connect-->
          <div class="z-10 flex w-full flex-row items-center justify-start px-5">
            <div class="h-12 w-12 flex-none border-opacity-10">
              <img :src="walletIcon" />
            </div>
            <div class="flex flex-col">
              <span
                class="flex-1 text-ellipsis pl-4 text-left text-sm font-bold text-neutrals-dark"
              >
                <!-- todo still have to finalize the text of the title-->
                TrustBloc Digital Identity Wallet'
              </span>
              <div class="flex flex-row items-center justify-center pl-4">
                <img src="@/assets/img/small-lock-icon.svg" />
                <span class="flex-1 text-ellipsis pl-1 text-left text-xs text-neutrals-medium">
                  {{ walletUrl }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex h-auto w-full justify-between bg-neutrals-magnolia py-4 px-5">
          <button id="cancelBtn" class="btn-outline" @click="cancel">Decline</button>
          <button id="didconnect" class="btn-primary" @click="connect">
            {{ buttonLabel }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { DIDConn } from '@/mixins';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';
import HeaderComponent from '@/components/Header/HeaderComponent.vue';
import { mapGetters } from 'vuex';

export default {
  components: {
    SpinnerIcon,
    HeaderComponent,
  },
  data() {
    return {
      errors: [],
      requestOrigin: '',
      loading: true,
      userCredentials: [],
      buttonLabel: 'Connect',
      govnVC: null,
      walletTitle: '',
    };
  },
  computed: {
    walletUrl() {
      return window.location.origin;
    },
    walletIcon() {
      return document.querySelector("link[rel~='icon']").href;
    },
    issuerName() {
      // TODO: issue-1055 read from external meta data
      return 'Issuer';
    },
  },
  created: function () {
    this.wallet = new DIDConn(
      this.getAgentInstance(),
      this.getCurrentUser().profile,
      this.getAgentOpts(),
      this.$parent.protocolHandler,
      this.getCredentialManifests()
    );

    this.requestOrigin = this.$parent.protocolHandler.requestor();
    this.userCredentials = this.wallet.userCredentials;
    this.govnVC = this.wallet.govnVC;
    this.buttonLabel = this.userCredentials.length > 0 ? 'Store & Connect' : 'Connect';

    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'getCredentialManifests']),
    cancel: function () {
      this.wallet.cancel();
    },
    connect: async function () {
      this.loading = true;
      this.errors = [];
      try {
        await this.wallet.connect(this.getCurrentUser().preference);
      } catch (e) {
        console.error('failed to connect', e);
        this.errors.push('failed to perform connection, please try again later');
      }

      this.loading = false;
    },
  },
};
</script>
<style scoped>
.chapi-container {
  width: 28rem;
}
.description {
  height: 25.5rem;
  width: 2.6rem;
}
</style>
