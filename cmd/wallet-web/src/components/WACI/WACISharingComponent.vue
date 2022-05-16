<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <SpinnerIcon />
  <span v-if="message.length" class="mt-8 text-base md:text-lg text-neutrals-dark">{{
    message
  }}</span>
</template>

<script>
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

export default {
  name: 'WACISharingComponent',
  components: {
    SpinnerIcon,
  },
  props: {
    message: {
      type: String,
      required: true,
    },
  },
};
</script>
