<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div v-bind="$attrs" class="relative flex w-full flex-col items-start justify-start">
    <div
      class="z-10 flex h-auto w-full flex-row items-center justify-start rounded-xl border bg-neutrals-white py-4 px-5"
      :class="
        credential?.styles.background.color !== '#fff'
          ? `border-neutrals-black border-opacity-10`
          : `border-neutrals-thistle`
      "
      :style="`background-color: ${credential?.styles.background.color}`"
    >
      <div class="h-12 w-12 flex-none border-opacity-10">
        <img :src="credentialIconSrc" />
      </div>
      <span
        class="flex-1 text-ellipsis pl-4 text-left text-sm font-bold"
        :style="`color: ${credential?.styles.text.color}`"
        >{{ credentialHeading }}</span
      >
    </div>
    <slot name="bannerBottomContainer" />
  </div>
  <slot name="credentialDetails" />
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';
import { getCredentialIcon } from '@/mixins';

export default {
  name: 'CredentialOverviewComponent',
  props: {
    credential: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const store = useStore();
    const getStaticAssetsUrl = () => store.getters.getStaticAssetsUrl;
    const credentialIconSrc = computed(() =>
      props?.credential?.styles?.thumbnail?.uri?.includes('https://')
        ? props?.credential?.styles?.thumbnail?.uri
        : getCredentialIcon(getStaticAssetsUrl(), props?.credential?.styles?.thumbnail?.uri)
    );
    const credentialHeading = computed(() =>
      props?.credential?.name?.length ? props?.credential?.name : props?.credential?.title
    );
    return { credentialIconSrc, credentialHeading };
  },
};
</script>
