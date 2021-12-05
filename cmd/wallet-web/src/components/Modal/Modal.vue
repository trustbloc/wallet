<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <teleport :to="target">
    <slot />
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
