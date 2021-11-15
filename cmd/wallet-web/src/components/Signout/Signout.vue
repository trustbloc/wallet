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
      focus-within:bg-gradient-to-r
      focus-within:from-neutrals-black
      focus-within:shadow-inner-outline-blue
      opacity-60
      focus-within:opacity-100
      hover:from-neutrals-black hover:bg-gradient-to-r hover:opacity-100
      flex flex-row
      bar
    "
  >
    <button
      class="flex flex-row justify-start items-center px-10 w-full h-16"
      type="button"
      @click="signout()"
      @keyup.enter="signout()"
    >
      <img src="@/assets/img/signout.svg" />
      <span class="ml-4 text-lg font-bold text-neutrals-white">{{ t('Signout.signout') }}</span>
    </button>
  </div>
</template>

<script>
import { CHAPIHandler } from '@/pages/mixins';
import { mapActions, mapGetters } from 'vuex';
import { useI18n } from 'vue-i18n';

export default {
  name: 'Signout',
  setup() {
    const { t } = useI18n();
    return { t };
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
    ...mapGetters(['getAgentOpts', 'isCHAPI']),
    signout: async function () {
      let actions = [this.signoutUser()];
      if (this.isCHAPI()) {
        actions.push(this.chapi.uninstall());
      }

      await Promise.all(actions);

      this.$router.push({ name: 'signin' });
    },
  },
};
</script>

<style scoped>
.bar:not(:focus-within):hover:before {
  position: absolute;
  content: '';
  display: block;
  background-color: theme('colors.primary.purple.hashita');
  height: theme('spacing.16');
  width: 4px;
}
</style>
