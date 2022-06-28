<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { computed, onBeforeMount, onMounted, onUnmounted, ref } from 'vue';
import { useStore } from 'vuex';
import getStartingLocale from '@/mixins/i18n/getStartingLocale.js';
import { updateI18nLocale } from '@/plugins/i18n';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

const store = useStore();
const loaded = ref(false);
// const isAgentInitialized = computed(() => store.getters['agent/isInitialized']);
// const initAgent = () => store.dispatch('agent/init');
// const initOpts = () => store.dispatch('initOpts');
// const loadUser = () => store.dispatch('loadUser');

// Get starting locale, set it in i18n and in the store
const startingLocale = getStartingLocale();
store.dispatch('setLocale', startingLocale);

onBeforeMount(async () => {
  await updateI18nLocale(startingLocale.id);
});

onMounted(() => {
  try {
    // load opts
    // await initOpts();
    // load user if already logged in
    // loadUser();
    // load agent if user already logged in and agent not initialized (scenario: page refresh)
    // if (store.getters.getCurrentUser && !isAgentInitialized.value) {
    //   await initAgent();
    // }
  } catch (e) {
    console.log('Could not initialize Vue App:', e);
  }
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
