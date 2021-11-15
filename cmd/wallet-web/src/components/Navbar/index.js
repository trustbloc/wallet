/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { reactive } from 'vue';

export const navbarStore = reactive({
  isNavbarOpen: false,
  currentPage: '',
});

export const navbarMutations = {
  toggleNavbar() {
    navbarStore.isNavbarOpen = !navbarStore.isNavbarOpen;
  },
  setCurrentPage(newPage) {
    navbarStore.currentPage = newPage;
  },
};
