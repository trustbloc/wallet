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
import OIDCSharePage from '@/pages/OIDCSharePage.vue';

export const OIDCShareLayoutStore = reactive({
  component: OIDCSharePage,
});

export const OIDCShareLayoutMutations = {
  setComponent(value) {
    OIDCShareLayoutStore.component = value;
  },
};

export default {
  name: 'OIDCShareLayout',
  setup() {
    const component = ref(OIDCShareLayoutStore.component);
    watch(
      () => OIDCShareLayoutStore.component,
      () => {
        component.value = OIDCShareLayoutStore.component;
      }
    );
    return { component };
  },
};
</script>
