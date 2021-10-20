<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <h3 class="mb-5 font-bold text-neutrals-dark">Wallet Operations</h3>
    <div id="#ivc">
      <textarea id="vcDataTextArea" v-model="interopData" rows="15" class="w-full" />
      <div class="flex flex-col justify-start items-start my-4">
        <span class="mb-2 font-bold">Sample requests:</span>
        <div class="w-full">
          <button
            id="store-vc-sample-1"
            class="m-1 btn-outline"
            @click="prefillRequest('vp', 'store')"
          >
            Store Presentation (Sample 1)
          </button>
          <button
            id="store-vc-sample-2"
            class="m-1 btn-outline"
            @click="prefillRequest('vp2', 'store')"
          >
            Store Presentation (Sample 2)
          </button>
          <button id="req-vp" class="m-1 btn-outline" @click="prefillRequest('getvp', 'get')">
            Request Presentation
          </button>
          <button id="sel-disclosure" class="m-1 btn-outline" @click="prefillRequest('bbs', 'get')">
            Selective Disclosure
          </button>
          <button
            id="multi-query-1"
            class="m-1 btn-outline"
            @click="prefillRequest('multiQ1', 'get')"
          >
            Multiple Query - 1
          </button>
          <button
            id="multi-query-2"
            class="m-1 btn-outline"
            @click="prefillRequest('multiQ2', 'get')"
          >
            Multiple Query - 2
          </button>
          <button
            id="multi-query-3"
            class="m-1 btn-outline"
            @click="prefillRequest('multiQ3', 'get')"
          >
            Multiple Query - 3
          </button>
          <button id="pexq" class="m-1 btn-outline" @click="prefillRequest('pexq', 'get')">
            Presentation Exchange Query
          </button>
          <button
            id="pexq-didcomm"
            class="m-1 btn-outline"
            @click="prefillRequest('pexq-didcomm', 'get')"
          >
            Presentation Exchange Query With DIDComm
          </button>
          <button
            id="pexq-didcomm-govnvc"
            class="m-1 btn-outline"
            @click="prefillRequest('pexq-didcomm-govnvc', 'get')"
          >
            Presentation Exchange Query With DIDComm & Governance VC
          </button>
          <button id="didauth" class="m-1 btn-outline" @click="prefillRequest('didauth', 'get')">
            DID Auth
          </button>
          <button id="didconn" class="m-1 btn-outline" @click="prefillRequest('didconn', 'get')">
            DID Connect
          </button>
          <button
            id="didconn-manifest"
            class="m-1 btn-outline"
            @click="prefillRequest('didconn-manifest', 'get')"
          >
            DID Connect with manifest
          </button>
          <button
            id="didconn-manifest-usrc"
            class="m-1 btn-outline"
            @click="prefillRequest('didconn-manifest-usrc', 'get')"
          >
            DID Connect with manifest and user credential
          </button>
          <button
            id="didconn-manifest-usrc-govvc"
            class="m-1 btn-outline"
            @click="prefillRequest('didconn-manifest-usrc-govvc', 'get')"
          >
            DID Connect with manifest, user credential and governance VC
          </button>
          <button class="m-1 btn-outline" @click="prefillRequest('waci-credential-share', 'get')">
            WACI Credential Share
          </button>
        </div>
      </div>

      <div class="flex flex-col justify-start items-start">
        <span class="mb-2 font-bold">Wallet Operations:</span>
        <div class="flex flex-row justify-start items-center space-x-2 w-full">
          <button id="store-btn" class="btn-primary" :disabled="disableStore" @click="store">
            Store
          </button>
          <button id="get-btn" class="btn-primary" :disabled="disableGet" @click="get">Get</button>
        </div>
      </div>

      <div v-if="responses.length" style="color: #0e9a00">
        <p v-for="response in responses" :key="response">{{ response }}</p>
      </div>
      <div v-if="errors.length" style="color: #fb4934">
        <b>Please correct the following error(s):</b>
        <ul>
          <li v-for="error in errors" :key="error">{{ error }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import { WebCredential } from 'credential-handler-polyfill/WebCredential.js';
import { getSample } from './webWalletSamples';

export default {
  data() {
    return {
      interopData: '',
      mode: '',
      errors: [],
      responses: [],
    };
  },
  computed: {
    disableStore() {
      return this.mode == 'get';
    },
    disableGet() {
      return this.mode == 'store';
    },
  },
  created: async function () {
    let opts = this.$store.getters.getAgentOpts;
    if (!opts) {
      this.errors.push('Please login to your webwallet before running this demo');
      return;
    }

    await this.$polyfill.loadOnce(opts.credentialMediatorURL);
  },
  methods: {
    clearResults: function () {
      this.errors.length = 0;
      this.responses.length = 0;
    },
    prefillRequest: function (id, mode) {
      this.interopData = JSON.stringify(getSample(id), null, 2);
      this.mode = mode;
    },
    store: async function () {
      this.clearResults();
      if (this.interopData.length == 0) {
        this.errors.push('Invalid presentation');
        return;
      }
      const webCredentialWrapper = new WebCredential(
        'VerifiablePresentation',
        JSON.parse(this.interopData)
      );
      const result = await navigator.credentials.store(webCredentialWrapper);
      console.log('Result received via store() request:', result);
      this.responses.push('Successfully stored verifiable presentation to wallet.');
    },
    get: async function () {
      this.clearResults();
      if (this.interopData.length == 0) {
        this.errors.push('Invalid query');
        return;
      }
      const result = await navigator.credentials.get(JSON.parse(this.interopData));
      if (!result) {
        this.errors.push('Failed to get result');
        return;
      }

      this.showResp(result.data);
      this.responses.push('Successfully got response from wallet.');
    },
    showResp: function (data) {
      if (typeof data == 'object') {
        this.interopData = JSON.stringify(data, null, 2);
      } else {
        this.responses.push('Warning: received unexpected string data type');
        this.interopData = data;
      }
    },
  },
};
</script>
