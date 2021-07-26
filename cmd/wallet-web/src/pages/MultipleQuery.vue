/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <div v-if="loading" class="w-screen" style="margin-left: 40%; margin-top: 20%; height: 200px">
    <div class="md-layout">
      <md-progress-spinner
        :md-diameter="100"
        class="md-accent"
        :md-stroke="10"
        md-mode="indeterminate"
      ></md-progress-spinner>
    </div>
  </div>
  <div v-else class="flex justify-center w-screen md-layout">
    <div class="max-w-screen-sm md-layout-item">
      <div
        class="flex flex-col items-start md-layout md-alignment-center-center"
        style="margin-top: 20px"
      >
        <div class="md-headline">Credential Presentation Requested</div>
        <div class="md-subheading">A credential presentation is been requested:</div>
      </div>

      <div style="margin: 10px"></div>

      <div class="md-layout md-alignment-center-center">By</div>

      <div class="md-layout md-alignment-center-center">
        <div style="padding-bottom: 10px">
          <governance :govn-v-c="govnVC" :request-origin="requestOrigin" :issuer="false" />
        </div>
      </div>

      <div style="margin: 20px"></div>

      <div v-if="errors.length">
        <b>Failed with following error(s):</b>
        <md-field style="margin-top: -15px">
          <ul>
            <li v-for="error in errors" :key="error" style="color: #9d0006">{{ error }}</li>
          </ul>
        </md-field>

        <md-button
          id="cancelBtnNrc"
          style="background-color: #9d0006 !important"
          class="md-cancel-text"
          @click="cancel"
        >
          Cancel
        </md-button>
      </div>

      <div
        v-if="reasons.length || presExchReasons.length"
        class="md-layout md-alignment-center-center reasons"
      >
        <ul>
          <md-card class="md-layout md-alignment-center-center" style="background: none !important">
            <md-card-expand>
              <md-card-actions md-alignment="space-between" style="background: none !important">
                <div class="md-subheading">Reason:</div>
                <md-card-expand-trigger>
                  <md-button class="md-icon-button">
                    <md-icon>keyboard_arrow_down</md-icon>
                  </md-button>
                </md-card-expand-trigger>
              </md-card-actions>

              <md-card-expand-content>
                <md-card-content>
                  <ul>
                    <li v-for="(reason, index) in reasons" :key="index">
                      <b>{{ reason }}</b>
                    </li>

                    <li v-for="(requirement, index) in presExchReasons" :key="index">
                      <b>{{ requirement.name }}</b
                      >: {{ requirement.purpose }}
                      <div>{{ requirement.rule }}</div>
                      <ul>
                        <li v-for="descriptor in requirement.descriptors" :key="descriptor.name">
                          <b>{{ descriptor.name }} </b>{{ descriptor.purpose }}
                          <ul>
                            <li v-for="constraint in descriptor.constraints" :key="constraint">
                              {{ constraint }}
                            </li>
                          </ul>
                        </li>
                      </ul>
                    </li>
                  </ul>
                </md-card-content>
              </md-card-expand-content>
            </md-card-expand>
          </md-card>
        </ul>
      </div>

      <md-card-content v-if="records.length" class="md-layout md-alignment-center-center card-list">
        <ul>
          <li v-for="(card, index) in records" :key="index">
            <transition name="flip">
              <div class="card" style="padding-bottom: 35px">
                <div class="cardContent">
                  <div class="cardHeader">
                    {{ card.title }}
                  </div>

                  <div class="cardBody">
                    <div class="cardDetailsL">
                      <md-icon class="md-size-4x">{{ card.icon }}</md-icon>
                    </div>
                    <div class="cardDetailsR">
                      <p>{{ card.description }}</p>
                      <div v-if="card.body">
                        The verifier can only access below information from your credential.
                        <div v-for="(subj, skey) in card.body" :key="skey">
                          <div v-if="displayContent(skey)" class="md-caption">
                            <b>{{ skey.replace('.', ' ') }} </b>: {{ subj }}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </transition>
          </li>
        </ul>
      </md-card-content>

      <div v-if="showShareCredential" class="md-layout md-alignment-center-center">
        <p class="md-body-1">
          By clicking Agree you will be sharing a unique identifier to
          <b style="color: #2e7d32">{{ requestOrigin }}</b
          >, the Credential content, and your digital signature.
          <a href="https://www.w3.org/TR/vc-data-model/#proofs-signatures" target="_blank"
            >Learn more</a
          >
        </p>

        <md-button
          id="share-credentials"
          class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100 col"
          style="background-color: #29a329 !important"
          @click="share"
        >
          Agree
        </md-button>
        <md-button
          id="cancelBtn"
          style="margin-left: 5px; background-color: #9d0006 !important"
          class="md-cancel-text"
          @click="cancel"
        >
          Cancel
        </md-button>
      </div>
    </div>
  </div>
