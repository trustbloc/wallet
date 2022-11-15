<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 * Copyright Avast Software. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import HeaderComponent from '@/components/Header/HeaderComponent.vue';
import FooterComponent from '@/components/Footer/FooterComponent.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import { QrcodeStream } from 'vue3-qrcode-reader';

// Local Variables
const initiateRequest = ref('');
const loading = ref(false);
const error = ref('');
const decodedString = ref('');
const loadingScanner = ref(false);

// Hooks
const router = useRouter();

// Methods
function initiateIssuanceFlow() {
  const initiateRequestDataUrl = new URL(initiateRequest.value);
  if (
    initiateRequestDataUrl.protocol.includes('openid-initiate-issuance') ||
    initiateRequestDataUrl.pathname.includes('initiate_issuance')
  ) {
    loading.value = true;
    router.push({ name: 'save' });
  }
  // For presentation initiate messages
  else if (
    initiateRequestDataUrl.protocol.includes('openid-vc:') &&
    initiateRequestDataUrl.searchParams.has('request_uri')
  ) {
    loading.value = true;
    router.push({
      name: 'openid4vc-share',
      query: { url: initiateRequest.value },
    });
  }
}
async function onInit(promise) {
  try {
    await promise;
  } catch (error) {
    if (error.name === 'NotAllowedError') {
      error.value = 'user denied camera access permission';
    } else if (error.name === 'NotFoundError') {
      error.value = 'no suitable camera device installed';
    } else if (error.name === 'NotSupportedError') {
      error.value = 'page is not served over HTTPS (or localhost)';
    } else if (error.name === 'NotReadableError') {
      error.value = 'maybe camera is already in use';
    } else if (error.name === 'OverconstrainedError') {
      error.value = 'did you requested the front camera although there is none?';
    } else if (error.name === 'StreamApiNotSupportedError') {
      error.value = 'browser seems to be lacking features';
    }
  }
}
function onDecode(decodedString) {
  window.location.replace(decodedString);
}
</script>

<template>
  <div
    class="mx-auto flex h-screen max-h-screen max-w-7xl grow flex-col items-center justify-start shadow-main-container"
  >
    <HeaderComponent :has-custom-gradient="true" :show-menu-dropdown="false">
      <template #gradientContainer>
        <div class="oval absolute h-15 bg-gradient-full" />
      </template>
    </HeaderComponent>
    <div class="grid h-full w-full grid-cols-2 divide-x">
      <div
        class="relative z-10 flex h-full w-full grow flex-col items-center justify-start overflow-hidden bg-neutrals-softWhite px-6 pt-32"
      >
        <h4>OpenID4VC Initiate</h4>
        <textarea
          id="initiateIssuanceRequest"
          v-model="initiateRequest"
          placeholder="Paste Initiate URL"
          class="h-auto w-full border-b"
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
          class="h-16 w-16"
          @click="loadingScanner = !loadingScanner"
        />
        <div
          v-if="loadingScanner"
          class="relative z-10 flex h-full w-full grow flex-col items-center justify-start overflow-hidden"
        >
          <p v-if="error">Error: {{ error }}</p>
          <p>Result: {{ decodedString }}</p>
          <qrcode-stream @init="onInit" @decode="onDecode" @detect="onDetect" />
        </div>
      </div>
    </div>
  </div>

  <FooterComponent
    class="sticky bottom-0 z-20 border-t border-neutrals-thistle bg-neutrals-magnolia"
  />
</template>

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
