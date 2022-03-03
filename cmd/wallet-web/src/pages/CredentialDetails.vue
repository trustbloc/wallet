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
      <flyout :tool-tip-label="t('CredentialDetails.toolTipLabel')">
        <template #button="{ toggleFlyoutMenu, setShowTooltip }">
          <button
            id="credential-details-flyout-button"
            class="
              w-11
              h-11
              bg-neutrals-white
              rounded-lg
              focus:border-neutrals-chatelle
              focus:ring-2
              focus:ring-primary-purple
              focus:ring-opacity-70
              focus-within:ring-offset-2
              outline-none
              border border-neutrals-chatelle
              hover:border-neutrals-mountainMist-light
            "
            @click="toggleFlyoutMenu()"
            @focus="setShowTooltip(false)"
            @mouseover="setShowTooltip(true)"
            @mouseout="setShowTooltip(false)"
          >
            <img alt="flyout menu icon" class="p-2" src="@/assets/img/more-icon.svg" />
          </button>
        </template>
        <template #menu>
          <flyout-menu>
            <flyout-button
              id="renameCredential"
              :text="t('CredentialDetails.renameCredential')"
              class="text-neutrals-medium"
            />
            <flyout-button id="moveCredential" :text="t('CredentialDetails.moveCredential')" />
            <flyout-button
              id="deleteCredential"
              :text="t('CredentialDetails.deleteCredential')"
              class="text-primary-vampire"
              @click="toggleDelete"
            >
            </flyout-button>
          </flyout-menu>
        </template>
      </flyout>
      <delete-credential :show="showModal" :credential-id="credential.id" />
    </div>
    <banner
      :styles="credential.styles"
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
          <td
            v-if="property.schema.format === 'image/png'"
            class="py-4 pr-6 pl-3 text-neutrals-dark break-words"
          >
            <img :src="property.value" class="w-20 h-20" />
          </td>
          <td v-else class="py-4 pr-6 pl-3 text-neutrals-dark break-words">
            {{ property.value }}
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
import Banner from '@/components/CredentialDetails/Banner.vue';
import Flyout from '@/components/Flyout/Flyout.vue';
import FlyoutMenu from '@/components/Flyout/FlyoutMenu.vue';
import FlyoutButton from '@/components/Flyout/FlyoutButton.vue';
import DeleteCredential from '@/components/CredentialDetails/DeleteCredentialModal.vue';

export default {
  name: 'CredentialDetails',
  components: {
    Banner,
    Flyout,
    FlyoutMenu,
    FlyoutButton,
    DeleteCredential,
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
    const { profile, username } = this.getCurrentUser();
    this.token = profile.token;
    const user = profile.user;
    this.username = username;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    try {
      const { id, issuanceDate, resolved } = await this.credentialManager.getCredentialMetadata(
        this.token,
        decode(this.$route.params.id)
      );
      this.credential = { id, issuanceDate, ...resolved[0] };
    } catch (e) {
      console.error('failed to fetch a credential:', e);
    }
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
  },
};
</script>
