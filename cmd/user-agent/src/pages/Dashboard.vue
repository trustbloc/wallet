/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="content">
            <div class="md-layout">
                <div>
                    <md-toolbar class="md-accent" md-elevation="1">
                        <i class="md-title" style="flex: 1; font-size: 20px">Welcome &nbsp;{{ username}}</i>
                        <logout/>
                    </md-toolbar>

                    <md-card md-with-hover>
                        <md-card-header data-background-color="green">
                            <h4 class="title">
                                <md-icon>content_paste</md-icon>
                                Your Stored Credentials
                            </h4>
                        </md-card-header>
                        <md-card-content>
                            <simple-table v-for="vc in verifiableCredential"
                                          :key=vc.name
                                          :name=credDisplayName(vc.type)
                                          :data=vc.credential>
                            </simple-table>
                        </md-card-content>
                    </md-card>
                </div>
            </div>
        </div>
    </div>

</template>

<script>
    import {SimpleTable} from "@/components";
    import {getCredentialType} from "@/pages/chapi/wallet";
    import Logout from "@/pages/chapi/Logout.vue";

    let vcData = [];
    async function fetchCredentials() {
        // Get the VC data
        for (let i = 0; i < vcData.length; i++) {
            try {
                let resp = await window.$aries.verifiable.getCredential({
                    id: vcData[i].id
                })
                vcData[i].credential = JSON.parse(resp.verifiableCredential)
            } catch (e) {
                console.error('get vc failed : errMsg=' + e)
            }
        }
    }

    export default {
        components: {
            Logout,
            SimpleTable,
        },
        beforeCreate: async function () {
            // Load the Credentials
            let aries = await this.$arieslib
            window.$webCredentialHandler = this.$webCredentialHandler
            window.$aries = aries
            await this.getCredentials(aries)
            this.username = this.$store.getters.getCurrentUser.username
        },
        methods: {
            getCredentials: async function (aries) {
                try {
                    let resp = await aries.verifiable.getCredentials()
                    if (!resp.result) {
                        console.log('no credentials exists')
                        return
                    }

                    vcData = resp.result
                    if (resp.result.length === 0) {
                        console.log('no credentials exists')
                    }

                } catch (e) {
                    console.error('get credentials failed : errMsg=' + e)
                }

                await fetchCredentials()
                this.verifiableCredential = vcData
            },
            credDisplayName: function (types) {
                return getCredentialType(types)
            }
        },
        data() {
            return {
                verifiableCredential: [],
                username: '',
            }
        }
    }
</script>

<style>
    .title {
        text-transform: capitalize;
    }

    .md-content {
        overflow: auto;
        padding: 1px;
        font-size: 6px;
        line-height: 16px;
    }

</style>


