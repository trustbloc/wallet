<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="flex sticky bottom-0 z-20 flex-row justify-center items-center pt-4 pr-5 pb-4 pl-5 w-full bg-neutrals-magnolia border-t border-neutrals-thistle"
  >
    <div class="flex flex-row flex-grow justify-between items-center w-full max-w-3xl">
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
