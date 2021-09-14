<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <!-- Todo Add Credential Card and Credential Details skeleton here -->
  <div
    v-if="type === 'vault'"
    class="flex flex-col lg:flex-row flex-wrap lg:gap-x-8 gap-y-8 px-6 md:px-0"
  >
    <vault-card v-for="id in skeletons" :key="`vault-skeleton-${id}`" />
  </div>
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
  data() {
    return {
      skeletons: [1, 2, 3],
    };
  },
};
</script>
