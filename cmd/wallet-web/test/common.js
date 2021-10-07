/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import * as Agent from '@trustbloc/agent-sdk-web';
import { DIDManager, WalletUser, CredentialManager } from '@trustbloc/wallet-sdk';
import Vuex from 'vuex';
import { createLocalVue } from '@vue/test-utils';

var uuid = require('uuid/v4');

export const DIDEXCHANGE_STATE_TOPIC = 'didexchange_states';
export const POST_STATE = 'post_state';
export const DIDEXCHANGE_STATE_REQUESTED = 'requested';

export const testConfig = window.__FIXTURES__
  ? JSON.parse(window.__FIXTURES__['test/fixtures/agent-config.json'])
  : {};
testConfig.walletUserPassphrase = uuid();

console.debug('test configuration:', JSON.stringify(testConfig, null, 2));

const { agentStartupOpts } = testConfig;

export const localVue = createLocalVue();
localVue.use(Vuex);

// loads testdata from fixtures and returns JSON parsed response.
export function getTestData(filename) {
  return JSON.parse(window.__FIXTURES__[`test/fixtures/testdata/${filename}`]);
}

// loadFrameworks loads agent instance
export async function loadFrameworks({ name = 'user-agent', logLevel = '' } = {}) {
  let agentOpts = JSON.parse(JSON.stringify(agentStartupOpts));
  agentOpts['indexedDB-namespace'] = `${name}db`;
  agentOpts['agent-default-label'] = `${name}-wallet-web`;

  if (logLevel) {
    agentOpts['log-level'] = logLevel;
  }

  return new Agent.Framework(agentOpts);
}

/**
 *  Wallet Test Setup
 */
export class Setup {
  constructor({ user = 'default_user' } = {}) {
    this.user = user;
    this.agent = {};
  }

  async loadAgent({ logLevel = '' } = {}) {
    this.agent = await loadFrameworks({ name: `${this.user}-agent`, logLevel });
  }

  destroyAgent() {
    if (this.agent) {
      this.agent.destroy();
    }
  }

  async createProfile() {
    if (!this.agent) {
      throw 'invalid operation, agent not initialized';
    }

    await new WalletUser({
      agent: this.agent,
      user: this.user,
    }).createWalletProfile({ localKMSPassphrase: testConfig.walletUserPassphrase });
  }

  async unlockWallet() {
    let { token } = await new WalletUser({
      agent: this.agent,
      user: this.user,
    }).unlock({ localKMSPassphrase: testConfig.walletUserPassphrase });

    this.token = token;
  }

  async createPreference() {
    let didManager = new DIDManager({ agent: this.agent, user: this.user });
    let docRes = await didManager.createOrbDID(this.token, {
      purposes: ['authentication', 'assertionMethod'],
    });

    this.preference = {
      controller: docRes.DIDDocument.id,
      proofType: 'Ed25519Signature2018',
    };

    await new WalletUser({ agent: this.agent, user: this.user }).savePreferences(
      this.token,
      this.preference
    );
  }

  async saveCredentials(...credentials) {
    await new CredentialManager({ agent: this.agent, user: this.user }).save(this.token, {
      credentials,
    });
  }

  getStateStore() {
    const user = this.user;
    const token = this.token;
    const agent = this.agent;
    const preference = this.preference;

    let store = {
      getters: {
        getCurrentUser(state) {
          return {
            username: user,
            profile: { user, token },
            preference,
          };
        },
        getAgentOpts(state) {
          return agentStartupOpts;
        },
        getCredentialManifestData(state) {
          return require('@/config/credentialDisplayData.js').default;
        },
      },
      modules: {
        agent: {
          namespaced: true,
          actions: {
            async init({ commit, rootState, state }) {},
            async destroy({ commit, state }) {},
          },
          getters: {
            getInstance(state) {
              return agent;
            },
          },
        },
      },
    };

    return new Vuex.Store(store);
  }
}

/**
 *  Mock credential handler for tests.
 */
export class MockCredentialHandler {
  constructor() {
    this.eventQueue = [];
  }

  setRequestEvent(event) {
    let respond;
    event.respondWith = async (promise) => {
      respond(promise);
    };

    this.eventQueue.push(event);

    // handle for event response
    return new Promise((resolve, reject) => {
      const timer = setTimeout(
        (_) => reject(new Error('timeout waiting for credential event response')),
        40000
      );
      respond = async (result) => {
        clearTimeout(timer);
        resolve(await result);
      };
    });
  }

  async receiveCredentialEvent() {
    let event = this.eventQueue.pop();

    return new Promise((resolve, reject) => {
      if (!event) {
        reject(event);
      }

      resolve(event);
    });
  }
}

// promise on resolving given function to true or timeout
export function promiseWhen(fn, timeout, interval) {
  function loop(resolve) {
    if (fn()) {
      resolve();
    }
    setTimeout(() => loop(resolve), interval ? interval : 300);
  }

  return new Promise(function (resolve, reject) {
    setTimeout(
      (_) => reject(new Error('timeout waiting for condition')),
      timeout ? timeout : 10000
    );
    loop(resolve);
  });
}
