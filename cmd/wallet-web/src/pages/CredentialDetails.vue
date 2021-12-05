<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <!-- TODO: add a loading state -->
  <div v-if="credential" class="flex flex-col justify-start items-start py-6 px-3">
    <div class="flex flex-row justify-between items-center mb-4 w-full">
      <div class="flex flex-grow">
        <h3 class="text-neutrals-dark">{{ t('CredentialDetails.heading') }}</h3>
      </div>
      <div
        class="inline-flex items-center bg-transparent rounded-lg border border-neutrals-chatelle"
      >
        <flyout-menu id="credFlyoutMenu" type="outline">
          <flyout-menu-list>
            <flyout-menu-button
              id="renameCredential"
              :text="t('CredentialDetails.renameCredential')"
              color="neutrals-medium"
            />
            <flyout-menu-button
              id="moveCredential"
              :text="t('CredentialDetails.moveCredential')"
              color="neutrals-medium"
            />
            <flyout-menu-button
              id="deleteCredential"
              :text="t('CredentialDetails.deleteCredential')"
              color="primary-vampire"
              @click="toggleDelete"
            >
            </flyout-menu-button>
          </flyout-menu-list>
        </flyout-menu>
        <delete-credential :show="showModal" :credential-id="credential.id" />
      </div>
    </div>
    <banner
      :brand-color="credential.brandColor"
      :icon="credential.icon"
      :title="credential.title"
      :issuance-date="credential.issuanceDate"
      :vault-name="vaultName"
    />
    <!-- List of Credential Details -->
    <div class="flex flex-col justify-start items-start mt-8 md:mt-2 w-full">
      <span class="mb-5 text-xl font-bold text-neutrals-dark">{{
        t('CredentialDetails.verifiedInformation')
      }}</span>
      <table class="w-full border-t border-neutrals-chatelle">
        <tr
          v-for="(property, index) of credential.properties"
          :key="index"
          class="border-b border-neutrals-thistle border-dotted"
        >
          <!-- TODO: Add the dropdown for the nested credentials 1016 -->
          <td class="py-4 pr-6 pl-3 text-neutrals-medium">{{ property.label }}</td>
          <td v-if="property.type != 'image'" class="py-4 pr-6 pl-3 text-neutrals-dark break-words">
            {{ property.value }}
          </td>
          <td
            v-if="property.type === 'image'"
            class="py-4 pr-6 pl-3 text-neutrals-dark break-words"
          >
            <img :src="property.value" class="w-20 h-20" />
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { mapGetters } from 'vuex';
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { decode } from 'js-base64';
import { getCredentialType, getCredentialDisplayData } from '@/utils/mixins';
import Banner from '@/components/CredentialDetails/Banner.vue';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu.vue';
import FlyoutMenuList from '@/components/FlyoutMenu/FlyoutMenuList.vue';
import FlyoutMenuButton from '@/components/FlyoutMenu/FlyoutMenuButton.vue';
import DeleteCredential from '@/components/CredentialDetails/DeleteCredentialModal.vue';

export default {
  name: 'CredentialDetails',
  components: {
    DeleteCredential,
    Banner,
    FlyoutMenu,
    FlyoutMenuList,
    FlyoutMenuButton,
  },
  setup() {
    const { t } = useI18n();
    const showModal = ref(false);

    function toggleDelete() {
      showModal.value = !showModal.value;
    }

    return {
      t,
      showModal,
      toggleDelete,
    };
  },
  data() {
    return {
      loading: true,
      credential: null,
      vaultName: this.$route.params.vaultName,
    };
  },
  created: async function () {
    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.username = this.getCurrentUser().username;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.credentialDisplayData = await this.getCredentialManifestData();
    try {
      this.credential = await this.fetchCredential(this.$route.params.fullCredId);
    } catch (e) {
      console.error('failed to fetch a credential:', e);
    }
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'getCredentialManifestData']),
    fetchCredential: async function (id) {
      const { content: credential } = await this.credentialManager.get(this.token, decode(id));
      const manifest = this.getManifest(credential);
      return this.getCredentialDisplayData(credential, manifest);
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
