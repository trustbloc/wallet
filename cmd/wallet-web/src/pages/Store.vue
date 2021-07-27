<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex justify-center w-screen content">
    <div class="max-w-screen-xl md-layout">
      <div class="md-layout-item">
        <form>
          <md-card>
            <md-card-header style="background-color: #00bcd4">
              <h3 class="title">
                <md-icon>fingerprint</md-icon>
                Credential
              </h3>
            </md-card-header>

            <md-card-content
              v-if="records.length"
              class="md-layout md-alignment-center-center card-list"
            >
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

            <md-card-content>
              <div v-if="errors.length">
                <b>Please correct the following error(s):</b>
                <ul>
                  <li v-for="error in errors" :key="error" style="color: #9d0006">{{ error }}</li>
                </ul>
              </div>
              <div class="md-layout md-alignment-center-center">
                <div>
                  <md-button id="cancelBtn" class="md-cancel-text" @click="cancel"
                    >Cancel
                  </md-button>
                </div>
                <div>
                  <md-button
                    id="storeVCBtn"
                    class="md-raised md-success"
                    :disabled="isDisabled"
                    @click="store"
                    >Confirm
                  </md-button>
                </div>
              </div>
            </md-card-content>
          </md-card>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { CHAPIEventHandler, getCredentialType, getVCIcon, isVPType } from './mixins';
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';

export default {
  data() {
    return {
      records: [],
      storeButton: true,
      subject: '',
      issuer: '',
      issuance: '',
      errors: [],
    };
  },
  computed: {
    isDisabled() {
      return this.storeButton;
    },
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
      this.records = data.verifiableCredential.map((vc) => {
        return {
          title: vc.name ? vc.name : getCredentialType(vc.type),
          description: vc.description,
          icon: getVCIcon(vc.type),
        };
      });
      console.log('this.records', JSON.stringify(this.records, null, 2));
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
  },
};
</script>

<style scoped>
.card {
  display: block;
  width: 360px;
  padding: 10px;
  background-color: whitesmoke;
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
</style>
