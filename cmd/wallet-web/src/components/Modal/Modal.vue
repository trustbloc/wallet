<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <teleport :to="target">
    <div
        class="
        flex
        overflow-y-auto
        fixed
        inset-0
        z-50
        justify-center
        items-center
        bg-neutrals-dark bg-opacity-50
      "
    >
      <div class="relative mx-6 lg:mx-auto flex-grow max-w-md bg-neutrals-white rounded-2xl modal-width">
      <slot />
      </div>
    </div>
  </teleport>
</template>

<script>
export default {
  name: 'Modal',
  props: {
    target: {
      type: String,
      default: 'body',
    },
  },
};
</script>
