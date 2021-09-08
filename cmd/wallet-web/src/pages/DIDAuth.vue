<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div v-if="loading">
    <Spinner />
  </div>

  <div v-else class="flex justify-center w-screen">
    <div>
      <Header />

      <div class="h-auto bg-gray-light rounded-b border border-neutrals-black chapi-container">
        <div v-if="errors.length">
          <b>Failed with following error(s):</b>
          <ul>
            <li v-for="error in errors" :key="error" class="text-primary-valencia">{{ error }}</li>
          </ul>
        </div>
        <!-- todo issue-1055 Read meta data from external urls -->
        <div class="p-5 text-neutrals-dark font-sm">Issuer wants to connect to your wallet.</div>

        <hr class="mx-5 border border-neutrals-thistle" />
        <div class="py-6 px-5">
          <div class="flex z-10 flex-row justify-start items-center p-5 w-full">
            <div class="flex-none w-12 h-12 border-opacity-10">
              <!-- todo issue-1055 Read meta data from external urls -->
              <img src="@/assets/img/generic-issuer-icon.svg" />
            </div>
            <div class="flex flex-col">
              <span
                class="flex-1 pl-4 font-bold text-left text-neutrals-dark overflow-ellipsis font-sm"
              >
                <!-- todo issue-1055 Read meta data from external urls -->
                Issuer
              </span>
              <div class="flex flex-row justify-center items-center pl-4">
                <img src="@/assets/img/small-lock-icon.svg" />
                <span class="flex-1 pl-1 text-left text-neutrals-medium overflow-ellipsis font-xs">
                  {{ requestOrigin }}
                </span>
              </div>
            </div>
          </div>
          <!-- Share icon -->
          <div class="flex flex-col pl-6 w-14 h-14 border-opacity-10">
            <img src="@/assets/img/share.svg" />
          </div>

          <!-- Wallet to connect-->
          <div class="flex z-10 flex-row justify-start items-center px-5 w-full">
            <div class="flex-none w-12 h-12 border-opacity-10">
              <img :src="walletIcon" />
            </div>
            <div class="flex flex-col">
              <span
                class="flex-1 pl-4 font-bold text-left text-neutrals-dark overflow-ellipsis font-sm"
              >
                TrustBloc Wallet
              </span>
              <div class="flex flex-row justify-center items-center pl-4">
                <img src="@/assets/img/small-lock-icon.svg" />
                <span class="flex-1 pl-1 text-left text-neutrals-medium overflow-ellipsis font-xs">
                  {{ walletUrl }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-between py-4 px-5 w-full h-auto bg-neutrals-magnolia">
          <button id="cancelBtn" class="btn-outline" @click="cancel">Decline</button>
          <button id="didauth" class="btn-primary" @click="authorize">Connect</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import Spinner from '@/components/Spinner/Spinner.vue';
import Header from '@/components/Header/Header.vue';
import { mapGetters } from 'vuex';

export default {
  components: {
    Spinner,
    Header,
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
    let { query } = this.protocolHandler.getEventData();

    let { user, token } = this.getCurrentUser().profile;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });

    try {
      let { results } = await this.credentialManager.query(
        token,
        Array.isArray(query) ? query : [query]
      );
      this.presentation = results[0];
    } catch (e) {
      console.error('failed to prepare DIDAuth response:', e);
      this.errors.push('failed to handle request, try again later.');
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
      let { domain, challenge } = this.protocolHandler.getEventData();

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
