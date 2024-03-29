<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Dashboard Layout -->
    <div
      v-if="breakpoints.xs || breakpoints.sm"
      class="flex w-screen flex-col items-center justify-start bg-neutrals-softWhite"
    >
      <HeaderComponent />
      <router-view
        v-if="!isNavbarOpen"
        class="h-auto min-h-screen w-screen bg-neutrals-softWhite"
      />
    </div>

    <!-- Desktop Dashboard Layout -->
    <div
      v-else
      class="mx-auto flex max-w-7xl grow flex-row justify-center bg-neutrals-softWhite shadow-main-container"
    >
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
      <router-view id="dashboard-content" class="flex grow flex-col md:py-12 md:px-16" />
    </div>
  </div>
</template>

<script>
import HeaderComponent from '@/components/Header/HeaderComponent.vue';
import NavbarComponent from '@/components/Navbar/NavbarComponent.vue';
import NavbarLinkComponent from '@/components/Navbar/NavbarLinkComponent.vue';
import { navbarStore } from '@/components/Navbar';
import useBreakpoints from '@/plugins/breakpoints.js';
import { useI18n } from 'vue-i18n';

export default {
  name: 'DashboardLayout',
  components: {
    HeaderComponent,
    NavbarComponent,
    NavbarLinkComponent,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      breakpoints: useBreakpoints(),
    };
  },
  computed: {
    isNavbarOpen() {
      return navbarStore.isNavbarOpen;
    },
  },
};
</script>
