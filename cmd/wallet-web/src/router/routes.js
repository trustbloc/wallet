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
    redirect: 'vaults',
    children: [
      {
        path: 'vaults',
        name: 'vaults',
        component: load('Vaults'),
        meta: { requiresAuth: true },
      },
      {
        path: 'credentials',
        name: 'credentials',
        component: load('Credentials'),
        meta: { requiresAuth: true },
      },
      {
        path: 'credentials/:id',
        name: 'credential-details',
        component: load('CredentialDetails'),
        meta: { requiresAuth: true },
      },
      {
        path: 'waci',
        name: 'waci',
        component: load('WACI'),
        meta: { requiresAuth: true, signin: true, disableCHAPI: true },
      },
      {
        path: 'did-management',
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
    ],
  },
  {
    path: 'StoreInWallet',
    name: 'chapi-store',
    component: load('Store'),
    meta: { blockNoAuth: true },
  },
  {
    path: 'GetFromWallet',
    name: 'chapi-get',
    component: load('Get'),
    meta: { blockNoAuth: true, isNavbarHidden: true },
  },
  {
    path: 'worker',
    name: 'chapi-worker',
    component: load('Worker'),
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
    path: 'needauth',
    name: 'block-no-auth',
    component: load('BlockNoAuth'),
  },
  {
    path: 'pathMatch(.*)*',
    name: 'PageNotFound',
    component: load('PageNotFound'),
  },
];
