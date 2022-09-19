<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { onBeforeMount, onUnmounted, computed } from 'vue';
import { useStore } from 'vuex';
import getStartingLocale from '@/mixins/i18n/getStartingLocale.js';
import { updateI18nLocale } from '@/plugins/i18n';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

// Local Variables
const startingLocale = getStartingLocale(); // Get starting locale, set it in i18n and in the store

// Hooks
const store = useStore();

// Store Getters
const userLoaded = computed(() => store.getters['isUserLoaded']);

store.dispatch('setLocale', startingLocale);

onBeforeMount(async () => {
  await updateI18nLocale(startingLocale.id);
});

onUnmounted(() => {
  if (store) {
    store.dispatch('agent/flushStore');
    console.log('store flushed', store);
  }
});
</script>
<template>
  <div class="h-screen w-screen bg-neutrals-mischka font-sans">
    <div v-if="!userLoaded" class="relative top-1/3 flex flex-col items-center justify-start">
      <SpinnerIcon />
      <span class="mt-6">Loading Wallet . . .</span>
    </div>
    <router-view v-else />
  </div>
</template>
