<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <button
    class="group credentialPreviewContainer inline-flex w-full items-center rounded-xl border py-4 pr-3 pl-5 text-sm font-bold outline-none focus-within:ring-2 focus-within:ring-offset-2 md:text-base"
    :class="
      styles.background.color !== '#fff'
        ? `border-neutrals-chatelle border-opacity-10 focus-within:ring-[#${styles.background.color}]`
        : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`
    "
    :style="`background-color: ${styles.background.color}`"
  >
    <div class="h-12 w-12 flex-none border-opacity-10">
      <img :src="credentialIconSrc" />
    </div>
    <div class="flex grow p-4">
      <span
        class="text-ellipsis text-left text-sm font-bold"
        :style="`color: ${styles.text.color}`"
      >
        {{ title }}
      </span>
    </div>
    <div
      :class="[
        `flex-none w-8 h-8 rounded-full`,
        styles.background.color !== '#fff'
          ? `bg-neutrals-black bg-opacity-25 group-hover:bg-opacity-60`
          : `bg-neutrals-thistle`,
      ]"
    >
      <div class="p-1">
        <img
          :src="
            styles.background.color !== '#fff'
              ? require('@/assets/img/credential--arrow-right-icon-light.svg')
              : require('@/assets/img/credential--arrow-right-icon.svg')
          "
        />
      </div>
    </div>
  </button>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';
import { getCredentialIcon } from '@/mixins';

export default {
  name: 'CredentialBannerComponent',
  props: {
    id: {
      type: String,
      required: true,
    },
    styles: {
      type: Object,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const store = useStore();
    const getStaticAssetsUrl = () => store.getters.getStaticAssetsUrl;
    const credentialIconSrc = computed(() =>
      props?.styles?.thumbnail?.uri?.includes('https://')
        ? props?.styles?.thumbnail?.uri
        : getCredentialIcon(getStaticAssetsUrl(), props?.styles?.thumbnail?.uri)
    );
    return { credentialIconSrc };
  },
};
</script>

<style scoped>
.credentialPreviewContainer:not(:focus-within):hover {
  box-shadow: 0px 4px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
