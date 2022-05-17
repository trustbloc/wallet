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
    <VaultCardComponent v-for="id in skeletons" :key="`vault-skeleton-${id}`" />
  </div>
</template>

<script>
import VaultCardComponent from './VaultCardComponent';
export default {
  components: {
    VaultCardComponent,
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
