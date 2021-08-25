/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Navbar from './Navbar.vue';
import NavbarLink from './NavbarLink.vue';

const NavbarStore = {
  navbarOpen: false,
  setNavbarOpen(value) {
    this.navbarOpen = value;
  },
};

const NavbarPlugin = {
  install(Vue) {
    Vue.mixin({
      created: function () {
        return {
          NavbarStore,
        };
      },
    });
    Object.defineProperty(Vue.prototype, '$navbar', {
      get() {
        return NavbarStore;
      },
    });
    Vue.component('Navbar', Navbar), Vue.component('NavbarLink', NavbarLink);
  },
};

export default NavbarPlugin;
