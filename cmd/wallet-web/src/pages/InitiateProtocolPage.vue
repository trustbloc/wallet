<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="mx-auto flex h-screen max-h-screen max-w-7xl grow flex-col items-center justify-start shadow-main-container"
  >
    <HeaderComponent :has-custom-gradient="true" :show-menu-dropdown="false">
      <template #gradientContainer>
        <div class="oval absolute h-15 bg-gradient-full" />
      </template>
    </HeaderComponent>
    <div class="grid grid-cols-2 w-full h-full divide-x">
      <div
        class="relative z-10 flex h-full w-full grow flex-col items-center justify-start overflow-hidden bg-neutrals-softWhite px-6 pt-32"
      >
        <h4>OpenID4VC Initiate</h4>
        <textarea
          id="initiateIssuanceRequest"
          v-model="initiateRequest"
          placeholder="Paste Initiate URL"
          class="w-full h-auto border-b"
        >
        </textarea>
        <StyledButtonComponent
          id="initiateFlow"
          class="mt-5"
          type="btn-primary"
          :loading="loading"
          @click="initiateIssuanceFlow()"
        >
          Initiate
        </StyledButtonComponent>
      </div>
      <div
        class="relative z-10 flex h-full w-full grow flex-col items-center justify-start overflow-hidden bg-neutrals-softWhite px-6 pt-32"
      >
        <h4>Scan a QR Code</h4>
        <img
          src="@/assets/img/qr-code-scan-icon.svg"
          class="w-16 h-16"
          @click="loadingScanner = !loadingScanner"
        />
        <div
          v-if="loadingScanner"
          class="relative z-10 flex h-full w-full grow flex-col items-center justify-start overflow-hidden"
        >
          <p v-if="error">Error: {{ error }}</p>
          <p>Result: {{ decodedString }}</p>
          <qrcode-stream @init="onInit" @decode="onDecode" @detect="onDetect"></qrcode-stream>
        </div>
      </div>
    </div>
  </div>

  <FooterComponent
    class="sticky bottom-0 z-20 border-t border-neutrals-thistle bg-neutrals-magnolia"
  />
</template>
<script>
import HeaderComponent from '@/components/Header/HeaderComponent.vue';
import FooterComponent from '@/components/Footer/FooterComponent.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import { QrcodeStream } from 'vue3-qrcode-reader';

export default {
  name: 'InitiateProtocolPage',
  components: {
    StyledButtonComponent,
    HeaderComponent,
    FooterComponent,
    QrcodeStream,
  },
  data() {
    return {
      initiateRequest: '',
      loading: false,
      error: '',
      decodedString: '',
      loadingScanner: false,
    };
  },
  methods: {
    initiateIssuanceFlow() {
      const initiateRequestDataUrl = new URL(this.initiateRequest);
      if (
        initiateRequestDataUrl.protocol.includes('openid-initiate-issuance') ||
        initiateRequestDataUrl.pathname.includes('initiate_issuance')
      ) {
        this.loading = true;
        this.$router.push({ name: 'save' });
      }
      // For presentation initiate messages
      if (initiateRequestDataUrl.protocol.includes('openid:')) {
        this.loading = true;
        this.$router.push({ name: 'share' });
      }
    },
    async onInit(promise) {
      try {
        await promise;
      } catch (error) {
        if (error.name === 'NotAllowedError') {
          this.error = 'user denied camera access permisson';
        } else if (error.name === 'NotFoundError') {
          this.error = 'no suitable camera device installed';
        } else if (error.name === 'NotSupportedError') {
          this.error = 'page is not served over HTTPS (or localhost)';
        } else if (error.name === 'NotReadableError') {
          this.error = 'maybe camera is already in use';
        } else if (error.name === 'OverconstrainedError') {
          this.error = 'did you requested the front camera although there is none?';
        } else if (error.name === 'StreamApiNotSupportedError') {
          this.error = 'browser seems to be lacking features';
        }
      }
    },
    onDecode(decodedString) {
      window.location.replace(decodedString);
    },
  },
};
</script>

<style scoped>
.oval {
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  filter: blur(50px);
  width: 15.625rem; /* 250px */
  top: 2.0625rem; /* 33px */
}
</style>
