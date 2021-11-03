<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div class="flex flex-col justify-start items-start py-6 px-3">
    <div class="flex flex-row justify-between items-center mb-4 w-full">
      <div class="flex flex-grow">
        <h3 class="text-neutrals-dark">{{ i18n.heading }}</h3>
      </div>
      <div
        class="inline-flex items-center rounded-lg border border-neutrals-chatelle bg-transparent"
      >
        <flyout-menu type="outline">
          <flyout-menu-list>
            <flyout-menu-button
              id="renameCredential"
              :text="i18n.renameCredential"
              color="neutrals-medium"
            />
            <flyout-menu-button
              id="moveCredential"
              :text="i18n.moveCredential"
              color="neutrals-medium"
            />
            <flyout-menu-button
              id="deleteCredential"
              :text="i18n.deleteCredential"
              color="primary-vampire"
              @click="openDeleteCredential"
            >
            </flyout-menu-button>
          </flyout-menu-list>
        </flyout-menu>
        <ModalRoot />
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
import { ModalBus } from '@/EventBus';
import ModalRoot from '@/components/Modal/Modal';
import DeleteCredential from '@/components/Modal/Credential/Delete';
import Banner from '@/components/CredentialDetails/Banner';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';
import FlyoutMenuList from '@/components/FlyoutMenu/FlyoutMenuList';
import FlyoutMenuButton from '@/components/FlyoutMenu/FlyoutMenuButton';

import { mapGetters } from 'vuex';

export default {
  name: 'CredentialDetails',
  components: {
    Banner,
    FlyoutMenu,
    FlyoutMenuList,
    FlyoutMenuButton,
    ModalRoot,
  },
  data: function () {
    return {
      modalAuth: false,
    };
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
  methods: {
    // Add modal bus event for other modal like rename, move credential etc
    openDeleteCredential() {
      ModalBus.$emit('open', {
        component: DeleteCredential,
        props: { credentialId: this.credential.id },
      });
    },
  },
};
</script>
