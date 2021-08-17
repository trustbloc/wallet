<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="content">
    <div class="content">
      <div class="md-layout">
        <div>
          <span v-if="loadingStatus === 'inprogress'">
            <skeleton-loader type="vault" />
          </span>
          <span v-else-if="loadingStatus === 'success'" id="dashboard-success-msg" class="px-24">
            <md-icon style="color: green" class="px-4">check_circle_outline</md-icon> Successfully
            setup your user for secured communication.
          </span>
          <span v-else-if="loadingStatus === 'failed'">
            <md-icon style="color: red" class="px-4">warning</md-icon>
            <b>Warning:</b> Failed to connect to server. Your wallet can not participate in secured
            communication.
          </span>
          <div v-if="cards.length" class="px-24 md-card" md-with-hover>
            <md-card-content>
              <ul class="credential-list">
                <li v-for="(card, index) in cards" :key="index" @click="toggleCard(card)">
                  <transition name="flip">
                    <div v-if="!card.flipped" :key="card.flipped" class="card">
                      <div class="cardContent">
                        <div class="cardHeader">
                          {{ credDisplayName(card.content) }}
                          <div v-for="cd in cjson" :key="cd.schema">
                            <div v-if="credDisplayName(card.content) === cd.schema">
                              <div class="cardBody">
                                <div class="cardDetailsL">
                                  <md-icon>{{ cd.icon }}</md-icon>
                                </div>
                                <div class="cardDetailsR">
                                  <p v-for="(value, key) in cd.credentialSubject" :key="key">
                                    {{ key }}: {{ value }}
                                  </p>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      <!--<json-modal :item="card.content" />-->
                    </div>
                    <div v-else :key="card.flipped" class="card">
                      <div class="cardContent cardBack">
                        <p>
                          Issuance Date:
                          {{
                            card.content.credentialSubject.issue_date ||
                            card.content.credentialSubject.issuedate ||
                            card.content.issuanceDate ||
                            'N/A'
                          }}
                        </p>
                        <p>
                          Expiration Date:
                          {{
                            card.content.credentialSubject.expiry_date ||
                            card.content.credentialSubject.cardexpires ||
                            card.content.expirationDate ||
                            'N/A'
                          }}
                        </p>
                      </div>
                    </div>
                  </transition>
                </li>
              </ul>
            </md-card-content>
          </div>
          <md-empty-state
            v-else
            md-icon="devices_other"
            :md-label="error"
            :md-description="errorDescription"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { getCredentialType } from '@/pages/mixins';
import { mapGetters } from 'vuex';
import SkeletonLoader from '../components/SkeletonLoader/SkeletonLoader';
import credentialDisplayData from '@/config/credentialDisplayData.json';

const filterBy = ['IssuerManifestCredential', 'GovernanceCredential'];
// TODO: issue-627 Add generic vue card for all the credentials to dynamically add support for all VC types.
export default {
  name: 'DashboardV2',
  components: {
    SkeletonLoader,
  },
  data() {
    return {
      cards: [],
      cjson: credentialDisplayData,
      username: '',
      agent: null,
      error: 'No stored credentials',
      errorDescription: "Your wallet is empty, there aren't any stored credentials to show.",
    };
  },
  computed: {
    loadingStatus() {
      return this.getCurrentUser().setupStatus;
    },
  },
  created: function () {
    let { user, token } = this.getCurrentUser().profile;
    this.username = this.getCurrentUser().username;

    let credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.fetchAllCredentials(credentialManager.getAll(token));
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts']),
    fetchAllCredentials: async function (getCredential) {
      let { contents } = await getCredential;
      console.log(`found ${Object.keys(contents).length} credentials`);
      const _filter = (id) => {
        return !contents[id].type.some((t) => filterBy.includes(t));
      };

      const _createCard = (id) => {
        return { content: contents[id], flipped: false };
      };

      this.cards = Object.keys(contents).filter(_filter).map(_createCard);

      console.log(`showing ${this.cards.length} credentials`);

      // Reading the values of the credentials and mapping it to the credential display data schemas
      this.cjson.forEach((obj) => {
        let flattened = {};
        for (let credentialKey in this.cards) {
          this.flatten(this.cards[credentialKey].content, flattened);
          for (let credentialContent in flattened) {
            if (
              obj.credentialSubject.hasOwnProperty(credentialContent) &&
              obj.schema === this.credDisplayName(this.cards[credentialKey].content)
            ) {
              obj.credentialSubject[credentialContent] = flattened[credentialContent];
            }
          }
        }
      });
    },
    credDisplayName: function (vc) {
      return vc.name ? vc.name : getCredentialType(vc.type);
    },
    toggleCard: function (card) {
      card.flipped = !card.flipped;
    },
    flatten: function (json, flattened) {
      for (let key in json) {
        if (json.hasOwnProperty(key)) {
          if (json[key] instanceof Object && json[key] != '') {
            this.flatten(json[key], flattened, key);
          } else {
            flattened[key] = json[key];
          }
        }
      }
    },
  },
};
</script>

<style scoped>
.md-card {
  display: inline-block !important;
  position: relative !important;
  width: 100% !important;
  margin: 25px 0 !important;
  overflow: unset !important;
  background: none !important;
  box-shadow: none !important;
  -webkit-box-shadow: none !important;
}

ul.credential-list {
  padding-left: 0;
  display: flex;
  flex-flow: row wrap;
}

li {
  list-style-type: none;
  padding: 10px 10px;
  transition: all 0.3s ease;
}

.card {
  display: block;
  width: 360px;
  height: 233px;
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

.cardBack {
  padding-top: 40px;
  color: rgba(0, 0, 0, 0.54);
}

li:hover {
  transform: scale(1.1);
}

.flip-enter-active {
  transition: all 0.4s ease;
}

.flip-leave-active {
  display: none;
}

.flip-enter,
.flip-leave {
  transform: rotateY(180deg);
  opacity: 0;
}

.md-dialog-container {
  width: 100% !important;
}
.cardBody {
  content: '""';
  display: table;
  clear: both;
  width: 100%;
  padding: 5px;
}

.cardDetailsL {
  float: left;
  width: 30%;
}

.cardDetailsR {
  float: left;
  width: 70%;
}

.cardDetailsR p {
  margin-bottom: 2px;
  color: rgba(0, 0, 0, 0.54);
}

.cardDetailsL i {
  font-size: 80px !important;
  padding-top: 20px;
  padding-left: 20px;
}
</style>
