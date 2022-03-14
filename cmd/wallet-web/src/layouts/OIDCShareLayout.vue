<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <keep-alive>
    <component
      :is="component"
      class="
        overflow-hidden
        relative
        z-10
        flex-grow
        justify-start
        items-start
        w-full
        h-full
        flex flex-col
        bg-neutrals-softWhite
      "
    />
  </keep-alive>
</template>

<script>
import { reactive, ref, watch } from 'vue';
import OIDCShare from '@/pages/OIDCShare.vue';

export const OIDCShareLayoutStore = reactive({
  component: OIDCShare,
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
