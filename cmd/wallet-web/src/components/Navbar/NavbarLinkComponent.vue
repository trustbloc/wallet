<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <li
    :class="[
      currentPage === id && `bg-gradient-to-r opacity-100 from-neutrals-black current-bar`,
      currentPage !== id && `focus-within:shadow-inner-outline-blue`,
      `
      justify-start
      items-center
      font-bold
      hover:bg-gradient-to-r hover:opacity-100 hover:from-neutrals-black
      focus-within:bg-gradient-to-r
      focus-within:opacity-100
      focus-within:from-neutrals-black
      opacity-60
      bar
      flex flex-row
    `,
    ]"
  >
    <a
      :id="id"
      tabindex="0"
      class="
        justify-start
        items-center
        px-10
        w-full
        h-16
        cursor-pointer
        focus:outline-none
        flex flex-row
      "
      v-bind="$attrs"
      @click="handleClick($attrs)"
      @keyup.enter="handleClick($attrs)"
    >
      <img v-if="icon" :src="iconUrl" class="w-8 h-8" />
      <div v-else class="w-8 h-8" />
      <span class="ml-4 text-lg text-neutrals-white">{{ heading }}</span>
    </a>
  </li>
</template>

<script>
import { navbarStore, navbarMutations } from '@/components/Navbar';

export default {
  name: 'NavbarLinkComponent',
  props: {
    id: {
      type: String,
      required: true,
    },
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
    currentPage() {
      return navbarStore.currentPage;
    },
  },
  methods: {
    toggleNavbar() {
      navbarMutations.toggleNavbar();
    },
    handleClick(attrs) {
      this.$router.push({ name: attrs.to.name });
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
.current-bar:before {
  position: absolute;
  content: '';
  display: block;
  background: theme('gradients.full');
  height: theme('spacing.16');
  width: 4px;
}
</style>