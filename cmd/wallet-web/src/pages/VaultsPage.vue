<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<script>
import { reactive } from 'vue';
export const vaultsStore = reactive({
  vaultsOutdated: false,
});

export const vaultsMutations = {
  setVaultsOutdated(value) {
    vaultsStore.vaultsOutdated = value;
  },
};
</script>

<script setup>
import { computed, onMounted, ref, watchEffect } from 'vue';
import { CollectionManager, CredentialManager, WalletUser } from '@trustbloc/wallet-sdk';
import { useStore } from 'vuex';
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

// Hooks
const { t } = useI18n();
const store = useStore();
const breakpoints = useBreakpoints();

// Local Variables
const numOfCreds = ref(0);
const vaults = ref([]);
const skippedLocally = ref(false);
const loading = ref(true);
const showDeleteModal = ref(false);
const showRenameModal = ref(false);
const selectedVaultId = ref('');
const authToken = ref(null);
const credentialManager = ref(null);
const collectionManager = ref(null);
const walletUser = ref(null);
const existingNames = computed(() => vaults.value.map((vault) => vault.name));

// Store Getters
const currentUser = computed(() => store.getters['getCurrentUser']);
const agentInstance = computed(() => store.getters['agent/getInstance']);

// Store Actions
const refreshUserPreference = () => store.dispatch('refreshUserPreference');

// Methods
function toggleDeleteModal(vaultId) {
  selectedVaultId.value = vaultId;
  showDeleteModal.value = !showDeleteModal.value;
}

function toggleRenameModal(vaultId) {
  selectedVaultId.value = vaultId;
  showRenameModal.value = !showRenameModal.value;
}

async function getNumOfCreds() {
  // Fetching all credentials
  const credentials = await credentialManager.value.getAllCredentialMetadata(authToken.value);
  numOfCreds.value = credentials.length;
}

async function fetchVaults() {
  vaults.value = [];
  // Fetching all vaults
  const { contents: rawVaults } = await collectionManager.value.getAll(authToken.value);

  const vaultsArray = Object.values(rawVaults);
  // For each vault get a number of credentials
  await Promise.all(
    vaultsArray.map(async (vault) => {
      // Fetching all credentials stored inside each vault
      // TODO: #1236 Revisit the solution to avoid getting all the credentials
      await credentialManager.value
        .getAllCredentialMetadata(authToken.value, { collection: vault.id })
        .then((credentials) => {
          vault['numOfCreds'] = credentials.length;
          vaults.value.push(vault);
        });
    })
  );
}

async function updateUserPreferences() {
  try {
    // Used to close the welcome message banner instantly in the UI
    skippedLocally.value = true;
    await walletUser.value.updatePreferences(authToken.value, { skipWelcomeMsg: true });
    await refreshUserPreference();
  } catch (e) {
    console.error('error updating user preferences', e);
  }
}

function handleRenameModalClose() {
  showRenameModal.value = false;
}

function handleDeleteModalClose() {
  showDeleteModal.value = false;
}

onMounted(async () => {
  const { profile } = currentUser.value;
  const { user, token } = profile;
  authToken.value = token;
  credentialManager.value = new CredentialManager({ agent: agentInstance.value, user });
  collectionManager.value = new CollectionManager({ agent: agentInstance.value, user });
  walletUser.value = new WalletUser({ agent: agentInstance.value, user });
  // Fetch all the credentials stored.
  // TODO: Issue-1250 Refactor to not to save credentials without vault ID.
  await Promise.all([getNumOfCreds(), fetchVaults()]);
  loading.value = false;
  watchEffect(async () => {
    if (vaultsStore.vaultsOutdated) {
      loading.value = true;
      await fetchVaults();
      vaultsMutations.setVaultsOutdated(false);
      loading.value = false;
    }
  });
});
</script>

<template>
  <div>
    <WelcomeBannerComponent
      v-if="!currentUser?.preference?.skipWelcomeMsg && !skippedLocally && !loading"
      id="welcome-banner-close-button"
      class="md:mb-10"
      @click="updateUserPreferences"
    >
      <div class="flex flex-col items-start justify-start">
        <div class="mb-2 inline-flex">
          <img src="@/assets/img/vault-icon-colored.svg" />
          <span
            class="flex items-center whitespace-nowrap pl-3 text-base font-bold text-neutrals-dark"
            >{{ t('Vaults.WelcomeBanner.AddAVault.heading') }}</span
          >
        </div>
        <span class="text-sm text-neutrals-medium">
          {{ t('Vaults.WelcomeBanner.AddAVault.message') }}
        </span>
      </div>
      <div class="flex flex-col items-start justify-start">
        <div class="mb-2 inline-flex">
          <img src="@/assets/img/credential-icon-colored.svg" />
          <span
            class="flex items-center whitespace-nowrap pl-3 text-base font-bold text-neutrals-dark"
            >{{ t('Vaults.WelcomeBanner.AddACredential.heading') }}</span
          >
        </div>
        <span class="text-sm text-neutrals-medium">
          {{ t('Vaults.WelcomeBanner.AddACredential.message') }}
        </span>
      </div>
    </WelcomeBannerComponent>
    <div class="p-6 md:p-0">
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
        class="mt-6 grid w-full grid-cols-1 gap-4 md:mt-8 lg:grid-cols-2 xl:grid-cols-3 xl:gap-8"
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
                class="h-8 w-8 rounded-full bg-neutrals-white hover:bg-neutrals-softWhite focus:bg-neutrals-mischka"
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
                class="h-8 w-8 rounded-full bg-neutrals-white hover:bg-neutrals-softWhite focus:bg-neutrals-mischka"
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
          class="order-last grid"
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
