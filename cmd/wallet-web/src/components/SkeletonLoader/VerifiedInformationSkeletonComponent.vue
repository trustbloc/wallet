<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex relative flex-row w-full sub-container">
    <div class="flex flex-col lg:flex-row justify-start lg:items-center w-full h-16">
      <div
        class="flex overflow-hidden relative mt-4 lg:my-8 ml-3 lg:ml-3 w-20 h-3 lg:h-5 bg-neutrals-white rounded-3xl skeleton-data"
      />
      <div
        class="flex overflow-hidden relative mb-4 mt-2 lg:my-8 ml-3 lg:mx-auto w-2/5 lg:w-48 h-4 lg:h-5 bg-neutrals-white rounded-3xl skeleton-data"
      />
    </div>
  </div>
</template>

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
