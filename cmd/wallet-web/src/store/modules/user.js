/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import * as Agent from '@trustbloc/agent-sdk-web';
import { WalletUser } from '@trustbloc/wallet-sdk';
import { toRaw } from 'vue';
import { getBootstrapData } from '@/mixins/gnap/gnap';
import { clearGnapStoreData, getGnapKeyPair } from '@/mixins/gnap/store';
import { RegisterWallet } from '@/mixins';

export const parseTIme = (ns) => parseInt(ns) * 60 * 10 ** 9;

export default {
  state: {
    username: null,
    setupStatus: null,
    profile: null,
    loggedIn: false,
    logInSuspended: false,
    loaded: false,
    chapi: false,
    selectedVaultId: null,
    selectedCredentialId: null,
    gnapAccessToken: null,
    gnapRequestAccessResp: null,
  },
  mutations: {
    setUser(state, val) {
      state.username = val;
      if (val !== null) {
        localStorage.setItem('user', val);
      }
    },
    setProfile(state, val) {
      state.profile = val;
      // only save profile.user to local storage, since profile.token is unique for each user session
      if (val?.user !== null) {
        localStorage.setItem('profile', JSON.stringify({ user: val.user }));
      }
    },
    setUserPreference(state, val) {
      state.preference = val;
      if (val !== null) {
        localStorage.setItem('preference', JSON.stringify(val));
      }
    },
    setUserSetupStatus(state, val) {
      state.setupStatus = val;
      if (val !== null) {
        localStorage.setItem('setupStatus', val);
      }
    },
    setCHAPI(state, val) {
      state.chapi = val;
      if (val !== null) {
        localStorage.setItem('chapi', JSON.stringify(val));
      }
    },
    setUserLoggedIn(state) {
      state.loggedIn = true;
    },
    setLogInSuspended(state) {
      state.logInSuspended = true;
    },
    setLoaded(state, val) {
      state.loaded = val;
    },
    setSelectedVaultId(state, val) {
      state.selectedVaultId = val;
      if (val !== null) {
        localStorage.setItem('selectedVaultId', val);
      }
    },
    setSelectedCredentialId(state, val) {
      state.selectedCredentialId = val;
      if (val !== null) {
        localStorage.setItem('selectedCredentialId', val);
      }
    },
    setAccessToken(state, val) {
      state.gnapAccessToken = val;
      if (val !== null) {
        localStorage.setItem('gnapAccessToken', val);
      } else {
        localStorage.removeItem('gnapAccessToken');
      }
    },
    setSubjectId(state, val) {
      state.gnapSubjectId = val;
      if (val !== null) {
        localStorage.setItem('gnapSubjectId', val);
      } else {
        localStorage.removeItem('gnapSubjectId');
      }
    },
    setGnapAccessReqResp(state, val) {
      state.gnapRequestAccessResp = val;
      if (val !== null) {
        localStorage.setItem('gnapRequestAccessResp', JSON.stringify(val));
      } else {
        localStorage.removeItem('gnapRequestAccessResp');
      }
    },
    clearUser(state) {
      state.username = null;
      state.setupStatus = null;
      state.profile = null;
      state.loggedIn = false;
      state.loaded = false;
      state.chapi = false;
      state.selectedVaultId = null;
      state.selectedCredentialId = null;
      state.gnapAccessToken = null;
      state.gnapSubjectId = null;
      state.gnapRequestAccessResp = null;

      localStorage.removeItem('user');
      localStorage.removeItem('setupStatus');
      localStorage.removeItem('profile');
      localStorage.removeItem('preference');
      localStorage.removeItem('chapi');
      localStorage.removeItem('selectedVaultId');
      localStorage.removeItem('selectedCredentialId');
      localStorage.removeItem('gnapAccessToken');
      localStorage.removeItem('gnapSubjectId');
      localStorage.removeItem('gnapRequestAccessResp');

      clearGnapStoreData();
    },
    loadUser(state) {
      state.username = localStorage.getItem('user');
      state.setupStatus = localStorage.getItem('setupStatus');
      state.profile = JSON.parse(localStorage.getItem('profile'));
      state.preference = JSON.parse(localStorage.getItem('preference'));
      state.chapi = JSON.parse(localStorage.getItem('chapi'));
      state.selectedVaultId = localStorage.getItem('selectedVaultId');
      state.selectedCredentialId = localStorage.getItem('selectedCredentialId');
      state.gnapAccessToken = localStorage.getItem('gnapAccessToken');
      state.gnapSubjectId = localStorage.getItem('gnapSubjectId');
      state.gnapRequestAccessResp = JSON.parse(localStorage.getItem('gnapRequestAccessResp'));
    },
  },
  actions: {
    async refreshUserPreference({ commit, state, rootGetters }, profile = state.profile) {
      if (!profile) {
        console.error('failed to refresh user preference, profile not found.');
        throw 'invalid operation, user profile not set';
      }

      let { user, token } = profile;

      let walletUser = new WalletUser({ agent: rootGetters['agent/getInstance'], user });
      try {
        let { content } = await walletUser.getPreferences(token);
        commit('setUserPreference', content);
      } catch (e) {
        console.error('user preference not found, may be user yet to get registered', e);
      }
    },
    async logout({ commit, dispatch, getters }) {
      await dispatch('agent/destroy');
      commit('clearUser');
    },
    loadUser({ commit }) {
      commit('loadUser');
    },
    startUserSetup({ commit }) {
      commit('setUserSetupStatus', 'inprogress');
    },
    completeUserSetup({ commit }, failure) {
      commit('setUserSetupStatus', failure ? 'failed' : 'success');
    },
    updateUser({ commit }, username) {
      commit('setUser', username);
    },
    updateUserProfile({ commit }, profile) {
      commit('setProfile', profile);
    },
    updateUserOnboard({ commit }) {
      commit('setUserLoggedIn');
    },
    updateLoginSuspended({ commit }) {
      commit('setLogInSuspended');
    },
    updateUserLoaded({ commit }, loaded) {
      commit('setLoaded', loaded);
    },
    activateCHAPI({ commit }) {
      commit('setCHAPI', true);
    },
    updateSelectedVaultId({ commit }, selectedVaultId) {
      commit('setSelectedVaultId', selectedVaultId);
    },
    updateSelectedCredentialId({ commit }, selectedCredentialId) {
      commit('setSelectedCredentialId', selectedCredentialId);
    },
    updateAccessToken({ commit }, gnapAccessToken) {
      commit('setAccessToken', gnapAccessToken);
    },
    updateSubjectId({ commit }, gnapSubjectId) {
      commit('setSubjectId', gnapSubjectId);
    },
    updateGnapReqAccessResp({ commit }, gnapRequestAccessResp) {
      commit('setGnapAccessReqResp', gnapRequestAccessResp);
    },
  },
  getters: {
    getCurrentUser(state) {
      return state.username
        ? {
            username: state.username,
            setupStatus: state.setupStatus,
            profile: toRaw(state.profile),
            preference: toRaw(state.preference),
          }
        : undefined;
    },
    isUserLoggedIn(state) {
      return state.loggedIn;
    },
    isLoginSuspended(state) {
      return state.logInSuspended;
    },
    isUserLoaded(state) {
      return state.loaded;
    },
    isCHAPI(state) {
      return state.chapi;
    },
    getSelectedVaultId(state) {
      return state.selectedVaultId;
    },
    getSelectedCredentialId(state) {
      return state.selectedCredentialId;
    },
    getGnapAccessToken(state) {
      return state.gnapAccessToken;
    },
    getGnapSubjectId(state) {
      return state.gnapSubjectId;
    },
    getGnapReqAccessResp(state) {
      return state.gnapRequestAccessResp;
    },
  },
  modules: {
    agent: {
      namespaced: true,
      state: {
        instance: null,
        notifiers: null,
        agentName: null,
      },
      mutations: {
        setInstance(state, { instance, user }) {
          state.instance = instance;
          state.agentName = user;
        },
        addNotifier(state, notifier) {
          if (state.notifiers) {
            state.notifiers.push(notifier);
          } else {
            state.notifiers = [notifier];
          }
        },
        startNotifier(state, notifier) {
          state.instance.startNotifier(notifier.callback, notifier.topics);
        },
        startAllNotifiers(state) {
          if (!state.notifiers) {
            return;
          }
          state.notifiers.forEach(function (notifier) {
            state.instance.startNotifier(notifier.callback, notifier.topics);
          });
        },
      },
      actions: {
        async init(
          { commit, rootState, state, rootGetters, dispatch },
          { accessToken, subjectId, newUser = false }
        ) {
          const username = rootState.user.username;

          // If agent is already initialized for current user, then return
          if (state.instance && state.agentName === username) {
            return;
          }

          // If user is not authenticated, then agent cannot be initialized
          if (!username && !subjectId && !accessToken) {
            throw new Error('Error initializing agent: user is not authenticated');
          }

          // If hubAuthURL is missing, then agent cannot be initialized
          const hubAuthURL = rootGetters.hubAuthURL;
          if (!hubAuthURL) throw new Error('Error initializing agent: hubAuthURL is missing');

          const { privateKey, kid, alg } = await getGnapKeyPair();
          const signingKey = await window.crypto.subtle.exportKey('jwk', privateKey);
          signingKey.kid = kid;
          signingKey.alg = alg;

          // Updating agentOpts with new user data
          const agentOpts = rootGetters.getAgentOpts;
          const profileOpts = rootGetters.getProfileOpts;

          Object.assign(agentOpts, {
            'agent-default-label': username,
            'db-namespace': username,
            'gnap-signing-jwk': JSON.stringify(signingKey),
            'gnap-access-token': accessToken,
            'gnap-user-subject': subjectId,
          });

          const newOpts = await getBootstrapData(agentOpts, hubAuthURL, dispatch, accessToken);

          if (newOpts?.newAgentOpts) {
            Object.assign(agentOpts, newOpts?.newAgentOpts);
            commit('updateAgentOpts', agentOpts, { root: true });
          }
          if (newOpts?.newProfileOpts) {
            Object.assign(profileOpts, newOpts?.newProfileOpts);
            commit('updateProfileOpts', profileOpts, { root: true });
          }
          // Initialize agent and update state accordingly
          const agent = await new Agent.Framework(agentOpts);
          commit('setInstance', { instance: agent, user: username });

          // For new users fetch newly created bootstrap data
          if (newUser) {
            const newOpts = await getBootstrapData(agentOpts, hubAuthURL, dispatch, accessToken);
            if (newOpts?.newAgentOpts) {
              Object.assign(agentOpts, newOpts?.newAgentOpts);
              commit('updateAgentOpts', agentOpts, { root: true });
            }
            if (newOpts?.newProfileOpts) {
              Object.assign(profileOpts, newOpts?.newProfileOpts);
              commit('updateProfileOpts', profileOpts, { root: true });
            }
          }

          await dispatch('unlockWallet');
          commit('startAllNotifiers');
        },
        async unlockWallet({ state, rootGetters, dispatch }) {
          // create wallet profile if it doesn't exist

          const profileOpts = rootGetters.getProfileOpts;

          const { user } = profileOpts.bootstrap.data;

          const walletUser = new WalletUser({ agent: state.instance, user });
          if (!(await walletUser.profileExists())) {
            const createOpts = profileCreationOpts(profileOpts);
            await walletUser.createWalletProfile(createOpts);
          }

          const { token } = await walletUser.unlock(profileUnlockOpts(profileOpts));

          await dispatch('refreshUserPreference', { user, token }, { root: true });

          const currentUser = rootGetters.getCurrentUser;

          if (!currentUser.preference) {
            // todo: should this be rootGetters.getCurrentUser.value.preference?
            dispatch('startUserSetup', '', { root: true });

            const agent = rootGetters['agent/getInstance'];

            await new RegisterWallet(agent, rootGetters.getAgentOpts).register(
              {
                name: currentUser.username,
                user: user,
                token: token,
              },
              () => dispatch('completeUserSetup', '', { root: true })
            );
          }

          await dispatch('updateUserProfile', { user, token }, { root: true });
        },
        async flushStore({ state }) {
          if (state.instance) {
            try {
              await state.instance.store.flush();
            } catch (e) {
              console.error('error flushing store', e);
            }
          }
        },
        async destroy({ commit, state, rootGetters, dispatch }) {
          const { user } = rootGetters.getCurrentUser.profile;

          const walletUser = new WalletUser({ agent: state.instance, user });
          try {
            await walletUser.lock();
          } catch (e) {
            console.error('error locking walletUser', e);
          }

          if (state.instance) {
            await dispatch('flushStore');
            await state.instance.destroy();
          }
          commit('setInstance', {});
        },
        addNotifier({ commit, state }, notifier) {
          commit('addNotifier', notifier);
          if (state.instance) {
            commit('startNotifier', notifier);
          }
        },
      },
      getters: {
        getInstance(state) {
          return state.instance;
        },
        isInitialized(state) {
          return state.instance != null;
        },
      },
    },
  },
};

