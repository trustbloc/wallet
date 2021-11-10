<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Dashboard Layout -->
    <keep-alive v-if="breakpoints.xs || breakpoints.sm">
      <div class="flex flex-col justify-start items-center w-screen bg-neutrals-softWhite">
        <Header />
        <keep-alive v-if="!isNavbarOpen">
          <dashboard-content class="w-screen h-screen bg-neutrals-softWhite" />
        </keep-alive>
      </div>
    </keep-alive>

    <!-- Desktop Dashboard Layout -->
    <keep-alive v-else>
      <div
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
            id="navbar-link-dashboard"
            :to="{ name: 'dashboard' }"
            :heading="i18n.credentials"
            icon="credentials.svg"
          />
          <!-- TODO: uncomment once vault component is ready and move this to the upper order -->
          <navbar-link
            id="navbar-link-vaults"
            :to="{ name: 'vaults' }"
            :heading="i18n.vaults"
            icon="vaults.svg"
          />
          <!-- TODO: uncomment once corresponding components are ready -->
          <!-- <navbar-link id="navbar-link-account" :to="{ name: 'account' }" :heading="i18n.account" icon="profile.svg" /> -->
          <!-- TODO: link to actual settings once implemented -->
          <navbar-link
            id="navbar-link-did-management"
            :to="{ name: 'did-management' }"
            :heading="i18n.settings"
            icon="settings.svg"
          />
        </navbar>
        <dashboard-content class="flex flex-col flex-grow py-12 px-16" />
      </div>
    </keep-alive>
  </div>
</template>

<script>
import DashboardContent from './DashboardContent.vue';
import Header from '@/components/Header/Header.vue';
import Navbar from '@/components/Navbar/Navbar.vue';
import NavbarLink from '@/components/Navbar/NavbarLink.vue';
import { navbarStore } from '@/components/Navbar';
import useBreakpoints from '@/plugins/breakpoints.js';

export default {
  name: 'DashboardLayout',
  components: {
    DashboardContent,
    Header,
    Navbar,
    NavbarLink,
  },
  data() {
    return {
      breakpoints: useBreakpoints(),
    };
  },
  computed: {
    i18n() {
      return this.$t('DashboardLayout');
    },
    isNavbarOpen() {
      return navbarStore.isNavbarOpen;
    },
  },
};
</script>
