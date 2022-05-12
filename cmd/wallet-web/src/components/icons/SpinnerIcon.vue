<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <svg class="animate-spin" width="48" height="48" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <linearGradient id="icon-spinner-a" x1="45.116%" y1="0%" x2="45.116%" y2="100%">
        <stop stop-color="#8631A0" offset="0%" />
        <stop stop-color="#360B4C" offset="100%" />
      </linearGradient>
    </defs>
    <g transform="translate(3 3)" stroke-width="6" fill="none" fill-rule="evenodd">
      <circle stroke="#C7C3C8" cx="21" cy="21" r="21" />
      <path
        d="M42 21C42 9.402 32.598 0 21 0"
        stroke="url(#icon-spinner-a)"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </g>
  </svg>
</template>
