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
  <div class="flex w-full flex-col items-start justify-start px-3 pt-10 md:px-0">
    <CredentialPreviewSkeletonComponent :is-tall="true" />
    <div
      class="sub-container relative mb-1 flex w-full flex-col rounded-b-xl bg-neutrals-white px-6 pt-12 pb-4 lg:flex-row lg:pt-14 lg:pb-6"
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
          class="skeleton-data relative mt-2 h-5 w-2/5 overflow-hidden rounded-3xl bg-neutrals-softWhite lg:w-32"
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
