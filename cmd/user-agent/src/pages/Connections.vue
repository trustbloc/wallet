/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="content">
    <div class="md-layout">
      <div class="md-layout-item">
        <div class="md-layout-item">
          <mediator title="Mediator"></mediator>
        </div>
        <div class="md-layout-item">
          <create-invitation></create-invitation>
        </div>
        <div class="md-layout-item">
          <receive-invitation></receive-invitation>
        </div>
      </div>
      <div class="md-layout-item">
        <md-card class="md-card-plain">
          <md-card-header data-background-color="green">
            <h4 class="title">Query Connections</h4>
            <md-button v-on:click="queryConnections"
                       class="md-icon-button md-dense md-raised md-info right refresh-connections">
              <md-icon>cached</md-icon>
            </md-button>
          </md-card-header>
          <md-card-content style="background-color: white;">
            <div class="text-center" v-if="allConnections.length===0">No connections</div>
            <md-content class="md-content-connections md-scrollbar">
              <md-list class="md-triple-line">
                <md-list-item v-for="conn in allConnections" :key="conn.id">
                  <div class="md-list-item-text">
                    <span v-if="conn.ConnectionID">ConnectionID: {{ conn.ConnectionID }}</span>
                    <span v-if="conn.State">State: {{ conn.State }}</span>
                    <span v-if="conn.ThreadID">ThreadID: {{ conn.ThreadID }}</span>
                    <span v-if="conn.MyDID">MyDID: {{ conn.MyDID }}</span>
                  </div>
                  <md-button v-if="canAcceptExchangeRequest(conn)"
                             v-on:click="acceptExchangeRequest(conn.ConnectionID)"
                             class="md-icon-button md-dense md-raised md-info right">
                    <md-icon>done</md-icon>
                  </md-button>
                </md-list-item>
              </md-list>
            </md-content>
          </md-card-content>
        </md-card>
      </div>
    </div>
  </div>
</template>

<script>
import {Mediator, CreateInvitation, ReceiveInvitation} from "@/components";
import {mapActions, mapGetters} from "vuex";

export default {
  components: {Mediator, CreateInvitation, ReceiveInvitation},
  methods: {
    ...mapActions(['queryConnections', 'acceptExchangeRequest', 'createInvitation']),
    canAcceptExchangeRequest: function (conn) {
      return conn.State === 'requested' && conn.Namespace === 'their'
    },
  },
  mounted() {
    // refreshes connections when component is mounted
    this.queryConnections()
  },
  computed: mapGetters(['isMediatorRegistered', 'allConnections']),
}
</script>

<style scoped>
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