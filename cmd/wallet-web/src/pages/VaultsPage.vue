<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <WelcomeBannerComponent
      v-if="!skipWelcomeMsg && !skippedLocally && !loading"
      id="welcome-banner-close-button"
      class="md:mb-10"
      @click="updateUserPreferences"
    >
      <div class="flex flex-col justify-start items-start">
        <div class="inline-flex mb-2">
          <img src="@/assets/img/vault-icon-colored.svg" />
          <span
            class="flex items-center pl-3 text-base font-bold text-neutrals-dark whitespace-nowrap"
            >{{ t('Vaults.WelcomeBanner.AddAVault.heading') }}</span
          >
        </div>
        <span class="text-sm text-neutrals-medium">
          {{ t('Vaults.WelcomeBanner.AddAVault.message') }}
        </span>
      </div>
      <div class="flex flex-col justify-start items-start">
        <div class="inline-flex mb-2">
          <img src="@/assets/img/credential-icon-colored.svg" />
          <span
            class="flex items-center pl-3 text-base font-bold text-neutrals-dark whitespace-nowrap"
            >{{ t('Vaults.WelcomeBanner.AddACredential.heading') }}</span
          >
        </div>
        <span class="text-sm text-neutrals-medium">
          {{ t('Vaults.WelcomeBanner.AddACredential.message') }}
        </span>
      </div>
    </WelcomeBannerComponent>
    <div class="md:p-0 py-6 px-6">
      <span
        v-if="breakpoints.xs || breakpoints.sm"
        class="w-full text-3xl font-bold text-neutrals-dark"
      >
        {{ t('Vaults.heading') }}
      </span>
      <h3 v-else class="w-full text-neutrals-dark">
        {{ t('Vaults.heading') }}
      </h3>
      <SkeletonLoaderComponent v-if="loading" type="Vault" />
      <div
        v-else
        id="vaults-loaded"
        class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4 xl:gap-8 mt-6 md:mt-8 w-full"
      >
        <VaultCardComponent
          id="all-vaults-button"
          color="pink"
          :num-of-creds="t('Vaults.foundCredentials', numOfCreds)"
          :name="t('Vaults.allVaults')"
        />
        <!--TODO: Issue-1198 Add flyout menu to default and other vaults create this vault -->
        <VaultCardComponent
          v-for="(vault, index) in vaults"
          :id="`vault-card-${vault.name.replaceAll(' ', '-')}`"
          :key="index"
          :name="vault.name"
          :num-of-creds="t('Vaults.foundCredentials', vault.numOfCreds)"
          :vault-id="vault.id"
        >
          <FlyoutComponent v-if="vaults.length === 1">
            <template #button="{ toggleFlyoutMenu }">
              <button
                :id="`vaults-flyout-menu-button-${vault.name.replaceAll(' ', '-')}`"
                class="w-8 h-8 bg-neutrals-white hover:bg-neutrals-softWhite focus:bg-neutrals-mischka rounded-full"
                @click="toggleFlyoutMenu()"
              >
                <img class="p-2" src="@/assets/img/more-icon.svg" />
              </button>
            </template>
            <template #menu>
              <FlyoutMenuComponent>
                <FlyoutButtonComponent
                  id="renameVault"
                  :text="t('Vaults.RenameModal.renameVault')"
                  class="text-neutrals-medium"
                  @click="toggleRenameModal(vault.id, vault.name)"
                />
              </FlyoutMenuComponent>
            </template>
          </FlyoutComponent>

          <FlyoutComponent v-else>
            <template #button="{ toggleFlyoutMenu }">
              <button
                :id="`vaults-flyout-menu-button-${vault.name.replaceAll(' ', '-')}`"
                class="w-8 h-8 bg-neutrals-white hover:bg-neutrals-softWhite focus:bg-neutrals-mischka rounded-full"
                @click="toggleFlyoutMenu()"
              >
                <img class="p-2" src="@/assets/img/more-icon.svg" />
              </button>
            </template>
            <template #menu>
              <FlyoutMenuComponent>
                <FlyoutButtonComponent
                  id="renameVault"
                  :text="t('Vaults.RenameModal.renameVault')"
                  class="text-neutrals-medium"
                  @click="toggleRenameModal(vault.id, vault.name)"
                />
                <FlyoutButtonComponent
                  id="delete-vault-flyout-button"
                  :text="t('Vaults.DeleteModal.deleteVault')"
                  class="text-primary-vampire"
                  @click="toggleDeleteModal(vault.id)"
                />
              </FlyoutMenuComponent>
            </template>
          </FlyoutComponent>
        </VaultCardComponent>
        <VaultCardComponent
          id="add-new-vault-button"
          type="addNew"
          :name="t('Vaults.AddModal.addVault')"
          :existing-names="existingNames"
          class="grid order-last"
        />
      </div>
    </div>
    <DeleteVaultComponent
      :show="showDeleteModal"
      :vault-id="selectedVaultId"
      @close="handleDeleteModalClose"
    />
    <RenameVaultComponent
      :show="showRenameModal"
      :vault-id="selectedVaultId"
      :existing-names="existingNames"
      @close="handleRenameModalClose"
    />
  </div>
