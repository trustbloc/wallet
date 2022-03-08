<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <router-link
    :class="[
      `group outline-none inline-flex items-center rounded-xl py-6 pl-5 pr-3 text-sm md:text-base font-bold border w-full h-20 md:h-24 focus-within:ring-2 focus-within:ring-offset-2 credentialPreviewContainer`,
      brandColor.length
        ? `bg-gradient-${brandColor} border-neutrals-black border-opacity-10 focus-within:ring-primary-${brandColor}`
        : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`,
    ]"
  >
    <div class="flex-none w-12 h-12 border-opacity-10">
      <img :src="credentialIcon" />
    </div>
    <div class="flex-grow p-4">
      <span
        :class="[
          `text-sm md:text-base font-bold text-left text-ellipsis`,
          brandColor.length ? `text-neutrals-white` : `text-neutrals-dark`,
        ]"
      >
        {{ title }}
      </span>
    </div>
    <div
      :class="[
        `flex-none w-8 h-8 rounded-full`,
        brandColor.length
          ? `bg-neutrals-black bg-opacity-25 group-hover:bg-opacity-60`
          : `bg-neutrals-thistle`,
      ]"
    >
      <div class="p-1">
        <img
          :src="
            brandColor.length
              ? require('@/assets/img/credential--arrow-right-icon-light.svg')
              : require('@/assets/img/credential--arrow-right-icon.svg')
          "
        />
      </div>
    </div>
  </router-link>
</template>

<script>
import { mapGetters } from 'vuex';
import { getCredentialIcon } from '@/mixins';

export default {
  name: 'CredentialPreview',
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
  data() {
    return {
      credentialIcon: this.getCredentialIcon(),
    };
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
