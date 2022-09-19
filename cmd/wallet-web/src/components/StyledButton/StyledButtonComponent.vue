<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <button
    :class="type"
    class="flex h-11 flex-row items-center justify-center rounded-lg px-8 text-base font-bold"
  >
    <svg
      v-if="getLoading"
      class="absolute animate-spin"
      width="24"
      height="24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <g transform="translate(2 2)" stroke="#FFF" stroke-width="3" fill="none" fill-rule="evenodd">
        <circle opacity=".3" cx="10" cy="10" r="10" />
        <path d="M20 10c0-5.523-4.477-10-10-10" stroke-linecap="round" stroke-linejoin="round" />
      </g>
    </svg>
    <div :class="getLoading && 'invisible'">
      <slot />
    </div>
  </button>
</template>

<script>
const BUTTON_TYPES = ['btn-primary', 'btn-inverse', 'btn-outline', 'btn-danger'];

export default {
  name: 'StyledButtonComponent',
  props: {
    loading: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      required: true,
      validator: function (value) {
        return BUTTON_TYPES.indexOf(value) !== -1;
      },
    },
  },
  computed: {
    getLoading() {
      return this.loading;
    },
  },
};
</script>

<style scoped></style>