</template>
<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import {
  CHAPIEventHandler,
  extractPresentationExchangeReasons,
  extractQueryReasons,
  flatCredentialSubject,
  normalizeQuery,
  getVCIcon,
  getCredentialType,
} from './mixins';
import { mapGetters } from 'vuex';
import Governance from './Governance.vue';

const nonDisplayContent = ['id', 'type'];

export default {
  components: {
    Governance,
  },
  data() {
    return {
      errors: [],
      requestOrigin: '',
      loading: true,
      allIcons: [
        'account_box',
        'contacts',
        'person',
        'person_outline',
        'card_membership',
        'portrait',
        'bento',
      ],
      records: [],
      govnVC: null,
      reasons: [],
      presExchReasons: [],
    };
  },
  created: async function () {
    this.loading = true;
    this.chapiHandler = new CHAPIEventHandler(this.$parent.credentialEvent);

    let { query } = this.chapiHandler.getEventData();
    query = normalizeQuery(query);

    let { user, token } = this.getCurrentUser().profile;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });

    try {
      let { results } = await this.credentialManager.query(token, query);
      this.presentation = results;

      this.prepareRecords(results);
      this.reasons = extractQueryReasons(query);
      this.presExchReasons = extractPresentationExchangeReasons(query);
    } catch (e) {
      this.errors.push('No credentials found matching requested criteria.');
      console.error('get credentials failed,:', e);
    }

    // TODO governance VC check

    this.requestOrigin = this.chapiHandler.getRequestor();

    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    handleError(e) {
      this.errors.push(e);
      this.loading = false;
    },
    prepareRecords: function (results) {
      let vcs = results.reduce((acc, val) => acc.concat(val.verifiableCredential), []);

      let _recordIt = (vc) => {
        let body;
        if (vc.proof && vc.proof.type == 'BbsBlsSignatureProof2020') {
          body = flatCredentialSubject(vc.credentialSubject);
        }

        return {
          title: vc.name ? vc.name : getCredentialType(vc.type),
          description: vc.description,
          icon: this.getVCIcon(vc.type),
          body,
        };
      };

      this.records = vcs.map(_recordIt);
    },
    async share() {
      this.loading = true;
      let { profile, preference } = this.getCurrentUser();
      let { controller, proofType, verificationMethod } = preference;
      let { domain, challenge } = this.chapiHandler.getEventData();

      let _present = async (presentation) => {
        return (
          await this.credentialManager.present(
            profile.token,
            { presentation },
            {
              controller,
              proofType,
              domain,
              challenge,
              verificationMethod,
            }
          )
        ).presentation;
      };

      let results = await Promise.all(this.presentation.map(_present));
      // typically single presentation, but some verifier queries might produce multiple presentation.
      if (results.length == 1) {
        this.chapiHandler.present(results[0]);
      } else {
        this.chapiHandler.present(results);
      }

      this.loading = false;
    },
    cancel() {
      this.chapiHandler.cancel();
    },
    getVCIcon(types) {
      return getVCIcon(getCredentialType(types));
    },
    displayContent(k) {
      let parts = k.split('.');
      return !nonDisplayContent.includes(parts[parts.length - 1]);
    },
  },
  computed: {
    showShareCredential() {
      return this.records.length > 0;
    },
  },
};
</script>
<style scoped>
.card {
  display: block;
  width: 360px;
  padding: 10px;
  background-color: #ffffff;
  border-radius: 7px;
  margin: 5px;
  text-align: center;
  line-height: 22px;
  cursor: pointer;
  position: relative;
  color: black;
  font-weight: 400;
  font-size: 16px;
  -webkit-box-shadow: 9px 10px 22px -8px rgba(209, 193, 209, 0.5);
  -moz-box-shadow: 9px 10px 22px -8px rgba(209, 193, 209, 0.5);
  box-shadow: 9px 10px 22px -8px rgba(209, 193, 209, 0.5);
  will-change: transform;
  user-select: none;
}

.card i {
  color: rgb(11, 151, 196) !important;
}

.cardContent {
  text-align: left;
}

.cardHeader {
  font-weight: 500;
  padding: 10px 15px;
}

.card-list li {
  list-style-type: none;
  padding: 10px 10px;
  transition: all 0.3s ease;
}

.card-list li:hover {
  transform: scale(1.1);
}

.reasons {
}

.reasons li {
  list-style: square;
  margin-left: 30px;
  list-style-type: '→';
}
</style>
