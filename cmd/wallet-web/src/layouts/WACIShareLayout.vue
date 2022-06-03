<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <keep-alive>
    <component
      :is="component"
      class="flex overflow-hidden relative z-10 flex-col flex-grow justify-start items-start w-full h-full bg-neutrals-softWhite"
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
