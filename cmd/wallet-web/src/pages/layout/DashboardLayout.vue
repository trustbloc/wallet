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
      class="flex flex-col justify-start items-center w-screen bg-neutrals-softWhite"
    >
      <Header />
      <dashboard-content v-if="!isNavbarOpen" class="w-screen h-screen bg-neutrals-softWhite" />
    </div>

    <!-- Desktop Dashboard Layout -->
    <div
      v-else
      class="
        flex
        justify-center
        mx-auto
        max-w-7xl
        bg-neutrals-softWhite
        shadow-main-container
        flex-row flex-grow
      "
    >
      <navbar>
        <navbar-link
          id="navbar-link-vaults"
          :to="{ name: 'vaults' }"
          :heading="t('DashboardLayout.vaults')"
          icon="vaults.svg"
        />
        <navbar-link
          id="navbar-link-dashboard"
          :to="{ name: 'dashboard' }"
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
      <dashboard-content class="flex flex-col flex-grow py-12 px-16" />
    </div>
  </div>
</template>

<script>
import DashboardContent from './DashboardContent.vue';
import Header from '@/components/Header/Header.vue';
import Navbar from '@/components/Navbar/Navbar.vue';
import NavbarLink from '@/components/Navbar/NavbarLink.vue';
import { navbarStore } from '@/components/Navbar';
import useBreakpoints from '@/plugins/breakpoints.js';
import { useI18n } from 'vue-i18n';

export default {
  name: 'DashboardLayout',
  components: {
    DashboardContent,
    Header,
    Navbar,
    NavbarLink,
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
