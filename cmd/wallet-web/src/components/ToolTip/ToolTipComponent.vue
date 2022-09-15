<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="absolute top-14 right-0 z-20 inline-flex justify-center rounded-lg bg-neutrals-dark py-1 px-5"
  >
    <div
      class="triangle absolute z-10 rotate-45 rounded-sm border-8 border-neutrals-dark bg-neutrals-dark"
    />
    <span class="z-20 whitespace-nowrap text-center text-sm text-neutrals-white">{{
      toolTipLabel
    }}</span>
  </div>
</template>

<script>
export default {
  name: 'ToolTipComponent',
  props: {
    toolTipLabel: {
      type: String,
      default: 'default',
    },
  },
};
</script>

<style scoped>
.triangle {
  top: -4px;
  right: 20px;
}
</style>
