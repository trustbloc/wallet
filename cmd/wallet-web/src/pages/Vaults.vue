<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="py-6 px-3">
    <div class="mb-8 w-full">
      <h3 class="text-neutrals-dark">{{ $t('Vaults.heading') }}</h3>
    </div>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 xl:gap-8">
      <vault-card
        color="pink"
        :num-of-creds="
          $tc('Vaults.foundCredentials', credentialsFound, {
            credentialLength: credentialsFound,
          })
        "
        :name="$t('Vaults.allVaults')"
      />
      <!--Issue-1198 Add flyout menu to default and other vaults -->
      <vault-card
        :num-of-creds="
          $tc('Vaults.foundCredentials', credentialsFound, {
            credentialLength: credentialsFound,
          })
        "
        :name="$t('Vaults.defaultVault')"
      />
      <!-- TODO: Issue-1194 Implement Add Vault -->
      <vault-card type="addNew" :name="$t('Vaults.addVault')" />
    </div>
  </div>
</template>

<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import useBreakpoints from '@/plugins/breakpoints';
import VaultCard from '@/components/VaultCard/VaultCard';

export default {
  name: 'Vaults',
  components: {
    VaultCard,
  },
  data() {
    return {
      credentialsFound: '',
    };
  },
  created: function () {
    const { user, token } = this.getCurrentUser().profile;
    this.username = this.getCurrentUser().username;
    const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.fetchAllCredentials(credentialManager.getAll(token));
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    fetchAllCredentials: async function (getCredential) {
      const { contents } = await getCredential;
      this.credentialsFound = Object.keys(contents).length;
    },
  },
};
</script>
