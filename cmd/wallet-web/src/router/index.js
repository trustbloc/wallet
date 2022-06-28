/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createRouter, createWebHistory } from 'vue-router';
import store from '@/store';
import { getGnapKeyPair, gnapContinue, gnapRequestAccess } from '@/mixins';
import routes from './routes';
import { SHA3 } from 'sha3';
import axios from 'axios';

const router = createRouter({
  history: createWebHistory(__webpack_public_path__),
  routes,
});

router.beforeEach(async (to, from) => {
  store.dispatch('agent/flushStore');
  const locale = store.getters.getLocale;
  console.log('to.path', to.path);
  console.log('store.getters.getGnapReqAccessResp', store.getters.getGnapReqAccessResp);

  if (to.path === '/gnap/redirect' && store.getters.getGnapReqAccessResp) {
    console.log('handling gnap/redirect');
    // need to dispatch initOpts here to initialize hubAuthURL
    // initOpts should've already been called in requestAccess handler block
    // TODO need to persist hubAuthURL in local storage to avoid calling initOpts here
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
      const gnapKeyPair = await getGnapKeyPair();
      const signer = { SignatureVal: gnapKeyPair };
      console.log('called gnapContinue with', {
        signer,
        gnapAuthServerURL,
        interactRef,
        token: gnapResp.continue_access_token.value,
      });
      const gnapContinueResp = await gnapContinue(
        signer,
        gnapAuthServerURL,
        interactRef,
        gnapResp.continue_access_token.value
      );
      console.log('continue resp', gnapContinueResp);
      const accessToken = gnapContinueResp.data.access_token[0].value;
      const subjectId = gnapContinueResp.data.subject.sub_ids[0].id;
      console.log('new accessToken', accessToken);
      console.log('new subjectId', subjectId);
      store.dispatch('updateSessionToken', accessToken);
      store.dispatch('updateSubjectId', subjectId);
      store.dispatch('updateUser', subjectId);
      // need to call initOpts here to initialize edv and kms opts from bootstrap data
      // TODO refactor initOpts to avoid calling it twice (separate static and bootstrap agent opts initialization)
      console.log('calling initOpts with access token from gnap/continue');
      await store.dispatch('initOpts');
      try {
        await store.dispatch('agent/init', { accessToken, subjectId });
      } catch (e) {
        console.error('error initializing agent in gnap flow:', e);
      }
      // continue access token should only be needed to complete initial auth flow
      // thus, once successful call to gnapContinue is made, we can delete it
      // later on, if we need to authenticate the same user again we just call requestAccess
      // if needed, it will return us a new continue access token to complete user authentication
      store.dispatch('updateGnapReqAccessResp', null);
      console.log('authenticated user, initialized agent, redirecting to dashboard');
      return { name: 'DashboardLayout' };
    }
    console.error('error authenticating user: invalid hash received');
    return false;
  } else if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (store.getters.getCurrentUser) {
      console.log(`store.getters['agent/isInitialized']`, store.getters['agent/isInitialized']);
      if (!store.getters['agent/isInitialized']) {
        console.log('agent not initialized');
        // try loading gnap access token and subject id from store
        const accessToken = store.getters.getGnapSessionToken;
        const subjectId = store.getters.getGnapSubjectId;
        console.log('accessToken', accessToken);
        console.log('subjectId', subjectId);
        // initialize agent opts
        await store.dispatch('initOpts');
        try {
          console.log('initializing agent for existing user');
          await store.dispatch('agent/init', { accessToken, subjectId });
        } catch (e) {
          console.error('error initializing agent for existing user:', e);
        }
        console.log('initialized agent for user, redirecting to', to);
        return;
      }
      return;
    } else if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      if (!store.getters['agent/isInitialized']) {
        console.log('agent not initialized');
        // try loading gnap access token and subject id from store
        const accessToken = store.getters.getGnapSessionToken;
        const subjectId = store.getters.getGnapSubjectId;
        console.log('accessToken', accessToken);
        console.log('subjectId', subjectId);
        // initialize agent opts
        await store.dispatch('initOpts');
        try {
          console.log('initializing agent for existing user');
          await store.dispatch('agent/init', { accessToken, subjectId });
        } catch (e) {
          console.error('error initializing agent for existing user:', e);
        }
        console.log('initialized agent for user, redirecting to', to);
        return;
      }
      return;
    } else {
      // user in not authenticated, initiate auth flow
      console.log('initializing authentication flow');
      await store.dispatch('initOpts');
      const { signin, disableCHAPI } = to.meta;
      console.log('TODO: redirect to auth to authenticate the user');
      // const interactRedirectUrl = await initiateGnapAuth(store, router);

      const gnapAccessTokens = await store.getters['getGnapAccessTokenConfig'];
      const gnapAuthServerURL = store.getters.hubAuthURL;
      const walletWebUrl = store.getters.walletWebUrl;
      const gnapKeyPair = await getGnapKeyPair();
      console.log('gnapKeyPair', gnapKeyPair);
      const signer = { SignatureVal: gnapKeyPair };
      console.log('signer', signer.SignatureVal);
      const clientNonceVal = (Math.random() + 1).toString(36).substring(7);

      console.log(
        'requestAccess with:',
        JSON.stringify(
          {
            ...signer,
            gnapAccessTokens,
            gnapAuthServerURL,
            walletWebUrl,
            clientNonceVal,
          },
          null,
          2
        )
      );

      const resp = await gnapRequestAccess(
        signer,
        gnapAccessTokens,
        gnapAuthServerURL,
        walletWebUrl,
        clientNonceVal
      );
      console.log('resp', resp);
      // If user have already signed in then just redirect to requested page
      if (resp.data.access_token) {
        const accessToken = resp.data.access_token[0].value;
        const subjectId = resp.data.subject.sub_ids[0].id;
        console.log('111 new accessToken', accessToken);
        console.log('111 new subjectId', subjectId);
        store.dispatch('updateSessionToken', accessToken);
        store.dispatch('updateSubjectId', subjectId);
        store.dispatch('updateUser', subjectId);
        if (!store.getters['agent/isInitialized']) {
          console.log('agent not initialized');
          // initialize agent opts
          await store.dispatch('initOpts');
          try {
            console.log('initializing agent for existing user');
            await store.dispatch('agent/init', { accessToken, subjectId });
          } catch (e) {
            console.error('error initializing agent for new user:', e);
          }
          console.log('initialized agent for new user, redirecting to', to);
          return;
        }
        return;
      }
      const respMetaData = {
        uri: resp.data.continue.uri,
        continue_access_token: resp.data.continue.access_token,
        finish: resp.data.interact.finish,
        clientNonceVal: clientNonceVal,
      };
      store.dispatch('updateGnapReqAccessResp', respMetaData);
      console.log('redirecting to interact url', resp.data.interact.redirect);
      window.location.href = resp.data.interact.redirect;
      return false;
    }
  } else if (to.matched.some((record) => record.meta.blockNoAuth)) {
    if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
      if (!store.getters['agent/isInitialized']) {
        console.log('agent not initialized');
        // initialize agent opts
        await store.dispatch('initOpts');
        // try loading gnap access token and subject id from store
        const accessToken = store.getters.getGnapSessionToken;
        const subjectId = store.getters.getGnapSubjectId;
        console.log('accessToken', accessToken);
        console.log('subjectId', subjectId);
        try {
          console.log('initializing agent for existing user');
          await store.dispatch('agent/init', { accessToken, subjectId });
        } catch (e) {
          console.error('error initializing agent for existing user:', e);
        }
        console.log('initialized agent for user, redirecting to', to);
        return;
      }
      return;
    } else {
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
      return;
    } else {
      return;
    }
  }
});

export default router;
