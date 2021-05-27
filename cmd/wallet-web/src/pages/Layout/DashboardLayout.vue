/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div v-if="$route.meta.showNav === false" class="main-panel">
    <dashboard-content></dashboard-content>
    <content-footer v-if="!$route.meta.hideFooter"></content-footer>
  </div>

  <div v-else class="wrapper" :class="{ 'nav-open': $sidebar.showSidebar }">
    <!--
     TODO - base path for sidebar links should be configurable: https://github.com/trustbloc/edge-agent/issues/374
    -->
    <side-bar>
      <mobile-menu slot="content"></mobile-menu>
      <sidebar-link to="/dashboard">
        <div class="py-8 text-lg rounded text-white">
        <div class="flex items-center">
          <img class="h-8 w-8 mr-2" src="@/assets/img/wallet.png" alt="">
          <span class="normal-case px-4 mt-2 text-2xl font-sans">Wallet</span>
        </div>
        </div>
      </sidebar-link>
      <sidebar-link to="/DIDManagement">
        <div class="py-2 text-lg rounded text-white">
          <div class="flex items-center">
        <md-icon class="h-8 w-8 mr-2" style="color:white;">flip_to_back</md-icon>
        <span class="normal-case inline-block px-4 mt-2 text-xl font-sans">Digital Identity Management</span>
          </div>
        </div>
      </sidebar-link>
    </side-bar>

    <div class="main-panel">
      <top-navbar></top-navbar>

      <dashboard-content></dashboard-content>

      <content-footer v-if="!$route.meta.hideFooter"></content-footer>
    </div>
  </div>
</template>

<script>
import TopNavbar from "./TopNavbar.vue";
import ContentFooter from "./ContentFooter.vue";
import DashboardContent from "./Content.vue";
import MobileMenu from "@/pages/Layout/MobileMenu.vue";
import {mapGetters} from "vuex";

export default {
  components: {
    TopNavbar,
    DashboardContent,
    ContentFooter,
    MobileMenu
  },
  data() {
    return {
    };
  },
  computed: mapGetters(['pendingConnectionsCount', "isDevMode"]),
};
</script>
