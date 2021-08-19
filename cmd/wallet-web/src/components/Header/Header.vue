<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="flex relative flex-row justify-center items-center p-6 w-screen h-full bg-gradient-dark"
  >
    <div
      class="absolute bottom-0 left-0 w-full h-full bg-no-repeat bg-cover bg-header-mobile-flare"
    />
    <button class="absolute left-6 z-10" @click="toggleNavbar">
      <img src="@/assets/img/menu-icon.svg" />
    </button>
    <Logo class="h-6 z-10" />
  </div>
</template>

<script>
import Logo from '@/components/Logo/Logo.vue';
export default {
  name: 'Header',
  components: {
    Logo,
  },
  methods: {
    toggleNavbar() {
      this.$navbar.setNavbarOpen(!this.$navbar.navbarOpen);
    },
  },
};
</script>