// options for creating wallet profile
function profileCreationOpts(agentOpts) {
  const { bootstrap, config } = agentOpts;

  let keyStoreURL, localKMSPassphrase, edvConfiguration;

  // webkms
  if (config.kmsType == 'webkms') {
    keyStoreURL = bootstrap.data.opsKeyStoreURL;
  }

  // local
  if (config.kmsType == 'local') {
    localKMSPassphrase = config.localKMSPassphrase;
  }

  // edv
  if (config.storageType == 'edv') {
    edvConfiguration = {
      serverURL: bootstrap.data.userEDVServer,
      vaultID: bootstrap.data.userEDVVaultID,
      encryptionKID: bootstrap.data.userEDVEncKID,
      macKID: bootstrap.data.userEDVMACKID,
    };
  }

  return { keyStoreURL, localKMSPassphrase, edvConfiguration };
}

function profileUnlockOpts(agentOpts) {
  const { bootstrap, userConfig, config } = agentOpts;

  let webKMSAuth, localKMSPassphrase, edvUnlocks;

  // webkms
  if (config.kmsType == 'webkms') {
    webKMSAuth = {
      gnapToken: userConfig.accessToken,
      secretShare: userConfig.walletSecretShare, // Not present
      capability: bootstrap.data.opsKMSCapability,
      authzKeyStoreURL: bootstrap.data.authzKeyStoreURL,
    };
  }

  // local
  if (config.kmsType == 'local') {
    localKMSPassphrase = config.localKMSPassphrase;
  }

  // edv
  if (config.storageType == 'edv') {
    edvUnlocks = {
      gnapToken: userConfig.accessToken,
      secretShare: userConfig.walletSecretShare,
      capability: bootstrap.data.edvCapability,
      authzKeyStoreURL: bootstrap.data.authzKeyStoreURL,
    };
  }

  return { webKMSAuth, localKMSPassphrase, edvUnlocks, expiry: parseTIme(bootstrap.tokenExpiry) };
}
