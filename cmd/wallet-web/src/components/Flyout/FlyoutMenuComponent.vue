<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div
    class="flyout-box absolute right-0 z-20 mt-2 w-48 rounded-lg bg-neutrals-white text-neutrals-medium shadow-xl"
  >
    <ul class="py-2">
      <slot />
    </ul>
  </div>
</template>

<script>
export default {
  name: 'FlyoutMenuComponent',
};
</script>
<style scoped>
.flyout-box {
  box-shadow: 0px 4px 20px 0px rgba(25, 12, 33, 0.2);
}
</style>
