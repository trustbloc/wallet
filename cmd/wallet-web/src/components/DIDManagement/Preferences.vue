<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-start py-8">
    <div class="flex flex-col justify-start items-start w-full">
      <label for="did-selector" class="mb-1 text-base">Update Identity:</label>
      <select
        v-if="allDIDs.length"
        id="did-selector"
        v-model="selectedDID"
        class="mb-5 w-full max-w-full border-b"
        name="Selected DID"
        @change="didSelected(selectedDID)"
      >
        <option v-for="did in allDIDs" :key="did.id" :value="did.id">{{ did.id }}</option>
      </select>
    </div>

    <div class="flex flex-col flex-grow justify-start items-start w-full">
      <label for="verification-method" class="mb-1 text-base">Key ID:</label>
      <select
        id="verification-method"
        v-model="verificationMethod"
        class="mb-5 w-full max-w-full border-b"
        name="Verification Method"
      >
        <option value="default">Use Default</option>
        <option v-for="keyId in keyIDs" :key="keyId" :value="keyId">
          {{ keyId }}
        </option>
      </select>
    </div>

    <div class="flex flex-col flex-grow justify-start items-start w-full">
      <label for="signature-type-selector" class="mb-1 text-base">Update Signature Type:</label>
      <div
        v-for="signatureType in allSignatureTypes"
        :key="signatureType.id"
        class="flex flex-row justify-start items-center"
      >
        <input
          :id="['signature-type-selector-', signatureType.id]"
          v-model="selectedSignType"
          :value="signatureType.id"
          type="radio"
          name="Signature Type"
          checked
        />
        <label :for="['signature-type-selector-', signatureType.id]" class="pl-1"
          >{{ signatureType.id }}
        </label>
      </div>
    </div>

    <styled-button
      id="update-btn"
      class="mt-5"
      :disabled="!preferencesChanged"
      :loading="loading"
      type="primary"
      @click="updatePreferences()"
    >
      <span>Update Preferences</span>
    </styled-button>

    <div v-if="updateSuccessful && !preferencesChanged">
      <label id="update-preferences-success">Updated preferences successfully.</label>
    </div>
  </div>
</template>

<script>
import { DIDManager, WalletUser } from '@trustbloc/wallet-sdk';
import { mapActions, mapGetters } from 'vuex';
import { getDIDVerificationMethod } from '@/pages/mixins';
import StyledButton from '@/components/StyledButton/StyledButton.vue';

export default {
  name: 'Preferences',
  components: {
    StyledButton,
  },
  model: {
    prop: 'allDIDs',
    event: 'update',
  },
  props: {
    allDIDs: {
      type: Array,
      required: true,
    },
  },
  emits: ['update:allDIDs'],
  data() {
    return {
      keyID: '',
      allSignatureTypes: [{ id: 'Ed25519Signature2018' }, { id: 'JsonWebSignature2020' }],
      selectedDID: '',
      selectedSignType: '',
      verificationMethod: '',
      updateSuccessful: false,
      loading: false,
      keyIDs: [],
      preference: {},
    };
  },
  computed: {
    preferencesChanged() {
      const { controller, proofType, verificationMethod } = this.preference;
      if (controller !== this.selectedDID) {
        return true;
      }
      if (proofType !== this.selectedSignType) {
        return true;
      }
      if (verificationMethod !== this.verificationMethod) {
        return true;
      }
      return false;
    },
  },
  created: async function () {
    const agent = this.getAgentInstance();
    const { user } = this.getCurrentUser().profile;
    this.walletUser = new WalletUser({ agent, user });
    this.didManager = new DIDManager({ agent, user });
    await Promise.all([this.listDIDs(), this.loadPreferences()]);
    this.updateVerificationMethod();
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    ...mapActions(['refreshUserPreference']),
    listDIDs: async function () {
      const { contents } = await this.didManager.getAllDIDs(this.getCurrentUser().profile.token);
      const newAllDIDs = Object.keys(contents).map(
        (k) => contents[k].didDocument || contents[k].DIDDocument
      );
      this.$emit('update:allDIDs', newAllDIDs);
      if (this.selectedDID === '') this.didSelected(newAllDIDs[0].id);
    },
    loadPreferences: async function () {
      const { content } = await this.walletUser.getPreferences(this.getCurrentUser().profile.token);
      this.selectedDID = content.controller;
      this.selectedSignType = content.proofType;
      this.preference = {
        ...content,
        verificationMethod: content.verificationMethod ? content.verificationMethod : 'default',
      };
      this.verificationMethod = content.verificationMethod ? content.verificationMethod : 'default';
    },
    updatePreferences() {
      this.updateSuccessful = false;
      this.loading = true;
      this.walletUser
        .updatePreferences(this.getCurrentUser().profile.token, {
          controller: this.selectedDID,
          proofType: this.selectedSignType,
          verificationMethod: this.verificationMethod !== 'default' ? this.verificationMethod : '',
        })
        .then(() => {
          this.refreshUserPreference();
        });

      this.preference.controller = this.selectedDID;
      this.preference.proofType = this.selectedSignType;
      this.preference.verificationMethod = this.verificationMethod;
      this.loading = false;
      this.updateSuccessful = true;
    },
    didSelected(did) {
      this.selectedDID = did;
      this.verificationMethod = 'default';
      this.updateVerificationMethod();
    },
    updateVerificationMethod() {
      this.keyIDs = getDIDVerificationMethod(this.allDIDs, this.selectedDID);
    },
  },
};
</script>
