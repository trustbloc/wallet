<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 * Copyright Avast Software. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { OpenID4VCShareLayoutMutations } from '@/layouts/OpenID4VCShareLayout.vue';
import OpenID4VCSharePage from '@/pages/OpenID4VCSharePage.vue';
import HeaderComponent from '@/components/Header/HeaderComponent.vue';
import FooterComponent from '@/components/Footer/FooterComponent.vue';

// Hooks
const { t } = useI18n();

// Local Variables
const selectedCredentialId = ref(OpenID4VCStore.selectedCredentialId);

// Watchers
watch(
  () => OpenID4VCStore.selectedCredentialId,
  () => {
    selectedCredentialId.value = OpenID4VCStore.selectedCredentialId;
  }
);

// Methods
function handleBackButtonClick() {
  OpenID4VCShareLayoutMutations.setComponent(OpenID4VCSharePage);
  OpenID4VCMutations.setSelectedCredentialId(null);
}
</script>

<template>
  <div
    class="mx-auto flex h-screen max-h-screen max-w-7xl grow flex-col items-center justify-start shadow-main-container"
  >
    <HeaderComponent :has-custom-gradient="true">
      <template v-if="selectedCredentialId" #leftButtonContainer>
        <button
          class="flex flex-row items-center justify-start outline-none focus:ring-2 focus:ring-primary-purple focus:ring-offset-2"
          @click="handleBackButtonClick"
        >
          <div class="rounded-full bg-neutrals-black">
            <img
              class="z-10 h-6 w-6 rotate-180"
              src="@/assets/img/credential--arrow-right-icon-light.svg"
            />
          </div>
          <span class="px-3 text-base font-bold text-neutrals-white">{{ t('WACI.back') }}</span>
        </button>
      </template>
      <template #gradientContainer>
        <div class="oval absolute h-15 bg-gradient-full" />
      </template>
    </HeaderComponent>
    <router-view
      class="relative z-10 flex h-full w-full grow flex-col items-start justify-start overflow-hidden bg-neutrals-softWhite"
    />

    <FooterComponent
      class="sticky bottom-0 z-20 border-t border-neutrals-thistle bg-neutrals-magnolia"
    />
  </div>
</template>

<script>
export const OpenID4VCStore = reactive({
  processedCredentials: [],
  selectedCredentialId: null,
});

export const OpenID4VCGetters = {
  getProcessedCredentialById(id) {
    return OpenID4VCStore.processedCredentials.find((credential) => credential.id === id);
  },
};

export const OpenID4VCMutations = {
  setProcessedCredentials(value) {
    OpenID4VCStore.processedCredentials = value;
  },
  setSelectedCredentialId(value) {
    OpenID4VCStore.selectedCredentialId = value;
  },
};
</script>

<style>
.oval {
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  filter: blur(50px);
  width: 15.625rem; /* 250px */
  top: 2.0625rem; /* 33px */
}
</style>
