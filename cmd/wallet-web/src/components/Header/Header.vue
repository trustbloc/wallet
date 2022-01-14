<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div :class="['w-full', isNavbarOpen ? 'h-screen' : 'h-auto']">
    <!-- Navbar Closed -->
    <div
      v-if="!isNavbarOpen"
      class="flex relative flex-col justify-center items-center p-6 w-full h-auto bg-gradient-dark"
    >
      <div class="flex absolute left-6 flex-row justify-start items-center">
        <slot name="leftButtonContainer" />
      </div>
      <slot v-if="hasCustomGradient" name="gradientContainer" />
      <div v-else class="absolute w-full bg-gradient-full opacity-40 oval oval-navbar-closed" />
      <div class="flex flex-row justify-center items-center">
        <button v-if="!isNavbarHidden" class="absolute left-6 z-10" @click="toggleNavbar">
          <img src="@/assets/img/menu-icon.svg" />
        </button>
        <Logo class="z-10 h-6 cursor-pointer" @click="$router.push({ name: 'vaults' })" />
      </div>
    </div>
    <!-- Navbar Open -->
    <div
      v-else
      class="flex relative flex-col justify-center items-center p-6 w-full h-full bg-gradient-dark"
    >
      <div class="absolute w-full opacity-40 bg-gradient-full oval oval-navbar-open" />
      <div class="flex flex-row justify-center items-center">
        <button class="absolute left-6 z-10" @click="toggleNavbar">
          <img src="@/assets/img/close.svg" />
        </button>
        <Logo
          class="z-10 h-6 cursor-pointer"
          @click="
            () => {
              $router.push({ name: 'vaults' });
              toggleNavbar();
            }
          "
        />
      </div>
      <navbar>
        <navbar-link
          id="navbar-link-vaults"
          :to="{ name: 'vaults' }"
          :heading="t('DashboardLayout.vaults')"
          icon="vaults.svg"
        />
        <navbar-link
          id="navbar-link-credentials"
          :to="{ name: 'credentials' }"
          :heading="t('DashboardLayout.credentials')"
          icon="credentials.svg"
        />
        <!-- TODO: uncomment once corresponding components are ready -->
        <!-- <navbar-link id="navbar-link-account" :to="{ name: 'account' }" :heading="t('DashboardLayout.account')" icon="profile.svg" /> -->
        <!-- TODO: link to actual settings once implemented -->
        <navbar-link
          id="navbar-link-did-management"
          :to="{ name: 'did-management' }"
          :heading="t('DashboardLayout.settings')"
          icon="settings.svg"
        />
      </navbar>
    </div>
  </div>
</template>

<script>
import Logo from '@/components/Logo/Logo.vue';
import Navbar from '@/components/Navbar/Navbar.vue';
import NavbarLink from '@/components/Navbar/NavbarLink.vue';
import { navbarStore, navbarMutations } from '@/components/Navbar';
import { useI18n } from 'vue-i18n';

export default {
  name: 'Header',
  components: {
    Logo,
    Navbar,
    NavbarLink,
  },
  props: {
    hasCustomGradient: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  computed: {
    isNavbarOpen() {
      return navbarStore.isNavbarOpen;
    },
    isNavbarHidden() {
      return this.$route.meta.isNavbarHidden;
    },
  },
  methods: {
    toggleNavbar() {
      navbarMutations.toggleNavbar();
    },
  },
};
</script>

<style scoped>
.oval {
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  filter: blur(50px);
}
.oval-navbar-closed {
  height: 5.875rem; /* 94px */
  top: -4rem;
}
.oval-navbar-open {
  height: 8.25rem; /* 132px */
  top: -5rem;
}
</style>
