<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start px-3 pt-10 w-full md:px-0">
    <div class="flex flex-col justify-start items-start w-full root-container">
      <div
        class="z-10 flex h-24 w-full flex-row items-center justify-start rounded-xl border py-6 px-5"
        :class="
          styles.background.color !== '#fff'
            ? `border-neutrals-black border-opacity-10`
            : `border-neutrals-thistle`
        "
        :style="`background-color: ${styles.background.color}`"
      >
        <div class="h-12 w-12 flex-none border-opacity-10">
          <img :src="credentialIconSrc" />
        </div>
        <span
          class="flex-1 pl-4 text-sm font-bold text-left text-ellipsis md:text-base"
          :style="`color: ${styles.text.color}`"
        >
          {{ title }}
        </span>
      </div>
      <div
        class="sub-container relative flex w-full flex-col items-start justify-start rounded-b-xl bg-neutrals-white px-6 pt-14 pb-6"
      >
        <table v-if="breakpoints.xs || breakpoints.sm" class="w-full text-left">
          <tr class="mb-1 flex flex-1 font-bold text-neutrals-dark">
            <th>
              {{ t('CredentialDetails.Banner.addedOn') }}
            </th>
          </tr>
          <tr class="mb-4 flex flex-1 text-neutrals-medium">
            <td>
              {{ addedOn }}
            </td>
          </tr>
          <tr class="mb-1 flex flex-1 font-bold text-neutrals-dark">
            <th>
              {{ t('CredentialDetails.Banner.expiresOn') }}
            </th>
          </tr>
          <tr class="mb-4 flex flex-1 text-neutrals-medium">
            <td>N/A</td>
          </tr>

          <tr class="mb-1 flex flex-1 font-bold text-neutrals-dark">
            <th>
              {{ t('CredentialDetails.Banner.lastUsed') }}
            </th>
          </tr>
          <tr class="mb-4 flex flex-1 text-neutrals-medium">
            <td>Never</td>
          </tr>

          <tr class="mb-1 flex flex-1 font-bold text-neutrals-dark">
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
import { computed } from 'vue';
import { useStore } from 'vuex';
import useBreakpoints from '@/plugins/breakpoints.js';
import { useI18n } from 'vue-i18n';
import { getCredentialIcon } from '@/mixins';

export default {
  name: 'CredentialDetailsBannerComponent',
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
  setup(props) {
    const store = useStore();
    const { t, locale } = useI18n();
    const breakpoints = useBreakpoints();
    const getStaticAssetsUrl = () => store.getters.getStaticAssetsUrl;
    const credentialIconSrc = computed(() =>
      props?.styles?.thumbnail?.uri?.includes('https://')
        ? props?.styles?.thumbnail?.uri
        : getCredentialIcon(getStaticAssetsUrl(), props?.styles?.thumbnail?.uri)
    );
    const addedOn = computed(() =>
      new Date(props?.issuanceDate).toLocaleDateString(locale, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    );
    return { addedOn, breakpoints, credentialIconSrc, t, locale };
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
