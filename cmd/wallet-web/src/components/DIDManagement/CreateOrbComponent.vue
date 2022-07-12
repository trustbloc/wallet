<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-start py-8">
    <div class="flex flex-col grow justify-start items-start w-full">
      <label for="select-key" class="mb-1 text-base">Key Type:</label>
      <select id="select-key" v-model="keyType" class="mb-5 w-full max-w-full border-b">
        <option value="ED25519">Ed25519</option>
        <option value="ECDSAP256IEEEP1363">P-256</option>
        <option value="ECDSAP384IEEEP1363">P-384</option>
        <option value="BLS12381G2">BLS12381G2</option>
      </select>
    </div>

    <div class="flex flex-col grow justify-start items-start w-full">
      <label for="select-signature-suite" class="mb-1 text-base">Signature Suite:</label>
      <select
        id="select-signature-suite"
        v-model="signType"
        class="mb-5 w-full max-w-full border-b"
      >
        <option value="Ed25519VerificationKey2018">Ed25519VerificationKey2018</option>
        <option value="JwsVerificationKey2020">JwsVerificationKey2020</option>
        <option value="Bls12381G2Key2020">Bls12381G2Key2020</option>
      </select>
    </div>

    <div class="flex flex-col grow justify-start items-start w-full">
      <label for="select-key-purpose" class="mb-1 text-base">Key Purpose:</label>
      <select id="select-key-purpose" v-model="purpose" class="mb-5 w-full max-w-full border-b">
        <option value="all">all</option>
        <option value="authentication">authentication</option>
        <option value="assertionMethod">assertionMethod</option>
      </select>
    </div>

    <button id="createDIDBtn" class="mt-5 btn-primary" @click="createDID">Create and Save</button>

    <div v-if="errors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in errors" :key="error">{{ error }}</li>
      </ul>
    </div>
    <div v-if="loading" class="mt-2 ml-4 h-48">
      <div>
        <SpinnerIcon />
      </div>
    </div>
    <div v-if="createDIDSuccess">
      <div style="color: green">
        <label id="create-did-success" class="text-sm">Saved your DID successfully.</label>
      </div>
    </div>
    <div>
      <span class="max-w-full break-words">{{ didDocTextArea }}</span>
    </div>
  </div>
</template>

<script>
import { DIDManager } from '@trustbloc/wallet-sdk';
import { mapActions, mapGetters } from 'vuex';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

export default {
  name: 'CreateOrbComponent',
  components: {
    SpinnerIcon,
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
      didDocTextArea: '',
      purpose: 'all',
      keyType: 'ED25519',
      signType: 'Ed25519VerificationKey2018',
      errors: [],
      loading: false,
      createDIDSuccess: false,
    };
  },
  created: function () {
    const agent = this.getAgentInstance();
    const { user } = this.getCurrentUser().profile;
    this.didManager = new DIDManager({ agent, user });
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts']),
    ...mapActions(['refreshUserPreference']),
    updateAllDIDs: async function () {
      const { contents } = await this.didManager.getAllDIDs(this.getCurrentUser().profile.token);
      const newAllDIDs = Object.keys(contents).map((k) => contents[k].didDocument);
      this.$emit('update:allDIDs', newAllDIDs);
    },
    createDID: async function () {
      this.errors.length = 0;
      this.createDIDSuccess = false;
      this.loading = true;

      let docRes;
      try {
        docRes = await this.didManager.createOrbDID(this.getCurrentUser().profile.token, {
          purposes: this.purpose == 'all' ? ['assertionMethod', 'authentication'] : [this.purpose],
          keyType: this.keyType,
          signatureType: this.signType,
        });
      } catch (e) {
        this.loading = false;
        this.didDocTextArea = `failed to create did: ${e.toString()}`;
        return;
      }

      this.didDocTextArea = `Created ${docRes.didDocument.id}`;
      this.createDIDSuccess = true;
      this.loading = false;
      await this.updateAllDIDs();
    },
  },
};
</script>
