<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Credentials Layout -->
    <div v-if="breakpoints.xs || breakpoints.sm" class="flex flex-col justify-start w-screen">
      <FlyoutComponent :tool-tip-label="t('Credentials.switchVaults')">
        <template #button="{ toggleFlyoutMenu, setShowTooltip }">
          <button
            id="credentials-flyout-menu-button-mobile"
            class="inline-flex justify-between items-center px-3 w-screen md:w-auto h-11 bg-neutrals-white md:rounded-lg border border-neutrals-chatelle focus:ring-2 focus:ring-primary-purple focus:ring-opacity-70 focus-within:ring-offset-2 outline-none hover:border-neutrals-mountainMist-light"
            @click="toggleFlyoutMenu()"
            @focus="setShowTooltip(false)"
            @mouseover="setShowTooltip(true)"
            @mouseout="setShowTooltip(false)"
          >
            <img class="w-6 h-6" src="@/assets/img/icons-sm--vault-icon.svg" />
            <span class="flex-grow pl-2 text-sm font-bold text-left text-neutrals-dark truncate">{{
              selectedVaultName
            }}</span>
            <img class="w-6 h-6" src="@/assets/img/icons-sm--chevron-down-icon.svg" />
          </button>
        </template>
        <template #menu>
          <FlyoutMenuComponent>
            <FlyoutButtonComponent
              id="flyout-menu-select-all-vaults"
              :text="t('Vaults.allVaults')"
              class="font-bold text-neutrals-dark"
              :class="!selectedVaultId && 'bg-neutrals-softWhite'"
            />
            <div class="my-1 mx-4 h-px bg-neutrals-thistle" />
            <FlyoutButtonComponent
              v-for="(vault, key) in allVaults"
              :id="`flyout-menu-select-${vault.id.slice(-5)}`"
              :key="key"
              :text="vault.name"
              class="text-neutrals-medium"
              :class="selectedVaultId === vault.id && 'bg-neutrals-softWhite'"
              @click="setSelectedVaultId(vault.id)"
            />
          </FlyoutMenuComponent>
        </template>
      </FlyoutComponent>
      <div class="items-start">
        <h3 class="mx-6 mb-5 font-bold text-neutrals-dark">{{ t('Credentials.credentials') }}</h3>
      </div>
    </div>
    <!-- Desktop Credentials Layout -->
    <div v-else class="flex justify-between items-center mb-8 w-full align-middle">
      <div class="flex flex-grow">
        <h3 class="m-0 font-bold text-neutrals-dark">{{ t('Credentials.credentials') }}</h3>
      </div>
      <FlyoutComponent :tool-tip-label="t('Credentials.switchVaults')">
        <template #button="{ toggleFlyoutMenu, setShowTooltip }">
          <button
            id="credentials-flyout-menu-button-desktop"
            class="inline-flex justify-between items-center px-3 w-screen md:w-auto h-11 bg-neutrals-white md:rounded-lg border border-neutrals-chatelle focus:ring-2 focus:ring-primary-purple focus:ring-opacity-70 focus-within:ring-offset-2 outline-none hover:border-neutrals-mountainMist-light"
            @click="toggleFlyoutMenu()"
            @focus="setShowTooltip(false)"
            @mouseover="setShowTooltip(true)"
            @mouseout="setShowTooltip(false)"
          >
            <img class="w-6 h-6" src="@/assets/img/icons-sm--vault-icon.svg" />
            <span class="flex-grow pl-2 text-sm font-bold text-left text-neutrals-dark truncate">{{
              selectedVaultName
            }}</span>
            <img class="w-6 h-6" src="@/assets/img/icons-sm--chevron-down-icon.svg" />
          </button>
        </template>
        <template #menu="{ toggleFlyoutMenu }">
          <FlyoutMenuComponent>
            <FlyoutButtonComponent
              id="flyout-menu-select-all-vaults"
              :text="t('Vaults.allVaults')"
              class="font-bold text-neutrals-dark"
              :class="!selectedVaultId && 'bg-neutrals-softWhite'"
              @click="
                () => {
                  toggleFlyoutMenu();
                  setSelectedVaultId(null);
                }
              "
            />
            <div class="my-1 mx-4 h-px bg-neutrals-thistle" />
            <FlyoutButtonComponent
              v-for="(vault, key) in allVaults"
              :id="`flyout-menu-select-${vault.id.slice(-5)}`"
              :key="key"
              :text="vault.name"
              class="text-neutrals-medium"
              :class="selectedVaultId && selectedVaultId === vault.id && 'bg-neutrals-softWhite'"
              @click="
                () => {
                  toggleFlyoutMenu();
                  setSelectedVaultId(vault.id);
                }
              "
            />
          </FlyoutMenuComponent>
        </template>
      </FlyoutComponent>
    </div>
    <!-- Loading State -->
    <SkeletonLoaderComponent v-if="loading" type="vault" />
    <!-- Error State -->
    <span v-else-if="userSetupStatus === 'failed'">
      <b>Warning:</b> Failed to connect to server. Your wallet can not participate in secured
      communication.
    </span>
    <!-- Main State -->
    <div v-else id="loaded-credentials-container">
      <div v-if="credentialsFound" class="mx-6 md:mx-0">
        <div v-for="(vault, key) in selectedVaults" :key="key">
          <div v-if="vault.credentials.length">
            <div class="md:mx-0 mb-5">
              <span class="text-xl font-bold text-neutrals-dark">{{ vault.name }}</span>
            </div>
            <ul class="grid grid-cols-1 xl:grid-cols-2 gap-4 xl:gap-8 my-8">
              <li v-for="(credential, index) in vault.credentials" :key="index">
                <CredentialPreviewComponent
                  :id="credential.id"
                  :to="{
                    name: 'credential-details',
                    params: {
                      id: credential.id,
                      vaultName: vault.name,
                    },
                  }"
                  :styles="credential.styles"
                  :title="credential.name"
                />
              </li>
            </ul>
          </div>
        </div>
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
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { computed } from 'vue';
import { mapGetters, useStore } from 'vuex';
import { encode } from 'js-base64';
import { useI18n } from 'vue-i18n';
import useBreakpoints from '@/plugins/breakpoints.js';
import CredentialPreviewComponent from '@/components/CredentialPreview/CredentialPreviewComponent.vue';
import SkeletonLoaderComponent from '@/components/SkeletonLoader/SkeletonLoaderComponent.vue';
import FlyoutComponent from '@/components/Flyout/FlyoutComponent.vue';
import FlyoutMenuComponent from '@/components/Flyout/FlyoutMenuComponent.vue';
import FlyoutButtonComponent from '@/components/Flyout/FlyoutButtonComponent.vue';

