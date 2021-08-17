<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <!-- Mobile Dashboard Layout -->
    <div class="flex md:hidden flex-col justify-start items-center">
      <!-- <header></header> -->
      <sidebar v-if="sidebarOpen">
        <!-- links -->
      </sidebar>
      <dashboard-content />
    </div>

    <!-- Desktop Dashboard Layout -->
    <div class="hidden md:flex flex-row justify-start items-center px-20">
      <sidebar>
        <!-- TODO: bring link to vault on top once the component is implemented -->
        <sidebar-link :to="'dashboard'" :heading="i18n.credentials" icon="credentials.svg" />
        <!-- TODO: uncomment once corresponding components are ready -->
        <!-- <sidebar-link :to="'vaults'" :heading="i18n.vaults" icon="vaults.svg" /> -->
        <!-- <sidebar-link :to="'account'" :heading="i18n.account" icon="profile.svg" /> -->
        <!-- TODO: link to actual settings once implemented -->
        <sidebar-link :to="'did-management'" :heading="i18n.settings" icon="settings.svg" />
      </sidebar>
      <dashboard-content />
    </div>
  </div>
</template>

<script>
import DashboardContent from './Content.vue';
import { mapGetters } from 'vuex';
import Sidebar from '@/components/SidebarPlugin/Sidebar.vue';
import SidebarLink from '@/components/SidebarPlugin/SidebarLink.vue';

export default {
  components: {
    DashboardContent,
    Sidebar,
    SidebarLink,
  },
  data() {
    return {
      sidebarOpen: false,
    };
  },
  computed: {
    ...mapGetters(['pendingConnectionsCount', 'isDevMode']),
    getSidebarOpen() {
      return this.sidebarOpen;
    },
    i18n() {
      return this.$t('DashboardLayout');
    },
  },
};
</script>
