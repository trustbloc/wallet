/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Root from './TheRoot';
import supportedLocales from '@/config/supportedLocales';

// Lazy load the component
function load(path) {
  return () => import(`@/${path}.vue`);
}

// Creates regex (en|fr)
function getLocaleRegex() {
  let reg = '';
  supportedLocales.forEach((locale, index) => {
    reg = `${reg}${locale.id}${index !== supportedLocales.length - 1 ? '|' : ''}`;
  });
  return `(${reg})`;
}

export default [
  {
    path: `/:locale${getLocaleRegex()}?`,
    component: Root,
    children: [
      {
        path: '',
        name: 'DashboardLayout',
        component: load('layouts/DashboardLayout'),
        redirect: 'vaults',
        children: [
          {
            path: 'vaults',
            name: 'vaults',
            component: load('pages/VaultsPage'),
            meta: { requiresAuth: true },
          },
          {
            path: 'credentials',
            name: 'credentials',
            component: load('pages/CredentialsPage'),
            meta: { requiresAuth: true },
          },
          {
            path: 'credentials/:id',
            name: 'credential-details',
            component: load('pages/CredentialDetailsPage'),
            meta: { requiresAuth: true },
          },
          {
            path: 'did-management',
            name: 'did-management',
            component: load('pages/DIDManagementPage'),
            meta: { requiresAuth: true },
          },
        ],
      },
      {
        path: 'waci',
        name: 'waci',
        component: load('layouts/WACILayout'),
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
            component: load('pages/WACIIssuePage'),
          },
        ],
      },
      {
        path: 'oidc',
        name: 'oidc',
        component: load('layouts/OIDCLayout'),
        meta: { requiresAuth: true, signin: true, disableCHAPI: true, isNavbarHidden: true },
        children: [
          {
            path: 'share',
            name: 'share',
            component: load('layouts/OIDCShareLayout'),
          },
          {
            path: 'save',
            name: 'save',
            component: load('layouts/OIDCSaveLayout'),
          },
          {
            path: 'initiate',
            name: 'initiate',
            component: load('layouts/OIDCInitiateLayout'),
          },
        ],
      },
      {
        path: 'StoreInWallet',
        name: 'chapi-store',
        component: load('pages/StorePage'),
        meta: { blockNoAuth: true },
      },
      {
        path: 'GetFromWallet',
        name: 'chapi-get',
        component: load('layouts/GetLayout'),
        meta: { blockNoAuth: true, isNavbarHidden: true },
      },
      {
        path: 'worker',
        name: 'chapi-worker',
        component: load('pages/WorkerPage'),
      },
      {
        path: 'loginhandle',
        name: 'ProviderPopup',
        component: load('pages/ProviderPopupPage'),
        props: (route) => ({ providerID: route.query.providerID }),
      },
      {
        path: 'signin',
        name: 'signin',
        component: load('pages/SigninPage'),
      },
      {
        path: 'signup',
        name: 'signup',
        component: load('pages/SignupPage'),
      },
      {
        path: 'needauth',
        name: 'block-no-auth',
        component: load('pages/BlockNoAuthPage'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: load('pages/NotFoundPage'),
  },
];
