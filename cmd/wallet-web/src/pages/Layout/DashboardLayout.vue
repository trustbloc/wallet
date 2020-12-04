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

    <side-bar
        :sidebar-item-color="sidebarBackground"
        :sidebar-background-image="sidebarBackgroundImage">
      <mobile-menu slot="content"></mobile-menu>
      <sidebar-link to="/dashboard">
        <md-icon>dashboard</md-icon>
        <p>Dashboard</p>
      </sidebar-link>
      <sidebar-link to="/DIDManagement">
        <md-icon>flip_to_back</md-icon>
        <p>Digital Identity Management</p>
      </sidebar-link>
      <sidebar-link to="/relationships">
        <md-icon>compare_arrows</md-icon>
        <p>Relationships</p>
        <md-badge v-if="pendingConnectionsCount" class="md-primary md-square" style="margin: 5px"
                  :md-content="pendingConnectionsCount"/>
      </sidebar-link>
      <sidebar-link v-if="isDevMode" to="/MyVC">
        <md-icon>border_outer</md-icon>
        <p>Presentation</p>
      </sidebar-link>
      <sidebar-link v-if="isDevMode" to="/connections">
        <md-icon>compare_arrows</md-icon>
        <p>Connections</p>
        <md-badge v-if="pendingConnectionsCount" class="md-primary md-square" style="margin: 5px"
                  :md-content="pendingConnectionsCount"/>
      </sidebar-link>
      <sidebar-link v-if="isDevMode" to="/issue-credential">
        <md-icon>note</md-icon>
        <p>Issue Credential</p>
      </sidebar-link>
      <sidebar-link v-if="isDevMode" to="/present-proof">
        <md-icon>security</md-icon>
        <p>Present Proof</p>
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
      sidebarBackground: "green",
      sidebarBackgroundImage: require("@/assets/img/sidebar-2.jpg")
    };
  },
  computed: mapGetters(['pendingConnectionsCount', "isDevMode"]),
};
</script>
