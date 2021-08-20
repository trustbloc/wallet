<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="">
    <div>
      <span v-if="loadingStatus === 'inprogress'">
        <skeleton-loader type="vault" />
      </span>
      <span v-else-if="loadingStatus === 'success'" id="dashboard-success-msg" class="px-8">
        <md-icon style="color: green" class="px-4">check_circle_outline</md-icon> Successfully setup
        your user for secured communication.
      </span>
      <span v-else-if="loadingStatus === 'failed'">
        <md-icon style="color: red" class="px-4">warning</md-icon>
        <b>Warning:</b> Failed to connect to server. Your wallet can not participate in secured
        communication.
      </span>
      <div v-if="cards.length">
        <ul class="grid grid-cols-2 gap-8">
          <li v-for="(card, index) in cards" :key="index" class="focus:ring-8 cursor-pointer">
            <!--Temporary solution to handle preview for the generic credential, this will be refactored issue-981-->
            <div class="flex justify-between credentialCard">
              <div class="flex flex-row flex-grow items-center">
                <div class="border-opacity-10 credentialLogoContainer">
                  <img src="@/assets/img/credential--generic-icon.svg" />
                </div>
                <div class="credentialHeader">
                  {{ credDisplayName(card.content) }}
                </div>
              </div>
              <div class="py-2">
                <div class="credentialArrowContainer">
                  <div class="credentialArrowLogo">
                    <button>
                      <img src="@/assets/img/credential--arrow-right-icon.svg" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>
      <md-empty-state
        v-else
        md-icon="devices_other"
        :md-label="error"
        :md-description="errorDescription"
      />
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
const images = require.context('@/assets/img/', false, /\.png$|\.jpg$|\.svg$/);
// TODO: issue-627 Add generic vue card for all the credentials to dynamically add support for all VC types.
export default {
  name: 'Dashboard',
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
      return this.getCurrentUser() ? this.getCurrentUser().setupStatus : null;
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
    loadImg(imgPath) {
      return images('./' + imgPath);
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

<style scoped></style>
