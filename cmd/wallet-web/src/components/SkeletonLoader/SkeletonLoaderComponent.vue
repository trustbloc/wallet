<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { ref } from 'vue';
import CredentialDetailsBannerSkeletonComponent from './CredentialDetailsBannerSkeletonComponent.vue';
import CredentialPreviewSkeletonComponent from './CredentialPreviewSkeletonComponent.vue';
import VaultCardSkeletonComponent from './VaultCardSkeletonComponent';
import VerifiedInformationSkeletonComponent from './VerifiedInformationSkeletonComponent.vue';

const props = defineProps({
  type: {
    type: String,
    default: 'vault',
  },
});
const numOfVaultSkeletons = ref(3);
const numOfRows = ref(4);
const numOfCredentialSkeletons = ref(2);
</script>

<template>
  <div
    v-if="type === 'Vault'"
    class="mt-6 flex flex-col flex-wrap gap-y-8 md:mt-8 lg:flex-row lg:gap-x-8"
  >
    <VaultCardSkeletonComponent v-for="id in numOfVaultSkeletons" :key="`vault-skeleton-${id}`" />
  </div>
  <div v-else-if="type === 'CredentialPreview'" class="flex w-full flex-col px-6 md:px-0">
    <div v-for="id in numOfCredentialSkeletons" :key="`cred-preview-skeleton-${id}`">
      <div
        class="relative mt-2 h-6 w-40 overflow-hidden rounded-3xl bg-neutrals-whiteLilac lg:h-5"
      />
      <div class="my-8 flex flex-col lg:flex-row">
        <CredentialPreviewSkeletonComponent class="mb-4 lg:mr-8 lg:mb-0" />
        <CredentialPreviewSkeletonComponent />
      </div>
    </div>
  </div>
  <div v-else-if="type === 'Flyout'" class="flex w-full px-6 md:w-40 md:px-0">
    <div
      class="relative mt-2 mb-3 h-6 w-full overflow-hidden rounded-3xl bg-neutrals-whiteLilac lg:mb-0"
    />
  </div>
  <div v-else-if="type === 'CredentialDetailsBanner'" class="flex w-full flex-row">
    <CredentialDetailsBannerSkeletonComponent />
  </div>
  <div v-else-if="type === 'VerifiedInformation'">
    <table class="w-full border-t border-neutrals-chatelle">
      <tr
        v-for="id in numOfRows"
        :key="`row-skeleton-${id}`"
        class="border-b border-dotted border-neutrals-thistle"
      >
        <VerifiedInformationSkeletonComponent />
      </tr>
    </table>
  </div>
</template>

<style>
.skeleton-data::before,
.skeleton-data::after {
  content: '';
  position: absolute;
  left: 0;
  width: 100%;
  height: 1px;
  background-image: linear-gradient(
    -90deg,
    rgba(255, 255, 255, 0) 15%,
    #867c8c 50%,
    rgba(255, 255, 255, 0) 85%
  );
  animation: shimmer 1.25s linear infinite;
}
.skeleton-data::before {
  top: 0;
}
.skeleton-data::after {
  bottom: 0;
}
@keyframes shimmer {
  from {
    transform: translateX(-100%);
  }
  to {
    transform: translateX(250%);
  }
}
</style>
