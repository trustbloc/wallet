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
                            <h4 class="title">Router</h4>
                        </md-card-header>
                        <md-card-content v-if="!routerConnID" style="background-color: white;">
                            <md-field>
                                <label>Router URL</label>
                                <md-input placeholder="http://router.example.com" v-model="routerURL"
                                          required></md-input>
                            </md-field>
                            <div style="display: flow-root">
                                <span class="error" v-if="routerRegisterError">{{ routerRegisterError }}</span>
                                <md-button v-bind:disabled="routerDisabledButton"
                                           class="md-button md-info md-square right"
                                           id='routerRegister'
                                           v-on:click="routerRegister">
                                    <b>Register</b>
                                </md-button>
                            </div>
                        </md-card-content>
                        <md-card-content v-if="routerConnID" style="background-color: white;">
                            <div class="md-layout router">
                                <div class="md-layout-item md-layout">
                                    <div class="md-layout-item router-done">
                                        <div>Router is registered</div>
                                    </div>
                                    <div class="md-layout-item">
                                        <md-button class="md-button md-danger md-square right"
                                                   id='routerUnregister'
                                                   v-on:click="routerUnregister">
                                            <b>Unregister</b>
                                        </md-button>
                                    </div>
                                </div>
                            </div>
                        </md-card-content>
                    </md-card>
                </div>
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
                                <md-list-item v-for="conn in connections" :key="conn.id">
                                    <div class="md-list-item-text">
                                        <span>ConnectionID: {{conn.ConnectionID}}</span>
                                        <span>State: {{conn.State}}</span>
                                        <span>MyDID: {{conn.MyDID}}</span>
                                    </div>
                                    <md-button v-if="canAcceptInvitation(conn)"
                                               v-on:click="acceptInvitation(conn.ConnectionID)"
                                               class="md-icon-button md-dense md-raised md-info right">
                                        <md-icon>done</md-icon>
                                    </md-button>
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
    import axios from "axios";

    export default {
        beforeCreate: async function () {
            window.$aries = await this.$arieslib
            await this.queryConnections()
            await this.isRouterRegistered()
        },
        methods: {
            waitFor: function (connection, state, callback) {
                return new Promise((resolve, reject) => {
                    const stop = window.$aries.startNotifier(notice => {
                        if (connection !== notice.payload.connection_id) {
                            return
                        }

                        if (notice.payload.state !== state) {
                            return
                        }

                        stop()

                        try {
                            callback().then(() => {
                                resolve()
                            })
                        } catch (err) {
                            reject(err)
                        }
                    }, ["all"])

                    setTimeout(() => {
                        stop()
                        reject(new Error("time out while waiting for connection"))
                    }, 5000)
                })
            },
            isRouterRegistered: async function () {
                let res = await window.$aries.mediator.getConnection().catch(err => {
                    if (!err.message.includes("router not registered")) {
                        throw err
                    }
                })

                if (res) {
                    this.routerConnID = res.connectionID
                }

                return this.routerConnID
            },
            routerRegister: async function () {
                this.routerRegisterError = ""
                let routerURL = this.routerURL.trim().replace(/\/$/, "");
                if (routerURL.length === 0) {
                    this.routerRegisterError = "Please fill in the field!"
                    return
                }

                this.routerDisabledButton = true

                try {
                    let invitation = await axios.post(routerURL + "/connections/create-invitation")
                    let connection = await window.$aries.didexchange.receiveInvitation(invitation.data.invitation)

                    await this.waitFor(connection.connection_id, 'invited', function () {
                        return window.$aries.didexchange.acceptInvitation({
                            id: connection.connection_id
                        })
                    })

                    await this.waitFor(connection.connection_id, 'completed', function () {
                        return window.$aries.mediator.register({"connectionID": connection.connection_id})
                    })
                } catch (e) {
                    this.routerRegisterError = e.message
                }

                this.routerDisabledButton = false
                if (await this.isRouterRegistered()) {
                    await this.queryConnections()
                }
            },
            routerUnregister: async function () {
                await window.$aries.mediator.unregister({id: this.routerConnID})
                this.routerConnID = ""
            },
            canAcceptInvitation: function (conn) {
                return conn.State === 'invited'
            },
            canAcceptExchangeRequest: function (conn) {
                return conn.State === 'requested' && conn.Namespace === 'their'
            },
            acceptInvitation: async function (id) {
                await window.$aries.didexchange.acceptInvitation({
                    id: id
                }).then(this.queryConnections)
            },
            acceptExchangeRequest: async function (id) {
                await window.$aries.didexchange.acceptExchangeRequest({
                    id: id
                }).then(this.queryConnections)
            },
            createInvitation: async function () {
                this.createInvitationError = ""
                if (!this.routerConnID) {
                    this.createInvitationError = "Please register a router first!"
                    return
                }

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
                    await this.queryConnections()
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
                routerURL: "",
                connections: [],
                routerConnID: "",
                routerDisabledButton: false,
                createInvitationError: "",
                receiveInvitationError: "",
                routerRegisterError: "",
            };
        },
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

    .router-done {
        display: flex;
        align-items: center;
    }

    .router .md-layout-item {
        padding: 0;
    }

    .md-content-connections {
        width: 100%;
        max-height: 500px;
    }
</style>