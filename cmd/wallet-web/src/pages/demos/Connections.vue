/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <div class="content">
    <div class="md-layout">
      <div class="md-layout-item">
        <div class="md-layout-item">
          <mediator title="Mediator" />
        </div>
        <div class="md-layout-item">
          <create-invitation />
        </div>
        <div class="md-layout-item">
          <receive-invitation />
        </div>
      </div>
      <div class="md-layout-item">
        <div class="md-layout-item">
          <connections
            :short="false"
            title="Connections"
            :count="pendingConnectionsCount"
            :connections="allConnections"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Connections, CreateInvitation, Mediator, ReceiveInvitation } from '@/components';
import { mapActions, mapGetters } from 'vuex';

export default {
  components: { Mediator, CreateInvitation, ReceiveInvitation, Connections },
  methods: {
    ...mapActions(['queryConnections', 'onDidExchangeState', 'loadMediatorState']),
    ...mapActions('agent', { addAgentNotifiers: 'addNotifier' }),
  },
  mounted() {
    this.addAgentNotifiers({ callback: this.onDidExchangeState, topics: ['didexchange_states'] });
    this.loadMediatorState();
    // refreshes connections when component is mounted
    this.queryConnections();
  },
  computed: mapGetters(['pendingConnectionsCount', 'allConnections']),
};
</script>
