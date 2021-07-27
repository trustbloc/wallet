/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export default {
  state: {
    locale: {
      id: '',
      base: '',
      name: '',
      translations: '',
    },
  },
  getters: {
    getLocale(state) {
      return state.locale;
    },
  },
  mutations: {
    setLocale(state, payload) {
      state.locale = payload;
    },
  },
  actions: {
    setLocale({ commit }, payload) {
      commit('setLocale', payload);
    },
  },
};
