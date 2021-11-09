<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-start py-8">
    <div class="flex flex-col flex-grow justify-start items-start mb-5 w-full">
      <label for="did-input">Enter Digital Identity:</label>
      <input id="did-input" v-model="didID" class="w-full border-b" required />
    </div>

    <div class="flex flex-col flex-grow justify-start items-start mb-5 w-full">
      <label>Select Key Format:</label>
      <div class="flex flex-row justify-start items-center">
        <input id="Base58" v-model="keyFormat" type="radio" value="Base58" />
        <label for="Base58" class="pl-1">Base58</label>
      </div>
      <div class="flex flex-row justify-start items-center">
        <input id="JWK" v-model="keyFormat" type="radio" value="JWK" />
        <label for="JWK" class="pl-1">JWK</label>
      </div>
    </div>

    <div class="flex flex-col flex-grow justify-start items-start mb-5 w-full">
      <label for="privateKeyStr">Enter Private Key (in JWK or Base58 format):</label>
      <input id="privateKeyStr" v-model="privateKeyStr" class="w-full border-b" required />
    </div>

    <div class="flex relative flex-col flex-grow justify-start items-start mb-5 w-full">
      <label for="keyID">Enter matching Key ID:</label>
      <input
        id="keyID"
        v-model="keyID"
        class="w-full border-b"
        placeholder="Enter key ID for above private key"
        required
      />
    </div>

    <div v-if="showImportKeyType" class="mb-5">
      <label for="importKeyType">Select Key Type:</label>
      <select id="importKeyType" v-model="importKeyType">
        <option value="ed25519verificationkey2018">Ed25519VerificationKey2018</option>
        <option value="bls12381g1key2020">Bls12381G1Key2020</option>
      </select>
    </div>

    <button id="saveDIDBtn" class="mt-5 btn-primary" @click="saveAnyDID">
      Resolve and Save Digital Identity
    </button>

    <div v-if="saveErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in saveErrors" :key="error">{{ error }}</li>
      </ul>
    </div>

    <div v-if="saveAnyDIDSuccess">
      <div>
        <label id="save-anydid-success">Saved your DID successfully.</label>
      </div>
    </div>

    <div>
      <span class="h-80">{{ anyDidDocTextArea }}</span>
    </div>
  </div>
</template>

<script>
import { DIDManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';

export default {
  name: 'ImportDid',
  props: {
    allDIDs: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      anyDidDocTextArea: '',
      keyType: 'ED25519',
      didID: '',
      privateKeyStr: '',
      keyID: '',
      saveErrors: [],
      loading: false,
      saveAnyDIDSuccess: false,
      importKeyType: '',
      keyFormat: '',
    };
  },
  computed: {
    showImportKeyType() {
      return this.keyFormat == 'Base58';
    },
  },
  created: function () {
    const agent = this.getAgentInstance();
    const { user } = this.getCurrentUser().profile;
    this.didManager = new DIDManager({ agent, user });
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    updateAllDIDs: async function () {
      const { contents } = await this.didManager.getAllDIDs(this.getCurrentUser().profile.token);
      const newAllDIDs = Object.keys(contents).map((k) => contents[k].didDocument);
      this.$emit('update:allDIDs', newAllDIDs);
    },
    saveAnyDID: async function () {
      this.saveErrors.length = 0;
      this.saveAnyDIDSuccess = false;
      this.anyDidDocTextArea = '';
      this.loading = true;

      if (this.didID.length == 0) {
        this.saveErrors.push('did id required.');
        return;
      }

      if (this.keyFormat.length == 0) {
        this.saveErrors.push('please select format of the key being imported.');
        return;
      }

      if (this.privateKeyStr.length == 0) {
        this.saveErrors.push('private key is required.');
        return;
      }

      if (this.keyID.length == 0) {
        this.saveErrors.push('key ID (verification method) matching private key is required.');
        return;
      }

      if (this.keyFormat == 'Base58' && this.importKeyType.length == 0) {
        this.saveErrors.push('key type of private key for importing base58 private keys.');
        return;
      }

      try {
        await this.didManager.importDID(this.getCurrentUser().profile.token, {
          did: this.didID,
          key: {
            keyType: this.importKeyType,
            privateKeyBase58: this.keyFormat == 'Base58' ? this.privateKeyStr : '',
            privateKeyJwk: this.keyFormat == 'JWK' ? JSON.parse(this.privateKeyStr) : undefined,
            keyID: this.keyID,
          },
        });
      } catch (e) {
        this.loading = false;
        this.anyDidDocTextArea = `failed to import did: ${e.toString()}`;
        return;
      }

      this.loading = false;
      this.saveAnyDIDSuccess = true;
      await this.updateAllDIDs();
    },
  },
};
</script>
