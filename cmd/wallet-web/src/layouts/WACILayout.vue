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
    <HeaderComponent :has-custom-gradient="true">
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
    </HeaderComponent>
    <keep-alive>
      <component
        :is="component"
        :protocol-handler="protocolHandler"
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
    <FooterComponent
      class="sticky bottom-0 z-20 border-t border-neutrals-thistle bg-neutrals-magnolia"
    />
  </div>
</template>

<script>
import { reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { extractOOBGoalCode } from '@trustbloc/wallet-sdk';
import { WACIRedirectHandler } from '@/mixins';
import { decode } from 'js-base64';
import { WACIShareLayoutMutations } from '@/layouts/WACIShareLayout.vue';
import WACIShareLayout from '@/layouts/WACIShareLayout.vue';
import WACISharePage from '@/pages/WACISharePage.vue';
import WACIIssuePage from '@/pages/WACIIssuePage.vue';
import HeaderComponent from '@/components/Header/HeaderComponent.vue';
import FooterComponent from '@/components/Footer/FooterComponent.vue';

export const WACIStore = reactive({
  processedCredentials: [],
  selectedCredentialId: null,
  protocolHandler: null,
});

export const WACIGetters = {
  getProcessedCredentialById(id) {
    return WACIStore.processedCredentials.find((credential) => credential.id === id);
  },
};

export const WACIMutations = {
  setProcessedCredentials(value) {
    WACIStore.processedCredentials = value;
  },
  setProtocolHandler(value) {
    WACIStore.protocolHandler = value;
  },
  setSelectedCredentialId(value) {
    WACIStore.selectedCredentialId = value;
  },
};

function findForm(query) {
  if (!query.oob) {
    // TODO [Issue#1325] should be redirected to standard error screen.
    throw 'access denied, oob invitation missing';
  }

  const invitation = JSON.parse(decode(query.oob));

  switch (extractOOBGoalCode(invitation)) {
    case 'streamlined-vc':
      return {
        component: WACIIssuePage,
        protocolHandler: new WACIRedirectHandler(invitation),
      };
    case 'streamlined-vp':
      return {
        component: WACIShareLayout,
        protocolHandler: new WACIRedirectHandler(invitation),
      };
    default:
      // TODO [Issue#1326] should throw error if goal-code is missing once other trustbloc components starts sending right goal codes.
      return {
        component: WACIShareLayout,
        protocolHandler: new WACIRedirectHandler(invitation),
      };
  }
}

export default {
  components: {
    HeaderComponent,
    FooterComponent,
  },
  setup() {
    const { t } = useI18n();
    const selectedCredentialId = ref(WACIStore.selectedCredentialId);
    watch(
      () => WACIStore.selectedCredentialId,
      () => {
        selectedCredentialId.value = WACIStore.selectedCredentialId;
      }
    );
    return { selectedCredentialId, t };
  },
  data() {
    return {
      component: null,
      protocolHandler: null,
    };
  },
  created: function () {
    const { component, protocolHandler } = findForm(this.$route.query);
    this.component = component;
    WACIMutations.setProtocolHandler(protocolHandler);
  },
  methods: {
    handleBackButtonClick() {
      WACIShareLayoutMutations.setComponent(WACISharePage);
      WACIMutations.setSelectedCredentialId(null);
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
