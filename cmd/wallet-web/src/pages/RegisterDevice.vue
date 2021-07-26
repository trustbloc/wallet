/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <div style="display: flex">
    <md-ripple>
      <md-button
        v-if="!registered"
        :md-ripple="false"
        class="md-primary"
        style="background-color: #4f424c !important; font-size: 16px"
        @click="registerDevice"
        >Register Device
      </md-button>
    </md-ripple>
    <div v-if="register === 'success'">
      <i class="md-title" style="font-size: 20px; color: white !important">
        <md-icon style="color: green; margin: 10px">check_circle</md-icon>
        Device successfully registered.
      </i>
    </div>
    <div v-if="register === 'failure'">
      <i class="md-title" style="font-size: 20px; color: white !important">
        <md-icon style="color: red; margin: 10px">report</md-icon>
        Device failed to register.
      </i>
    </div>
  </div>
</template>

<script>
import { DeviceRegister } from '@trustbloc/wallet-sdk';
import { mapActions, mapGetters } from 'vuex';
export default {
  data() {
    return {
      register: 'none',
      registered: false,
    };
  },
  created: async function () {
    this.deviceRegister = new DeviceRegister(this.getAgentOpts()['edge-agent-server']);
    if (this.$cookies.isKey('registerSuccess')) {
      this.registered = true;
    }
  },
  methods: {
    ...mapActions({ logoutUser: 'logout' }),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getAgentOpts']),
    registerDevice: async function () {
      this.register = await this.deviceRegister.register();
      if (this.register == 'success') {
        this.$cookies.set('registerSuccess', 'success');
        this.registered = true;
      }
    },
  },
};
</script>
