<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

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
import { reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { OIDCShareLayoutMutations } from '@/layouts/OIDCShareLayout.vue';
import OIDCSharePage from '@/pages/OIDCSharePage.vue';
import HeaderComponent from '@/components/Header/HeaderComponent.vue';
import FooterComponent from '@/components/Footer/FooterComponent.vue';

export const OIDCStore = reactive({
  processedCredentials: [],
  selectedCredentialId: null,
});

export const OIDCGetters = {
  getProcessedCredentialById(id) {
    return OIDCStore.processedCredentials.find((credential) => credential.id === id);
  },
};

export const OIDCMutations = {
  setProcessedCredentials(value) {
    OIDCStore.processedCredentials = value;
  },
  setSelectedCredentialId(value) {
    OIDCStore.selectedCredentialId = value;
  },
};

export default {
  components: {
    HeaderComponent,
    FooterComponent,
  },
  setup() {
    const { t } = useI18n();
    const selectedCredentialId = ref(OIDCStore.selectedCredentialId);
    watch(
      () => OIDCStore.selectedCredentialId,
      () => {
        selectedCredentialId.value = OIDCStore.selectedCredentialId;
      }
    );
    return { selectedCredentialId, t };
  },
  methods: {
    handleBackButtonClick() {
      OIDCShareLayoutMutations.setComponent(OIDCSharePage);
      OIDCMutations.setSelectedCredentialId(null);
    },
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
