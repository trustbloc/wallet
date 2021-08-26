<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <li
    class="
      justify-start
      items-center
      font-bold
      hover:bg-gradient-to-r hover:opacity-100 hover:from-neutrals-black
      focus-within:bg-gradient-to-r
      focus-within:opacity-100
      focus-within:from-neutrals-black
      focus-within:shadow-inner-outline-blue
      opacity-60
      bar
      flex flex-row
    "
  >
    <a
      tabindex="0"
      class="flex flex-row justify-start items-center px-10 w-full h-16 cursor-pointer"
      v-bind="$attrs"
      @click="handleClick($attrs.to)"
    >
      <img v-if="icon" :src="iconUrl" class="w-8 h-8" />
      <div v-else class="w-8 h-8" />
      <span class="ml-4 text-lg text-neutrals-white">{{ heading }}</span>
    </a>
  </li>
</template>

<script>
import { navbarMutations } from '@/components/Navbar';

export default {
  name: 'NavbarLink',
  props: {
    heading: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      required: true,
    },
  },
  computed: {
    iconUrl() {
      return require(`@/assets/img/${this.icon}`);
    },
  },
  methods: {
    toggleNavbar() {
      navbarMutations.toggleNavbar();
    },
    handleClick(to) {
      this.$router.push(to);
      this.toggleNavbar();
    },
  },
};
</script>

<style scoped>
.bar:not(:focus-within):hover:before {
  position: absolute;
  content: '';
  display: block;
  background-color: theme('colors.primary.purple.hashita');
  height: theme('spacing.16');
  width: 4px;
}
</style>
