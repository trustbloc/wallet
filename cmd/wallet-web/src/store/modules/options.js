/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { toRaw } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import axios from 'axios';

const agentOptsLocation = (l) => `${l}/walletconfig/agent`;
const credentialMediator = (url) =>
  url
    ? `${url}?origin=${encodeURIComponent(window.location.origin)}${__webpack_public_path__}/`
    : undefined;

let defaultAgentStartupOpts = {
  assetsPath: `${__webpack_public_path__}agent-js-worker/assets`,
  'outbound-transport': ['ws', 'http'],
  'transport-return-route': 'all',
  'http-resolver-url': ['web@http://localhost:9080/1.0/identifiers'],

  'agent-default-label': 'demo-wallet-web',
  'auto-accept': true,
  'log-level': 'debug',
  'indexed-db-namespace': 'agent',
  // default backend server url
  'edge-agent-server': 'https://localhost:9099',
  walletWebUrl: 'https://localhost:9098',
  // remote JSON-LD context provider urls
  'context-provider-url': [],

  'bloc-domain': 'testnet.orb.local',
  walletMediatorURL: '',
  blindedRouting: false,
  credentialMediatorURL: '',
  'storage-type': `indexedDB`,
  'edv-server-url': '',
  'edv-vault-id': '',
  'edv-capability': '',
  'authz-key-store-url': '',
  'kms-type': `local`,
  localKMSPassphrase: `demo`,
  'use-edv-cache': false,
  'edv-clear-cache': '',
  'ops-kms-capability': '',
  'use-edv-batch': false,
  'cache-size': 100,
  'edv-batch-size': 0,
  'did-anchor-origin': 'origin',
  'sidetree-token': '',
  'hub-auth-url': '',
  staticAssetsUrl: '',
  'unanchored-din-max-life-time': 0,
  'media-type-profiles': ['didcomm/aip2;env=rfc587', 'didcomm/v2'],
  'key-type': 'ecdsap256ieee1363',
  'key-agreement-type': 'p256kw',
  'web-socket-read-limit': 0,
  'kms-server-url': '',
};

