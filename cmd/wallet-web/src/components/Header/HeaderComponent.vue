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
      class="relative flex h-auto w-full flex-col items-center justify-center bg-gradient-dark p-6"
    >
      <div class="absolute left-6 flex flex-row items-center justify-start">
        <slot name="leftButtonContainer" />
      </div>
      <slot v-if="hasCustomGradient" name="gradientContainer" />
      <div v-else class="oval oval-navbar-closed absolute w-full bg-gradient-full opacity-40" />
      <div class="flex flex-row items-center justify-center">
        <button
          v-if="!isNavbarHidden && showMenuDropdown"
          class="absolute left-6 z-10"
          @click="toggleNavbar"
        >
          <img src="@/assets/img/menu-icon.svg" />
        </button>
        <LogoComponent class="z-10 w-40 cursor-pointer" @click="$router.push({ name: 'vaults' })" />
      </div>
    </div>
    <!-- Navbar Open -->
    <div
      v-else
      class="relative flex h-full w-full flex-col items-center justify-center bg-gradient-dark p-6"
    >
      <div class="oval oval-navbar-open absolute w-full bg-gradient-full opacity-40" />
      <div class="flex flex-row items-center justify-center">
        <button class="absolute left-6 z-10" @click="toggleNavbar">
          <img src="@/assets/img/close.svg" />
        </button>
        <LogoComponent
          class="z-10 w-40 cursor-pointer"
          @click="
            () => {
              $router.push({ name: 'vaults' });
              toggleNavbar();
            }
          "
        />
      </div>
      <NavbarComponent>
        <NavbarLinkComponent
          id="navbar-link-vaults"
          :to="{ name: 'vaults' }"
          :heading="t('DashboardLayout.vaults')"
          icon="vaults.svg"
        />
        <NavbarLinkComponent
          id="navbar-link-credentials"
          :to="{ name: 'credentials' }"
          :heading="t('DashboardLayout.credentials')"
          icon="credentials.svg"
        />
        <!-- TODO: uncomment once corresponding components are ready -->
        <!-- <navbar-link id="navbar-link-account" :to="{ name: 'account' }" :heading="t('DashboardLayout.account')" icon="profile.svg" /> -->
        <!-- TODO: link to actual settings once implemented -->
        <NavbarLinkComponent
          id="navbar-link-did-management"
          :to="{ name: 'did-management' }"
          :heading="t('DashboardLayout.settings')"
          icon="settings.svg"
        />
      </NavbarComponent>
    </div>
  </div>
</template>

<script>
import LogoComponent from '@/components/Logo/LogoComponent.vue';
import NavbarComponent from '@/components/Navbar/NavbarComponent.vue';
import NavbarLinkComponent from '@/components/Navbar/NavbarLinkComponent.vue';
import { navbarStore, navbarMutations } from '@/components/Navbar';
import { useI18n } from 'vue-i18n';

export default {
  name: 'HeaderComponent',
  components: {
    LogoComponent,
    NavbarComponent,
    NavbarLinkComponent,
  },
  props: {
    hasCustomGradient: {
      type: Boolean,
      default: false,
    },
    showMenuDropdown: {
      type: Boolean,
      default: true,
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
