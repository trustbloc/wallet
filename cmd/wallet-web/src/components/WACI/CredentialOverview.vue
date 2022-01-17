<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex relative flex-col justify-start items-start my-5 my:mt-6 w-full root-container">
    <div
      :class="[
        `flex
          z-10
          flex-row
          justify-start
          items-center
          py-4
          px-5
          w-full
          h-auto
          bg-neutrals-white
          rounded-xl
          border `,
        credential.brandColor.length
          ? `bg-gradient-${credential.brandColor} border-neutrals-black border-opacity-10`
          : `bg-neutrals-white border-neutrals-thistle`,
      ]"
    >
      <div class="flex-none w-12 h-12 border-opacity-10">
        <img :src="credentialIconSrc" />
      </div>
      <span
        :class="[
          `flex-1 pl-4 font-bold text-left text-sm overflow-ellipsis`,
          credential.brandColor.length ? `text-neutrals-white` : `text-neutrals-dark`,
        ]"
        >{{ credential.title }}</span
      >
    </div>
    <div
      class="
        absolute
        justify-start
        items-start
        px-5
        pb-3
        w-full
        rounded-b-xl
        pt-13
        bg-neutrals-white
        flex flex-row
        sub-container
      "
    >
      <span class="flex text-sm font-bold text-neutrals-dark">
        {{ t('CredentialDetails.Banner.vault') }}
      </span>
      <span class="flex ml-3 text-sm text-neutrals-medium">
        {{ credential.vaultName }}
      </span>
    </div>
  </div>
  <!-- Credential Details -->
  <div class="flex flex-col justify-start items-start w-full">
    <span class="py-4 pl-3 text-xl font-bold text-neutrals-dark">{{
      t('WACI.Share.whatIsShared')
    }}</span>
    <table class="w-full border-t border-neutrals-chatelle">
      <tr
        v-for="(property, key) of credential.properties"
        :key="key"
        class="border-b border-dotted border-neutrals-thistle"
      >
        <td class="py-4 pr-6 pl-3 w-2/5 whitespace-nowrap text-neutrals-medium">
          {{ property.label }}
        </td>
        <td
          v-if="property.type === 'image'"
          class="py-4 pr-6 pl-3 w-3/5 break-words text-neutrals-dark"
        >
          <img :src="property.value" class="w-20 h-20" />
        </td>
        <td v-else class="py-4 pr-6 pl-3 w-3/5 break-words text-neutrals-dark">
          {{ property.value }}
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';
import { mapGetters } from 'vuex';
import { getCredentialIcon } from '@/mixins';

export default {
  name: 'CredentialOverview',
  props: {
    credential: {
      type: Object,
      required: true,
    },
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  computed: {
    credentialIconSrc() {
      return this.getCredentialIcon();
    },
  },
  methods: {
    ...mapGetters(['getStaticAssetsUrl']),
    // Get credential icon based on docker configuration
    getCredentialIcon: function () {
      return getCredentialIcon(this.getStaticAssetsUrl(), this.credential.icon);
    },
  },
};
</script>

<style scoped>
.root-container {
  padding-bottom: 2.5rem;
}
.sub-container {
  top: 2.5rem;
  left: 0;
  /* TODO: replace with tailwind shadow once defined in config */
  box-shadow: 0px 2px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
