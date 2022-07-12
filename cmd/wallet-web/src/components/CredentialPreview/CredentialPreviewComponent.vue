<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <router-link
    :id="id"
    class="group inline-flex items-center py-6 pr-3 pl-5 w-full h-20 text-sm font-bold rounded-xl border outline-none focus-within:ring-2 focus-within:ring-offset-2 md:h-24 md:text-base credentialPreviewContainer"
    :class="
      styles.background.color !== '#fff'
        ? `border-neutrals-black border-opacity-10 notWhiteCredentialPreview`
        : `border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`
    "
    :style="credentialStyles"
  >
    <div class="flex-none w-12 h-12 border-opacity-10">
      <img :src="credentialIconSrc" />
    </div>
    <div class="grow p-4">
      <span
        class="text-sm font-bold text-left text-ellipsis md:text-base"
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
  name: 'CredentialPreviewComponent',
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
  computed: {
    credentialStyles() {
      return {
        'background-color': this.styles?.background?.color,
        '--focus-color': this.styles?.background?.color,
      };
    },
  },
};
</script>

<style scoped>
.credentialPreviewContainer:not(:focus-within):hover {
  box-shadow: 0px 4px 12px 0px rgba(25, 12, 33, 0.1);
}

.notWhiteCredentialPreview:focus {
  outline: 2px solid var(--focus-color);
  outline-offset: 2px;
}
</style>
