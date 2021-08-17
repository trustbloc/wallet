<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="
      justify-start
      items-center
      w-full
      h-16
      opacity-60
      hover:from-neutrals-black hover:bg-gradient-to-r hover:opacity-100
      flex flex-row
      bar
    "
  >
    <button
      class="flex flex-row justify-start items-center px-10 w-full"
      type="button"
      @click="signout"
    >
      <img src="@/assets/img/signout.svg" />
      <span class="ml-4 text-lg font-bold text-neutrals-white">{{ i18n.signout }}</span>
    </button>
  </div>
</template>

<script>
import { CHAPIHandler } from '@/pages/mixins';
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'Signout',
  computed: {
    i18n() {
      return this.$t('Signout');
    },
  },
  created: function () {
    this.chapi = new CHAPIHandler(
      this.$polyfill,
      this.$webCredentialHandler,
      this.getAgentOpts().credentialMediatorURL
    );
  },
  methods: {
    ...mapActions({ signoutUser: 'logout' }),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getAgentOpts']),
    signout: async function () {
      await Promise.all([this.chapi.uninstall(), this.signoutUser()]);
      this.$router.push({ name: 'signin', params: this.$route.params });
    },
  },
};
</script>

<style scoped>
.bar:hover:before {
  position: absolute;
  content: "''";
  display: block;
  background-color: theme('colors.primary.purple.hashita');
  height: theme('spacing.16');
  width: 4px;
}
</style>
