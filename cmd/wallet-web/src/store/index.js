/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import Vuex from 'vuex';
import connections from './modules/connections';
import mode from './modules/mode';
import mediator from './modules/mediator';
import credentials from './modules/credentials';
import presentation from './modules/presentation';
import user from './modules/user';
import options from './modules/options';
import locale from './modules/locale';
import sharedMutations from 'vuex-shared-mutations';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: { connections, mode, mediator, credentials, presentation, user, options, locale },
  plugins: [sharedMutations({ predicate: ['setUserLoggedIn', 'setLocale', 'setLogInSuspended'] })],
});
