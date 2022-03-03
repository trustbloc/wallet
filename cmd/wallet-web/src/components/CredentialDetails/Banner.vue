<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start px-3 md:px-0 pt-10 w-full">
    <div class="flex flex-col justify-start items-start w-full root-container">
      <div
        class="
          flex
          z-10
          flex-row
          justify-start
          items-center
          py-6
          px-5
          w-full
          h-24
          rounded-xl
          border
        "
        :class="
          styles.background.color !== '#fff'
            ? `border-neutrals-black border-opacity-10`
            : `border-neutrals-thistle`
        "
        :style="`background-color: ${styles.background.color}`"
      >
        <div class="flex-none w-12 h-12 border-opacity-10">
          <img :src="credentialIcon" />
        </div>
        <span
          class="flex-1 pl-4 text-sm md:text-base font-bold text-left text-ellipsis"
          :style="`color: ${styles.text.color}`"
        >
          {{ title }}
        </span>
      </div>
      <div
        class="
          relative
          justify-start
          items-start
          px-6
          pt-14
          pb-6
          w-full
          bg-neutrals-white
          rounded-b-xl
          flex flex-col
          sub-container
        "
      >
        <table v-if="breakpoints.xs || breakpoints.sm" class="w-full text-left">
          <tr class="flex flex-1 mb-1 font-bold text-neutrals-dark">
            <th>
              {{ t('CredentialDetails.Banner.addedOn') }}
            </th>
          </tr>
          <tr class="flex flex-1 mb-4 text-neutrals-medium">
            <td>
              {{ addedOn }}
            </td>
          </tr>
          <tr class="flex flex-1 mb-1 font-bold text-neutrals-dark">
            <th>
              {{ t('CredentialDetails.Banner.expiresOn') }}
            </th>
          </tr>
          <tr class="flex flex-1 mb-4 text-neutrals-medium">
            <td>N/A</td>
          </tr>

          <tr class="flex flex-1 mb-1 font-bold text-neutrals-dark">
            <th>
              {{ t('CredentialDetails.Banner.lastUsed') }}
            </th>
          </tr>
          <tr class="flex flex-1 mb-4 text-neutrals-medium">
            <td>Never</td>
          </tr>

          <tr class="flex flex-1 mb-1 font-bold text-neutrals-dark">
            <th>
              {{ t('CredentialDetails.Banner.vault') }}
            </th>
          </tr>
          <tr class="flex flex-1 text-neutrals-medium">
            <td>{{ vaultName }}</td>
          </tr>
        </table>
        <table v-else class="w-full text-left">
          <thead class="font-bold text-neutrals-dark">
            <tr class="flex">
              <th class="flex-1">{{ t('CredentialDetails.Banner.addedOn') }}</th>
              <th class="flex-1">{{ t('CredentialDetails.Banner.expiresOn') }}</th>
              <th class="flex-1">{{ t('CredentialDetails.Banner.lastUsed') }}</th>
              <th class="flex-1">{{ t('CredentialDetails.Banner.vault') }}</th>
            </tr>
          </thead>
          <tbody class="text-neutrals-medium">
            <tr class="flex">
              <td class="flex-1">
                {{ addedOn }}
              </td>
              <td class="flex-1">N/A</td>
              <td class="flex-1">Never</td>
              <td class="flex-1">{{ vaultName }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import useBreakpoints from '@/plugins/breakpoints.js';
import { mapGetters } from 'vuex';
import { getCredentialIcon } from '@/mixins';
import { useI18n } from 'vue-i18n';

export default {
  name: 'Banner',
  props: {
    styles: {
      type: Object,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    issuanceDate: {
      type: String,
      required: true,
    },
    vaultName: {
      type: String,
      required: true,
    },
  },
  setup() {
    const { t, locale } = useI18n();
    return { t, locale };
  },
  data() {
    return {
      breakpoints: useBreakpoints(),
      credentialIcon: this.getCredentialIcon(),
    };
  },
  computed: {
    addedOn() {
      return new Date(this.issuanceDate).toLocaleDateString(this.locale, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      });
    },
  },
  methods: {
    ...mapGetters(['getStaticAssetsUrl']),
    getCredentialIcon: function () {
      return getCredentialIcon(this.getStaticAssetsUrl(), this.styles.thumbnail.uri);
    },
  },
};
</script>

<style scoped>
.root-container {
  height: calc(100% - 2.5rem);
}
.sub-container {
  top: -2.5rem;
  left: 0;
  /* TODO: replace with tailwind shadow once defined in config */
  box-shadow: 0px 2px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