export default {
  actions: {
    async initOpts(
      { commit, getters, dispatch },
      { location = window.location.origin, accessToken } = {}
    ) {
      let agentOpts = {};

      const profileOpts = getters.getProfileOpts;
      if (accessToken) {
        Object.assign(profileOpts, {
          userConfig: { accessToken },
        });
      }

      let readCredentialManifests;

      if (process.env.NODE_ENV === 'production') {
        // call service to get the agent opts
        await axios
          .get(agentOptsLocation(location))
          .then((resp) => {
            agentOpts = resp.data;
          })
          .catch((err) => {
            console.log('error fetching start up options - using default options : errMsg=', err);
          });

        agentOpts['http-resolver-url'] = agentOpts['http-resolver-url'].split(',');
        agentOpts['context-provider-url'] = agentOpts['context-provider-url']
          ? agentOpts['context-provider-url'].split(',')
          : [];
        agentOpts['media-type-profiles'] = agentOpts['media-type-profiles']
          ? agentOpts['media-type-profiles'].split(',')
          : ['didcomm/aip2;env=rfc587', 'didcomm/v2'];

        readCredentialManifests = readManifests(agentOpts['staticAssetsUrl']);

        Object.assign(profileOpts, {
          config: {
            storageType: agentOpts['storage-type'],
            kmsType: agentOpts['kms-type'],
            localKMSPassphrase: agentOpts.localKMSPassphrase,
          },
        });
      } else {
        // strictly, for dev mode only

        // if you clear browser data, no user profile would be found locally, so we generate new user each time
        // thus, we create a new wallet each time too, so we can't access existing user's data
        let user = uuidv4();

        dispatch('loadUser');
        if (getters.getCurrentUser) {
          const { profile } = getters.getCurrentUser;
          user = profile ? profile.user : user;
        }

        // dev mode agent opts
        agentOpts.walletMediatorURL = 'https://localhost:10093';
        agentOpts['hub-auth-url'] = 'https://localhost:8044';

        Object.assign(profileOpts, {
          bootstrap: {
            data: {
              user,
              tokenExpiry: '10',
            },
          },
          config: {
            storageType: defaultAgentStartupOpts['storage-type'],
            kmsType: defaultAgentStartupOpts['kms-type'],
            localKMSPassphrase: defaultAgentStartupOpts.localKMSPassphrase,
          },
        });

        readCredentialManifests = readManifests();
      }

      const optValue = (opt) => (opt in agentOpts ? agentOpts[opt] : defaultAgentStartupOpts[opt]);

      commit('updateAgentOpts', {
        assetsPath: defaultAgentStartupOpts['assetsPath'],
        'outbound-transport': defaultAgentStartupOpts['outbound-transport'],
        'transport-return-route': defaultAgentStartupOpts['transport-return-route'],
        'http-resolver-url': optValue('http-resolver-url'),
        'agent-default-label': optValue('agent-default-label'),
        'auto-accept': optValue('auto-accept'),
        'log-level': optValue('log-level'),
        'indexed-db-namespace': optValue('indexed-db-namespace'),
        'edge-agent-server': optValue('edge-agent-server'),
        'context-provider-url': optValue('context-provider-url'),
        'bloc-domain': optValue('bloc-domain'),
        walletMediatorURL: optValue('walletMediatorURL'),
        credentialMediatorURL: credentialMediator(optValue('credentialMediatorURL')),
        blindedRouting: optValue('blindedRouting'),
        'storage-type': optValue('storage-type'),
        'edv-server-url': optValue('edv-server-url'),
        'edv-vault-id': optValue('edv-vault-id'),
        'edv-capability': optValue('edv-capability'),
        'authz-key-store-url': optValue('authz-key-store-url'),
        'user-config': optValue('user-config'),
        'use-edv-cache': optValue('use-edv-cache'),
        'edv-clear-cache': optValue('edv-clear-cache'),
        'kms-type': optValue('kms-type'),
        'ops-key-store-url': optValue('ops-key-store-url'),
        'edv-ops-kid-url': optValue('edv-ops-kid-url'),
        'edv-hmac-kid-url': optValue('edv-hmac-kid-url'),
        'ops-kms-capability': optValue('ops-kms-capability'),
        'use-edv-batch': optValue('use-edv-batch'),
        'edv-batch-size': optValue('edv-batch-size'),
        'unanchored-din-max-life-time': optValue('unanchored-din-max-life-time'),
        'cache-size': optValue('cache-size'),
        'did-anchor-origin': optValue('did-anchor-origin'),
        'sidetree-token': optValue('sidetree-token'),
        'hub-auth-url': optValue('hub-auth-url'),
        walletWebUrl: optValue('walletWebUrl'),
        staticAssetsUrl: optValue('staticAssetsUrl'),
        'media-type-profiles': optValue('media-type-profiles'),
        'key-type': optValue('key-type'),
        'key-agreement-type': optValue('key-agreement-type'),
        'web-socket-read-limit': optValue('web-socket-read-limit'),
        'kms-server-url': optValue('kms-server-url'),
      });

      commit('updateProfileOpts', profileOpts);

      const manifests = await readCredentialManifests;
      commit('updateCredentialManifests', manifests);
    },
  },
  mutations: {
    updateAgentOpts(state, opts) {
      state.agentOpts = opts;
    },
    updateProfileOpts(state, opts) {
      state.profileOpts = opts;
    },
    updateCredentialManifests(state, manifests) {
      state.credentialManifests = manifests;
    },
  },
  state: {
    agentOpts: {},
    profileOpts: {},
    credentialManifests: {},
  },
  getters: {
    getAgentOpts(state) {
      return toRaw(state.agentOpts);
    },
    getProfileOpts(state) {
      return state.profileOpts;
    },
    agentDefaultLabel(state) {
      return state.agentOpts['agent-default-label'];
    },
    serverURL(state) {
      return state.agentOpts['edge-agent-server'];
    },
    hubAuthURL(state) {
      return state.agentOpts['hub-auth-url'];
    },
    walletWebUrl(state) {
      return state.agentOpts['walletWebUrl'];
    },
    getStaticAssetsUrl(state) {
      return state.agentOpts['staticAssetsUrl'];
    },
    getCredentialManifests(state) {
      return toRaw(state.credentialManifests);
    },
    async getGnapAccessTokenConfig(state) {
      const staticUrl = state.agentOpts['staticAssetsUrl'];
      if (staticUrl) {
        const response = await axios.get(`${staticUrl}/config/gnap-access-token.json`);
        return response.data;
      }

      return require('@/config/gnap-access-token.json');
    },
  },
};

const readManifests = async (staticUrl) => {
  if (staticUrl) {
    try {
      const response = await axios.get(`${staticUrl}/config/credential-output-descriptors.json`);
      return response.data;
    } catch (e) {
      console.warn(e);
      console.warn(
        'unable to read credential manifests from remote location, switching to default manifests'
      );
    }
  }

  return require('@/config/credential-output-descriptors.json');
};
