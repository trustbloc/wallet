<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
const props = defineProps({
  isTall: {
    type: Boolean,
    required: false,
    default: false,
  },
});
</script>

<template>
  <div class="flex w-full flex-col items-start justify-start px-0">
    <div
      class="z-10 flex w-full flex-row items-center justify-between rounded-xl bg-neutrals-whiteLilac pr-8 pl-5 lg:h-24 lg:py-6"
      :class="isTall ? 'h-24 py-6' : 'h-20 py-5'"
    >
      <div class="flex w-full flex-row items-center">
        <div
          class="skeleton-data relative w-12 flex-none overflow-hidden rounded-lg border-opacity-10 bg-neutrals-softWhite lg:h-12"
          :class="isTall ? 'h-12' : 'h-10'"
        />
        <div
          class="skeleton-data relative mt-0 ml-4 w-40 overflow-hidden rounded-3xl bg-neutrals-softWhite lg:h-5"
          :class="isTall ? 'h-6' : 'h-5'"
        />
      </div>
      <div
        class="skeleton-data relative h-8 w-8 overflow-hidden rounded-full bg-neutrals-softWhite"
      />
    </div>
  </div>
</template>
