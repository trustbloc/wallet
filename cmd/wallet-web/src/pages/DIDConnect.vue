<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div v-if="loading" class="w-screen" style="margin-left: 40%; margin-top: 20%; height: 200px">
    <div class="md-layout">
      <md-progress-spinner
        :md-diameter="100"
        class="md-accent"
        :md-stroke="10"
        md-mode="indeterminate"
      ></md-progress-spinner>
    </div>
  </div>

  <div v-else class="flex justify-center w-screen md-layout">
    <div class="max-w-screen-md md-layout-item">
      <md-card class="md-card-plain">
        <md-card-header>
          <h4 class="title">Connect your wallet</h4>
        </md-card-header>

        <md-card-content style="background-color: white">
          <div v-if="errors.length">
            <b>Failed with following error(s):</b>
            <ul>
              <li v-for="error in errors" :key="error" style="color: #9d0006">{{ error }}</li>
            </ul>
          </div>

          <md-card-content class="viewport">
            This issuer would like to connect to your wallet for secured communication.
            <governance :govn-v-c="govnVC" :request-origin="requestOrigin" />
          </md-card-content>

          <md-card-content v-if="userCredentials.length" class="viewport">
            Here are the credentials being sent to your wallet,

            <md-list class="md-double-line">
              <md-list-item v-for="credential in userCredentials" :key="credential">
                <md-icon class="md-primary md-size-2x">perm_identity</md-icon>

                <div class="md-list-item-text">
                  <span>{{
                    credential.name ? credential.name : 'Credential name not provided'
                  }}</span>
                  <span>{{ credential.description }}</span>
                </div>
              </md-list-item>
            </md-list>
          </md-card-content>

          <md-divider></md-divider>

          <md-card-content class="md-layout md-alignment-center-center">
            <md-button
              id="didconnect"
              style="margin-right: 5%"
              class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
              @click="connect"
              >{{ buttonLabel }}
            </md-button>
            <md-button id="cancelBtn" class="md-cancel-text" @click="cancel"> Cancel </md-button>
          </md-card-content>
        </md-card-content>
      </md-card>
    </div>
  </div>
</template>
<script>
import { DIDConn } from './mixins';
import Governance from './Governance.vue';
import { mapGetters } from 'vuex';

export default {
  components: { Governance },
  data() {
    return {
      errors: [],
      requestOrigin: '',
      loading: true,
      userCredentials: [],
      buttonLabel: 'Connect',
      govnVC: null,
    };
  },
  created: function () {
    this.wallet = new DIDConn(
      this.getAgentInstance(),
      this.getCurrentUser().profile,
      this.getAgentOpts(),
      this.$parent.credentialEvent
    );

    this.requestOrigin = this.wallet.chapiHandler.getRequestor();
    this.userCredentials = this.wallet.userCredentials;
    this.govnVC = this.wallet.govnVC;
    this.buttonLabel = this.userCredentials.length > 0 ? 'Store & Connect' : 'Connect';

    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getAgentOpts']),
    cancel: function () {
      this.wallet.cancel();
    },
    connect: async function () {
      this.loading = true;
      this.errors = [];
      try {
        await this.wallet.connect(this.getCurrentUser().preference);
      } catch (e) {
        console.error('failed to connect', e);
        this.errors.push('failed to perform connection, please try again later');
      }

      this.loading = false;
    },
  },
};
</script>
