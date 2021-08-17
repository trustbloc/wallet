/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import SidebarV2 from './SideBarV2.vue';
import SidebarLinkV2 from './SidebarLinkV2.vue';

const SidebarStoreV2 = {
  showSidebar: false,
  displaySidebar(value) {
    this.showSidebar = value;
  },
};

const SidebarPluginV2 = {
  install(Vue) {
    Vue.mixin({
      created: function () {
        return {
          SidebarStoreV2,
        };
      },
    });
    Object.defineProperty(Vue.prototype, '$sidebarV2', {
      get() {
        return SidebarStoreV2;
      },
    });
    Vue.component('SideBarV2', SidebarV2), Vue.component('SidebarLinkV2', SidebarLinkV2);
  },
};

export default SidebarPluginV2;
