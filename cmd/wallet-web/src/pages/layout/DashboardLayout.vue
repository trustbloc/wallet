<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Dashboard Layout -->
    <div class="flex md:hidden flex-col justify-start items-center w-screen bg-neutrals-softWhite">
      <Header />
      <keep-alive v-if="!isNavbarOpen">
        <dashboard-content class="h-screen bg-neutrals-softWhite" />
      </keep-alive>
    </div>

    <!-- Desktop Dashboard Layout -->
    <div
      class="
        hidden
        md:flex
        justify-center
        mx-auto
        max-w-7xl
        bg-neutrals-softWhite
        shadow-main-container
        flex-row flex-grow
      "
    >
      <navbar>
        <!-- TODO: bring link to vault on top once the component is implemented -->
        <navbar-link to="dashboard" :heading="i18n.credentials" icon="credentials.svg" />
        <!-- TODO: uncomment once corresponding components are ready -->
        <!-- <navbar-link to="vaults" :heading="i18n.vaults" icon="vaults.svg" /> -->
        <!-- <navbar-link to="account" :heading="i18n.account" icon="profile.svg" /> -->
        <!-- TODO: link to actual settings once implemented -->
        <navbar-link to="did-management" :heading="i18n.settings" icon="settings.svg" />
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
import { navbarStore, navbarMutations } from '@/components/Navbar';

export default {
  name: 'DashboardLayout',
  components: {
    DashboardContent,
    Header,
    Navbar,
    NavbarLink,
  },
  computed: {
    i18n() {
      return this.$t('DashboardLayout');
    },
    isNavbarOpen() {
      return navbarStore.isNavbarOpen;
    },
  },
  methods: {
    toggleNavbar() {
      navbarMutations.toggleNavbar();
    },
  },
};
</script>
