<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <button :id="id" :class="[`block font-bold text-${color}`]" @click="$emit('click', $event)">
    {{ text }}
    <slot />
  </button>
</template>

<script>
export default {
  name: 'FlyoutMenuButton',
  props: {
    text: {
      type: String,
      required: true,
    },
    color: {
      type: String,
      default: 'neutrals-medium',
    },
    id: {
      type: String,
      required: true,
    },
  },
  emits: ['click'],
};
</script>
