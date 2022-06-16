<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <!-- Loading State -->
  <div v-if="loading" class="flex flex-col justify-start items-start py-6 px-3">
    <div class="flex flex-row justify-between items-center mb-4 w-full">
      <h3 class="text-neutrals-dark">{{ t('CredentialDetails.heading') }}</h3>
    </div>
    <SkeletonLoaderComponent type="CredentialDetailsBanner" />
    <div class="flex flex-col justify-start items-start mt-8 md:mt-2 w-full">
      <span class="mb-5 text-xl font-bold text-neutrals-dark">{{
        t('CredentialDetails.verifiedInformation')
      }}</span>
      <SkeletonLoaderComponent class="w-full" type="VerifiedInformation" />
    </div>
  </div>
  <div v-else-if="credential" class="flex flex-col justify-start items-start py-6 px-3">
    <div class="flex flex-row justify-between items-center mb-4 w-full">
      <div class="flex flex-grow">
        <h3 class="text-neutrals-dark">{{ t('CredentialDetails.heading') }}</h3>
      </div>
      <FlyoutComponent :tool-tip-label="t('CredentialDetails.toolTipLabel')">
        <template #button="{ toggleFlyoutMenu, setShowTooltip }">
          <button
            id="credential-details-flyout-button"
            class="w-11 h-11 bg-neutrals-white rounded-lg border border-neutrals-chatelle focus:border-neutrals-chatelle focus:ring-2 focus:ring-primary-purple focus:ring-opacity-70 focus-within:ring-offset-2 outline-none hover:border-neutrals-mountainMist-light"
            @click="toggleFlyoutMenu()"
            @focus="setShowTooltip(false)"
            @mouseover="setShowTooltip(true)"
            @mouseout="setShowTooltip(false)"
          >
            <img alt="flyout menu icon" class="p-2" src="@/assets/img/more-icon.svg" />
          </button>
        </template>
        <template #menu>
          <FlyoutMenuComponent>
            <FlyoutButtonComponent
              id="renameCredential"
              :text="t('CredentialDetails.RenameModal.renameCredential')"
              class="text-neutrals-medium"
              @click="toggleRenameModal"
            />
            <FlyoutButtonComponent
              id="moveCredential"
              :text="t('CredentialDetails.MoveModal.moveCredential')"
            />
            <FlyoutButtonComponent
              id="deleteCredential"
              :text="t('CredentialDetails.DeleteModal.deleteCredential')"
              class="text-primary-vampire"
              @click="toggleDeleteModal"
            >
            </FlyoutButtonComponent>
          </FlyoutMenuComponent>
        </template>
      </FlyoutComponent>
      <DeleteCredentialComponent
        :show="showDeleteModal"
        :credential-id="credential.id"
        @close="handleDeleteModalClose"
      />
      <RenameCredentialComponent
        :show="showRenameModal"
        :credential-id="credential.id"
        :name="credential.name"
        :vault-name="vaultName"
        @close="handleRenameModalClose"
      />
    </div>
    <CredentialDetailsBannerComponent
      :styles="credential.styles"
      :title="credential.name"
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
import { reactive, ref, watchEffect } from 'vue';
import { useI18n } from 'vue-i18n';
import { mapGetters } from 'vuex';
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { decode } from 'js-base64';
import CredentialDetailsBannerComponent from '@/components/CredentialDetails/CredentialDetailsBannerComponent.vue';
import FlyoutComponent from '@/components/Flyout/FlyoutComponent.vue';
import FlyoutMenuComponent from '@/components/Flyout/FlyoutMenuComponent.vue';
import FlyoutButtonComponent from '@/components/Flyout/FlyoutButtonComponent.vue';
import DeleteCredentialComponent from '@/components/CredentialDetails/DeleteCredentialModalComponent.vue';
import RenameCredentialComponent from '@/components/CredentialDetails/RenameCredentialModalComponent.vue';
import SkeletonLoaderComponent from '@/components/SkeletonLoader/SkeletonLoaderComponent.vue';

export const credentialStore = reactive({
  credentialOutdated: false,
});

export const credentialMutations = {
  setCredentialOutdated(value) {
    credentialStore.credentialOutdated = value;
  },
};
export default {
  name: 'CredentialDetailsPage',
  components: {
    CredentialDetailsBannerComponent,
    FlyoutComponent,
    FlyoutMenuComponent,
    FlyoutButtonComponent,
    DeleteCredentialComponent,
    RenameCredentialComponent,
    SkeletonLoaderComponent,
  },
  setup() {
    const { t } = useI18n();
    const showDeleteModal = ref(false);
    const showRenameModal = ref(false);
    const selectedCredentialId = ref('');

    function toggleDeleteModal() {
      showDeleteModal.value = !showDeleteModal.value;
    }

    function toggleRenameModal(credentialId) {
      selectedCredentialId.value = credentialId;
      showRenameModal.value = !showRenameModal.value;
    }

    return {
      t,
      showDeleteModal,
      showRenameModal,
      selectedCredentialId,
      toggleDeleteModal,
      toggleRenameModal,
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
    await this.fetchCredential();
    this.loading = false;
    watchEffect(async () => {
      if (credentialStore.credentialOutdated) {
        this.loading = true;
        await this.fetchCredential();
        credentialMutations.setCredentialOutdated(false);
        this.loading = false;
      }
    });
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    fetchCredential: async function () {
      const { profile } = this.getCurrentUser();
      this.token = profile.token;
      const user = profile.user;
      this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
      try {
        const { id, name, issuanceDate, resolved } =
          await this.credentialManager.getCredentialMetadata(
            this.token,
            decode(this.$route.params.id)
          );
        this.credential = { id, issuanceDate, name, ...resolved[0] };
      } catch (e) {
        console.error('failed to fetch a credential:', e);
      }
    },
    handleRenameModalClose: function () {
      this.showRenameModal = false;
    },
    handleDeleteModalClose: function () {
      this.showDeleteModal = false;
    },
  },
};
</script>
