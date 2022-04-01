<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="
      flex-grow
      justify-start
      items-center
      mx-auto
      max-w-7xl
      h-screen
      max-h-screen
      shadow-main-container
      flex flex-col
    "
  >
    <Header :has-custom-gradient="true">
      <template v-if="selectedCredentialId" #leftButtonContainer>
        <button
          class="
            justify-start
            items-center
            focus:ring-2 focus:ring-primary-purple focus:ring-offset-2
            outline-none
            flex flex-row
          "
          @click="handleBackButtonClick"
        >
          <div class="rounded-full bg-neutrals-black">
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
    </Header>
    <keep-alive>
      <component
        :is="component"
        class="
          overflow-hidden
          relative
          z-10
          flex-grow
          justify-start
          items-start
          w-full
          h-full
          flex flex-col
          bg-neutrals-softWhite
        "
      />
    </keep-alive>
    <Footer class="sticky bottom-0 z-20 border-t border-neutrals-thistle bg-neutrals-magnolia" />
  </div>
</template>

<script>
import { reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { OIDCShareLayoutMutations } from '@/layouts/OIDCShareLayout.vue';
import OIDCShareLayout from '@/layouts/OIDCShareLayout.vue';
import OIDCShare from '@/pages/OIDCShare.vue';
import OIDCSaveLayout from '@/layouts/OIDCSaveLayout.vue';
import Header from '@/components/Header/Header.vue';
import Footer from '@/components/Footer/Footer.vue';
import { sendCredentialAuthorizeRequest, readOpenIDConfiguration } from '@/mixins';

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
    Header,
    Footer,
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
  data() {
    return {
      component: null,
    };
  },
  created: function () {
    this.decideFlow(this.$route.path);
  },
  methods: {
    handleBackButtonClick() {
      OIDCShareLayoutMutations.setComponent(OIDCShare);
      OIDCMutations.setSelectedCredentialId(null);
    },
    async decideFlow(path) {
      if (path === '/oidc/initiate') {
        const configuration = await readOpenIDConfiguration(this.$route.query.issuer);
        sendCredentialAuthorizeRequest(
          configuration,
          this.$route.query,
          `${location.protocol}//${location.host}/oidc/save`
        );
      } else if (path === '/oidc/save') {
        this.component = OIDCSaveLayout;
      } else if (path === '/oidc/share') {
        this.component = OIDCShareLayout;
      } else {
        // TODO error should be thrown, for now by default switch to OIDC share flow issue #1619
        this.component = OIDCShareLayout;
      }
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
