<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <SpinnerIcon />
  <span v-if="message.length" class="mt-8 text-base text-neutrals-dark md:text-lg">{{
    message
  }}</span>
</template>

<script>
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

export default {
  name: 'WACILoadingComponent',
  components: {
    SpinnerIcon,
  },
  props: {
    message: {
      type: String,
      default: '',
    },
  },
};
</script>
