<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div v-bind="$attrs" class="flex relative flex-col justify-start items-start w-full">
    <div
      class="flex z-10 flex-row justify-start items-center py-4 px-5 w-full h-auto bg-neutrals-white rounded-xl border"
      :class="
        credential?.styles.background.color !== '#fff'
          ? `border-neutrals-black border-opacity-10`
          : `border-neutrals-thistle`
      "
      :style="`background-color: ${credential?.styles.background.color}`"
    >
      <div class="flex-none w-12 h-12 border-opacity-10">
        <img :src="credentialIconSrc" />
      </div>
      <span
        class="flex-1 pl-4 text-sm font-bold text-left text-ellipsis"
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
