<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 * Copyright Avast Software. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script>
import { reactive } from 'vue';
import OpenID4VCSharePage from '@/pages/OpenID4VCSharePage.vue';

export const OpenID4VCShareLayoutStore = reactive({
  component: OpenID4VCSharePage,
});

export const OpenID4VCShareLayoutMutations = {
  setComponent(value) {
    OpenID4VCShareLayoutStore.component = value;
  },
};
</script>

<script setup>
import { ref, watch } from 'vue';

const component = ref(OpenID4VCShareLayoutStore.component);

watch(
  () => OpenID4VCShareLayoutStore.component,
  () => {
    component.value = OpenID4VCShareLayoutStore.component;
  }
);
</script>

<template>
  <keep-alive>
    <component
      :is="component"
      class="relative z-10 flex h-full w-full grow flex-col items-start justify-start overflow-hidden bg-neutrals-softWhite"
    />
  </keep-alive>
</template>
