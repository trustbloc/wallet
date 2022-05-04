<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <IconSpinner />
  <span v-if="message.length" class="mt-8 text-base md:text-lg text-neutrals-dark">{{
    message
  }}</span>
</template>

<script>
import IconSpinner from '@/components/icons/IconSpinner.vue';

export default {
  name: 'WACILoading',
  components: {
    IconSpinner,
  },
  props: {
    message: {
      type: String,
      default: '',
    },
  },
};
</script>
