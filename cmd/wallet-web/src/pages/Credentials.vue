<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Credentials Layout -->
    <div v-if="breakpoints.xs || breakpoints.sm" class="flex flex-col justify-start w-screen">
      <flyout :tool-tip-label="t('Credentials.switchVaults')">
        <template #button="{ toggleFlyoutMenu, setShowTooltip }">
          <button
            id="credentials-flyout-menu-button-mobile"
            class="
              inline-flex
              justify-between
              items-center
              px-3
              w-screen
              md:w-auto
              h-11
              bg-neutrals-white
              md:rounded-lg
              focus:ring-primary-purple focus:ring-opacity-70 focus:ring-2
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
            <img class="w-6 h-6" src="@/assets/img/icons-sm--vault-icon.svg" />
            <span class="flex-grow pl-2 text-sm font-bold text-left text-neutrals-dark truncate">{{
              t('CredentialDetails.allVaultLabel')
            }}</span>
            <img class="w-6 h-6" src="@/assets/img/icons-sm--chevron-down-icon.svg" />
          </button>
        </template>
        <template #menu>
          <flyout-menu>
            <flyout-button
              id="flyout-menu-select-all-vaults"
              :text="t('Vaults.allVaults')"
              class="font-bold text-neutrals-dark"
              :class="!selectedVault.length && 'bg-neutrals-softWhite'"
            />
            <div class="my-1 mx-4 h-px bg-neutrals-thistle" />
            <flyout-button
              v-for="(vault, key) in vaults"
              :id="`flyout-menu-select-${vault.id.slice(-5)}`"
              :key="key"
              :text="vault.name"
              class="text-neutrals-medium"
            />
          </flyout-menu>
        </template>
      </flyout>
      <div class="items-start">
        <h3 class="mx-6 mb-5 font-bold text-neutrals-dark">{{ t('Credentials.credentials') }}</h3>
      </div>
    </div>
    <!-- Desktop Credentials Layout -->
    <div v-else class="flex justify-between items-center mb-8 w-full align-middle">
      <div class="flex flex-grow">
        <h3 class="m-0 font-bold text-neutrals-dark">{{ t('Credentials.credentials') }}</h3>
      </div>
      <flyout :tool-tip-label="t('Credentials.switchVaults')">
        <template #button="{ toggleFlyoutMenu, setShowTooltip }">
          <button
            id="credentials-flyout-menu-button-desktop"
            class="
              inline-flex
              justify-between
              items-center
              px-3
              w-screen
              md:w-auto
              h-11
              bg-neutrals-white
              md:rounded-lg
              focus:ring-primary-purple focus:ring-opacity-70 focus:ring-2
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
            <img class="w-6 h-6" src="@/assets/img/icons-sm--vault-icon.svg" />
            <span class="flex-grow pl-2 text-sm font-bold text-left text-neutrals-dark truncate">{{
              t('CredentialDetails.allVaultLabel')
            }}</span>
            <img class="w-6 h-6" src="@/assets/img/icons-sm--chevron-down-icon.svg" />
          </button>
        </template>
        <template #menu>
          <flyout-menu>
            <flyout-button
              id="flyout-menu-select-all-vaults"
              :text="t('Vaults.allVaults')"
              class="font-bold text-neutrals-dark"
              :class="!selectedVault.length && 'bg-neutrals-softWhite'"
            />
            <div class="my-1 mx-4 h-px bg-neutrals-thistle" />
            <flyout-button
              v-for="(vault, key) in vaults"
              :id="`flyout-menu-select-${vault.id.slice(-5)}`"
              :key="key"
              :text="vault.name"
              class="text-neutrals-medium"
            />
          </flyout-menu>
        </template>
      </flyout>
    </div>
    <!-- Loading State -->
    <skeleton-loader v-if="loading" type="vault" />
    <!-- Error State -->
    <!-- Main State -->
    <div v-else id="loaded-credentials-container">
      <div v-if="credentialsFound" class="mx-6 md:mx-0">
        <div v-for="(vault, key) in vaults" :key="key">
          <div v-if="vault.credentials.length">
            <div class="md:mx-0 mb-5">
              <span class="text-xl font-bold text-neutrals-dark">{{ vault.name }}</span>
            </div>
            <ul class="grid grid-cols-1 xl:grid-cols-2 gap-4 xl:gap-8 my-8">
              <li v-for="(credential, index) in vault.credentials" :key="index">
                <credential-preview
                  :id="credential.id"
                  :to="{
                    name: 'credential-details',
                    params: {
                      id: credential.id,
                      vaultName: vault.name,
                    },
                  }"
                  :brand-color="credential.brandColor"
                  :icon="credential.icon"
                  :title="credential.title"
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
import {CollectionManager, CredentialManager} from '@trustbloc/wallet-sdk';
import {mapGetters} from 'vuex';
import {useI18n} from 'vue-i18n';
import {getCredentialDisplayData, getCredentialType} from '@/utils/mixins';
import useBreakpoints from '@/plugins/breakpoints.js';
import CredentialPreview from '@/components/CredentialPreview/CredentialPreview.vue';
import SkeletonLoader from '@/components/SkeletonLoader/SkeletonLoader.vue';
import Flyout from '@/components/Flyout/Flyout.vue';
import FlyoutMenu from '@/components/Flyout/FlyoutMenu.vue';
import FlyoutButton from '@/components/Flyout/FlyoutButton.vue';

const filterBy = ['IssuerManifestCredential', 'GovernanceCredential'];
export default {
  name: 'Credentials',
  components: {
    CredentialPreview,
    SkeletonLoader,
    Flyout,
    FlyoutMenu,
    FlyoutButton,
  },
  setup() {
    const breakpoints = useBreakpoints();
    const { t } = useI18n();
    return { breakpoints, t };
  },
  data() {
    return {
      loading: true,
      vaults: [],
      credentialsFound: false,
      selectedVault: '',
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
    this.credentialDisplayData = await this.getCredentialManifestData();
    await this.fetchVaults(this.$route.params.vaultId);
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts', 'getCredentialManifestData']),
    fetchCredentials: async function (vaultId) {
      try {
        const { contents } = await this.credentialManager.getAll(this.token, {
          collectionID: vaultId,
        });
        if (!contents) return;

        const _filter = (id) => {
          return !contents[id].type.some((t) => filterBy.includes(t));
        };

        const credentials = Object.keys(contents)
          .filter(_filter)
          .map((id) => contents[id]);

        return credentials.map((credential) => {
          const manifest = this.getManifest(credential);
          return this.getCredentialDisplayData(credential, manifest);
        });
      } catch (e) {
        console.error(`failed to fetch credentials for vault ${vaultId}:`, e);
      }
    },
    // Fetching vaults (specific one, if vaultId is provided, otherwise all of them)
    fetchVaults: async function (vaultId) {
      this.vaults = [];
      if (vaultId) {
        try {
          const { content: vault } = await this.collectionManager.get(this.token, vaultId);
          this.vaults = [vault];
        } catch (e) {
          console.error(`failed to fetch vault ${vaultId}:`, e);
        }
      } else {
        try {
          const { contents: rawVaults } = await this.collectionManager.getAll(this.token);
          const vaults = Object.values(rawVaults);
          this.vaults = vaults;
        } catch (e) {
          console.error('failed to fetch vaults:', e);
        }
      }

      this.vaults = await Promise.all(
        this.vaults.map(async (vault) => {
          const credentials = await this.fetchCredentials(vault.id);
          if (credentials && credentials.length) this.credentialsFound = true;
          return { ...vault, credentials };
        })
      );
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
