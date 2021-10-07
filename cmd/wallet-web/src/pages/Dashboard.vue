<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Dashboard Layout -->
    <div v-if="breakpoints.xs || breakpoints.sm" class="flex flex-col justify-start w-screen">
      <div class="bg-neutrals-white border border-neutrals-chatelle">
        <flyout-menu />
      </div>
      <div class="items-start">
        <h3 class="mx-6 mb-5 font-bold text-neutrals-dark">{{ i18n.credentials }}</h3>
      </div>
    </div>
    <!-- Desktop Dashboard Layout -->
    <div v-else class="flex justify-between items-center mb-8 w-full align-middle">
      <div class="flex flex-grow">
        <h3 class="m-0 font-bold text-neutrals-dark">{{ i18n.credentials }}</h3>
      </div>
      <div
        class="
          inline-flex
          items-center
          bg-neutrals-white
          rounded-lg
          border border-neutrals-chatelle
        "
      >
        <flyout-menu />
      </div>
    </div>
    <skeleton-loader v-if="loadingStatus === 'inProgress'" type="vault" />
    <span v-else-if="loadingStatus === 'failed'">
      <b>Warning:</b> Failed to connect to server. Your wallet can not participate in secured
      communication.
    </span>
    <div v-else>
      <span id="dashboard-success-msg" class="px-8" />
      <div v-if="processedCredentials.length" class="mx-6 md:mx-0">
        <div class="md:mx-0 mb-5">
          <span class="text-xl font-bold text-neutrals-dark">{{ i18n.defaultVault }}</span>
        </div>
        <ul class="grid grid-cols-1 xl:grid-cols-2 gap-4 xl:gap-8 my-8">
          <li v-for="(processedCredential, index) in processedCredentials" :key="index">
            <!-- todo  load credential images from external source -->
            <credential-preview
              :id="processedCredential.id"
              :brand-color="processedCredential.brandColor"
              :icon="processedCredential.icon"
              :title="processedCredential.title"
            />
          </li>
        </ul>
      </div>
      <div
        v-else
        class="py-8 px-6 mx-auto rounded-lg border border-neutrals-thistle nocredentialCard"
      >
        <div class="flex justify-center">
          <img src="@/assets/img/icons-md--credentials-icon.svg" />
        </div>
        <div class="flex justify-center">
          <span class="text-base font-bold text-neutrals-medium"> {{ i18n.error }}</span>
        </div>
        <div class="flex justify-center">
          <span class="text-base text-neutrals-medium"> {{ i18n.description }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { getCredentialType, getCredentialDisplayData } from '@/pages/mixins';
import { mapActions, mapGetters } from 'vuex';
import CredentialPreview from '@/components/CredentialPreview/CredentialPreview';
import SkeletonLoader from '@/components/SkeletonLoader/SkeletonLoader';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';
import useBreakpoints from '@/plugins/breakpoints.js';

const filterBy = ['IssuerManifestCredential', 'GovernanceCredential'];
export default {
  name: 'Dashboard',
  components: {
    CredentialPreview,
    SkeletonLoader,
    FlyoutMenu,
  },
  data() {
    return {
      processedCredentials: [],
      username: '',
      agent: null,
      breakpoints: useBreakpoints(),
      credentialDisplayData: '',
    };
  },
  computed: {
    loadingStatus() {
      return this.getCurrentUser() ? this.getCurrentUser().setupStatus : null;
    },
    i18n() {
      return this.$t('Credentials');
    },
  },
  created: async function () {
    let { user, token } = this.getCurrentUser().profile;
    this.username = this.getCurrentUser().username;
    let credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.fetchAllCredentials(credentialManager.getAll(token));
    this.credentialDisplayData = await this.getCredentialManifestData();
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'getCredentialManifestData']),
    ...mapActions(['updateProcessedCredentials']),
    fetchAllCredentials: async function (getCredential) {
      const { contents } = await getCredential;
      console.log(`found ${Object.keys(contents).length} credentials`);
      const _filter = (id) => {
        return !contents[id].type.some((t) => filterBy.includes(t));
      };

      const credentials = Object.keys(contents)
        .filter(_filter)
        .map((id) => contents[id]);

      console.debug(`showing ${credentials.length} credentials`);

      credentials.map((credential) => {
        const manifest = this.getManifest(credential);
        const processedCredential = this.getCredentialDisplayData(credential, manifest);
        this.processedCredentials.push(processedCredential);
      });

      this.updateProcessedCredentials(this.processedCredentials);

      console.debug(`showing ${this.processedCredentials.length} credentials`);
    },
    getCredentialType: function (vc) {
      return getCredentialType(vc.type);
    },
    getCredentialDisplayData: function (vc, manifestCredential) {
      return getCredentialDisplayData(vc, manifestCredential);
    },
    getManifest: function (credential) {
      const currentCredentialType = this.getCredentialType(credential);
      return (
        this.credentialDisplayData[currentCredentialType] || this.credentialDisplayData.fallback
      );
    },
  },
};
</script>
