<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col items-center justify-start py-16 px-5">
    <img src="@/assets/img/icons-circle-check.svg" />
    <span class="mt-5 mb-3 text-center text-xl font-bold text-neutrals-dark">{{ heading }}</span>
    <span class="text-center text-lg text-neutrals-medium">{{ message }}</span>
    <StyledButtonComponent v-bind="$attrs" type="btn-primary" class="mt-6">
      {{ buttonLabel }}
    </StyledButtonComponent>
  </div>
</template>

<script>
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';

export default {
  name: 'WACISuccessComponent',
  components: {
    StyledButtonComponent,
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