</template>

<script>
import { computed, reactive, ref, watchEffect } from 'vue';
import { CollectionManager, CredentialManager, WalletUser } from '@trustbloc/wallet-sdk';
import { mapActions, mapGetters, useStore } from 'vuex';
import { useI18n } from 'vue-i18n';
import useBreakpoints from '@/plugins/breakpoints.js';
import DeleteVaultComponent from '@/components/Vaults/DeleteVaultModalComponent';
import FlyoutComponent from '@/components/Flyout/FlyoutComponent';
import FlyoutMenuComponent from '@/components/Flyout/FlyoutMenuComponent';
import FlyoutButtonComponent from '@/components/Flyout/FlyoutButtonComponent';
import RenameVaultComponent from '@/components/Vaults/RenameVaultModalComponent';
import VaultCardComponent from '@/components/Vaults/VaultCardComponent';
import WelcomeBannerComponent from '@/components/Vaults/VaultWelcomeBannerComponent.vue';
import SkeletonLoaderComponent from '@/components/SkeletonLoader/SkeletonLoaderComponent.vue';

export const vaultsStore = reactive({
  vaultsOutdated: false,
});

export const vaultsMutations = {
  setVaultsOutdated(value) {
    vaultsStore.vaultsOutdated = value;
  },
};

export default {
  name: 'VaultsPage',
  components: {
    RenameVaultComponent,
    VaultCardComponent,
    FlyoutComponent,
    FlyoutMenuComponent,
    FlyoutButtonComponent,
    DeleteVaultComponent,
    WelcomeBannerComponent,
    SkeletonLoaderComponent,
  },
  setup() {
    const breakpoints = useBreakpoints();
    const { t } = useI18n();
    const store = useStore();
    const skipWelcomeMsg = computed(() => store?.state?.user?.preference?.skipWelcomeMsg);
    const showDeleteModal = ref(false);
    const showRenameModal = ref(false);
    const selectedVaultId = ref('');
    function toggleDeleteModal(vaultId) {
      selectedVaultId.value = vaultId;
      showDeleteModal.value = !showDeleteModal.value;
    }

    function toggleRenameModal(vaultId) {
      selectedVaultId.value = vaultId;
      showRenameModal.value = !showRenameModal.value;
    }

    return {
      breakpoints,
      t,
      showDeleteModal,
      showRenameModal,
      selectedVaultId,
      skipWelcomeMsg,
      toggleDeleteModal,
      toggleRenameModal,
    };
  },
  data() {
    return {
      numOfCreds: 0,
      vaults: [],
      skippedLocally: false,
      loading: true,
    };
  },
  computed: {
    existingNames() {
      return this.vaults.map((vault) => vault.name);
    },
  },
  created: async function () {
    const { profile, username } = this.getCurrentUser();
    const { user, token } = profile;
    this.token = token;
    this.username = username;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    this.walletUser = new WalletUser({ agent: this.getAgentInstance(), user });
    // Fetch all the credentials stored.
    // TODO: Issue-1250 Refactor to not to save credentials without vault ID.
    await this.getNumOfCreds();
    await this.fetchVaults();
    this.loading = false;
    watchEffect(async () => {
      if (vaultsStore.vaultsOutdated) {
        this.loading = true;
        await this.fetchVaults();
        vaultsMutations.setVaultsOutdated(false);
        this.loading = false;
      }
    });
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    ...mapActions(['refreshUserPreference']),
    getNumOfCreds: async function () {
      // Fetching all credentials
      const credentials = await this.credentialManager.getAllCredentialMetadata(this.token);
      this.numOfCreds = credentials.length;
    },
    fetchVaults: async function () {
      this.vaults = [];
      // Fetching all vaults
      const { contents: rawVaults } = await this.collectionManager.getAll(this.token);
      console.log(`found ${Object.keys(rawVaults).length} vaults`);

      const vaults = Object.values(rawVaults);
      // For each vault get a number of credentials
      await Promise.all(
        vaults.map(async (vault) => {
          // Fetching all credentials stored inside each vault
          // TODO: #1236 Revisit the solution to avoid getting all the credentials
          await this.credentialManager
            .getAllCredentialMetadata(this.token, { collection: vault.id })
            .then((credentials) => {
              vault['numOfCreds'] = credentials.length;
              this.vaults.push(vault);
            });
        })
      );
    },
    updateUserPreferences: async function () {
      try {
        // Used to close the welcome message banner instantly in the UI
        this.skippedLocally = true;
        await this.walletUser.updatePreferences(this.token, { skipWelcomeMsg: true });
        this.refreshUserPreference();
      } catch (e) {
        console.error('error updating user preferences', e);
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
<style>
.card-list {
  display: grid;
  grid-gap: 1em;
}

.card-item {
  padding: 2em;
}
</style>
