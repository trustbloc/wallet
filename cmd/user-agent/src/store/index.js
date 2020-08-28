/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import Vuex from 'vuex';
import connections from './modules/connections'
import mode from './modules/mode'
import init from './modules/init'
import mediator from './modules/mediator'
import credentials from './modules/credentials'
import presentation from './modules/presentation'

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {init, connections, mode, mediator, credentials, presentation}
})