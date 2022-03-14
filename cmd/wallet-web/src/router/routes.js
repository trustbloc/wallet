/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Lazy load the component
function load(path) {
  return () => import(`@/${path}.vue`);
}

export default [
  {
    path: '',
    name: 'DashboardLayout',
    component: load('layouts/DashboardLayout'),
    redirect: 'vaults',
    children: [
      {
        path: 'vaults',
        name: 'vaults',
        component: load('pages/Vaults'),
        meta: { requiresAuth: true },
      },
      {
        path: 'credentials',
        name: 'credentials',
        component: load('pages/Credentials'),
        meta: { requiresAuth: true },
      },
      {
        path: 'credentials/:id',
        name: 'credential-details',
        component: load('pages/CredentialDetails'),
        meta: { requiresAuth: true },
      },
      {
        path: 'did-management',
        name: 'did-management',
        component: load('pages/DIDManagement'),
        meta: { requiresAuth: true },
      },
    ],
  },
  {
    path: 'waci',
    name: 'waci',
    component: load('layouts/WACI'),
    meta: { requiresAuth: true, signin: true, disableCHAPI: true, isNavbarHidden: true },
    children: [
      {
        path: 'share',
        name: 'share',
        component: load('layouts/WACIShareLayout'),
      },
      {
        path: 'issue',
        name: 'issue',
        component: load('pages/WACIIssue'),
      },
    ],
  },
  {
    path: 'oidc-share',
    name: 'oidc-share',
    component: load('layouts/OIDC'),
    meta: { requiresAuth: true, signin: true, disableCHAPI: true, isNavbarHidden: true },
    children: [
      {
        path: 'share',
        name: 'share',
        component: load('layouts/OIDCShareLayout'),
      },
    ],
  },
  {
    path: 'StoreInWallet',
    name: 'chapi-store',
    component: load('pages/Store'),
    meta: { blockNoAuth: true },
  },
  {
    path: 'GetFromWallet',
    name: 'chapi-get',
    component: load('layouts/Get'),
    meta: { blockNoAuth: true, isNavbarHidden: true },
  },
  {
    path: 'worker',
    name: 'chapi-worker',
    component: load('pages/Worker'),
  },
  {
    path: 'loginhandle',
    name: 'loginhandle',
    component: load('pages/LoginHandle'),
    props: (route) => ({ provider: route.query.provider }),
  },
  {
    path: 'signin',
    name: 'signin',
    component: load('pages/Signin'),
  },
  {
    path: 'signup',
    name: 'signup',
    component: load('pages/Signup'),
  },
  {
    path: 'needauth',
    name: 'block-no-auth',
    component: load('pages/BlockNoAuth'),
  },
  {
    path: 'pathMatch(.*)*',
    name: 'NotFound',
    component: load('pages/NotFound'),
  },
];
