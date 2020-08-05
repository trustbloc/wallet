/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="content">
    <div class="md-layout">
      <div class="md-layout-item">
        <div class="md-layout-item">
          <md-card class="md-card-plain">
            <md-card-header data-background-color="green">
              <h4 class="title">Make New Friends By Sharing This Invitation!</h4>
              <md-button v-on:click="copyInvitationToClipboard"
                         class="md-icon-button md-dense md-raised md-info right refresh-connections">
                <md-icon>content_copy</md-icon>
              </md-button>
            </md-card-header>
            <md-card-content style="background-color: white;">
              <div style="display: flow-root;text-align: center;">
                <div style="line-break: anywhere;">{{ generatedInvitation }}</div>
              </div>
              <input type="hidden" id="created-invitation" :value="generatedInvitation">
            </md-card-content>
          </md-card>
        </div>
        <div class="md-layout-item">
          <md-card class="md-card-plain">
            <md-card-header data-background-color="green">
              <h4 class="title">Have an invitation? Put it below and hit the Connect button.</h4>
            </md-card-header>
            <md-card-content style="background-color: white;">
              <md-field>
                <label>Invitation</label>
                <md-textarea v-model="inboundInvitation"
                             required></md-textarea>
              </md-field>
              <div style="display: flow-root">
                <span class="error" v-if="inboundInvitationError">{{ inboundInvitationError }}</span>
                <span class="success" v-if="inboundInvitationSuccess">{{ inboundInvitationSuccess }}</span>
                <md-button class="md-button md-info md-square right"
                           id='receiveInvitation'
                           v-on:click="receiveInvitation">
                  <b>Connect</b>
                </md-button>
              </div>
            </md-card-content>
          </md-card>
        </div>
      </div>
      <div class="md-layout-item">
        <div class="md-layout-item" v-if="pendingRequests().length!==0">
          <md-card class="md-card-plain">
            <md-card-header data-background-color="green">
              <h4 class="title">Pending requests</h4>
              <div style="display: -webkit-inline-box;position: absolute">
                <md-badge class="md-primary md-square" style="float: left;position: relative;margin-top: 5px;"
                          :md-content="pendingRequests().length"/>
              </div>
              <md-button v-on:click="queryConnections"
                         class="md-icon-button md-dense md-raised md-info right refresh-connections">
                <md-icon>cached</md-icon>
              </md-button>
            </md-card-header>
            <md-card-content style="background-color: white;">
              <md-content class="md-content-connections md-scrollbar">
                <md-list class="md-triple-line">
                  <md-list-item v-for="conn in pendingRequests()" :key="conn.id">
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
        <div class="md-layout-item" v-if="completedConnections().length!==0">
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
                  <md-list-item v-for="conn in completedConnections()" :key="conn.id">
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
export default {
  methods: {
    pendingRequests: function () {
      return this.connections.filter(conn => conn.State === 'requested' && conn.Namespace === 'their')
    },
    completedConnections: function () {
      return this.connections.filter(conn => conn.State === 'completed')
    },
    acceptExchangeRequest: async function (id) {
      await window.$aries.didexchange.acceptExchangeRequest({
        id: id
      }).then(this.queryConnections)
    },
    createInvitation: async function () {
      // creates invitation through the out-of-band protocol
      let res = await window.$aries.outofband.createInvitation({label: "User"})
      // encodes invitation to base64 string
      return window.btoa(JSON.stringify(res.invitation));
    },
    copyInvitationToClipboard: function () {
      let inv = document.querySelector('#created-invitation')
      inv.setAttribute('type', 'text')
      inv.select()

      document.execCommand('copy');

      inv.setAttribute('type', 'hidden')
      window.getSelection().removeAllRanges()
    },
    receiveInvitation: async function () {
      let inv = this.inboundInvitation;
      // clears user-friendly messages
      this.inboundInvitationError = ""
      this.inboundInvitationSuccess = ""

      // checks whether the invitation was provided by the user
      if (inv.trim().length === 0) {
        this.inboundInvitationError = "Please fill in the field!"
        return
      }

      let invitation;
      try {
        // parses invitation that user has provided
        invitation = JSON.parse(window.atob(inv.trim()))
      } catch (_) {
        this.inboundInvitationError = "Please make sure you are providing a JSON invitation"
        return
      }

      try {
        // accepts invitation thought out-of-band protocol
        let res = await window.$aries.outofband.acceptInvitation({
          my_label: "agent",
          invitation: invitation,
        })
        // shows connection id that has been received
        this.inboundInvitationSuccess = `Your connection ID is ${res['connection_id']}`
        // refreshes connections
        await this.queryConnections()
        // clears provided invitation
        this.inboundInvitation = ""
      } catch (e) {
        this.inboundInvitationError = e.message
      }
    },
    queryConnections: async function () {
      // retrieves all agent connections
      let res = await window.$aries.didexchange.queryConnections()
      if (res.results) {
        // sets connections
        this.connections = res.results
      }
    },
  },
  data() {
    return {
      generatedInvitation: "",
      inboundInvitation: "",
      connections: [],
      inboundInvitationError: "",
      inboundInvitationSuccess: "",
    };
  },
  async mounted() {
    // sets aries instance globally
    window.$aries = await this.$arieslib
    // generates new invitation, it would be nice to use public invitation here
    this.generatedInvitation = await this.createInvitation()
    // load connections when component is mounted
    await this.queryConnections()

    let _this = this;
    // update connections when on update_relationships event
    this.$root.$on('update_relationships', async function () {
      await _this.queryConnections()
    })
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

.error {
  color: red;
}

.success {
  color: green;
}

.md-content-connections {
  width: 100%;
  max-height: 500px;
}
</style>