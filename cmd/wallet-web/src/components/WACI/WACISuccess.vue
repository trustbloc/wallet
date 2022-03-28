<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
    <img src="@/assets/img/icons-circle-check.svg" />
    <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{ heading }}</span>
    <span class="text-lg text-center text-neutrals-medium">{{ message }}</span>
    <styled-button v-bind="$attrs" type="btn-primary" class="mt-6">
      {{ buttonLabel }}
    </styled-button>
  </div>
</template>

<script>
import StyledButton from '@/components/StyledButton/StyledButton.vue';

export default {
  name: 'WACISuccess',
  components: {
    StyledButton,
  },
  props: {
    heading: {
      type: String,
      required: true,
    },
    message: {
      type: String,
      required: true,
    },
    buttonLabel: {
      type: String,
      required: true,
    },
  },
};
</script>
