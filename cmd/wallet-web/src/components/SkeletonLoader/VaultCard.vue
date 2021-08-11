<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div class="flex flex-col md:flex-row px-6 space-y-5 md:space-y-0 md:space-x-8">
    <div class="skeleton">
      <div class="skeleton-photo"></div>
      <div class="space-y-3 col-1">
        <div class="w-48 h-5 md:h-6 skeleton-data top"></div>
        <div class="w-32 h-3 md:h-4 skeleton-data bottom"></div>
      </div>
    </div>
    <div class="skeleton">
      <div class="skeleton-photo"></div>
      <div class="space-y-3 col-1">
        <div class="w-48 h-5 md:h-6 skeleton-data top"></div>
        <div class="w-32 h-3 md:h-4 skeleton-data bottom"></div>
      </div>
    </div>
    <div class="skeleton">
      <div class="skeleton-photo"></div>
      <div class="space-y-3 col-1">
        <div class="w-48 h-5 md:h-6 skeleton-data top"></div>
        <div class="w-32 h-3 md:h-4 skeleton-data bottom"></div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    type: {
      type: String,
      default: 'vault',
    },
  },
};
</script>
