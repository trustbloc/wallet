/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="content">
    <mediator v-if="!isMediatorRegistered" title="Please, set up a mediator to proceed with this page!"></mediator>
    <div class="md-layout" v-if="isMediatorRegistered">
      <div class="md-layout-item">
        <div class="md-layout-item">
          <public-invitation></public-invitation>
        </div>
        <div class="md-layout-item">
          <receive-invitation title="Have an invitation? Put it below and hit the Connect button."
                              type="base64"></receive-invitation>
        </div>
      </div>
      <div class="md-layout-item">
        <div class="md-layout-item" v-if="pendingConnectionsCount">
          <md-card class="md-card-plain">
            <md-card-header data-background-color="green">
              <h4 class="title">Pending requests</h4>
              <div style="display: -webkit-inline-box; position: absolute">
                <md-badge class="md-primary md-square" style="float: left;position: relative;margin-top: 5px;"
                          :md-content="pendingConnectionsCount"/>
              </div>
              <md-button v-on:click="queryConnections"
                         class="md-icon-button md-dense md-raised md-info right refresh-connections">
                <md-icon>cached</md-icon>
              </md-button>
            </md-card-header>
            <md-card-content style="background-color: white;">
              <md-content class="md-content-connections md-scrollbar">
                <md-list class="md-triple-line">
                  <md-list-item v-for="conn in pendingConnections" :key="conn.id">
                    <div class="md-list-item-text">
                      <span v-if="conn.TheirDID">{{ conn.TheirDID }}</span>
                    </div>
                    <md-button v-on:click="acceptExchangeRequest(conn.ConnectionID)"
                               class="md-icon-button md-dense md-raised md-info right">
                      <md-icon>done</md-icon>
                    </md-button>
                  </md-list-item>
                </md-list>
              </md-content>
            </md-card-content>
          </md-card>
        </div>
        <div class="md-layout-item" v-if="completedConnectionsCount">
          <md-card class="md-card-plain">
            <md-card-header data-background-color="green">
              <h4 class="title">List of agents you have a connection with</h4>
              <md-button v-on:click="queryConnections"
                         class="md-icon-button md-dense md-raised md-info right refresh-connections">
                <md-icon>cached</md-icon>
              </md-button>
            </md-card-header>
            <md-card-content style="background-color: white;">
              <md-content class="md-content-connections md-scrollbar">
                <md-list class="md-triple-line">
                  <md-list-item v-for="conn in completedConnections" :key="conn.id">
                    <div class="md-list-item-text">
                      <span v-if="conn.TheirLabel">{{ conn.TheirLabel }}</span>
                    </div>
                  </md-list-item>
                </md-list>
              </md-content>
            </md-card-content>
          </md-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {mapGetters, mapActions} from 'vuex'
import {Mediator, PublicInvitation, ReceiveInvitation} from "@/components";

export default {
  components: {Mediator, PublicInvitation, ReceiveInvitation},
  methods: mapActions(['queryConnections', 'acceptExchangeRequest']),
  computed: mapGetters([
    'pendingConnections', 'completedConnections',
    'pendingConnectionsCount', 'completedConnectionsCount',
    'isMediatorRegistered'
  ]),
  mounted() {
    // refreshes connections when component is mounted
    this.queryConnections()
  }
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