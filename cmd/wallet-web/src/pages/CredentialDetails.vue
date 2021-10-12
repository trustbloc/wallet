<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div class="flex flex-col justify-start items-start py-6 px-3">
    <div class="flex flex-row justify-between items-center mb-4 w-full">
      <div class="flex flex-grow">
        <!-- TODO: remove m-0 and font-bold once vue-material is scrapped -->
        <h3 class="m-0 font-bold text-neutrals-dark">{{ i18n.heading }}</h3>
      </div>
      <div
        class="inline-flex items-center rounded-lg border border-neutrals-chatelle bg-transparent"
      >
        <flyout-menu type="outline" />
      </div>
    </div>
    <banner
      :brand-color="credential.brandColor"
      :icon="credential.icon"
      :title="credential.title"
      :issuance-date="credential.issuanceDate"
    />
    <!-- List of Credential Details -->
    <div class="flex flex-col justify-start items-start mt-8 md:mt-2 w-full">
      <span class="mb-5 text-xl font-bold text-neutrals-dark">{{ i18n.verifiedInformation }}</span>
      <table class="w-full border-t border-neutrals-chatelle">
        <tr
          v-for="(property, index) of credential.properties"
          :key="index"
          class="border-b border-neutrals-thistle border-dotted"
        >
          <!-- TODO: Add the dropdown for the nested credentials 1016 -->
          <td class="py-4 pr-6 pl-3 text-neutrals-medium">{{ property.label }}</td>
          <td v-if="property.type != 'image'" class="py-4 pr-6 pl-3 text-neutrals-dark break-words">
            {{ property.value }}
          </td>
          <td
            v-if="property.type === 'image'"
            class="py-4 pr-6 pl-3 text-neutrals-dark break-words"
          >
            <img :src="property.value" class="w-20 h-20" />
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import Banner from '@/components/CredentialDetails/Banner';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';
import { mapGetters } from 'vuex';

export default {
  name: 'CredentialDetails',
  components: {
    Banner,
    FlyoutMenu,
  },
  computed: {
    i18n() {
      return this.$t('CredentialDetails');
    },
    ...mapGetters(['getProcessedCredentialByID']),
    credential() {
      return this.getProcessedCredentialByID(this.$route.params.id);
    },
  },
};
</script>
