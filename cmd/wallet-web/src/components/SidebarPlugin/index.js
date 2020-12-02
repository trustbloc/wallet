/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Sidebar from "./SideBar.vue";
import SidebarLink from "./SidebarLink.vue";

const SidebarStore = {
  showSidebar: false,
  displaySidebar(value) {
    this.showSidebar = value;
  }
};

const SidebarPlugin = {
  install(Vue) {
      Vue.mixin({
        created: function () {
          return{
            SidebarStore
          }
        }
      });
    Object.defineProperty(Vue.prototype, "$sidebar", {
      get() {
        return SidebarStore;
      }
    });
    Vue.component("side-bar", Sidebar),
    Vue.component("sidebar-link", SidebarLink)
  }
};

export default SidebarPlugin;
