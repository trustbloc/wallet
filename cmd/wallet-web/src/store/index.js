/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createStore } from 'vuex';
import connections from './modules/connections';
import mode from './modules/mode';
import mediator from './modules/mediator';
import credentials from './modules/credentials';
import presentation from './modules/presentation';
import user from './modules/user';
import options from './modules/options';
import locale from './modules/locale';
import sharedMutations from 'vuex-shared-mutations';

const store = createStore({
  modules: { connections, mode, mediator, credentials, presentation, user, options, locale },
  plugins: [
    sharedMutations({
      predicate: ['setUserLoggedIn', 'setLocale', 'setLogInSuspended'],
    }),
  ],
});

export default store;
