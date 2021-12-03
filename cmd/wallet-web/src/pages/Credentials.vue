<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Credentials Layout -->
    <div v-if="breakpoints.xs || breakpoints.sm" class="flex flex-col justify-start w-screen">
      <div class="bg-neutrals-white border border-neutrals-chatelle">
        <flyout-menu />
      </div>
      <div class="items-start">
        <h3 class="mx-6 mb-5 font-bold text-neutrals-dark">{{ t('Credentials.credentials') }}</h3>
      </div>
    </div>
    <!-- Desktop Credentials Layout -->
    <div v-else class="flex justify-between items-center mb-8 w-full align-middle">
      <div class="flex flex-grow">
        <h3 class="m-0 font-bold text-neutrals-dark">{{ t('Credentials.credentials') }}</h3>
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
    <skeleton-loader v-if="loading" type="vault" />
    <span v-else-if="loadingStatus === 'failed'">
      <b>Warning:</b> Failed to connect to server. Your wallet can not participate in secured
      communication.
    </span>
    <div v-else id="loaded-credentials-container">
      <div v-if="processedCredentials.length" class="mx-6 md:mx-0">
        <div class="md:mx-0 mb-5">
          <span class="text-xl font-bold text-neutrals-dark">{{ selectedVault.name }}</span>
        </div>
        <ul class="grid grid-cols-1 xl:grid-cols-2 gap-4 xl:gap-8 my-8">
          <li v-for="(processedCredential, index) in processedCredentials" :key="index">
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
          <span class="text-base font-bold text-neutrals-medium">
            {{ t('Credentials.error') }}</span
          >
        </div>
        <div class="flex justify-center">
          <span class="text-base text-neutrals-medium"> {{ t('Credentials.description') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { CredentialManager, CollectionManager } from '@trustbloc/wallet-sdk';
import { getCredentialType, getCredentialDisplayData } from '@/utils/mixins';
import { mapActions, mapGetters } from 'vuex';
import CredentialPreview from '@/components/CredentialPreview/CredentialPreview';
import SkeletonLoader from '@/components/SkeletonLoader/SkeletonLoader';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';
import useBreakpoints from '@/plugins/breakpoints.js';
import { useI18n } from 'vue-i18n';

const filterBy = ['IssuerManifestCredential', 'GovernanceCredential'];
export default {
  name: 'Credentials',
  components: {
    CredentialPreview,
    SkeletonLoader,
    FlyoutMenu,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      processedCredentials: [],
      username: '',
      agent: null,
      breakpoints: useBreakpoints(),
      credentialDisplayData: '',
      selectedVault: { id: this.$route.params.vaultId, name: '' },
      loading: true,
    };
  },
  computed: {
    loadingStatus() {
      return this.getCurrentUser() ? this.getCurrentUser().setupStatus : null;
    },
  },
  created: async function () {
    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.username = this.getCurrentUser().username;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    this.fetchCredentials(
      this.credentialManager.getAll(token, { collectionID: this.selectedVault.id })
    );
    this.credentialDisplayData = await this.getCredentialManifestData();
    if (this.selectedVault.id) await this.fetchSelectedVault();
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'getCredentialManifestData']),
    ...mapActions(['updateProcessedCredentials']),
    fetchCredentials: async function () {
      const { contents } = await this.credentialManager.getAll(this.token, {
        collectionID: this.$route.params.vaultId,
      });
      const _filter = (id) => {
        return !contents[id].type.some((t) => filterBy.includes(t));
      };

      const credentials = Object.keys(contents)
        .filter(_filter)
        .map((id) => contents[id]);

      credentials.map((credential) => {
        const manifest = this.getManifest(credential);
        const processedCredential = this.getCredentialDisplayData(credential, manifest);
        this.processedCredentials.push(processedCredential);
      });

      this.updateProcessedCredentials(this.processedCredentials);

      console.debug(`showing ${this.processedCredentials.length} credentials`);
    },
    fetchSelectedVault: async function () {
      const {
        content: { id, name },
      } = await this.collectionManager.get(this.token, this.$route.params.vaultId);
      this.selectedVault = { id, name };
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
