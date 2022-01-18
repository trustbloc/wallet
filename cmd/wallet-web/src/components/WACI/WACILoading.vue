<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <Spinner />
  <span v-if="message.length" class="mt-8 text-base md:text-lg text-neutrals-dark">{{
    message
  }}</span>
</template>

<script>
import Spinner from '@/components/Spinner/Spinner.vue';

export default {
  name: 'WACILoading',
  components: {
    Spinner,
  },
  props: {
    message: {
      type: String,
      default: '',
    },
  },
};
</script>
