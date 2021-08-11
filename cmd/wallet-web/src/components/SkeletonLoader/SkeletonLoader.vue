<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <!-- Todo Add Credential Card and Credential Details skeleton here -->
  <vault-card v-if="type === 'vault'" />
</template>

<script>
import VaultCard from './VaultCard';
export default {
  components: {
    VaultCard,
  },
  props: {
    type: {
      type: String,
      default: 'vault',
    },
  },
};
</script>
