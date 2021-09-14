<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <!-- Loading State -->
  <div v-if="loading" class="flex justify-center items-start w-screen h-screen">
    <div
      class="
        flex
        justify-center
        items-center
        w-full
        max-w-md
        h-80
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
      "
    >
      <Spinner />
    </div>
  </div>
  <div v-else-if="sharing" class="flex justify-center items-start w-screen h-screen">
    <div
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-80
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
        flex flex-col
      "
    >
      <Spinner />
      <span class="mt-8 text-base text-neutrals-dark">{{
        $t('CHAPI.Share.sharingCredential')
      }}</span>
    </div>
  </div>
  <div v-else-if="showCredentialsMissing" class="flex justify-center items-start w-screen h-screen">
    <div
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-auto
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
        flex flex-col
      "
    >
      <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
        <img src="@/assets/img/icons-warning.svg" />
        <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
          $t('CHAPI.Share.CredentialsMissing.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          $t('CHAPI.Share.CredentialsMissing.body')
        }}</span>
      </div>
      <div
        class="
          justify-center
          items-center
          pt-4
          pr-5
          pb-4
          pl-5
          w-full
          bg-neutrals-magnolia
          flex flex-row
          border-t border-neutrals-thistle
        "
      >
        <button id="share-credentials-ok-btn" class="btn-outline" @click="cancel">
          {{ $t('CHAPI.Share.CredentialsMissing.ok') }}
        </button>
      </div>
    </div>
  </div>
  <div v-else-if="errors.length" class="flex justify-center items-start w-screen h-screen">
    <div
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-auto
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
        flex flex-col
      "
    >
      <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
        <img src="@/assets/img/icons-warning.svg" />
        <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
          $t('CHAPI.Share.Error.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          $t('CHAPI.Share.Error.body')
        }}</span>
      </div>
      <div
        class="
          justify-center
          items-center
          pt-4
          pr-5
          pb-4
          pl-5
          w-full
          bg-neutrals-magnolia
          flex flex-row
          border-t border-neutrals-thistle
        "
      >
        <button id="share-credentials-ok-btn" class="btn-primary" @click="cancel">
          {{ $t('CHAPI.Share.Error.tryAgain') }}
        </button>
      </div>
    </div>
  </div>
  <!-- Main Container -->
  <div v-else class="flex justify-center items-start w-screen h-screen">
    <div class="w-full max-w-md bg-gray-light md:border md:border-t-0 border-neutrals-black">
      <div class="p-5">
        <!-- Heading -->
        <div class="flex flex-row justify-start items-start mb-4 w-full">
          <div class="flex-none w-12 h-12 border-opacity-10">
            <!-- todo issue-1055 Read meta data from external urls -->
            <img src="@/assets/img/generic-issuer-icon.svg" />
          </div>
          <div class="flex flex-col pl-3">
            <span
              class="flex-1 mb-1 text-sm font-bold text-left text-neutrals-dark overflow-ellipsis"
            >
              <!-- todo issue-1055 Read meta data from external urls -->
              Requestor
            </span>
            <div class="flex flex-row justify-center items-center">
              <img src="@/assets/img/small-lock-icon.svg" />
              <span class="flex-1 pl-1 text-xs text-left text-neutrals-medium overflow-ellipsis">
                {{ requestOrigin }}
              </span>
            </div>
          </div>
        </div>

        <!-- Errors -->
        <div v-if="errors.length" class="mb-4">
          <b>Failed with following error(s):</b>
          <ul>
            <li v-for="error in errors" :key="error" style="color: #9d0006">{{ error }}</li>
          </ul>
        </div>

        <span class="text-neutrals-dark">{{
          $tc('CHAPI.Share.headline', processedCredentials.length, { issuer: 'Requestor' })
        }}</span>

        <!-- Credentials Preview -->
        <div
          v-if="processedCredentials.length"
          class="flex flex-col justify-start items-center mt-6 mb-6 w-full"
        >
          <ul class="space-y-5 w-full">
            <li v-for="(credential, index) in processedCredentials" :key="index">
              <credential-preview
                :id="credential.id"
                :brand-color="credential.brandColor"
                :icon="credential.icon"
                :title="credential.title"
              />
            </li>
          </ul>
        </div>
      </div>

      <div
        class="
          justify-between
          items-center
          pt-4
          pr-5
          pb-4
          pl-5
          w-full
          bg-neutrals-magnolia
          flex flex-row
          border-t border-neutrals-thistle
        "
      >
        <button id="cancelBtn" class="btn-outline" @click="cancel">
          {{ $t('CHAPI.Share.decline') }}
        </button>
        <button id="share-credentials" class="btn-primary" @click="share">
          {{ $t('CHAPI.Share.share') }}
        </button>
      </div>
    </div>
  </div>
</template>
<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { normalizeQuery, getCredentialType, getCredentialDisplayData } from './mixins';
import { mapGetters } from 'vuex';
import Spinner from '@/components/Spinner/Spinner.vue';
import CredentialPreview from '@/components/CredentialPreview/CredentialPreview.vue';
import credentialDisplayData from '@/config/credentialDisplayData';

export default {
  components: {
    CredentialPreview,
    Spinner,
  },
  data() {
    return {
      errors: [],
      requestOrigin: '',
      loading: true,
      sharing: false,
      processedCredentials: [],
    };
  },
  computed: {
    showCredentialsMissing() {
      return this.processedCredentials.length === 0;
    },
  },
  created: async function () {
    this.loading = true;
    this.protocolHandler = this.$parent.protocolHandler;
    const query = normalizeQuery(this.protocolHandler.getEventData().query);
    const { user, token } = this.getCurrentUser().profile;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });

    try {
      const { results } = await this.credentialManager.query(token, query);
      this.presentation = results;
      const credentials = results.reduce((acc, val) => acc.concat(val.verifiableCredential), []);
      credentials.map((credential) => {
        const manifest = this.getManifest(credential);
        const processedCredential = this.getCredentialDisplayData(credential, manifest);
        this.processedCredentials.push(processedCredential);
      });
    } catch (e) {
      this.errors.push('No credentials found matching requested criteria.');
      console.error('get credentials failed,:', e);
    }

    // TODO: governance VC check
    this.requestOrigin = this.protocolHandler.requestor();
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    handleError(e) {
      this.errors.push(e);
      this.loading = false;
    },
    async share() {
      this.sharing = true;
      const { profile, preference } = this.getCurrentUser();
      const { controller, proofType, verificationMethod } = preference;
      const { domain, challenge } = this.protocolHandler.getEventData();

      const _present = async (presentation) => {
        return (
          await this.credentialManager.present(
            profile.token,
            { presentation },
            {
              controller,
              proofType,
              domain,
              challenge,
              verificationMethod,
            }
          )
        ).presentation;
      };
      // TODO: match expected format
      const results = await Promise.all(this.presentation.map(_present));
      // typically single presentation, but some verifier queries might produce multiple presentation.
      if (results.length === 1) {
        this.protocolHandler.present(results[0]);
      } else {
        this.protocolHandler.present(results);
      }

      this.sharing = false;
    },
    cancel() {
      this.protocolHandler.cancel();
    },
    getCredentialType: function (vc) {
      return getCredentialType(vc.type);
    },
    getCredentialDisplayData: function (vc, manifestCredential) {
      return getCredentialDisplayData(vc, manifestCredential);
    },
    getManifest: function (credential) {
      const currentCredentialType = this.getCredentialType(credential);
      return credentialDisplayData[currentCredentialType] || credentialDisplayData.fallback;
    },
  },
};
</script>
