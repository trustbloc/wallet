/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createRouter, createWebHistory } from 'vue-router';
import store from '@/store';
import { getGnapKeyPair, gnapContinue, gnapRequestAccess } from '@/mixins';
import routes from './routes';
import { SHA3 } from 'sha3';
import { HTTPSigner } from '@trustbloc/wallet-sdk';
import { CHAPIHandler } from '@/mixins';
import useBreakpoints from '@/plugins/breakpoints.js';

const router = createRouter({
  history: createWebHistory(__webpack_public_path__),
  routes,
});

router.beforeEach(async (to, from) => {
  store.dispatch('agent/flushStore');
  const locale = store.getters.getLocale;

  if (to.path === '/gnap/redirect') {
    if (store.dispatch('loadUser') && store.getters.getGnapReqAccessResp) {
      await store.dispatch('initOpts');
      const gnapResp = store.getters.getGnapReqAccessResp;
      const params = new URL(document.location).searchParams;
      const hash = params.get('hash');
      const interactRef = params.get('interact_ref');
      const data = gnapResp.clientNonceVal + '\n' + gnapResp.finish + '\n' + interactRef + '\n';
      const shaHash = new SHA3(512);
      shaHash.update(data);
      let hashB64 = shaHash.digest({ format: 'base64' });
      hashB64 = hashB64.replace(/\+/g, '-').replace(/\//g, '_').replace(/\=+$/, '');
      if (hash === hashB64) {
        const gnapAuthServerURL = store.getters.hubAuthURL;
        const gnapKeyPair = await getGnapKeyPair('key1', 'ES256');
        const signer = new HTTPSigner({
          authorization: gnapResp.continue_access_token.value,
          signingKey: gnapKeyPair,
        });

        const gnapContinueResp = await gnapContinue(
          signer,
          gnapAuthServerURL,
          interactRef,
          gnapResp.continue_access_token.value
        );

        const accessToken = gnapContinueResp.data.access_token[0].value;
        const subjectId = gnapContinueResp.data.subject.sub_ids[0].id;
        store.dispatch('updateAccessToken', accessToken);
        store.dispatch('updateSubjectId', subjectId);
        store.dispatch('updateUser', subjectId);

        await store.dispatch('initOpts', { accessToken });

        const breakpoints = useBreakpoints();
        const enableCHAPI =
          store.getters['getEnableCHAPI'] &&
          !to.params?.disableCHAPI &&
          (to.name === 'chapi-worker' || (!breakpoints.xs && !breakpoints.sm));
        if (enableCHAPI) {
          const polyfill = await import('credential-handler-polyfill');
          const webCredentialHandler = await import('web-credential-handler');
          const chapi = new CHAPIHandler(
            polyfill,
            webCredentialHandler,
            store.getters['getCredentialMediatorURL']
          );
          chapi
            .install(subjectId)
            .then(() => store.dispatch('activateCHAPI'))
            .catch((e) => {
              console.error(e);
            });
        }
        try {
          await store.dispatch('agent/init', { accessToken, subjectId, newUser: true });
        } catch (e) {
          console.error('error initializing agent in gnap flow:', e);
        }
        // continue access token should only be needed to complete initial auth flow
        // thus, once successful call to gnapContinue is made, we can delete it
        // later on, if we need to authenticate the same user again we just call requestAccess
        // if needed, it will return us a new continue access token to complete user authentication
        store.dispatch('updateGnapReqAccessResp', null);
        store.dispatch('updateUserLoaded', true);

        return store.getters.getTargetPage || 'DashboardLayout';
      }
      console.error('error authenticating user: invalid hash received');
      return false;
    }
  } else if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (store.getters.getCurrentUser) {
      if (!store.getters['agent/isInitialized']) {
        const accessToken = store.getters.getGnapAccessToken;
        const subjectId = store.getters.getGnapSubjectId;

        await store.dispatch('initOpts', { accessToken });

        try {
          await store.dispatch('agent/init', { accessToken, subjectId });
        } catch (e) {
          console.error('error initializing agent for existing user:', e);
        }
      }
      store.dispatch('updateUserLoaded', true);
      return;
    } else if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      if (!store.getters['agent/isInitialized']) {
        const accessToken = store.getters.getGnapAccessToken;
        const subjectId = store.getters.getGnapSubjectId;

        await store.dispatch('initOpts', { accessToken });

        try {
          await store.dispatch('agent/init', { accessToken, subjectId });
        } catch (e) {
          console.error('error initializing agent for existing user:', e);
        }
      }
      store.dispatch('updateUserLoaded', true);
      return;
    } else {
      // User is not authenticated, request access from auth server
      await store.dispatch('initOpts');

      const gnapAccessTokens = await store.getters['getGnapAccessTokenConfig'];
      const gnapAuthServerURL = store.getters.hubAuthURL;
      const walletWebUrl = store.getters.walletWebUrl;
      const gnapKeyPair = await getGnapKeyPair('key1', 'ES256');

      const signer = new HTTPSigner({
        signingKey: gnapKeyPair,
      });

      const clientNonceVal = (Math.random() + 1).toString(36).substring(7);

      const resp = await gnapRequestAccess(
        signer,
        gnapAccessTokens,
        gnapAuthServerURL,
        walletWebUrl,
        clientNonceVal
      );

      // If user is not new, just make sure agent is initialized and proceed to the requested page
      if (resp.data.access_token) {
        const accessToken = resp.data.access_token[0].value;
        const subjectId = resp.data.subject.sub_ids[0].id;
        store.dispatch('updateAccessToken', accessToken);
        store.dispatch('updateSubjectId', subjectId);
        store.dispatch('updateUser', subjectId);
        if (!store.getters['agent/isInitialized']) {
          await store.dispatch('initOpts', { accessToken });

          const breakpoints = useBreakpoints();
          const enableCHAPI =
            store.getters['getEnableCHAPI'] &&
            !to.params?.disableCHAPI &&
            (to.name === 'chapi-worker' || (!breakpoints.xs && !breakpoints.sm));
          if (enableCHAPI) {
            const polyfill = await import('credential-handler-polyfill');
            const webCredentialHandler = await import('web-credential-handler');
            const chapi = new CHAPIHandler(
              polyfill,
              webCredentialHandler,
              store.getters['getCredentialMediatorURL']
            );
            chapi
              .install(subjectId)
              .then(() => store.dispatch('activateCHAPI'))
              .catch((e) => {
                console.error(e);
              });
          }

          try {
            await store.dispatch('agent/init', { accessToken, subjectId });
          } catch (e) {
            console.error('error initializing agent for new user:', e);
          }
        }
        store.dispatch('updateUserLoaded', true);
        return;
      }

      // If user is new, continue with the sign up flow and do not show the requested page at this step
      const respMetaData = {
        uri: resp.data.continue.uri,
        continue_access_token: resp.data.continue.access_token,
        finish: resp.data.interact.finish,
        clientNonceVal: clientNonceVal,
      };
      store.dispatch('updateGnapReqAccessResp', respMetaData);
      store.dispatch('updateTargetPage', to);
      window.location.href = resp.data.interact.redirect;
      return false;
    }
  } else if (to.matched.some((record) => record.meta.blockNoAuth)) {
    if (store.getters.getCurrentUser) {
      if (!store.getters['agent/isInitialized']) {
        const accessToken = store.getters.getGnapAccessToken;
        const subjectId = store.getters.getGnapSubjectId;

        await store.dispatch('initOpts', { accessToken });

        try {
          await store.dispatch('agent/init', { accessToken, subjectId });
        } catch (e) {
          console.error('error initializing agent for existing user:', e);
        }
      }
      store.dispatch('updateUserLoaded', true);
      return;
    } else if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      if (!store.getters['agent/isInitialized']) {
        const accessToken = store.getters.getGnapAccessToken;
        const subjectId = store.getters.getGnapSubjectId;

        await store.dispatch('initOpts', { accessToken });

        try {
          await store.dispatch('agent/init', { accessToken, subjectId });
        } catch (e) {
          console.error('error initializing agent for existing user:', e);
        }
      }
      store.dispatch('updateUserLoaded', true);
      return;
    } else {
      store.dispatch('updateUserLoaded', true);
      store.dispatch('updateTargetPage', to);
      return {
        name: 'block-no-auth',
        params: {
          ...router.currentRoute._value.params,
          locale: locale.base,
          redirect: '/',
        },
      };
    }
  } else {
    if (to.params.locale && to.params.locale !== locale.id) {
      router.replace({
        name: to.params.name,
        params: {
          ...router.currentRoute._value.params,
          ...to.params,
          locale: locale.base,
        },
        query: to.query,
      });
    }
    store.dispatch('updateUserLoaded', true);
    return;
  }
});

export default router;
