<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { computed, inject, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useStore } from 'vuex';
import { CHAPIHandler } from '@/mixins';
import useBreakpoints from '@/plugins/breakpoints.js';

// Local Variables
const disableCHAPI = ref(false);
const loaded = ref(false);

// Hooks
const route = useRoute();
const store = useStore();
const breakpoints = useBreakpoints();
const polyfill = inject('polyfill');
const webCredentialHandler = inject('webCredentialHandler');

// Store Getters
const currentUser = computed(() => store.getters['getCurrentUser']);
const agentOpts = computed(() => store.getters['getAgentOpts']);

// Store Actions
const activateCHAPI = () => store.dispatch('activateCHAPI');

onMounted(async () => {
  // if intended target doesn't require CHAPI.
  disableCHAPI.value = route.params.disableCHAPI || false;
  try {
    if (!breakpoints.xs && !breakpoints.sm && !disableCHAPI.value) {
      const chapi = new CHAPIHandler(
        polyfill,
        webCredentialHandler,
        agentOpts.value.credentialMediatorURL
      );
      await chapi.install(currentUser.value.username);
      activateCHAPI();
    }
  } catch (e) {
    console.log('Could not initialize Vue App:', e);
  }
  loaded.value = true;
});
</script>
<template>
  <router-view />
</template>
