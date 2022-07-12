<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { onBeforeMount, onMounted, onUnmounted, ref } from 'vue';
import { useStore } from 'vuex';
import getStartingLocale from '@/mixins/i18n/getStartingLocale.js';
import { updateI18nLocale } from '@/plugins/i18n';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

// Local Variables
const startingLocale = getStartingLocale(); // Get starting locale, set it in i18n and in the store
const loaded = ref(false);

// Hooks
const store = useStore();

store.dispatch('setLocale', startingLocale);

onBeforeMount(async () => {
  await updateI18nLocale(startingLocale.id);
});

onMounted(() => {
  loaded.value = true;
});

onUnmounted(() => {
  if (store) {
    store.dispatch('agent/flushStore');
    console.log('store flushed', store);
  }
});
</script>
<template>
  <div class="w-screen h-screen font-sans bg-neutrals-mischka">
    <div v-if="!loaded" class="flex relative top-1/3 flex-col justify-start items-center">
      <SpinnerIcon />
      <span class="mt-6">Loading Wallet . . .</span>
    </div>
    <router-view v-else />
  </div>
</template>
