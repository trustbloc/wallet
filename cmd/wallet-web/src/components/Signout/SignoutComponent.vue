<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { computed, inject, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { CHAPIHandler } from '@/mixins';
import { useI18n } from 'vue-i18n';

// Hooks
const router = useRouter();
const store = useStore();
const { t } = useI18n();
const polyfill = inject('polyfill');
const webCredentialHandler = inject('webCredentialHandler');

// Store Getters
const agentOpts = computed(() => store.getters['getAgentOpts']);
const isCHAPI = computed(() => store.getters['isCHAPI']);

// Store Actions
const signoutUser = () => store.dispatch('logout');

// Methods
async function signout() {
  const actions = [signoutUser()];
  if (isCHAPI.value) {
    actions.push(chapi.value.uninstall());
  }

  await Promise.all(actions);

  router.push({ name: 'signin' });
}

const chapi = ref(
  new CHAPIHandler(polyfill, webCredentialHandler, agentOpts.value.credentialMediatorURL)
);
</script>

<template>
  <div
    class="
      justify-start
      items-center
      focus-within:bg-gradient-to-r
      opacity-60
      focus-within:opacity-100
      hover:from-neutrals-black hover:bg-gradient-to-r hover:opacity-100
      focus-within:from-neutrals-black focus-within:shadow-inner-outline-blue
      flex flex-row
      bar
    "
  >
    <button
      id="signout-button"
      class="flex flex-row justify-start items-center px-10 w-full h-16 focus:outline-none"
      type="button"
      @click="signout()"
      @keyup.enter="signout()"
    >
      <img src="@/assets/img/signout.svg" />
      <span class="ml-4 text-lg font-bold text-neutrals-white">{{ t('Signout.signout') }}</span>
    </button>
  </div>
</template>

<style scoped>
.bar:not(:focus-within):hover:before {
  position: absolute;
  content: '';
  display: block;
  background-color: theme('colors.primary.purple.hashita');
  height: theme('spacing.16');
  width: 4px;
}
</style>