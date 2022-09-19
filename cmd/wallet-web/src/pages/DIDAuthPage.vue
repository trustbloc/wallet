<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div v-if="loading">
    <SpinnerIcon />
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
        <!-- TODO: issue-1055 Read meta data from external urls -->
        <div class="p-5 text-sm text-neutrals-dark">Issuer wants to connect to your wallet.</div>

        <hr class="mx-5 border border-neutrals-thistle" />
        <div class="py-6 px-5">
          <div class="z-10 flex w-full flex-row items-center justify-start p-5">
            <div class="h-12 w-12 flex-none border-opacity-10">
              <!-- TODO: issue-1055 Read meta data from external urls -->
              <img src="@/assets/img/generic-issuer-icon.svg" />
            </div>
            <div class="flex flex-col">
              <span
                class="flex-1 text-ellipsis pl-4 text-left text-sm font-bold text-neutrals-dark"
              >
                <!-- TODO: issue-1055 Read meta data from external urls -->
                Issuer
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
                TrustBloc Wallet
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
          <button id="didauth" class="btn-primary" @click="authorize">Connect</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { toRaw } from 'vue';
import { CredentialManager } from '@trustbloc/wallet-sdk';
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
      issuers: [{ id: 0, name: 'Select Identity' }],
      errors: [],
      requestOrigin: '',
      loading: true,
      govnVC: null,
    };
  },
  computed: {
    walletUrl() {
      return window.location.origin;
    },
    walletIcon() {
      return document.querySelector("link[rel~='icon']").href;
    },
  },
  created: async function () {
    this.protocolHandler = this.$parent.protocolHandler;
    const eventData = this.protocolHandler.getEventData();
    const query = toRaw(eventData.query);

    const { user, token } = this.getCurrentUser().profile;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });

    try {
      let { results } = await this.credentialManager.query(
        token,
        Array.isArray(query) ? query : [query]
      );
      this.presentation = results[0];
    } catch (e) {
      console.error('failed to prepare DIDAuth response:', e);
      // TODO: https://github.com/trustbloc/wallet/issues/1067 DIDAuth query failure
      // this.errors.push('failed to handle request, try again later.');
      this.loading = false;
    }

    this.requestOrigin = this.protocolHandler.requestor();
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    cancel: function () {
      this.protocolHandler.cancel();
    },
    authorize: async function () {
      this.loading = true;

      let { profile, preference } = this.getCurrentUser();
      let { controller, proofType, verificationMethod } = preference;
      let { domain, challenge } = toRaw(this.protocolHandler.getEventData());

      let { presentation } = await this.credentialManager.present(
        profile.token,
        { presentation: this.presentation },
        {
          controller,
          proofType,
          domain,
          challenge,
          verificationMethod,
        }
      );

      this.protocolHandler.present(presentation);

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
