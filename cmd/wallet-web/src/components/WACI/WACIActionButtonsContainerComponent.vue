<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="sticky bottom-0 z-20 flex w-full flex-row items-center justify-center border-t border-neutrals-thistle bg-neutrals-magnolia py-4 px-5"
  >
    <div class="flex w-full max-w-3xl grow flex-row items-center justify-between">
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
