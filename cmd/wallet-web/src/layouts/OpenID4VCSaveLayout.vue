<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 * Copyright Avast Software. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script>
import { reactive } from 'vue';
import OpenID4VCSavePage from '@/pages/OpenID4VCSavePage.vue';

export const OpenID4VCSaveLayoutStore = reactive({
  component: OpenID4VCSavePage,
});
</script>

<script setup>
import { ref } from 'vue';

const component = ref(OpenID4VCSaveLayoutStore.component);
</script>

<template>
  <keep-alive>
    <component
      :is="component"
      class="relative z-10 flex h-full w-full grow flex-col items-start justify-start overflow-hidden bg-neutrals-softWhite"
    />
  </keep-alive>
</template>
