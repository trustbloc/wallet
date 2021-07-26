<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex justify-center w-full">
    <button
      type="button"
      class="w-full border-2 border-solid border-red-50 md-button logout-button"
      @click="logout"
    >
      Log Out
    </button>
  </div>
</template>

<script>
import { CHAPIHandler } from './mixins';
import { mapActions, mapGetters } from 'vuex';

export default {
  created: function () {
    this.chapi = new CHAPIHandler(
      this.$polyfill,
      this.$webCredentialHandler,
      this.getAgentOpts().credentialMediatorURL
    );
  },
  methods: {
    ...mapActions({ logoutUser: 'logout' }),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getAgentOpts']),
    logout: async function () {
      await this.chapi.uninstall();
      await this.logoutUser();
      this.$router.push('/signup');
    },
  },
};
</script>

<style scoped>
.logout-button {
  height: 48px;
  font-size: 16px; /* remove once global styles for new designs are defined */
  font-weight: bold;
  border: 2px solid;
  border-image-slice: 1;
  border-image-source: linear-gradient(to left, #3f5fd3, #743ad5, #d53a9d, #cd3a67);
}
/*--Remove this once vue-material css is removed */
.md-button {
  text-transform: none !important;
  background: transparent !important;
  font-size: large;
  font-family: sans-serif;
}
.md-button:hover {
  background: transparent !important;
}
.md-button:active:after {
  background: transparent !important;
}
</style>
