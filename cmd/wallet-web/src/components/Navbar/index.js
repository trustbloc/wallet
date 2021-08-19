/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';

export const navbarStore = Vue.observable({
  isNavbarOpen: false,
});

export const navbarMutations = {
  toggleNavbar() {
    navbarStore.isNavbarOpen = !navbarStore.isNavbarOpen;
  },
};
