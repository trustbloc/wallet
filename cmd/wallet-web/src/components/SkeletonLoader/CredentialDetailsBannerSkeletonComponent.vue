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
  <div class="flex flex-col justify-start items-start px-3 md:px-0 pt-10 w-full">
    <CredentialPreviewSkeletonComponent :isTall="true" />
    <div
      class="flex relative flex-col lg:flex-row px-6 pt-12 lg:pt-14 pb-4 lg:pb-6 mb-1 w-full bg-neutrals-white rounded-b-xl sub-container"
    >
      <div
        v-for="id in numOfSkeletons"
        :key="`cred-details-skeleton-${id}`"
        class="justify-start items-center lg:w-1/2"
      >
        <div
          class="overflow-hidden relative mt-2 w-20 h-4 bg-neutrals-softWhite rounded-3xl skeleton-data"
        />
        <div
          class="overflow-hidden relative mt-2 w-2/5 lg:w-32 h-5 bg-neutrals-softWhite rounded-3xl skeleton-data"
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
