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
    <div class="px-6 w-64 h-auto bg-neutrals-white rounded-xl">
      <div class="pt-5">
        <div class="flex justify-center items-center w-12 h-12 bg-gradient-pink rounded-full">
          <img class="w-6 h-5" src="@/assets/img/vaults.svg" alt="Vault Icon" />
        </div>
      </div>
      <div class="pb-4 space-y-1">
        <span class="block pt-4 text-lg font-bold"> {{ $t('Vaults.allVaults') }}</span>
        <span class="block text-sm text-neutrals-medium">
          {{
            $tc('Vaults.foundCredentials', credentialsFound, {
              credentialLength: credentialsFound,
            })
          }}
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import useBreakpoints from '@/plugins/breakpoints';

export default {
  name: 'Vaults',
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
