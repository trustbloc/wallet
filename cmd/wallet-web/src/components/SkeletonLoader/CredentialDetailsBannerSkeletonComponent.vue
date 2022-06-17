<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start px-3 md:px-0 pt-10 w-full items-start">
    <div
      class="flex z-10 flex-row justify-between items-center py-6 pl-5 pr-8 w-full h-24 bg-neutrals-whiteLilac rounded-xl"
    >
      <div class="flex flex-row items-center w-full">
        <div
          class="overflow-hidden relative flex-none w-12 h-12 rounded-lg border-opacity-10 bg-neutrals-softWhite skeleton-data"
        />
        <div
          class="overflow-hidden relative mt-0 ml-4 w-40 h-6 lg:h-5 bg-neutrals-softWhite rounded-3xl skeleton-data"
        />
      </div>
      <div
        class="overflow-hidden relative w-8 h-8 bg-neutrals-softWhite rounded-full skeleton-data"
      />
    </div>
    <div
      class="flex relative flex-col lg:flex-row px-6 pt-12 lg:pt-14 pb-4 lg:pb-6 mb-1 w-full bg-neutrals-white rounded-b-xl sub-container"
    >
      <div
        v-for="id in skeletons"
        :key="`cred-details-skeleton-${id}`"
        class="flex-column justify-start items-center lg:w-1/2"
      >
        <div
          class="overflow-hidden relative mt-2 w-20 h-4 bg-neutrals-softWhite rounded-3xl skeleton-data"
        />
        <div
          class="overflow-hidden relative mt-2 w-2/5 lg:w-32 h-5 bg-neutrals-softWhite rounded-3xl skeleton-data"
        />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      skeletons: [1, 2, 3, 4],
    };
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
.sub-container {
  top: -2.5rem;
  left: 0;
  box-shadow: 0px 2px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
