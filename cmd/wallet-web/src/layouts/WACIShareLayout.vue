<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <keep-alive>
    <component
      :is="component"
      class="relative z-10 flex h-full w-full grow flex-col items-start justify-start overflow-hidden bg-neutrals-softWhite"
    />
  </keep-alive>
</template>

<script>
import { reactive, ref, watch } from 'vue';
import WACISharePage from '@/pages/WACISharePage.vue';

export const WACIShareLayoutStore = reactive({
  component: WACISharePage,
});

export const WACIShareLayoutMutations = {
  setComponent(value) {
    WACIShareLayoutStore.component = value;
  },
};

export default {
  name: 'WACIShareLayout',
  setup() {
    const component = ref(WACIShareLayoutStore.component);
    watch(
      () => WACIShareLayoutStore.component,
      () => {
        component.value = WACIShareLayoutStore.component;
      }
    );
    return { component };
  },
};
</script>
