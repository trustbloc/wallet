/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Sidebar on the right. Used as a local plugin in DashboardLayout.vue
import Sidebar from './components/SidebarPlugin';

// asset imports
import VueMaterial from 'vue-material';
import 'vue-material/dist/vue-material.css';
import './assets/scss/material-dashboard.scss';

// library auto imports
import 'es6-promise';
/**
 * This is the main Light Bootstrap Dashboard Vue plugin where dashboard related plugins are registerd.
 */
export default {
  install(Vue) {
    Vue.use(Sidebar);
    Vue.use(VueMaterial);
  },
};
