/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import * as Agent from '@trustbloc/agent-sdk-web';
import { WalletUser } from '@trustbloc/wallet-sdk';
import { toRaw } from 'vue';
import axios from 'axios';
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
    chapi: false,
    selectedVaultId: null,
    selectedCredentialId: null,
    gnapSessionToken: null,
    gnapRequestAccessResp: null,
  },
  mutations: {
    setUser(state, val) {
      console.log('setUser to', val);
      state.username = val;
      if (val !== null) {
        localStorage.setItem('user', val);
      }
    },
    setProfile(state, val) {
      console.log('setProfile to', val);
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
      console.log('set setupStatus to', state.setupStatus);
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
    setSessionToken(state, val) {
      state.gnapSessionToken = val;
      if (val !== null) {
        localStorage.setItem('gnapSessionToken', JSON.stringify(val));
      } else {
        localStorage.removeItem('gnapSessionToken');
      }
    },
    setSubjectId(state, val) {
      state.gnapSubjectId = val;
      if (val !== null) {
        localStorage.setItem('gnapSubjectId', JSON.stringify(val));
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
      state.chapi = false;
      state.selectedVaultId = null;
      state.selectedCredentialId = null;
      state.gnapSessionToken = null;
      state.gnapSubjectId = null;

      localStorage.removeItem('user');
      localStorage.removeItem('setupStatus');
      localStorage.removeItem('profile');
      localStorage.removeItem('preference');
      localStorage.removeItem('chapi');
      localStorage.removeItem('selectedVaultId');
      localStorage.removeItem('selectedCredentialId');
      localStorage.removeItem('gnapSessionToken');
      localStorage.removeItem('gnapSubjectId');

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
      state.gnapSessionToken = localStorage.getItem('gnapSessionToken');
      state.gnapSubjectId = localStorage.getItem('gnapSubjectId');
      console.log('loadUser state:', state);
    },
  },
  actions: {
    async refreshUserPreference({ commit, state, rootGetters }, profile = state.profile) {
      if (!profile) {
        console.error('failed to refresh user preference, profile not found.');
        throw 'invalid operation, user profile not set';
      }

      let { user, token } = profile;

      console.log('user preference for', user, 'with token', token);

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
    activateCHAPI({ commit }) {
      commit('setCHAPI', true);
    },
    updateSelectedVaultId({ commit }, selectedVaultId) {
      commit('setSelectedVaultId', selectedVaultId);
    },
    updateSelectedCredentialId({ commit }, selectedCredentialId) {
      commit('setSelectedCredentialId', selectedCredentialId);
    },
    updateSessionToken({ commit }, gnapSessionToken) {
      commit('setSessionToken', gnapSessionToken);
    },
    updateSubjectId({ commit }, gnapSubjectId) {
      commit('setSubjectId', gnapSubjectId);
    },
    updateGnapReqAccessResp({ commit }, gnapRequestAccessResp) {
      console.log('called updateGnapReqAccessResp(null);');
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
    isCHAPI(state) {
      return state.chapi;
    },
    getSelectedVaultId(state) {
      return state.selectedVaultId;
    },
    getSelectedCredentialId(state) {
      return state.selectedCredentialId;
    },
    getGnapSessionToken(state) {
      state.gnapSessionToken = JSON.parse(localStorage.getItem('gnapSessionToken'));
      return state.gnapSessionToken;
    },
    getGnapSubjectId(state) {
      state.gnapSubjectId = JSON.parse(localStorage.getItem('gnapSubjectId'));
      return state.gnapSubjectId;
    },
    getGnapReqAccessResp(state) {
      state.gnapRequestAccessResp = JSON.parse(localStorage.getItem('gnapRequestAccessResp'));
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
        async init({ commit, rootState, state, rootGetters, dispatch }, gnapOpts = {}) {
          if (state.instance && state.agentName == rootState.user.username) {
            return;
          }

          if (!rootState.user.username && !gnapOpts.subjectId) {
            console.error('user should be logged in to initialize agent instance');
            throw 'invalid user state';
          }

          const gnapKeyPair = await getGnapKeyPair();
          const signingKey = await window.crypto.subtle.exportKey('jwk', gnapKeyPair.privateKey);

          let opts = {};
          Object.assign(opts, rootGetters.getAgentOpts, {
            'agent-default-label': rootState.user.username,
            'db-namespace': rootState.user.username,
            'gnap-signing-jwk': JSON.stringify(signingKey),
            'gnap-access-token': gnapOpts?.accessToken || '',
            'gnap-user-subject': gnapOpts?.subjectId || '',
          });

          console.log('all agent opts', opts);

          let profileOpts = rootGetters.getProfileOpts;
          console.log('profileOpts (in agent/init)', profileOpts);

          Object.assign(profileOpts, { userConfig: { accessToken: gnapOpts?.accessToken || '' } });

          console.log('profileOpts (in agent/init)', profileOpts);
          console.log('rootState (in agent/init)', rootState);

          await axios
            .get(rootGetters.hubAuthURL + '/gnap/bootstrap', {
              headers: { Authorization: `GNAP ${gnapOpts?.accessToken}` },
            })
            .then((resp) => {
              console.log(
                'received bootstrap resp after initializing agent' + JSON.stringify(resp, null, 2)
              );
              let { data } = resp;

              // TODO to be removed after universal wallet migration
              if (opts.storageType === 'edv') {
                const edvVaultURL = data.data.edvVaultURL;

                console.log('User EDV Vault URL is: ' + edvVaultURL);

                const edvVaultID = data.data.userEDVVaultID;

                console.log('User EDV Vault ID is: ' + edvVaultID);

                opts.edvVaultID = edvVaultID;
                // TODO this property is not returned from the bootstrap data - remove if not needed
                opts.edvCapability = data.data.edvCapability;
              }

              // TODO to be removed after universal wallet migration
              if (opts.kmsType === 'webkms') {
                opts.opsKeyStoreURL = data.data.opsKeyStoreURL;
                opts.edvOpsKIDURL = data.data.edvOpsKIDURL;
                opts.edvHMACKIDURL = data.data.edvHMACKIDURL;

                console.log('ops key store url : ', opts.opsKeyStoreURL);
                console.log('edv ops key url : ', opts.edvOpsKIDURL);
                console.log('edv ops key url : ', opts.edvHMACKIDURL);
              }

              // TODO to be removed after universal wallet migration
              // TODO this property is not returned from the bootstrap data - remove if not needed
              opts.authzKeyStoreURL = data.data.authzKeyStoreURL;
              // TODO this property is not returned from the bootstrap data - remove if not needed
              // it is not even defined in the defaultAgentStartupOpts
              // opts.userConfig = data.data.userConfig;
              // TODO this property is not returned from the bootstrap data - remove if not needed
              opts.opsKMSCapability = data.data.opsKMSCapability;

              Object.assign(profileOpts, { bootstrap: data });

              commit('updateAgentOpts', opts, { root: true });
              commit('updateProfileOpts', profileOpts, { root: true });
            })
            .catch((err) => {
              console.log('error fetching bootstrap data BEFORE Agent Init: errMsg=', err);
              console.log(
                "Note: If you haven't logged in yet and you just got a 403 error, then it's expected"
              );

              // http 400 denotes expired cookie at server - logout the user and make user to signin
              if (err.response && err.response.status === 400) {
                dispatch('logout');
              }
            });

          try {
            console.log('Agent.Framework opts: ' + JSON.stringify(opts, null, 2));
            let agent = await new Agent.Framework(opts);

            // TODO to be moved from here to 'loadOIDCUser' in case server based universal wallet.
            // if (process.env.NODE_ENV === 'production') {
            await axios
              .get(rootGetters.hubAuthURL + '/gnap/bootstrap', {
                headers: { Authorization: `GNAP ${gnapOpts?.accessToken}` },
              })
              .then((resp) => {
                console.log(
                  'received bootstrap resp after initializing agent' + JSON.stringify(resp, null, 2)
                );
                let { data } = resp;

                // TODO to be removed after universal wallet migration
                if (opts.storageType === 'edv') {
                  const edvVaultURL = data.data.edvVaultURL;

                  console.log('User EDV Vault URL is: ' + edvVaultURL);

                  const edvVaultID = data.data.userEDVVaultID;

                  console.log('User EDV Vault ID is: ' + edvVaultID);

                  opts.edvVaultID = edvVaultID;
                  // TODO this property is not returned from the bootstrap data - remove if not needed
                  opts.edvCapability = data.data.edvCapability;
                }

                // TODO to be removed after universal wallet migration
                if (opts.kmsType === 'webkms') {
                  opts.opsKeyStoreURL = data.data.opsKeyStoreURL;
                  opts.edvOpsKIDURL = data.data.edvOpsKIDURL;
                  opts.edvHMACKIDURL = data.data.edvHMACKIDURL;

                  console.log('ops key store url : ', opts.opsKeyStoreURL);
                  console.log('edv ops key url : ', opts.edvOpsKIDURL);
                  console.log('edv ops key url : ', opts.edvHMACKIDURL);
                }

                // TODO to be removed after universal wallet migration
                // TODO this property is not returned from the bootstrap data - remove if not needed
                opts.authzKeyStoreURL = data.data.authzKeyStoreURL;
                // TODO this property is not returned from the bootstrap data - remove if not needed
                // it is not even defined in the defaultAgentStartupOpts
                // opts.userConfig = data.data.userConfig;
                // TODO this property is not returned from the bootstrap data - remove if not needed
                opts.opsKMSCapability = data.data.opsKMSCapability;

                Object.assign(profileOpts, { bootstrap: data });

                commit('updateAgentOpts', opts, { root: true });
                commit('updateProfileOpts', profileOpts, { root: true });
              })
              .catch((err) => {
                console.log('error fetching bootstrap data: errMsg=', err);
                console.log(
                  "Note: If you haven't logged in yet and you just got a 403 error, then it's expected"
                );

                // http 400 denotes expired cookie at server - logout the user and make user to signin
                if (err.response && err.response.status === 400) {
                  dispatch('logout');
                }
              });
            commit('setInstance', { instance: agent, user: rootState?.user.username });
            console.log('[user.js] getInstance', state.instance);
            await dispatch('unlockWallet');
            commit('startAllNotifiers');
          } catch (e) {
            console.error(e);
          }
        },
        async unlockWallet({ state, rootGetters, dispatch }) {
          // create wallet profile if it doesn't exist

          let profileOpts = rootGetters.getProfileOpts;

          console.log('profile opts', JSON.stringify(profileOpts, null, 2));
          console.log('opts bootstrap', profileOpts.bootstrap);

          let { user } = profileOpts.bootstrap.data;

          console.log('unlocking wallet for user ', user);
          console.log('agent', state.instance);

          console.log('initializing WalletUser with:', {
            agent: rootGetters['agent/getInstance'],
            user,
          });
          let walletUser = new WalletUser({ agent: state.instance, user });
          if (!(await walletUser.profileExists())) {
            console.log('wallet for this profile did not exist');
            let createOpts = profileCreationOpts(profileOpts);

            console.log('profile create opts', createOpts);

            await walletUser.createWalletProfile(createOpts);
          }

          console.log('unlocking wallet with:', profileUnlockOpts(profileOpts));
          let { token } = await walletUser.unlock(profileUnlockOpts(profileOpts));
          console.log('unlocked. user, token', user, token);

          console.log('refreshing user preference');
          await dispatch('refreshUserPreference', { user, token }, { root: true });

          const currentUser = rootGetters.getCurrentUser;
          console.log('unlock wallet currentUser', currentUser);

          if (!currentUser.preference) {
            // todo: should this be rootGetters.getCurrentUser.value.preference?
            dispatch('startUserSetup', '', { root: true });

            console.log('AAAAstate.instance', state.instance);

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
          console.log('updated preferences in state', rootGetters.getCurrentUser);
        },
        async flushStore({ state }) {
          console.debug('flushing store', state.instance);
          if (state.instance) {
            try {
              await state.instance.store.flush();
            } catch (e) {
              console.error('error flushing store', e);
            }
            console.debug('flushed store.');
          }
        },
        async destroy({ commit, state, rootGetters, dispatch }) {
          let { user } = rootGetters.getCurrentUser.profile;
          console.log('signing user out', user);

          let walletUser = new WalletUser({ agent: state.instance, user });
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
function profileCreationOpts(opts) {
  console.log('profileCreationOpts', opts);
  let { bootstrap, config } = opts;
  let keyStoreURL, localKMSPassphrase, edvConfiguration;
  // webkms
  if (config.kmsType == 'webkms') {
    keyStoreURL = bootstrap.data.opsKeyStoreURL;
  }

  // local
  if (config.kmsType == 'local') {
    localKMSPassphrase = config.localKMSScret;
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

function profileUnlockOpts(opts) {
  console.log('profileUnlockOpts', opts);
  let { bootstrap, userConfig, config } = opts;

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
    localKMSPassphrase = config.localKMSScret;
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
