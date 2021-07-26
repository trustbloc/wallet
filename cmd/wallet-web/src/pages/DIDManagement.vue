<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="content">
    <div class="px-4">
      <div class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
        <md-card class="md-card-plain">
          <md-card-content>
            <md-tabs class="md-success" md-alignment="left">
              <md-tab id="tab-home" md-label="Digital Identity Preference" md-icon="how_to_reg">
                <div class="md-layout-item md-layout md-gutter">
                  <div class="md-layout-item md-medium-size-50 md-xsmall-size-75 md-size-75">
                    <md-card-content>
                      <md-list>
                        <md-subheader
                          ><label>
                            <md-icon>how_to_reg</md-icon>
                            Update Identity:
                          </label>
                        </md-subheader>

                        <md-list-item v-for="did in allDIDs" :key="did.id">
                          <md-checkbox
                            v-model="selectedDID"
                            :value="did.id"
                            @change="didSelected(did.id)"
                            >{{ did.id }}
                          </md-checkbox>
                        </md-list-item>
                      </md-list>

                      <md-divider></md-divider>

                      <md-list>
                        <md-subheader
                          ><label>
                            <md-icon>aspect_ratio</md-icon>
                            Key ID:
                          </label>
                        </md-subheader>

                        <md-field>
                          <select
                            id="verification-method"
                            v-model="verificationMethod"
                            name="verification-method"
                          >
                            <option value="default">Use Default</option>
                            <option v-for="keyID in keyIDs" :key="keyID" :value="keyID">
                              {{ keyID }}
                            </option>
                          </select>
                        </md-field>
                      </md-list>

                      <md-divider></md-divider>

                      <md-list>
                        <md-subheader>
                          <label>
                            <md-icon>vpn_lock</md-icon>
                            Update Signature Type:
                          </label>
                        </md-subheader>

                        <md-list-item>
                          <md-checkbox v-model="selectedSignType" value="Ed25519Signature2018"
                            >Ed25519Signature2018
                          </md-checkbox>
                        </md-list-item>

                        <md-list-item>
                          <md-checkbox v-model="selectedSignType" value="JsonWebSignature2020"
                            >JsonWebSignature2020
                          </md-checkbox>
                        </md-list-item>
                      </md-list>

                      <md-divider></md-divider>

                      <md-button
                        class="md-button md-info md-square"
                        :disabled="!preferencesChanged"
                        @click="updatePreferences"
                        >Update Preferences
                      </md-button>
                    </md-card-content>
                  </div>
                </div>
              </md-tab>

              <md-tab
                id="tab-home-1"
                md-label="Create TrustBloc Digital Identity"
                md-icon="add_box"
              >
                <div class="md-layout-item md-layout md-gutter">
                  <div class="md-layout-item">
                    <md-card md-alignment="left">
                      <md-card-header data-background-color="green">
                        <h4 class="title"><b>Create New Trustbloc Digital Identity</b></h4>
                      </md-card-header>
                      <md-card-content>
                        <md-label>
                          <md-icon>vpn_key</md-icon>
                          Key Type:
                        </md-label>
                        <md-field>
                          <select
                            id="selectKey"
                            v-model="keyType"
                            style="color: grey"
                            md-alignment="left"
                          >
                            <option value="ED25519">Ed25519</option>
                            <option value="ECDSAP256IEEEP1363">P-256</option>
                            <option value="ECDSAP384IEEEP1363">P-384</option>
                            <option value="BLS12381G2">BLS12381G2</option>
                          </select>
                        </md-field>

                        <md-label>
                          <md-icon>lock</md-icon>
                          Signature Suite:
                        </md-label>

                        <md-field>
                          <select id="signKey" v-model="signType" style="color: grey">
                            <option value="Ed25519VerificationKey2018">
                              Ed25519VerificationKey2018
                            </option>
                            <option value="JwsVerificationKey2020">JwsVerificationKey2020</option>
                            <option value="Bls12381G2Key2020">Bls12381G2Key2020</option>
                          </select>
                        </md-field>

                        <md-label>
                          <md-icon>design_services</md-icon>
                          Key Purpose:
                        </md-label>

                        <md-field>
                          <select id="puporseKey" v-model="purpose" style="color: grey">
                            <option value="all">all</option>
                            <option value="authentication">authentication</option>
                            <option value="assertionMethod">assertionMethod</option>
                          </select>
                        </md-field>

                        <md-button
                          id="createDIDBtn"
                          class="md-button md-info md-square"
                          @click="createDID"
                        >
                          <b>Create and Save</b>
                        </md-button>

                        <div v-if="errors.length">
                          <b>Please correct the following error(s):</b>
                          <ul>
                            <li v-for="error in errors" :key="error">{{ error }}</li>
                          </ul>
                        </div>
                        <div
                          v-if="loading"
                          style="margin-left: 40%; margin-top: 20%; height: 200px"
                        >
                          <div class="md-layout">
                            <md-progress-spinner
                              :md-diameter="100"
                              class="md-primary"
                              :md-stroke="10"
                              md-mode="indeterminate"
                            ></md-progress-spinner>
                          </div>
                        </div>
                        <div v-if="createDIDSuccess">
                          <div class="md-layout-item md-size-100" style="color: green">
                            <label id="create-did-success" class="md-helper-text"
                              >Saved your DID successfully.</label
                            >
                          </div>
                        </div>
                        <md-field>
                          <md-textarea v-model="didDocTextArea" readonly style="min-height: 300px">
                          </md-textarea>
                        </md-field>
                      </md-card-content>
                    </md-card>
                  </div>
                </div>
              </md-tab>

              <md-tab id="tab-pages" md-label="Import Any Digital Identity" md-icon="upload_file">
                <md-card class="md-card-plain">
                  <md-card-header data-background-color="green">
                    <h4 class="title">Import Any Digital Identity</h4>
                  </md-card-header>
                  <md-card-content>
                    <div class="md-layout-item md-size-100">
                      <md-icon>line_style</md-icon>
                      <label class="md-helper-text">Enter Digital Identity</label>
                      <md-field maxlength="5">
                        <md-input id="did" v-model="didID" required></md-input>
                      </md-field>
                    </div>

                    <div class="md-layout-item md-size-100">
                      <md-icon>vpn_key</md-icon>
                      <label class="md-helper-text">Select Key Format</label>
                      <div></div>
                      <md-checkbox v-model="keyFormat" value="Base58">Base58</md-checkbox>
                      <md-checkbox v-model="keyFormat" value="JWK">JWK</md-checkbox>
                      <md-field style="margin-top: -25px"></md-field>
                    </div>

                    <div class="md-layout-item md-size-100">
                      <md-icon>vpn_key</md-icon>
                      <label class="md-helper-text"
                        >Enter Private Key (in JWK or Base58 format)</label
                      >
                      <md-field maxlength="5">
                        <md-input id="privateKeyStr" v-model="privateKeyStr" required></md-input>
                      </md-field>
                    </div>

                    <div class="md-layout-item md-size-100">
                      <md-icon
                        >aspect_ratio
                        <md-tooltip md-direction="top"
                          >Enter key ID for above private key
                        </md-tooltip>
                      </md-icon>
                      <label class="md-helper-text">Enter matching Key ID</label>
                      <md-field maxlength="5">
                        <md-input id="keyID" v-model="keyID" required></md-input>
                      </md-field>
                    </div>

                    <div v-if="showImportKeyType" class="md-layout-item md-size-100">
                      <md-icon>style</md-icon>
                      <label class="md-helper-text">Select Key Type</label>
                      <md-field>
                        <select
                          id="importKeyType"
                          v-model="importKeyType"
                          style="color: grey"
                          md-alignment="left"
                        >
                          <option value="ed25519verificationkey2018">
                            Ed25519VerificationKey2018
                          </option>
                          <option value="bls12381g1key2020">Bls12381G1Key2020</option>
                        </select>
                      </md-field>
                    </div>

                    <md-button
                      id="saveDIDBtn"
                      class="
                        md-button
                        md-success
                        md-square
                        md-theme-default
                        md-large-size-100
                        md-size-100
                      "
                      @click="saveAnyDID"
                      >Resolve and Save Digital Identity
                    </md-button>
                    <div v-if="saveErrors.length">
                      <b>Please correct the following error(s):</b>
                      <ul>
                        <li v-for="error in saveErrors" :key="error">{{ error }}</li>
                      </ul>
                    </div>

                    <div v-if="saveAnyDIDSuccess">
                      <div class="md-layout-item md-size-100" style="color: green">
                        <label id="save-anydid-success" class="md-helper-text"
                          >Saved your DID successfully.</label
                        >
                      </div>
                    </div>

                    <md-field>
                      <md-textarea v-model="anyDidDocTextArea" readonly style="min-height: 360px">
                      </md-textarea>
                    </md-field>
                  </md-card-content>
                </md-card>
              </md-tab>
            </md-tabs>
          </md-card-content>
        </md-card>
      </div>
    </div>
  </div>
