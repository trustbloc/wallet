/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import Vuex from 'vuex';
import connections from './modules/connections'
import mode from './modules/mode'
import mediator from './modules/mediator'
import credentials from './modules/credentials'
import presentation from './modules/presentation'
import user from './modules/user'
import options from './modules/options'
import loginhandle from './modules/loginhandle'
//import createPersistedState from 'vuex-persistedstate'
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {connections, mode, mediator, credentials, presentation, user, options, loginhandle,},
    plugins: [createPersistedState({
        storage: window.sessionStorage,
    })],
})
