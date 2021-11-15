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
        :num-of-creds="
          t('Vaults.foundCredentials', credentialsFound, {
            credentialLength: credentialsFound,
          })
        "
        :name="t('Vaults.allVaults')"
      />
      <!--TODO: Issue-1198 Add flyout menu to default and other vaults -->
      <vault-card
        :num-of-creds="
          t('Vaults.foundCredentials', credentialsFound, {
            credentialLength: credentialsFound,
          })
        "
        :name="t('Vaults.defaultVault')"
      />
      <vault-card type="addNew" :name="t('Vaults.addVault')" class="grid order-last" />
      <div v-for="(vault, index) in vaults" :key="index">
        <!-- TODO: Issue-1215 Add credentials found in the vault -->
        <vault-card :name="vault.name" class="grid order-last" />
      </div>
    </div>
  </div>
</template>

<script>
import { CredentialManager, CollectionManager } from '@trustbloc/wallet-sdk';
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
      credentialsFound: '',
      vaults: [],
    };
  },
  created: function () {
    const { user, token } = this.getCurrentUser().profile;
    this.username = this.getCurrentUser().username;
    const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.fetchAllCredentials(credentialManager.getAll(token));

    const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    this.fetchAllVaults(token, collectionManager);
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    fetchAllCredentials: async function (getCredential) {
      const { contents } = await getCredential;
      this.credentialsFound = Object.keys(contents).length;
    },
    fetchAllVaults: async function (token, collectionManager) {
      const { contents } = await collectionManager.getAll(token);
      console.log(`found ${Object.keys(contents).length} vaults`);
      this.vaults = Object.values(contents).map((vault) => vault);
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