</template>

<script>
import { DIDManager, WalletUser } from '@trustbloc/wallet-sdk';
import { mapActions, mapGetters } from 'vuex';
import { getDIDVerificationMethod } from './mixins';

export default {
  created: async function () {
    let agent = this.getAgentInstance();
    let { user } = this.getCurrentUser().profile;

    this.walletUser = new WalletUser({ agent, user });
    this.didManager = new DIDManager({ agent, user });

    await Promise.all([this.listDIDs(), this.loadPreferences()]);

    this.updateVerificationMethod();
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts']),
    ...mapActions(['refreshUserPreference']),
    listDIDs: async function () {
      let { contents } = await this.didManager.getAllDIDs(this.getCurrentUser().profile.token);
      this.allDIDs = Object.keys(contents).map((k) => contents[k].DIDDocument);
    },
    loadPreferences: async function () {
      let { content } = await this.walletUser.getPreferences(this.getCurrentUser().profile.token);

      this.selectedDID = content.controller;
      this.selectedSignType = content.proofType;
      this.preference = content;
      this.verificationMethod = content.verificationMethod ? content.verificationMethod : 'default';
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

      this.didDocTextArea = `Created ${docRes.DIDDocument.id}`;
      this.createDIDSuccess = true;
      this.loading = false;
      this.listDIDs();
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
      this.listDIDs();
    },
    updatePreferences() {
      this.walletUser
        .updatePreferences(this.getCurrentUser().profile.token, {
          controller: this.selectedDID,
          proofType: this.selectedSignType,
          verificationMethod: this.verificationMethod != 'default' ? this.verificationMethod : '',
        })
        .then(() => {
          this.refreshUserPreference();
        });

      this.preference.controller = this.selectedDID;
      this.preference.proofType = this.selectedSignType;
      this.preference.verificationMethod = this.verificationMethod;
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
  data() {
    return {
      didDocTextArea: '',
      anyDidDocTextArea: '',
      purpose: 'all',
      keyType: 'ED25519',
      signType: 'Ed25519VerificationKey2018',
      didID: '',
      privateKeyStr: '',
      keyID: '',
      errors: [],
      saveErrors: [],
      loading: false,
      createDIDSuccess: false,
      saveAnyDIDSuccess: false,
      allDIDs: {},
      importKeyType: '',
      keyFormat: '',
      allSignatureTypes: [{ id: 'Ed25519Signature2018' }, { id: 'JsonWebSignature2020' }],
      selectedDID: '',
      selectedSignType: '',
      verificationMethod: '',
      keyIDs: [],
      preference: {},
    };
  },
  computed: {
    preferencesChanged() {
      let { controller, proofType, verificationMethod } = this.preference;

      if (controller != this.selectedDID) {
        return true;
      }

      if (proofType != this.selectedSignType) {
        return true;
      }

      if (
        verificationMethod != (this.verificationMethod == 'default' ? '' : this.verificationMethod)
      ) {
        return true;
      }

      return false;
    },
    showImportKeyType() {
      return this.keyFormat == 'Base58';
    },
  },
};
</script>
<style></style>
