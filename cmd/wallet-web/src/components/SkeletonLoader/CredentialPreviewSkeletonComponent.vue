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
  <div class="flex flex-col justify-start items-start px-0 w-full">
    <div
      class="flex z-10 flex-row justify-between items-center pr-8 pl-5 w-full bg-neutrals-whiteLilac rounded-xl lg:py-6 lg:h-24"
      :class="isTall ? 'h-24 py-6' : 'h-20 py-5'"
    >
      <div class="flex flex-row items-center w-full">
        <div
          class="overflow-hidden relative flex-none w-12 bg-neutrals-softWhite rounded-lg border-opacity-10 lg:h-12 skeleton-data"
          :class="isTall ? 'h-12' : 'h-10'"
        />
        <div
          class="overflow-hidden relative mt-0 ml-4 w-40 bg-neutrals-softWhite rounded-3xl lg:h-5 skeleton-data"
          :class="isTall ? 'h-6' : 'h-5'"
        />
      </div>
      <div
        class="overflow-hidden relative w-8 h-8 bg-neutrals-softWhite rounded-full skeleton-data"
      />
    </div>
  </div>
</template>
