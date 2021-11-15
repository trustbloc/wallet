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
        :num-of-creds="t('Vaults.foundCredentials', numofCreds)"
        :name="t('Vaults.allVaults')"
      />
      <!--TODO: Issue-1198 Add flyout menu to default and other vaults create this vault-->
      <vault-card type="addNew" :name="t('Vaults.addVault')" class="grid order-last" />
      <div v-for="(vault, index) in vaults" :key="index">
        <vault-card
          :name="vault.name"
          :num-of-creds="t('Vaults.foundCredentials', vault.numofCreds)"
          class="grid order-last"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import VaultCard from '@/components/Vaults/VaultCard';
import { useI18n } from 'vue-i18n';

export default {
  name: 'Vaults',
  components: {
    VaultCard,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      numofCreds: 0,
      vaults: [],
    };
  },
  created: async function () {
    const { user, token } = this.getCurrentUser().profile;
    this.username = this.getCurrentUser().username;
    const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    // Fetch all the credentials stored.
    // TODO: Issue-1250 Refactor to not to save credentials without vault ID.
    await this.fetchAllCredentials(credentialManager.getAll(token));

    const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    this.fetchAllVaultsAndCredentials(token, collectionManager, credentialManager);
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    fetchAllCredentials: async function (getCredential) {
      const { contents } = await getCredential;
      this.numofCreds = Object.keys(contents).length;
    },
    fetchAllVaultsAndCredentials: async function (token, collectionManager, credentialManager) {
      // Fetch the list of the vaults/collections.
      const { contents } = await collectionManager.getAll(token);
      console.log(`found ${Object.keys(contents).length} vaults`);

      const vaultCred = Object.values(contents);
      vaultCred.forEach(async (vault) => {
        // Fetch the credentials stored inside the specific vault
        // TODO: #1236 Revisit the solution to avoid getting all the credentials
        await credentialManager.getAll(token, { collectionID: vault.id }).then((result) => {
          vault['numofCreds'] = Object.values(result.contents).length;
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
