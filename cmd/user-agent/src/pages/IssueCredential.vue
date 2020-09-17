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
                        <div v-if="dialogContent === 'acceptRequestForm'">
                            <md-field>
                                <label>JSON IssueCredential</label>
                                <md-textarea v-model="issueCredential"
                                             required></md-textarea>
                            </md-field>
                            <span class="error" v-if="acceptRequestError">{{ acceptRequestError }}</span>
                        </div>
                        <div v-else-if="dialogContent === 'acceptCredentialForm'">
                            <md-field>
                                <label>Credential names coma separated e.g name1,name2</label>
                                <md-textarea v-model="issueCredentialNames"
                                             required></md-textarea>
                            </md-field>
                            <span class="error" v-if="acceptCredentialError">{{ acceptCredentialError }}</span>
                        </div>
                        <pre v-else>{{ dialogContent }}</pre>
                    </div>
                </div>
            </div>
            <md-dialog-actions v-if="dialogContent === 'acceptRequestForm'">
                <md-button class="md-button md-info md-square" v-on:click="acceptRequest(selectedAction,true)">Send
                </md-button>
            </md-dialog-actions>
            <md-dialog-actions v-else-if="dialogContent === 'acceptCredentialForm'">
                <md-button class="md-button md-info md-square" v-on:click="acceptCredential(selectedAction,true)">Send
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
                            <h4 class="title">Send Offer</h4>
                        </md-card-header>
                        <md-card-content style="background-color: white;">
                            <md-field>
                                <label>JSON OfferCredential</label>
                                <md-textarea v-model="offerCredential"
                                             required></md-textarea>
                            </md-field>
                            <div class="md-layout md-gutter">
                                <div class="md-layout-item">
                                    <md-field>
                                        <md-select v-model="offerCredentialConnection"
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
                                <span class="error" v-if="sendOfferError">{{ sendOfferError }}</span>
                                <span class="success" v-if="sendOfferSuccess">{{ sendOfferSuccess }}</span>
                                <md-button :disabled="connections.length === 0"
                                           class="md-button md-info md-square right"
                                           id='receiveInvitation'
                                           v-on:click="sendOffer">
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
                                        <span>PIID: {{action.PIID}}</span>
                                    </div>
                                    <md-button v-if="isOfferCredential(action)" v-on:click="acceptOffer(action)"
                                               class="md-icon-button md-dense md-raised md-info right">
                                        <md-icon>done</md-icon>
                                    </md-button>
                                    <md-button v-if="isOfferCredential(action)" v-on:click="declineOffer(action)"
                                               class="md-icon-button md-dense md-raised md-danger right">
                                        <md-icon>close</md-icon>
                                    </md-button>
                                    <md-button v-if="isRequestCredential(action)" v-on:click="acceptRequest(action)"
                                               class="md-icon-button md-dense md-raised md-info right">
                                        <md-icon>done</md-icon>
                                    </md-button>
                                    <md-button v-if="isRequestCredential(action)" v-on:click="declineRequest(action)"
                                               class="md-icon-button md-dense md-raised md-danger right">
                                        <md-icon>close</md-icon>
                                    </md-button>

                                    <md-button v-if="isIssueCredential(action)" v-on:click="acceptCredential(action)"
                                               class="md-icon-button md-dense md-raised md-info right">
                                        <md-icon>done</md-icon>
                                    </md-button>
                                    <md-button v-if="isIssueCredential(action)" v-on:click="declineCredential(action)"
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
    import {mapGetters, mapActions} from 'vuex'

    export default {
        created: async function () {
            this.addAriesNotifiers({callback:this.onIssueCredentialState, topics:["issue-credential_states"]})
            this.aries = this.getAriesInstance()

            await this.queryConnections()
            await this.refreshActions()
        },
        methods: {
            ...mapGetters('aries', {getAriesInstance: 'getInstance'}),
            ...mapActions('aries', {addAriesNotifiers: 'addNotifier'}),
            ...mapActions(['onIssueCredentialState']),
            isOfferCredential: function (action) {
                return action.Msg['@type'].endsWith('/offer-credential')
            },
            isRequestCredential: function (action) {
                return action.Msg['@type'].endsWith('/request-credential')
            },
            isIssueCredential: function (action) {
                return action.Msg['@type'].endsWith('/issue-credential')
            },
            acceptCredential: async function (action, form) {
                if (!form) {
                    this.selectedAction = action
                    this.showDialog("Accept Credential", "acceptCredentialForm")
                    return
                }


                this.acceptCredentialError = ""
                if (this.issueCredentialNames.trim().length === 0) {
                    this.acceptCredentialError = "Please fill in the field!"
                    return
                }

                try {
                    await this.aries.issuecredential.acceptCredential({
                        piid: action.PIID,
                        names: this.issueCredentialNames.split(','),
                    })
                } catch (e) {
                    this.showDialog("Accept Credential", e.message)
                    return
                }

                this.showDialog("Accept Credential", "Accepted!")
                await this.refreshActions()
            },
            declineCredential: async function (action) {
                try {
                    await this.aries.issuecredential.declineCredential({
                        piid: action.PIID,
                    })
                } catch (e) {
                    this.showDialog("Decline Credential", e.message)
                    return
                }

                this.showDialog("Decline Credential", "Declined!")
                await this.refreshActions()
            },
            acceptRequest: async function (action, form) {
                if (!form) {
                    this.selectedAction = action
                    this.showDialog("Accept Request", "acceptRequestForm")
                    return
                }


                this.acceptRequestError = ""
                if (this.issueCredential.trim().length === 0) {
                    this.acceptRequestError = "Please fill in the field!"
                    return
                }

                let credential;
                try {
                    credential = JSON.parse(this.issueCredential.trim())
                } catch (_) {
                    this.acceptRequestError = "Please make sure you are providing a JSON Credential"
                    return
                }

                try {
                    await this.aries.issuecredential.acceptRequest({
                        piid: action.PIID,
                        issue_credential: credential,
                    })
                } catch (e) {
                    this.showDialog("Accept Request", e.message)
                    return
                }

                this.showDialog("Accept Request", "Accepted!")
                await this.refreshActions()
            },
            declineRequest: async function (action) {
                try {
                    await this.aries.issuecredential.declineRequest({
                        piid: action.PIID,
                    })
                } catch (e) {
                    this.showDialog("Decline Request", e.message)
                    return
                }

                this.showDialog("Decline Request", "Declined!")
                await this.refreshActions()
            },
            refreshActions: async function () {
                let res = await this.aries.issuecredential.actions()
                this.actions = res.actions
            },
            declineOffer: async function (action) {
                try {
                    await this.aries.issuecredential.declineOffer({
                        piid: action.PIID,
                    })
                } catch (e) {
                    this.showDialog("Decline Offer", e.message)
                    return
                }

                this.showDialog("Decline Offer", "Declined!")
                await this.refreshActions()
            },
            acceptOffer: async function (action) {
                try {
                    await this.aries.issuecredential.acceptOffer({
                        piid: action.PIID,
                    })
                } catch (e) {
                    this.showDialog("Accept Offer", e.message)
                    return
                }

                this.showDialog("Accept Offer", "Accepted!")
                await this.refreshActions()
            },
            sendOffer: async function () {
                this.sendOfferError = ""
                if (this.offerCredential.trim().length === 0) {
                    this.sendOfferError = "Please provide the JSON OfferCredential!"
                    return
                }

                if (!this.offerCredentialConnection) {
                    this.sendOfferError = "Please select a connection."
                    return
                }

                let connID = this.offerCredentialConnection
                let conn = this.connections.find(function (conn) {
                    return conn.ConnectionID === connID
                })

                let offerCredential;
                try {
                    offerCredential = JSON.parse(this.offerCredential.trim())
                } catch (_) {
                    this.sendOfferError = "Please make sure you are providing a JSON OfferCredential"
                    return
                }

                try {
                    await this.aries.issuecredential.sendOffer({
                        my_did: conn.MyDID,
                        their_did: conn.TheirDID,
                        offer_credential: offerCredential,
                    })
                    this.sendOfferSuccess = `Your offer was sent successfully!`
                } catch (e) {
                    this.sendOfferError = e.message
                }
            },
            queryConnections: async function () {
                try {
                    let res = await this.aries.didexchange.queryConnections()
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
                selectedAction: {},
                connections: [],
                actions: [],
                offerCredential: "",
                issueCredential: "",
                offerCredentialConnection: "",
                sendOfferError: "",
                acceptRequestError: "",
                acceptCredentialError: "",
                issueCredentialNames: "",
                sendOfferSuccess: "",
            };
        },
    }
</script>

<style>
    .offer-credential-connection {
        max-width: 346px !important;
    }

    .offer-credential-connection .md-list-item {
        margin: 0 5px !important;
    }

    .offer-credential-connection .md-ripple > span {
        left: initial !important;
        height: initial !important;
        position: absolute !important;
    }
</style>

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

