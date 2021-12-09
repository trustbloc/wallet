<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <button
    :id="id"
    class="flex justify-start items-center py-2 px-4 w-full text-base font-bold"
    @click="$emit('click', $event)"
  >
    {{ text }}
  </button>
</template>

<script>
export default {
  name: 'FlyoutButton',
  props: {
    text: {
      type: String,
      required: true,
    },
    id: {
      type: String,
      required: true,
    },
  },
  emits: ['click'],
};
</script>
