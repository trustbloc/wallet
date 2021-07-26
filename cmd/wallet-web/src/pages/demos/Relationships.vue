<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="content relationships">
    <mediator
      v-if="!isMediatorRegistered"
      title="Please, set up a mediator to proceed with this page!"
    />
    <div v-if="isMediatorRegistered" class="md-layout">
      <div class="md-layout-item">
        <div v-if="associatedCredentialsCount" class="md-layout-item">
          <associated-credentials
            title="Associated credentials"
            :credentials="associatedCredentials"
          />
        </div>
        <div v-if="associatedPresentationsCount" class="md-layout-item">
          <associated-presentation
            title="Associated presentations"
            :presentations="associatedPresentations"
          />
        </div>
        <div class="md-layout-item">
          <public-invitation />
        </div>
        <div class="md-layout-item">
          <receive-invitation
            title="Have an invitation? Put it below and hit the Connect button."
            type="base64"
          />
        </div>
      </div>
      <div class="md-layout-item">
        <div v-if="pendingConnectionsCount" class="md-layout-item">
          <connections
            title="Pending requests"
            :count="pendingConnectionsCount"
            :connections="pendingConnections"
          />
        </div>
        <div v-if="completedConnectionsCount" class="md-layout-item">
          <connections
            title="List of agents you have a connection with"
            :connections="completedConnections"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import {
  Mediator,
  PublicInvitation,
  ReceiveInvitation,
  Connections,
  AssociatedCredentials,
  AssociatedPresentation,
} from '@/components';

export default {
  components: {
    Mediator,
    PublicInvitation,
    ReceiveInvitation,
    Connections,
    AssociatedCredentials,
    AssociatedPresentation,
  },
  methods: mapActions(['queryConnections', 'loadMediatorState']),
  computed: mapGetters([
    'pendingConnections',
    'completedConnections',
    'pendingConnectionsCount',
    'completedConnectionsCount',
    'isMediatorRegistered',
    'associatedCredentials',
    'associatedCredentialsCount',
    'associatedPresentations',
    'associatedPresentationsCount',
  ]),
  mounted() {
    // load mediator registration details if already registered
    this.loadMediatorState();
    // refreshes connections when component is mounted
    this.queryConnections();
  },
};
</script>

<style>
.relationships .md-list.md-triple-line .md-list-item-content {
  min-height: 30px !important;
}
</style>
