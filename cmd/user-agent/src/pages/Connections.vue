/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <md-dialog :md-active.sync="dialog">
            <md-dialog-title>{{ dialogTitle }}</md-dialog-title>
            <div class="content">
                <div class="md-layout">
                    <div class="md-layout-item">
                        <pre>{{ dialogContent }}</pre>
                    </div>
                </div>
            </div>
            <md-dialog-actions>
                <md-button class="md-button md-info md-square" @click="dialog = false">Close</md-button>
            </md-dialog-actions>
        </md-dialog>
        <div class="md-layout">
            <div class="md-layout-item">
                <div class="md-layout-item">
                    <md-card class="md-card-plain">
                        <md-card-header data-background-color="green">
                            <h4 class="title">Create Invitation</h4>
                        </md-card-header>
                        <md-card-content style="background-color: white;">
                            <md-field>
                                <label>Alias</label>
                                <md-input v-model="alias" required></md-input>
                            </md-field>
                            <div style="display: flow-root">
                                <span class="error" v-if="createInvitationError">{{ createInvitationError }}</span>
                                <md-button class="md-button md-info md-square right"
                                           id='createInvitation'
                                           v-on:click="createInvitation">
                                    <b>Create</b>
                                </md-button>
                            </div>
                        </md-card-content>
                    </md-card>
                </div>
                <div class="md-layout-item">
                    <md-card class="md-card-plain">
                        <md-card-header data-background-color="green">
                            <h4 class="title">Receive Invitation</h4>
                        </md-card-header>
                        <md-card-content style="background-color: white;">
                            <md-field>
                                <label>JSON Invitation</label>
                                <md-textarea v-model="invitation"
                                             required></md-textarea>
                            </md-field>
                            <div style="display: flow-root">
                                <span class="error" v-if="receiveInvitationError">{{ receiveInvitationError }}</span>
                                <md-button class="md-button md-info md-square right"
                                           id='receiveInvitation'
                                           v-on:click="receiveInvitation">
                                    <b>Receive</b>
                                </md-button>
                            </div>
                        </md-card-content>
                    </md-card>
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
                        <div class="text-center" v-if="connections.length===0">No connections</div>
                        <md-content class="md-content-connections md-scrollbar">
                            <md-list class="md-triple-line">
                                <md-list-item v-for="connection in connections" :key="connection.id">
                                    <div class="md-list-item-text">
                                        <span>ConnectionID: {{connection.ConnectionID}}</span>
                                        <span>State: {{connection.State}}</span>
                                        <span>MyDID: {{connection.MyDID}}</span>
                                    </div>
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
    export default {
        beforeCreate: async function () {
            window.$aries = await this.$arieslib
            await this.queryConnections()
        },
        methods: {
            createInvitation: async function () {
                this.createInvitationError = ""
                if (this.alias.trim().length === 0) {
                    this.createInvitationError = "Please fill in the field!"
                    return
                }

                try {
                    let res = await window.$aries.didexchange.createInvitation({alias: this.alias.trim()})
                    this.showDialog(
                        "Invitation created!",
                        JSON.stringify(res, null, 2)
                    )
                } catch (e) {
                    this.createInvitationError = e.message
                }
            },
            receiveInvitation: async function () {
                this.receiveInvitationError = ""
                if (this.invitation.trim().length === 0) {
                    this.receiveInvitationError = "Please fill in the field!"
                    return
                }

                let invitation;
                try {
                    invitation = JSON.parse(this.invitation.trim())
                } catch (_) {
                    this.receiveInvitationError = "Please make sure you are providing a JSON invitation"
                    return
                }

                try {
                    let res = await window.$aries.didexchange.receiveInvitation(invitation)
                    this.showDialog(
                        "Invitation received!",
                        `Your connection ID is ${res['connection_id']}`
                    )
                } catch (e) {
                    this.receiveInvitationError = e.message
                }
            },
            queryConnections: async function () {
                try {
                    let res = await window.$aries.didexchange.queryConnections()
                    if (res.results) {
                        this.connections = res.results
                    }
                } catch (e) {
                    console.error(e.message)
                }
            },
            showDialog(title, content) {
                this.dialogTitle = title
                this.dialogContent = content
                this.dialog = true
            },
        },
        data() {
            return {
                dialog: false,
                dialogTitle: "",
                dialogContent: "",
                invitation: "",
                alias: "",
                connections: [],
                createInvitationError: "",
                receiveInvitationError: "",
            };
        },
    }
</script>

<style>
    .refresh-connections {
        top: -28px;
    }

    .right {
        float: right;
    }

    .error {
        color: red;
    }

    .md-content-connections {
        width: 100%;
        max-height: 500px;
    }
</style>