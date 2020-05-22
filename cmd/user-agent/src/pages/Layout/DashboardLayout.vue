/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="wrapper" :class="{ 'nav-open': $sidebar.showSidebar }">

    <side-bar
      :sidebar-item-color="sidebarBackground"
      :sidebar-background-image="sidebarBackgroundImage">
      <mobile-menu slot="content"></mobile-menu>
      <sidebar-link to="/dashboard">
        <md-icon>dashboard</md-icon>
        <p>Dashboard</p>
      </sidebar-link>
      <sidebar-link to="/MyVC">
        <md-icon>border_outer</md-icon>
        <p>Presentation</p>
      </sidebar-link>
      <sidebar-link to="/DIDManagement">
        <md-icon>flip_to_back</md-icon>
        <p>DID Management</p>
      </sidebar-link>
      <sidebar-link to="/connections">
        <md-icon>compare_arrows</md-icon>
        <p>Connections</p>
      </sidebar-link>
      <sidebar-link to="/issue-credential">
        <md-icon>note</md-icon>
        <p>Issue Credential</p>
      </sidebar-link>
      <sidebar-link to="/present-proof">
        <md-icon>security</md-icon>
        <p>Present Proof</p>
      </sidebar-link>
    </side-bar>

    <div class="main-panel">
      <top-navbar></top-navbar>

      <dashboard-content> </dashboard-content>

      <content-footer v-if="!$route.meta.hideFooter"></content-footer>
    </div>
  </div>
</template>

<script>
import TopNavbar from "./TopNavbar.vue";
import ContentFooter from "./ContentFooter.vue";
import DashboardContent from "./Content.vue";
import MobileMenu from "@/pages/Layout/MobileMenu.vue";

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
  }
};
</script>
