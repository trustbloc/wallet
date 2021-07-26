/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export default {
  actions: {
    loadMode({ commit }) {
      commit('updateDevMode', localStorage.getItem('devMode') === 'true');
    },
  },
  mutations: {
    updateDevMode(state, val) {
      state.devMode = val;
      localStorage.setItem('devMode', val);
    },
  },
  state: {
    devMode: false,
  },
  getters: {
    isDevMode(state) {
      return state.devMode;
    },
  },
};
