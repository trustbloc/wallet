<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { computed, onBeforeMount } from 'vue';
import { useStore } from 'vuex';
import * as polyfill from 'credential-handler-polyfill';
import * as webCredentialHandler from 'web-credential-handler';

// Hooks
const store = useStore();
const opts = computed(() => store.getters['getAgentOpts']);

onBeforeMount(async () => {
  try {
    await polyfill.loadOnce(opts.value ? opts.value.credentialMediatorURL : undefined);
  } catch (e) {
    console.error('Error in loadOnce:', e);
  }
  return webCredentialHandler.activateHandler({
    mediatorOrigin: opts.value.credentialMediatorURL,
    get(event) {
      console.debug('Received get() event:', event.event);
      return { type: 'redirect', url: `${__webpack_public_path__}GetFromWallet` };
    },
    store(event) {
      console.debug('Received store() event:', event.event);
      return { type: 'redirect', url: `${__webpack_public_path__}StoreInWallet` };
    },
  });
});
</script>

<template>
  <div></div>
</template>
