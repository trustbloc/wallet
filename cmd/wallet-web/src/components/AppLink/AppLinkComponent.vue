<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <slot />
  <a v-if="isExternalLink" v-bind="$attrs" :href="to" target="_blank" />
  <router-link v-else v-slot="{ isActive }" v-bind="$props" custom class="">
    <a
      v-bind="$attrs"
      :class="isActive ? activeClass : inactiveClass"
      class="absolute top-0 left-0 z-0 h-full w-full cursor-pointer"
    />
  </router-link>
</template>

<script>
import { RouterLink } from 'vue-router';

export default {
  name: 'AppLinkComponent',
  inheritAttrs: false,
  props: {
    ...RouterLink.props,
    inactiveClass: {
      type: String,
      default: '',
    },
  },
  computed: {
    isExternalLink() {
      return typeof this.to === 'string' && this.to.startsWith('http');
    },
  },
};
</script>
