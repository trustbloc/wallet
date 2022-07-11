<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { ref } from 'vue';
import CredentialPreviewSkeletonComponent from './CredentialPreviewSkeletonComponent.vue';
const numOfSkeletons = ref(4);
</script>
<template>
  <div class="flex flex-col justify-start items-start px-3 pt-10 w-full md:px-0">
    <CredentialPreviewSkeletonComponent :is-tall="true" />
    <div
      class="flex relative flex-col px-6 pt-12 pb-4 mb-1 w-full bg-neutrals-white rounded-b-xl lg:flex-row lg:pt-14 lg:pb-6 sub-container"
    >
      <div
        v-for="id in numOfSkeletons"
        :key="`cred-details-skeleton-${id}`"
        class="items-center justify-start lg:w-1/2"
      >
        <div
          class="skeleton-data relative mt-2 h-4 w-20 overflow-hidden rounded-3xl bg-neutrals-softWhite"
        />
        <div
          class="overflow-hidden relative mt-2 w-2/5 h-5 bg-neutrals-softWhite rounded-3xl lg:w-32 skeleton-data"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.sub-container {
  top: -2.5rem;
  left: 0;
  box-shadow: 0px 2px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
