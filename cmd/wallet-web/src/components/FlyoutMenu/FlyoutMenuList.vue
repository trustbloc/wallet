<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div
    class="
      absolute
      right-0
      z-20
      mt-2
      w-max
      text-neutrals-medium
      bg-neutrals-white
      rounded-lg
      shadow-xl
      flyout-box
    "
  >
    <ul class="p-4 space-y-2 text-base">
      <slot />
    </ul>
  </div>
</template>

<script>
export default {
  name: 'FlyoutMenuList',
};
</script>
<style scoped>
.flyout-box {
  box-shadow: 0px 4px 20px 0px rgba(25, 12, 33, 0.2);
}
</style>
