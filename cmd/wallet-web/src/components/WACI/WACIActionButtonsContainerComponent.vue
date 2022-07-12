<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="flex sticky bottom-0 z-20 flex-row justify-center items-center py-4 px-5 w-full bg-neutrals-magnolia border-t border-neutrals-thistle"
  >
    <div class="flex flex-row grow justify-between items-center w-full max-w-3xl">
      <slot name="leftButton" />
      <slot name="rightButton" />
    </div>
  </div>
</template>

<script>
export default {
  name: 'WACIActionButtonsContainerComponent',
};
</script>
