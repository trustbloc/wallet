<template>
    <div class="content">
        <md-dialog :md-active.sync="dialog">
            <md-dialog-title>{{ dialogTitle }}</md-dialog-title>
            <div class="content">
                <div class="md-layout">
                    <div class="md-layout-item">
                        <div v-if="dialogContent === 'acceptRequestPresentationForm'">
                            <md-field>
                                <label>JSON Presentation</label>
                                <md-textarea v-model="presentation"
                                             required></md-textarea>
                            </md-field>
                            <span class="error" v-if="acceptRequestPresentationError">{{ acceptRequestPresentationError }}</span>
                        </div>
                        <div v-else-if="dialogContent === 'acceptPresentationForm'">
                            <md-field>
                                <label>Presentation names coma separated e.g name1,name2</label>
                                <md-textarea v-model="presentationNames"
                                             required></md-textarea>
                            </md-field>
                            <span class="error" v-if="acceptPresentationError">{{ acceptPresentationError }}</span>
                        </div>
                        <pre v-else>{{ dialogContent }}</pre>
                    </div>
                </div>
            </div>
            <md-dialog-actions v-if="dialogContent === 'acceptRequestPresentationForm'">
                <md-button class="md-button md-info md-square"
                           v-on:click="acceptRequestPresentation(selectedAction,true)">Send
                </md-button>
            </md-dialog-actions>
            <md-dialog-actions v-else-if="dialogContent === 'acceptPresentationForm'">
                <md-button class="md-button md-info md-square" v-on:click="acceptPresentation(selectedAction,true)">Send
                </md-button>
            </md-dialog-actions>
            <md-dialog-actions v-else>
                <md-button class="md-button md-info md-square" @click="dialog = false">Close</md-button>
            </md-dialog-actions>
        </md-dialog>

        <div class="md-layout">
            <div class="md-layout-item">
                <div class="md-layout-item">
                    <md-card class="md-card-plain">
                        <md-card-header data-background-color="green">
                            <h4 class="title">Send Request Presentation</h4>
                        </md-card-header>
                        <md-card-content style="background-color: white;">
                            <md-field>
                                <label>JSON RequestPresentation</label>
                                <md-textarea v-model="requestPresentation"
                                             required></md-textarea>
                            </md-field>
                            <div class="md-layout md-gutter">
                                <div class="md-layout-item">
                                    <md-field>
                                        <md-select v-model="requestPresentationConnection"
                                                   md-class="offer-credential-connection"
                                                   :disabled="connections.length === 0">
                                            <md-option v-for="conn in connections" :key="conn.id"
                                                       :value="conn.ConnectionID">
                                                {{conn.TheirLabel}}
                                            </md-option>
                                        </md-select>
                                    </md-field>
                                </div>
                            </div>
                            <div style="display: flow-root">
                                <span class="error"
                                      v-if="sendRequestPresentationError">{{ sendRequestPresentationError }}</span>
                                <span class="success"
                                      v-if="sendRequestPresentationSuccess">{{ sendRequestPresentationSuccess }}</span>
                                <md-button :disabled="connections.length === 0"
                                           class="md-button md-info md-square right"
                                           id='receiveInvitation'
                                           v-on:click="sendRequestPresentation">
                                    <b>Send</b>
                                </md-button>
                            </div>
                        </md-card-content>
                    </md-card>
                </div>
            </div>
            <div class="md-layout-item">
                <md-card class="md-card-plain">
                    <md-card-header data-background-color="green">
                        <h4 class="title">Actions</h4>
                        <md-button v-on:click="refreshActions"
                                   class="md-icon-button md-dense md-raised md-info right refresh-connections">
                            <md-icon>cached</md-icon>
                        </md-button>
                    </md-card-header>
                    <md-card-content style="background-color: white;">
                        <div class="text-center" v-if="actions.length===0">No actions</div>
                        <md-content class="md-content-actions md-scrollbar">
                            <md-list class="md-triple-line">
                                <md-list-item v-for="action in actions" :key="action.id">
                                    <div class="md-list-item-text">
                                        <span>PIID: {{action.piid}}</span>
                                    </div>
                                    <md-button v-if="isRequestPresentation(action)"
                                               v-on:click="acceptRequestPresentation(action)"
                                               class="md-icon-button md-dense md-raised md-info right">
                                        <md-icon>done</md-icon>
                                    </md-button>
                                    <md-button v-if="isRequestPresentation(action)"
                                               v-on:click="declineRequestPresentation(action)"
                                               class="md-icon-button md-dense md-raised md-danger right">
                                        <md-icon>close</md-icon>
                                    </md-button>

                                    <md-button v-if="isPresentation(action)"
                                               v-on:click="acceptPresentation(action)"
                                               class="md-icon-button md-dense md-raised md-info right">
                                        <md-icon>done</md-icon>
                                    </md-button>
                                    <md-button v-if="isPresentation(action)"
                                               v-on:click="declinePresentation(action)"
                                               class="md-icon-button md-dense md-raised md-danger right">
                                        <md-icon>close</md-icon>
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
    export default {
        beforeCreate: async function () {
            window.$aries = await this.$arieslib
            await this.queryConnections()
            await this.refreshActions()
        },
        methods: {
            refreshActions: async function () {
                let res = await window.$aries.presentproof.actions()
                this.actions = res.actions
                console.log(this.actions)
            },
            isRequestPresentation: function (action) {
                return action.msg['@type'].endsWith('/request-presentation')
            },
            isPresentation: function (action) {
                return action.msg['@type'].endsWith('/presentation')
            },
            acceptPresentation: async function (action, form) {
                if (!form) {
                    this.selectedAction = action
                    this.showDialog("Accept Presentation", "acceptPresentationForm")
                    return
                }

                this.acceptCredentialError = ""
                if (this.presentationNames.trim().length === 0) {
                    this.acceptCredentialError = "Please fill in the field!"
                    return
                }

                try {
                    await window.$aries.presentproof.acceptPresentation({
                        piid: action.piid,
                        names: this.presentationNames.split(','),
                    })
                } catch (e) {
                    this.showDialog("Accept Presentation", e.message)
                    return
                }

                this.showDialog("Accept Presentation", "Accepted!")
                await this.refreshActions()
            },
            acceptRequestPresentation: async function (action, form) {
                if (!form) {
                    this.selectedAction = action
                    this.showDialog("Accept Request Presentation", "acceptRequestPresentationForm")
                    return
                }

                this.acceptRequestPresentationError = ""
                if (this.presentation.trim().length === 0) {
                    this.acceptRequestPresentationError = "Please fill in the field!"
                    return
                }

                let presentation;
                try {
                    presentation = JSON.parse(this.presentation.trim())
                } catch (_) {
                    this.acceptRequestPresentationError = "Please make sure you are providing a JSON Presentation"
                    return
                }

                try {
                    await window.$aries.presentproof.acceptRequestPresentation({
                        piid: action.piid,
                        presentation: presentation,
                    })
                } catch (e) {
                    this.showDialog("Accept Request Presentation", e.message)
                    return
                }

                this.showDialog("Accept Request Presentation", "Accepted!")
                await this.refreshActions()

            },
            declineRequestPresentation: async function (action) {
                try {
                    await window.$aries.presentproof.declineRequestPresentation({
                        piid: action.piid,
                    })
                } catch (e) {
                    this.showDialog("Decline Request Presentation", e.message)
                    return
                }

                this.showDialog("Decline Request Presentation", "Declined!")
                await this.refreshActions()
            },
            declinePresentation: async function (action) {
                try {
                    await window.$aries.presentproof.declinePresentation({
                        piid: action.piid,
                    })
                } catch (e) {
                    this.showDialog("Decline Presentation", e.message)
                    return
                }

                this.showDialog("Decline Presentation", "Declined!")
                await this.refreshActions()
            },
            sendRequestPresentation: async function () {
                this.sendRequestPresentationError = ""
                if (this.requestPresentation.trim().length === 0) {
                    this.sendRequestPresentationError = "Please provide the JSON RequestPresentation!"
                    return
                }

                if (!this.requestPresentationConnection) {
                    this.sendRequestPresentationError = "Please select a connection."
                    return
                }

                let connID = this.requestPresentationConnection
                let conn = this.connections.find(function (conn) {
                    return conn.ConnectionID === connID
                })

                let requestPresentation;
                try {
                    requestPresentation = JSON.parse(this.requestPresentation.trim())
                } catch (_) {
                    this.sendRequestPresentationError = "Please make sure you are providing a JSON RequestPresentation"
                    return
                }

                try {
                    await window.$aries.presentproof.sendRequestPresentation({
                        my_did: conn.MyDID,
                        their_did: conn.TheirDID,
                        request_presentation: requestPresentation,
                    })
                    this.sendRequestPresentationSuccess = `Your offer was sent successfully!`
                } catch (e) {
                    this.sendRequestPresentationError = e.message
                }
            },
            queryConnections: async function () {
                try {
                    let res = await window.$aries.didexchange.queryConnections()
                    if (res.results) {
                        this.connections = res.results.filter(function (conn) {
                            return conn.State === "completed";
                        });
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
                actions: [],
                connections: [],
                selectedAction: {},
                requestPresentation: "",
                requestPresentationConnection: "",
                presentation: "",
                acceptRequestPresentationError: "",
                sendRequestPresentationError: "",
                sendRequestPresentationSuccess: "",
                presentationNames: "",
                acceptPresentationError: "",
            };
        },
    }
</script>

<style scoped>
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
</style>