<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <welcome-banner
      v-if="!skipWelcomeMsg && !skippedLocally"
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
    </welcome-banner>
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
      <div
        class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4 xl:gap-8 mt-6 md:mt-8 w-full"
      >
        <vault-card
          id="all-vaults-button"
          color="pink"
          :num-of-creds="t('Vaults.foundCredentials', numOfCreds)"
          :name="t('Vaults.allVaults')"
        />
        <!--TODO: Issue-1198 Add flyout menu to default and other vaults create this vault -->
        <vault-card
          v-for="(vault, index) in vaults"
          :id="`vault-card-${vault.name.replaceAll(' ', '-')}`"
          :key="index"
          :name="vault.name"
          :num-of-creds="t('Vaults.foundCredentials', vault.numOfCreds)"
          :vault-id="vault.id"
        >
          <flyout v-if="vaults.length === 1">
            <template #button="{ toggleFlyoutMenu }">
              <button
                id="vaults-flyout-menu-button-default"
                class="
                  w-8
                  h-8
                  bg-neutrals-white
                  hover:bg-neutrals-softWhite
                  focus:bg-neutrals-mischka
                  rounded-full
                "
                @click="toggleFlyoutMenu()"
              >
                <img class="p-2" src="@/assets/img/more-icon.svg" />
              </button>
            </template>
            <template #menu>
              <flyout-menu>
                <flyout-button
                  id="renameVault"
                  :text="t('Vaults.RenameModal.renameVault')"
                  class="text-neutrals-medium"
                  @click="toggleRenameModal(vault.id, vault.name)"
                />
              </flyout-menu>
            </template>
          </flyout>

          <flyout v-else>
            <template #button="{ toggleFlyoutMenu }">
              <button
                id="vaults-flyout-menu-button-default"
                class="
                  w-8
                  h-8
                  bg-neutrals-white
                  hover:bg-neutrals-softWhite
                  focus:bg-neutrals-mischka
                  rounded-full
                "
                @click="toggleFlyoutMenu()"
              >
                <img class="p-2" src="@/assets/img/more-icon.svg" />
              </button>
            </template>
            <template #menu>
              <flyout-menu>
                <flyout-button
                  id="renameVault"
                  :text="t('Vaults.RenameModal.renameVault')"
                  class="text-neutrals-medium"
                  @click="toggleRenameModal(vault.id, vault.name)"
                />
                <flyout-button
                  :id="`delete-vault-${vault.id.slice(-5)}`"
                  :text="t('Vaults.DeleteModal.deleteVault')"
                  class="text-primary-vampire"
                  @click="toggleDeleteModal(vault.id)"
                />
              </flyout-menu>
            </template>
          </flyout>
        </vault-card>
        <vault-card
          id="add-new-vault-button"
          type="addNew"
          :name="t('Vaults.AddModal.addVault')"
          :existing-names="existingNames"
          class="grid order-last"
        />
      </div>
    </div>
    <delete-vault :show="showDeleteModal" :vault-id="selectedVaultId" />
    <rename-vault
      :show="showRenameModal"
      :vault-id="selectedVaultId"
      :existing-names="existingNames"
    />
  </div>
</template>

<script>
import { computed, reactive, ref, watchEffect } from 'vue';
import { useStore } from 'vuex';
import { CollectionManager, CredentialManager, WalletUser } from '@trustbloc/wallet-sdk';
import { mapGetters, mapActions } from 'vuex';
import { useI18n } from 'vue-i18n';
import useBreakpoints from '@/plugins/breakpoints.js';
import DeleteVault from '@/components/Vaults/DeleteVaultModal';
import Flyout from '@/components/Flyout/Flyout';
import FlyoutMenu from '@/components/Flyout/FlyoutMenu';
import FlyoutButton from '@/components/Flyout/FlyoutButton';
import RenameVault from '@/components/Vaults/RenameVaultModal';
import VaultCard from '@/components/Vaults/VaultCard';
import WelcomeBanner from '@/components/Vaults/WelcomeBanner.vue';

export const vaultsStore = reactive({
  vaultsOutdated: false,
});

export const vaultsMutations = {
  setVaultsOutdated(value) {
    vaultsStore.vaultsOutdated = value;
  },
};

export default {
  name: 'Vaults',
  components: {
    RenameVault,
    VaultCard,
    Flyout,
    FlyoutMenu,
    FlyoutButton,
    DeleteVault,
    WelcomeBanner,
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
    watchEffect(async () => {
      if (vaultsStore.vaultsOutdated) {
        await this.fetchVaults();
        vaultsMutations.setVaultsOutdated(false);
      }
    });
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    ...mapActions(['refreshUserPreference']),
    getNumOfCreds: async function () {
      // Fetching all credentials
      const { contents: credentials } = await this.credentialManager.getAll(this.token);
      this.numOfCreds = Object.keys(credentials).length;
    },
    fetchVaults: async function () {
      this.vaults = [];
      // Fetching all vaults
      const { contents: rawVaults } = await this.collectionManager.getAll(this.token);
      console.log(`found ${Object.keys(rawVaults).length} vaults`);

      const vaults = Object.values(rawVaults);
      // For each vault get a number of credentials
      vaults.forEach(async (vault) => {
        // Fetching all credentials stored inside each vault
        // TODO: #1236 Revisit the solution to avoid getting all the credentials
        await this.credentialManager
          .getAll(this.token, { collectionID: vault.id })
          .then(({ contents: credentials }) => {
            vault['numOfCreds'] = Object.keys(credentials).length;
            this.vaults.push(vault);
          });
      });
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
