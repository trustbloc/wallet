<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Dashboard Layout -->
    <div class="flex md:hidden flex-col justify-start w-screen">
      <!-- Todo move this to resuable component issue-985 -->
      <div class="bg-neutrals-white border border-neutrals-chatelle">
        <flyout-menu />
      </div>
      <div class="items-start">
        <h3 class="mx-6 mb-5 font-bold text-neutals-dark">{{ i18n.credentials }}</h3>
      </div>
    </div>
    <!-- Desktop Dashboard Layout -->
    <!-- Todo move this to seprate component as required -->
    <div class="hidden md:flex justify-between items-center mb-8 w-full align-middle">
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
    <div class="mx-6 md:mx-0 mb-5">
      <span class="font-bold font-xl text-neutals-dark">{{ i18n.defaultvault }}</span>
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
      <ul class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-8 mx-6 md:mx-0">
        <li v-for="(card, index) in cards" :key="index">
          <!--Temporary solution to handle preview for the generic credential, this will be refactored issue-981-->
          <div class="inline-flex items-center w-full credentialCard">
            <div class="flex-none w-12 h-12 border-opacity-10">
              <img src="@/assets/img/credential--generic-icon.svg" />
            </div>
            <div class="flex-grow p-4 text-left text-neutrals-dark overflow-ellipsi">
              {{ credDisplayName(card.content) }}
            </div>
            <div class="flex-none credentialArrowContainer">
              <div class="p-1">
                <button>
                  <img src="@/assets/img/credential--arrow-right-icon.svg" />
                </button>
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
