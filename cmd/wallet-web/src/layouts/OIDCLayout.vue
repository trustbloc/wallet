<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="flex flex-col flex-grow justify-start items-center mx-auto max-w-7xl h-screen max-h-screen shadow-main-container"
  >
    <HeaderComponent :has-custom-gradient="true">
      <template v-if="selectedCredentialId" #leftButtonContainer>
        <button
          class="flex flex-row justify-start items-center focus:ring-2 focus:ring-primary-purple focus:ring-offset-2 outline-none"
          @click="handleBackButtonClick"
        >
          <div class="bg-neutrals-black rounded-full">
            <img
              class="z-10 w-6 h-6 transform rotate-180"
              src="@/assets/img/credential--arrow-right-icon-light.svg"
            />
          </div>
          <span class="px-3 text-base font-bold text-neutrals-white">{{ t('WACI.back') }}</span>
        </button>
      </template>
      <template #gradientContainer>
        <div class="absolute h-15 bg-gradient-full oval" />
      </template>
    </HeaderComponent>
    <router-view
      class="flex overflow-hidden relative z-10 flex-col flex-grow justify-start items-start w-full h-full bg-neutrals-softWhite"
    />

    <FooterComponent
      class="sticky bottom-0 z-20 bg-neutrals-magnolia border-t border-neutrals-thistle"
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
