/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Lazy load the component
function load(component) {
  return () => import(/* webpackChunkName: "[request]" */ `@/pages/${component}.vue`);
}

export default [
  {
    path: '',
    name: 'DashboardLayout',
    component: load('layout/DashboardLayout'),
    redirect: 'dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: load('Dashboard'),
        meta: { requiresAuth: true },
      },
      {
        path: 'logout',
        name: 'logout',
        component: load('Logout'),
      },
      {
        path: 'WebWallet',
        name: 'web-wallet',
        component: load('demos/WebWallet'),
      },
      {
        path: 'DIDManagement',
        name: 'did-management',
        component: load('DIDManagement'),
        meta: { requiresAuth: true },
      },
      {
        path: 'connections',
        name: 'connections',
        component: load('demos/Connections'),
        meta: { requiresAuth: true },
      },
      {
        path: 'relationships',
        name: 'relationships',
        component: load('demos/Relationships'),
        meta: { requiresAuth: true },
      },
      {
        path: 'issue-credential',
        name: 'issue-credential',
        component: load('demos/IssueCredential'),
        meta: { requiresAuth: true },
      },
      {
        path: 'present-proof',
        name: 'present-proof',
        component: load('demos/PresentProof'),
        meta: { requiresAuth: true },
      },
      {
        path: 'worker',
        name: 'chapi-worker',
        component: load('Worker'),
        meta: { showNav: false, hideFooter: true },
      },
      {
        path: 'StoreInWallet',
        name: 'chapi-store',
        component: load('Store'),
        meta: { blockNoAuth: true, showNav: false, hideFooter: true },
      },
      {
        path: 'GetFromWallet',
        name: 'chapi-get',
        component: load('Get'),
        meta: { blockNoAuth: true, showNav: false, hideFooter: true },
      },
    ],
  },
  // Temporary path for the new dashboard
  {
    path: 'v2',
    name: 'DashboardLayoutV2',
    component: load('layout/DashboardLayoutV2'),
    redirect: 'dashboardV2',
    children: [
      {
        path: 'dashboard',
        name: 'dashboardV2',
        component: load('DashboardV2'),
        meta: { requiresAuth: true },
      },
      {
        path: 'logout',
        name: 'logoutV2',
        component: load('Logout'),
      },
      {
        path: 'WebWallet',
        name: 'web-walletV2',
        component: load('demos/WebWallet'),
      },
      {
        path: 'DIDManagement',
        name: 'did-managementV2',
        component: load('DIDManagement'),
        meta: { requiresAuth: true },
      },
      {
        path: 'connections',
        name: 'connectionsV2',
        component: load('demos/Connections'),
        meta: { requiresAuth: true },
      },
      {
        path: 'relationships',
        name: 'relationshipsV2',
        component: load('demos/Relationships'),
        meta: { requiresAuth: true },
      },
      {
        path: 'issue-credential',
        name: 'issue-credentialV2',
        component: load('demos/IssueCredential'),
        meta: { requiresAuth: true },
      },
      {
        path: 'present-proof',
        name: 'present-proofV2',
        component: load('demos/PresentProof'),
        meta: { requiresAuth: true },
      },
      {
        path: 'worker',
        name: 'chapi-workerV2',
        component: load('Worker'),
        meta: { showNav: false, hideFooter: true },
      },
      {
        path: 'StoreInWallet',
        name: 'chapi-storeV2',
        component: load('Store'),
        meta: { blockNoAuth: true, showNav: false, hideFooter: true },
      },
      {
        path: 'GetFromWallet',
        name: 'chapi-getV2',
        component: load('Get'),
        meta: { blockNoAuth: true, showNav: false, hideFooter: true },
      },
    ],
  },
  {
    path: 'loginhandle',
    name: 'loginhandle',
    component: load('LoginHandle'),
    props: (route) => ({ provider: route.query.provider }),
  },
  {
    path: 'signin',
    name: 'signin',
    component: load('Signin'),
  },
  {
    path: 'signup',
    name: 'signup',
    component: load('Signup'),
  },
  {
    path: '*',
    name: 'PageNotFound',
    component: load('PageNotFound'),
  },
  {
    path: 'needauth',
    name: 'block-no-auth',
    component: load('BlockNoAuth'),
  },
];
