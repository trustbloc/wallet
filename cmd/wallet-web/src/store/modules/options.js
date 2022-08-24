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
  enableDIDComm: false,
  enableCHAPI: false,
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

      const chooseOptValue = (opt) =>
        opt in agentOpts ? agentOpts[opt] : defaultAgentStartupOpts[opt];

      commit('updateAgentOpts', {
        assetsPath: defaultAgentStartupOpts['assetsPath'],
        'outbound-transport': defaultAgentStartupOpts['outbound-transport'],
        'transport-return-route': defaultAgentStartupOpts['transport-return-route'],
        'http-resolver-url': chooseOptValue('http-resolver-url'),
        'agent-default-label': chooseOptValue('agent-default-label'),
        'auto-accept': chooseOptValue('auto-accept'),
        'log-level': chooseOptValue('log-level'),
        'indexed-db-namespace': chooseOptValue('indexed-db-namespace'),
        'edge-agent-server': chooseOptValue('edge-agent-server'),
        'context-provider-url': chooseOptValue('context-provider-url'),
        'bloc-domain': chooseOptValue('bloc-domain'),
        walletMediatorURL: chooseOptValue('walletMediatorURL'),
        credentialMediatorURL: credentialMediator(chooseOptValue('credentialMediatorURL')),
        blindedRouting: chooseOptValue('blindedRouting'),
        'storage-type': chooseOptValue('storage-type'),
        'edv-server-url': chooseOptValue('edv-server-url'),
        'edv-vault-id': chooseOptValue('edv-vault-id'),
        'edv-capability': chooseOptValue('edv-capability'),
        'authz-key-store-url': chooseOptValue('authz-key-store-url'),
        'user-config': chooseOptValue('user-config'),
        'use-edv-cache': chooseOptValue('use-edv-cache'),
        'edv-clear-cache': chooseOptValue('edv-clear-cache'),
        'kms-type': chooseOptValue('kms-type'),
        'ops-key-store-url': chooseOptValue('ops-key-store-url'),
        'edv-ops-kid-url': chooseOptValue('edv-ops-kid-url'),
        'edv-hmac-kid-url': chooseOptValue('edv-hmac-kid-url'),
        'ops-kms-capability': chooseOptValue('ops-kms-capability'),
        'use-edv-batch': chooseOptValue('use-edv-batch'),
        'edv-batch-size': chooseOptValue('edv-batch-size'),
        'unanchored-din-max-life-time': chooseOptValue('unanchored-din-max-life-time'),
        'cache-size': chooseOptValue('cache-size'),
        'did-anchor-origin': chooseOptValue('did-anchor-origin'),
        'sidetree-token': chooseOptValue('sidetree-token'),
        'hub-auth-url': chooseOptValue('hub-auth-url'),
        walletWebUrl: chooseOptValue('walletWebUrl'),
        staticAssetsUrl: chooseOptValue('staticAssetsUrl'),
        'media-type-profiles': chooseOptValue('media-type-profiles'),
        'key-type': chooseOptValue('key-type'),
        'key-agreement-type': chooseOptValue('key-agreement-type'),
        'web-socket-read-limit': chooseOptValue('web-socket-read-limit'),
        'kms-server-url': chooseOptValue('kms-server-url'),
        enableDIDComm: chooseOptValue('enableDIDComm'),
        enableCHAPI: chooseOptValue('enableCHAPI'),
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
    getCredentialMediatorURL(state) {
      return state.agentOpts['credentialMediatorURL'];
    },
    getEnableDIDComm(state) {
      return state.agentOpts['enableDIDComm'];
    },
    getEnableCHAPI(state) {
      return state.agentOpts['enableCHAPI'];
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
