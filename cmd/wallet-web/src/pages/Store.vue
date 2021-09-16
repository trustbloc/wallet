<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex justify-center items-start w-screen h-screen">
    <div class="pt-5 bg-neutrals-softWhite rounded-b border border-neutrals-black chapi-container">
      <span class="px-5 text-xl font-bold text-neutrals-dark">Save credential</span>
      <div v-if="records.length" class="flex flex-col justify-center px-5">
        <ul class="grid grid-cols-1 gap-4 my-8">
          <li v-for="(record, index) in records" :key="index" @click="toggleDetails(index)">
            <div
              :class="[
                `group inline-flex items-center rounded-xl p-5 text-sm md:text-base font-bold border w-full h-20 md:h-24 focus-within:ring-2 focus-within:ring-offset-2 credentialPreviewContainer`,
                record.brandColor.length
                  ? `bg-gradient-${record.brandColor} border-neutrals-black border-opacity-10 focus-within:ring-primary-${record.brandColor}`
                  : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`,
              ]"
            >
              <div class="flex-none w-12 h-12 border-opacity-10">
                <img :src="require(`@/assets/img/${record.icon}`)" />
              </div>
              <div class="flex-grow p-4">
                <span
                  :class="[
                    `text-sm md:text-base font-bold text-left overflow-ellipsis`,
                    record.brandColor.length ? `text-neutrals-white` : `text-neutrals-dark`,
                  ]"
                >
                  {{ record.title }}
                </span>
              </div>
            </div>
            <!-- TODO refactor this solution, if only 1 credential present then display detail by default -->
            <div
              v-if="showDetails || records.length === 1"
              :class="index == active || records.length === 1 ? activeClass : 'hidden'"
              class="flex flex-col justify-start items-start mt-5 md:mt-6 w-full details"
            >
              <!-- todo populate with dynamic vault list -->
              <div
                class="
                  justify-start
                  items-start
                  px-4
                  mb-8
                  w-full
                  bg-neutrals-lilacSoft
                  rounded-t-lg
                  flex flex-col flex-grow
                  border-b border-neutrals-dark
                "
              >
                <label for="select-key" class="mb-1 text-sm font-bold text-neutrals-dark"
                  >Select Vault</label
                >
                <select
                  v-model="selectedDefault"
                  class="mb-1 w-full max-w-full text-base text-neutrals-dark bg-neutrals-lilacSoft"
                >
                  <option>Default Vault</option>
                </select>
              </div>

              <span class="py-3 text-base font-bold text-neutrals-dark">Verified Information</span>

              <!-- todo move this to resuable components -->
              <table class="w-full border-t border-neutrals-chatelle">
                <tr
                  v-for="(property, index) of record.properties"
                  :key="index"
                  class="border-b border-neutrals-thistle border-dotted"
                >
                  <td class="py-4 pr-6 pl-3 text-neutrals-medium">{{ property.label }}</td>
                  <td
                    v-if="property.type != 'image'"
                    class="py-4 pr-6 pl-3 text-neutrals-dark break-words"
                  >
                    {{ property.value }}
                  </td>
                  <td
                    v-if="property.type === 'image'"
                    class="py-4 pr-6 pl-3 text-neutrals-dark break-words"
                  >
                    <img :src="property.value" class="w-20 h-20" />
                  </td>
                </tr>
              </table>
            </div>
          </li>
        </ul>
      </div>
      <div v-if="errors.length">
        <b>Please correct the following error(s):</b>
        <ul>
          <!-- todo implement error as per ux design -->
          <li v-for="error in errors" :key="error" class="text-sm text-primary-valencia">
            {{ error }}
          </li>
        </ul>
      </div>
      <div class="flex justify-between p-5 w-full h-auto bg-neutrals-magnolia footerContainer">
        <button id="cancelBtn" class="btn-outline" @click="cancel">Decline</button>
        <button id="storeVCBtn" class="btn-primary" @click="store">Save</button>
      </div>
    </div>
  </div>
</template>

<script>
import {
  CHAPIEventHandler,
  getCredentialType,
  getCredentialDisplayData,
  getVCIcon,
  isVPType,
} from './mixins';
import { CredentialManager } from '@trustbloc/wallet-sdk';
import credentialDisplayData from '@/config/credentialDisplayData';
import { mapGetters } from 'vuex';

export default {
  data() {
    return {
      records: [],
      storeButton: true,
      subject: '',
      issuer: '',
      issuance: '',
      activeClass: 'is-visible',
      active: null,
      showDetails: false,
      errors: [],
      selectedDefault: 'Default Vault',
    };
  },
  computed: {
    isDisabled() {
      return this.storeButton;
    },
    // Todo issue 1075 add ii8n support : failing UI Test
  },
  created: async function () {
    // Load the Credentials
    this.credentialEvent = new CHAPIEventHandler(
      await this.$webCredentialHandler.receiveCredentialEvent()
    );
    let { dataType, data } = this.credentialEvent.getEventData();

    if (!isVPType(dataType)) {
      this.errors.push(`unknown credential data type '${dataType}'`);
      return;
    }

    let { user } = this.getCurrentUser().profile;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });

    // prepare cards
    this.prepareCards(data);
    this.presentation = data;

    // enable send vc button once loaded
    this.storeButton = false;
  },
  methods: {
    ...mapGetters(['getCurrentUser']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    prepareCards: function (data) {
      data.verifiableCredential.map((vc) => {
        const manifest = this.getManifest(vc);
        const record = this.getCredentialDisplayData(vc, manifest);
        this.records.push(record);
      });
      console.log('records', JSON.stringify(this.records, null, 2));
    },
    store: function () {
      this.errors.length = 0;
      let { token } = this.getCurrentUser().profile;

      this.credentialManager
        .save(token, { presentation: this.presentation })
        .then(() => {
          this.credentialEvent.done();
        })
        .catch((e) => {
          console.error(e);
          this.errors.push(`failed to save credential`);
        });
    },
    cancel: function () {
      this.credentialEvent.cancel();
    },
    getCredentialType: function (vc) {
      return getCredentialType(vc.type);
    },
    getCredentialDisplayData: function (vc, manifestCredential) {
      return getCredentialDisplayData(vc, manifestCredential);
    },
    getManifest: function (credential) {
      const currentCredentialType = this.getCredentialType(credential);
      return credentialDisplayData[currentCredentialType] || credentialDisplayData.fallback;
    },
    toggleDetails(i) {
      this.active = i;
      this.showDetails = !this.showDetails;
    },
  },
};
</script>

<style scoped>
.chapi-container {
  width: 28rem;
}
.credentialPreviewContainer:not(:focus-within):hover {
  box-shadow: 0px 4px 12px 0px rgba(25, 12, 33, 0.1);
}
.footerContainer {
  box-shadow: inset 0px 1px 0px 0px rgb(219, 215, 220);
}
</style>
