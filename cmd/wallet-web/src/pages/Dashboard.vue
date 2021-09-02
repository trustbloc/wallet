<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Dashboard Layout -->
    <div v-if="breakpoints.xs || breakpoints.sm" class="flex flex-col justify-start w-screen">
      <!-- Todo move this to reusable component issue-985 -->
      <div class="bg-neutrals-white border border-neutrals-chatelle">
        <flyout-menu />
      </div>
      <div class="items-start">
        <h3 class="mx-6 mb-5 font-bold text-neutrals-dark">{{ i18n.credentials }}</h3>
      </div>
    </div>
    <!-- Desktop Dashboard Layout -->
    <div v-else class="flex justify-between items-center mb-8 w-full align-middle">
      <div class="flex flex-grow">
        <h3 class="m-0 font-bold text-neutrals-dark">{{ i18n.credentials }}</h3>
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
    <div v-if="processedCredentials.length" class="mx-6 md:mx-0">
      <div class="md:mx-0 mb-5">
        <span class="text-xl font-bold text-neutrals-dark">{{ i18n.defaultVault }}</span>
      </div>
      <ul class="grid grid-cols-1 xl:grid-cols-2 gap-4 xl:gap-8 my-8">
        <li v-for="(processedCredential, index) in processedCredentials" :key="index">
          <router-link
            :to="{ name: 'credential-details', params: { id: processedCredential.id } }"
            class="inline-flex items-center credentialCard"
          >
            <div class="flex-none w-12 h-12 border-opacity-10">
              <img :src="require(`@/assets/img/${processedCredential.icon}`)" />
            </div>
            <div class="flex-grow p-4">
              <span
                class="
                  text-sm
                  md:text-base
                  font-bold
                  overflow-ellipsis
                  text-left text-neutrals-dark
                "
              >
                {{ processedCredential.title }}
              </span>
            </div>
            <div class="flex-none credentialArrowContainer">
              <div class="p-1">
                <img src="@/assets/img/credential--arrow-right-icon.svg" />
              </div>
            </div>
          </router-link>
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
import { getCredentialType, getCredentialDisplayData } from '@/pages/mixins';
import { mapGetters } from 'vuex';
import SkeletonLoader from '@/components/SkeletonLoader/SkeletonLoader';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';
import credentialDisplayData from '@/config/credentialDisplayData';
import useBreakpoints from '@/plugins/breakpoints.js';

const filterBy = ['IssuerManifestCredential', 'GovernanceCredential'];
// TODO: issue-627 Add generic vue card for all the credentials to dynamically add support for all VC types.
export default {
  name: 'Dashboard',
  components: {
    SkeletonLoader,
    FlyoutMenu,
  },
  data() {
    return {
      processedCredentials: [],
      username: '',
      agent: null,
      breakpoints: useBreakpoints(),
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
      const { contents } = await getCredential;
      console.log(`found ${Object.keys(contents).length} credentials`);
      const _filter = (id) => {
        return !contents[id].type.some((t) => filterBy.includes(t));
      };

      const credentials = Object.keys(contents)
        .filter(_filter)
        .map((id) => contents[id]);

      console.debug(`showing ${credentials.length} credentials`);

      credentials.map((credential) => {
        const manifest = this.getManifest(credential);
        const processedCredential = this.getCredentialDisplayData(credential, manifest);
        this.processedCredentials.push(processedCredential);
      });
      console.debug(`showing ${this.processedCredentials.length} credentials`);
    },
    getCredentialType: function (vc) {
      return getCredentialType(vc.type);
    },
    credDisplayName: function (vc) {
      return vc.name ? vc.name : getCredentialType(vc.type);
    },
    getCredentialDisplayData: function (vc, manifestCredential) {
      return getCredentialDisplayData(vc, manifestCredential);
    },
    getManifest: function (credential) {
      const currentCredentialType = this.getCredentialType(credential);
      return credentialDisplayData[currentCredentialType] || credentialDisplayData.fallback;
    },
  },
};
</script>
