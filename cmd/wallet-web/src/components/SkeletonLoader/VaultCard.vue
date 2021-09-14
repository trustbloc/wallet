<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div
    class="
      lg:flex-col
      p-4
      lg:p-6 lg:pt-5
      pr-14
      w-full
      lg:w-64
      bg-neutrals-whiteLilac
      rounded-xl
      flex flex-row
    "
  >
    <div class="mr-3 w-12 h-12 bg-neutrals-softWhite rounded-full"></div>
    <div class="flex-1 space-y-3">
      <div
        class="
          overflow-hidden
          relative
          mt-0
          ml-3
          lg:mt-3 lg:ml-0
          w-full
          h-5
          lg:h-6
          bg-neutrals-softWhite
          rounded-3xl
          skeleton-data
        "
      ></div>
      <div
        class="
          overflow-hidden
          relative
          mt-0
          ml-3
          lg:mt-3 lg:ml-0
          w-1/2
          h-3
          lg:h-4
          bg-neutrals-softWhite
          rounded-3xl
          skeleton-data
        "
      ></div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    type: {
      type: String,
      default: 'vault',
    },
  },
};
</script>

<style scoped>
.skeleton-data::before,
.skeleton-data::after {
  content: '';
  position: absolute;
  left: 0;
  width: 100%;
  height: 1px;
  background-image: linear-gradient(
    -90deg,
    rgba(255, 255, 255, 0) 15%,
    #867c8c 50%,
    rgba(255, 255, 255, 0) 85%
  );
  animation: shimmer 1.25s linear infinite;
}
.skeleton-data::before {
  top: 0;
}
.skeleton-data::after {
  bottom: 0;
}
@keyframes shimmer {
  from {
    transform: translateX(-100%);
  }
  to {
    transform: translateX(250%);
  }
}
</style>
