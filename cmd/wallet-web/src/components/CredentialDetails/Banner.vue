<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="px-3 md:px-0 w-full">
    <div class="flex relative flex-col justify-start items-start w-full h-48">
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
          bg-neutrals-white
          rounded-xl
          border border-neutrals-thistle
        "
      >
        <div class="flex-none w-12 h-12 border-opacity-10">
          <img src="@/assets/img/credential--generic-icon.svg" />
        </div>
        <span class="flex-1 pl-4 font-bold text-left text-neutrals-dark overflow-ellipsis">{{
          title
        }}</span>
      </div>
      <div
        class="
          absolute
          justify-start
          items-start
          px-6
          pt-11
          pb-6
          w-full
          bg-neutrals-white
          rounded-b-xl
          flex flex-col
          sub-container
        "
      >
        <table class="w-full text-left">
          <thead class="font-bold text-neutrals-dark">
            <tr class="flex">
              <th class="flex-1">{{ i18n.addedOn }}</th>
              <th class="flex-1">{{ i18n.expiresOn }}</th>
              <th class="flex-1">{{ i18n.lastUsed }}</th>
              <th class="flex-1">{{ i18n.vault }}</th>
            </tr>
          </thead>
          <tbody class="text-neutrals-medium">
            <tr class="flex">
              <td class="flex-1">
                {{ addedOn }}
              </td>
              <td class="flex-1">N/A</td>
              <td class="flex-1">Never</td>
              <td class="flex-1">Default Vault</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Banner',
  props: {
    title: {
      type: String,
      required: true,
    },
    issuanceDate: {
      type: String,
      required: true,
    },
  },
  computed: {
    i18n() {
      return this.$t('CredentialDetails.Banner');
    },
    addedOn() {
      return new Date(this.issuanceDate).toLocaleDateString(this.$i18n.locale, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      });
    },
  },
};
</script>

<style scoped>
.sub-container {
  top: 4.75rem; /* 76px */
  left: 0;
  height: 7.3125rem; /* 117px */
  /* TODO: replace with tailwind shadow once defined in config */
  box-shadow: 0px 2px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
