/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Sidebar from './Sidebar.vue';
import SidebarLink from './SidebarLink.vue';

const SidebarStore = {
  showSidebar: false,
  displaySidebar(value) {
    this.showSidebar = value;
  },
};

const SidebarPlugin = {
  install(Vue) {
    Vue.mixin({
      created: function () {
        return {
          SidebarStore,
        };
      },
    });
    Object.defineProperty(Vue.prototype, '$sidebar', {
      get() {
        return SidebarStore;
      },
    });
    Vue.component('Sidebar', Sidebar), Vue.component('SidebarLink', SidebarLink);
  },
};

export default SidebarPlugin;
