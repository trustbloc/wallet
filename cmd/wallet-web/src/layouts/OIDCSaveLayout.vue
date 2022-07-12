<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <keep-alive>
    <component
      :is="component"
      class="flex overflow-hidden relative z-10 flex-col grow justify-start items-start w-full h-full bg-neutrals-softWhite"
    />
  </keep-alive>
</template>

<script>
import { reactive, ref } from 'vue';
import OIDCSavePage from '@/pages/OIDCSavePage.vue';

export const OIDCSaveLayoutStore = reactive({
  component: OIDCSavePage,
});

export default {
  name: 'OIDCSaveLayout',
  setup() {
    const component = ref(OIDCSaveLayoutStore.component);
    return { component };
  },
};
</script>
