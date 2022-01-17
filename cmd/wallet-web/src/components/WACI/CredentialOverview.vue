<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div v-bind="$attrs" class="flex relative flex-col justify-start items-start w-full">
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
          border`,
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
    <slot name="bannerBottomContainer" />
  </div>
  <slot name="credentialDetails" />
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
