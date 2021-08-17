<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <button
    type="button"
    class="flex flex-row justify-start items-center w-full h-16"
    @click="signout"
  >
    <img src="@/assets/img/signout.svg" />
    <span class="ml-4 text-lg text-neutrals-white">{{ i18n.signout }}</span>
  </button>
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