const filterBy = ['IssuerManifestCredential', 'GovernanceCredential'];
export default {
  name: 'CredentialsPage',
  components: {
    CredentialPreviewComponent,
    SkeletonLoaderComponent,
    FlyoutComponent,
    FlyoutMenuComponent,
    FlyoutButtonComponent,
  },
  setup() {
    const store = useStore();
    const agentInstance = computed(() => store.getters['agent/getInstance']);
    const selectedVaultId = computed(() => store.getters.getSelectedVaultId);
    const currentUser = computed(() => store.getters.getCurrentUser);
    const userSetupStatus = currentUser.value ? currentUser.value.setupStatus : null;
    const getCredentialManifestData = () => store.getters.getCredentialManifestData;
    const updateSelectedVaultId = (selectedVaultId) =>
      store.dispatch('updateSelectedVaultId', selectedVaultId);
    const breakpoints = useBreakpoints();
    const { t } = useI18n();
    return {
      agentInstance,
      breakpoints,
      currentUser,
      getCredentialManifestData,
      selectedVaultId,
      t,
      userSetupStatus,
      updateSelectedVaultId,
    };
  },
  data() {
    return {
      loading: true,
      allVaults: [], // vaults to display in the flyout
      selectedVaults: [], // vaults to display in the main view along with credentials stored in each
      credentialsFound: false,
      selectedVaultName: null,
      credential: [],
    };
  },
  created: async function () {
    const { profile, username } = this.getCurrentUser();
    this.token = profile.token;
    const user = profile.user;
    this.username = username;
    this.credentialManager = new CredentialManager({ agent: this.agentInstance, user });
    this.collectionManager = new CollectionManager({ agent: this.agentInstance, user });
    await this.fetchVaults();
    this.loading = false;
  },
  methods: {
    ...mapGetters(['getCurrentUser']),
    fetchCredentials: async function (vaultId) {
      this.loading = true;
      try {
        const metadataList = await this.credentialManager.getAllCredentialMetadata(this.token, {
          collection: vaultId,
        });
        return metadataList.map((credential) => ({
          id: encode(credential.id),
          name: credential.name,
          ...credential.resolved[0],
        }));
      } catch (e) {
        console.error(`failed to fetch credentials for vault ${vaultId}:`, e);
      }
      this.loading = false;
    },
    // Fetching vaults (specific one, if vaultId is provided, otherwise all of them)
    fetchVaults: async function () {
      this.loading = true;
      this.credentialsFound = false;
      this.selectedVaults = [];
      this.allVaults = [];
      // Fetch all vaults
      try {
        const { contents: rawVaults } = await this.collectionManager.getAll(this.token);
        const vaults = Object.values(rawVaults);
        this.allVaults = vaults;
        // If user did not select any vault then we display all vaults in the main view
        if (!this.selectedVaultId) this.selectedVaults = vaults;
      } catch (e) {
        console.error('failed to fetch vaults:', e);
      }
      // If selected, fetch only these vaults to display in the main view
      if (this.selectedVaultId) {
        try {
          const { content: vault } = await this.collectionManager.get(
            this.token,
            this.selectedVaultId
          );
          this.selectedVaults = [vault];
        } catch (e) {
          console.error(`failed to fetch vault ${this.selectedVaultId}:`, e);
        }
      }

      // Determine which name to display in the flyout as selected vault
      this.selectedVaultName =
        this.selectedVaultId && this.selectedVaults.length
          ? this.selectedVaults[0].name
          : this.t('CredentialDetails.allVaultLabel');

      // Fetch and save all credentials stored inside each of the vaults to be displayed in the main view
      this.selectedVaults = await Promise.all(
        this.selectedVaults.map(async (vault) => {
          const credentials = await this.fetchCredentials(vault.id);
          // If any of the vaults selected to be displayed has at least one credential, we will render it
          // Otherwise, we will display an empty state
          if (credentials && credentials.length) this.credentialsFound = true;
          return { ...vault, credentials };
        })
      );
      this.loading = false;
    },
    setSelectedVaultId: async function (id) {
      if (this.selectedVaultId !== id) {
        this.updateSelectedVaultId(id);
        await this.fetchVaults();
      }
    },
  },
};
</script>
