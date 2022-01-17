<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <button
    :class="[
      `group inline-flex items-center rounded-xl p-5 text-sm md:text-base font-bold border w-full h-20 md:h-24 focus-within:ring-2 focus-within:ring-offset-2 credentialPreviewContainer`,
      credential.brandColor.length
        ? `bg-gradient-${credential.brandColor} border-neutrals-chatelle border-opacity-10 focus-within:ring-primary-${credential.brandColor}`
        : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`,
    ]"
    @click="toggleDetails(credential)"
  >
    <div class="flex-none w-12 h-12 border-opacity-10">
      <img :src="credentialIconSrc" />
    </div>
    <div class="flex flex-grow p-4">
      <span
        :class="[
          `text-sm md:text-base font-bold text-left overflow-ellipsis`,
          credential.brandColor.length ? `text-neutrals-white` : `text-neutrals-dark`,
        ]"
      >
        {{ credential.title }}
      </span>
    </div>
  </button>
</template>

<script>
import { mapGetters } from 'vuex';
import { getCredentialIcon } from '@/mixins';

export default {
  name: 'CredentialBanner',
  props: {
    id: {
      type: String,
      required: true,
    },
    brandColor: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
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
      return getCredentialIcon(this.getStaticAssetsUrl(), this.icon);
    },
  },
};
</script>

<style scoped>
.credentialPreviewContainer:not(:focus-within):hover {
  box-shadow: 0px 4px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
