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
      <side-bar v-if="sidebarOpen">
        <!-- links -->
        <!-- <footer></footer> -->
      </side-bar>
    </div>

    <!-- Desktop Dashboard Layout -->
    <div class="hidden md:flex flex-row justify-start items-center px-20">
      <side-bar-v2>
        <sidebar-link-v2 :to="{ name: 'dashboardV2' }" :heading="i18n.vaults" icon="vaults.svg" />
        <sidebar-link-v2 :to="{ name: '' }" :heading="i18n.credentials" icon="credentials.svg" />
        <sidebar-link-v2 :to="{ name: '' }" :heading="i18n.account" icon="profile.svg" />
        <sidebar-link-v2
          :to="{ name: 'did-managementV2' }"
          :heading="i18n.settings"
          icon="settings.svg"
        />
      </side-bar-v2>
      <dashboard-content />
    </div>
  </div>
</template>

<script>
import DashboardContent from './Content.vue';
import { mapGetters } from 'vuex';
import SideBarV2 from '@/components/SidebarPluginV2/SideBarV2.vue';
import SidebarLinkV2 from '@/components/SidebarPluginV2/SidebarLinkV2.vue';

export default {
  components: {
    DashboardContent,
    SideBarV2,
    SidebarLinkV2,
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
