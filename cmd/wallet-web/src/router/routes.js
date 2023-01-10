/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import supportedLocales from '@/config/supportedLocales';
import store from '@/store';

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
            path: 'credentials/:id/:vaultName',
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
        beforeEnter: async () => {
          if (store.getters.getEnableDIDComm === undefined) await store.dispatch('initOpts');
          if (!store.getters.getEnableDIDComm) {
            store.dispatch('updateUserLoaded', true);
            return {
              name: 'NotFound',
              params: {
                title:
                  'The page you requested requires DIDComm, but your Wallet configuration does not support it at the moment.',
                message: 'Please, contact your administrator.',
              },
            };
          }
        },
        meta: {
          requiresAuth: true,
          signin: true,
          disableCHAPI: true,
          isNavbarHidden: true,
          requiresDIDComm: true,
        },
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
        redirect: 'NotFound',
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
        path: 'openid4vc',
        name: 'openid4vc',
        component: load('layouts/OpenID4VCLayout'),
        redirect: 'NotFound',
        meta: { requiresAuth: true, signin: true, disableCHAPI: true, isNavbarHidden: true },
        children: [
          {
            path: 'save',
            name: 'openid4vc-save',
            component: load('layouts/OpenID4VCSaveLayout'),
          },
          {
            path: 'share',
            name: 'openid4vc-share',
            component: load('layouts/OpenID4VCShareLayout'),
          },
        ],
      },
      {
        path: 'StoreInWallet',
        name: 'chapi-store',
        component: load('pages/StorePage'),
        beforeEnter: async () => {
          if (store.getters.getEnableCHAPI === undefined) await store.dispatch('initOpts');
          if (!store.getters.getEnableCHAPI) {
            store.dispatch('updateUserLoaded', true);
            return {
              name: 'NotFound',
              params: {
                title:
                  'The page you requested requires CHAPI, but your Wallet configuration does not support it at the moment.',
                message: 'Please, contact your administrator.',
              },
            };
          }
        },
        meta: { blockNoAuth: true, requiresCHAPI: true },
      },
      {
        path: 'GetFromWallet',
        name: 'chapi-get',
        component: load('layouts/GetLayout'),
        beforeEnter: async () => {
          if (store.getters.getEnableCHAPI === undefined) await store.dispatch('initOpts');
          if (!store.getters.getEnableCHAPI) {
            store.dispatch('updateUserLoaded', true);
            return {
              name: 'NotFound',
              params: {
                title:
                  'The page you requested requires CHAPI, but your Wallet configuration does not support it at the moment.',
                message: 'Please, contact your administrator.',
              },
            };
          }
        },
        meta: { blockNoAuth: true, isNavbarHidden: true, requiresCHAPI: true },
      },
      {
        path: 'worker',
        name: 'chapi-worker',
        component: load('pages/WorkerPage'),
        beforeEnter: async () => {
          if (store.getters.getEnableCHAPI === undefined) await store.dispatch('initOpts');
          if (!store.getters.getEnableCHAPI) {
            store.dispatch('updateUserLoaded', true);
            return {
              name: 'NotFound',
              params: {
                title:
                  'The page you requested requires CHAPI, but your Wallet configuration does not support it at the moment.',
                message: 'Please, contact your administrator.',
              },
            };
          }
        },
        meta: { requiresCHAPI: true },
      },
      {
        path: 'needauth',
        name: 'block-no-auth',
        component: load('pages/BlockNoAuthPage'),
      },
      {
        path: 'initiateProtocol',
        name: 'Initiate Protocol',
        component: load('pages/InitiateProtocolPage'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: load('pages/NotFoundPage'),
    props: true,
  },
];
