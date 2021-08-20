<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <div>
      <!-- Todo move this to seprate component as required -->
      <div class="flex flex-row justify-between items-center align-middle">
        <div class="flex flex-grow">
          <h3 class="m-0 font-bold text-neutals-dark">{{ i18n.credentials }}</h3>
        </div>
        <div
          class="
            inline-flex
            items-center
            bg-neutrals-white
            rounded-lg
            border border-neutrals-chatelle
          "
        >
          <flyout-menu />
        </div>
      </div>
      <span v-if="loadingStatus === 'inprogress'">
        <skeleton-loader type="vault" />
      </span>
      <span v-else-if="loadingStatus === 'success'" id="dashboard-success-msg" class="px-8"> </span>
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
      <div
        v-else-if="loadingStatus === 'success'"
        class="py-8 px-6 mx-auto rounded-lg border border-neutrals-thistle nocredentialCard"
      >
        <div class="flex justify-center">
          <img src="@/assets/img/icons-md--credentials-icon.svg" />
        </div>
        <div class="flex justify-center">
          <span class="text-base font-bold text-neutrals-medium"> {{ i18n.error }}</span>
        </div>
        <div class="flex justify-center">
          <span class="text-base text-neutrals-medium"> {{ i18n.description }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { getCredentialType } from '@/pages/mixins';
import { mapGetters } from 'vuex';
import SkeletonLoader from '@/components/SkeletonLoader/SkeletonLoader';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';
import credentialDisplayData from '@/config/credentialDisplayData.json';

const filterBy = ['IssuerManifestCredential', 'GovernanceCredential'];
const images = require.context('@/assets/img/', false, /\.png$|\.jpg$|\.svg$/);
// TODO: issue-627 Add generic vue card for all the credentials to dynamically add support for all VC types.
export default {
  name: 'Dashboard',
  components: {
    SkeletonLoader,
    FlyoutMenu,
  },
  data() {
    return {
      cards: [],
      cjson: credentialDisplayData,
      username: '',
      agent: null,
    };
  },
  computed: {
    loadingStatus() {
      return this.getCurrentUser() ? this.getCurrentUser().setupStatus : null;
    },
    i18n() {
      return this.$t('Credentials');
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
