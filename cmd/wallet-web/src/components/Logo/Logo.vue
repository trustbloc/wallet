/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <a v-bind="$attrs" class="flex justify-start items-center border-0">
    <img class="w-48" :src="logoUrl" alt="logo" />
  </a>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  name: 'Logo',
  inheritAttrs: false,
  data() {
    return {
      logoUrl: this.getLogoUrl(),
    };
  },
  methods: {
    ...mapGetters(['getStaticAssetsUrl']),
    // Get logo url based on docker configuration
    getLogoUrl: function () {
      let staticAssetsUrl = this.getStaticAssetsUrl();
      if (staticAssetsUrl) {
        return (this.logoUrl = `${staticAssetsUrl}/images/logo.svg`);
      }
      return (this.logoUrl = `${require('@/assets/img/logo.svg')}`);
    },
  },
};
</script>
