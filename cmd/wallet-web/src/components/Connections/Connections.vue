/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <md-card class="md-card-plain">
    <md-card-header data-background-color="green">
      <h4 class="title">
        {{ title }}
        <md-badge
          v-if="count"
          class="md-primary md-square"
          style="position: relative; margin-top: 5px"
          :md-content="count"
        />
      </h4>
      <div class="title-btn-right">
        <md-button class="md-icon-button md-dense md-raised md-info" @click="queryConnections">
          <md-icon>cached</md-icon>
        </md-button>
      </div>
    </md-card-header>
    <md-card-content class="white">
      <div v-if="connections.length === 0" class="text-center">No connections</div>
      <md-content class="md-content-connections md-scrollbar">
        <md-list class="md-triple-line">
          <md-list-item v-for="conn in connections" :key="conn.id">
            <div class="md-list-item-text">
              <span v-if="short && conn.TheirLabel">{{ conn.TheirLabel }}</span>
              <span v-if="short && !conn.TheirLabel">{{ conn.TheirDID }}</span>
              <span v-if="!short && conn.ConnectionID">ConnectionID: {{ conn.ConnectionID }}</span>
              <span v-if="!short && conn.State">State: {{ conn.State }}</span>
              <span v-if="!short && conn.ThreadID">ThreadID: {{ conn.ThreadID }}</span>
              <span v-if="!short && conn.MyDID">MyDID: {{ conn.MyDID }}</span>
            </div>
            <md-button
              v-if="canAcceptExchangeRequest(conn)"
              class="md-icon-button md-dense md-raised md-info right"
              @click="acceptExchangeRequest(conn.ConnectionID)"
            >
              <md-icon>done</md-icon>
            </md-button>
          </md-list-item>
        </md-list>
      </md-content>
    </md-card-content>
  </md-card>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'Connections',
  props: {
    title: {
      type: String,
      default: 'Connections',
    },
    short: {
      type: Boolean,
      default: true,
    },
    count: {
      type: Number,
      default: 0,
    },
    connections: {
      type: Array,
    },
  },
  methods: {
    ...mapActions(['acceptExchangeRequest', 'queryConnections']),
    canAcceptExchangeRequest: function (conn) {
      return conn.State === 'requested' && conn.Namespace === 'their';
    },
  },
};
</script>

<style scoped>
.white {
  background-color: white;
}

.title {
  width: calc(100% - 42px);
  display: -webkit-inline-box;
}

.title-btn-right {
  right: 20px;
  position: absolute;
  display: -webkit-inline-box;
}

.md-list {
  width: 100%;
  display: inline-block;
  vertical-align: center;
}

.title {
  display: -webkit-inline-box;
}

.right {
  float: right;
}

.md-content-connections {
  width: 100%;
  max-height: 500px;
}
</style>
