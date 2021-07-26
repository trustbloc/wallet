/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import DashboardLayout from '@/pages/layout/DashboardLayout.vue';
import Dashboard from '@/pages/Dashboard.vue';
import Login from '@/pages/Login.vue';
import Signin from '@/pages/Signin.vue';
import Signup from '@/pages/Signup.vue';
import Logout from '@/pages/Logout.vue';
import StoreInWallet from '@/pages/Store.vue';
import GetFromWallet from '@/pages/Get.vue';
import WalletWorker from '@/pages/Worker.vue';
import WebWallet from '@/pages/demos/WebWallet.vue';
import DIDManagement from '@/pages/DIDManagement.vue';
import Connections from '@/pages/demos/Connections.vue';
import Relationships from '@/pages/demos/Relationships.vue';
import IssueCredential from '@/pages/demos/IssueCredential.vue';
import PresentProof from '@/pages/demos/PresentProof.vue';
import NotFound from '@/pages/PageNotFound';
import BlockNoAuth from '@/pages/BlockNoAuth.vue';

const routes = [
  {
    path: `${__webpack_public_path__}`,
    component: DashboardLayout,
    // name: "dashboard",
    redirect: { name: 'dashboard' },
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: Dashboard,
        meta: { requiresAuth: true },
      },
      {
        path: 'logout',
        name: 'logout',
        component: Logout,
      },
      {
        path: 'WebWallet',
        name: 'web-wallet',
        component: WebWallet,
      },
      {
        path: 'DIDManagement',
        name: 'did-management',
        component: DIDManagement,
        meta: { requiresAuth: true },
      },
      {
        path: 'connections',
        name: 'connections',
        component: Connections,
        meta: { requiresAuth: true },
      },
      {
        path: 'relationships',
        name: 'relationships',
        component: Relationships,
        meta: { requiresAuth: true },
      },
      {
        path: 'issue-credential',
        name: 'issue-credential',
        component: IssueCredential,
        meta: { requiresAuth: true },
      },
      {
        path: 'present-proof',
        name: 'present-proof',
        component: PresentProof,
        meta: { requiresAuth: true },
      },
      {
        path: 'worker',
        name: 'chapi-worker',
        component: WalletWorker,
        meta: { showNav: false, hideFooter: true },
      },
      {
        path: 'StoreInWallet',
        name: 'chapi-store',
        component: StoreInWallet,
        meta: { blockNoAuth: true, showNav: false, hideFooter: true },
      },
      {
        path: 'GetFromWallet',
        name: 'chapi-get',
        component: GetFromWallet,
        meta: { blockNoAuth: true, showNav: false, hideFooter: true },
      },
    ],
  },
  // This will be deleted once Sign Up and Sign In pages are ready
  {
    path: '/login',
    name: 'login',
    component: Login,
  },
  {
    path: '/signin',
    name: 'signin',
    component: Signin,
  },
  // signUp page as per new designs
  {
    path: '/signup',
    name: 'signup',
    component: Signup,
  },
  {
    path: '*',
    component: NotFound,
  },
  {
    path: '/needauth',
    name: 'block-no-auth',
    component: BlockNoAuth,
  },
];

export default routes;
