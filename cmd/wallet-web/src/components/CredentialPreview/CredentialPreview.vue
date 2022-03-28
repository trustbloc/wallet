<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <router-link
    class="
      group
      inline-flex
      items-center
      py-6
      pr-3
      pl-5
      w-full
      h-20
      md:h-24
      text-sm
      md:text-base
      font-bold
      rounded-xl
      border
      focus-within:ring-2 focus-within:ring-offset-2
      outline-none
      credentialPreviewContainer
    "
    :class="
      styles.background.color !== '#fff'
        ? `border-neutrals-black border-opacity-10 focus-within:ring-primary-${styles?.background?.color}`
        : `border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`
    "
    :style="`background-color: ${styles?.background?.color}`"
  >
    <div class="flex-none w-12 h-12 border-opacity-10">
      <img :src="credentialIconSrc" />
    </div>
    <div class="flex-grow p-4">
      <span
        class="text-sm md:text-base font-bold text-left text-ellipsis"
        :style="`color: ${styles?.text?.color}`"
      >
        {{ title }}
      </span>
    </div>
    <div
      :class="[
        `flex-none w-8 h-8 rounded-full`,
        styles?.background?.color !== '#fff'
          ? `bg-neutrals-black bg-opacity-25 group-hover:bg-opacity-60`
          : `bg-neutrals-light bg-opacity-25`,
      ]"
    >
      <div class="p-1">
        <img
          :src="
            styles?.background?.color !== '#fff'
              ? require('@/assets/img/credential--arrow-right-icon-light.svg')
              : require('@/assets/img/credential--arrow-right-icon.svg')
          "
        />
      </div>
    </div>
  </router-link>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';
import { getCredentialIcon } from '@/mixins';

export default {
  name: 'CredentialPreview',
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
