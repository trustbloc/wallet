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
    <VaultCardComponent v-for="id in vaultSkeletons" :key="`vault-skeleton-${id}`" />
  </div>
  <div v-else-if="type === 'CredentialDetailsBanner'" class="flex flex-row w-full">
    <CredentialDetailsBannerSkeletonComponent />
  </div>
  <div v-else-if="type === 'VerifiedInformation'">
    <table class="w-full border-t border-neutrals-chatelle">
      <tr
        v-for="id in rowSkeletons"
        :key="`row-skeleton-${id}`"
        class="border-b border-neutrals-thistle border-dotted"
      >
        <VerifiedInformationSkeletonComponent />
      </tr>
    </table>
  </div>
</template>

<script>
import VaultCardComponent from './VaultCardComponent';
import CredentialDetailsBannerSkeletonComponent from './CredentialDetailsBannerSkeletonComponent.vue';
import VerifiedInformationSkeletonComponent from './VerifiedInformationSkeletonComponent.vue';
export default {
  components: {
    VaultCardComponent,
    CredentialDetailsBannerSkeletonComponent,
    VerifiedInformationSkeletonComponent,
  },
  props: {
    type: {
      type: String,
      default: 'vault',
    },
  },
  data() {
    return {
      vaultSkeletons: [1, 2, 3],
      rowSkeletons: [1, 2, 3, 4],
    };
  },
};
</script>
