<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="py-6 px-3">
    <div class="mb-8 w-full">
      <h3 class="text-neutrals-dark">{{ t('Vaults.heading') }}</h3>
    </div>
    <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4 xl:gap-8 w-full">
      <vault-card
        color="pink"
        :num-of-creds="t('Vaults.foundCredentials', numOfCreds)"
        :name="t('Vaults.allVaults')"
      />
      <!--TODO: Issue-1198 Add flyout menu to default and other vaults create this vault -->
      <vault-card
        v-for="(vault, index) in vaults"
        :key="index"
        :name="vault.name"
        :num-of-creds="t('Vaults.foundCredentials', vault.numOfCreds)"
        :vault-id="vault.id"
      >
        <slot v-if="vault.name == 'Default Vault'">
          <flyout-menu type="vault">
            <flyout-menu-list>
              <flyout-menu-button
                id="renameVault"
                :text="t('Vaults.renameVault')"
                color="neutrals-medium"
              />
            </flyout-menu-list>
          </flyout-menu>
        </slot>
        <slot v-else>
          <flyout-menu id="vaultFlyoutMenu" type="vault">
            <flyout-menu-list>
              <flyout-menu-button
                id="renameVault"
                :text="t('Vaults.renameVault')"
                color="neutrals-medium"
              />
              <flyout-menu-button
                id="deleteVault"
                :text="t('Vaults.deleteVault')"
                color="primary-vampire"
              />
            </flyout-menu-list>
          </flyout-menu>
        </slot>
      </vault-card>

      <vault-card type="addNew" :name="t('Vaults.addVault')" class="grid order-last" />
    </div>
  </div>
</template>

<script>
import { reactive, watchEffect } from 'vue';
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import VaultCard from '@/components/Vaults/VaultCard';
import { useI18n } from 'vue-i18n';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';
import FlyoutMenuList from '@/components/FlyoutMenu/FlyoutMenuList';
import FlyoutMenuButton from '@/components/FlyoutMenu/FlyoutMenuButton';

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
    VaultCard,
    FlyoutMenu,
    FlyoutMenuList,
    FlyoutMenuButton,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      numOfCreds: 0,
      vaults: [],
    };
  },
  created: async function () {
    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.username = this.getCurrentUser().username;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
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
