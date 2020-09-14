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
    import {mapGetters, mapActions} from 'vuex'

    let vcData = [];
    async function fetchCredentials(aries) {
        // Get the VC data
        for (let i = 0; i < vcData.length; i++) {
            try {
                let resp = await aries.verifiable.getCredential({
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
        created: async function () {
            await this.initAries()
            // Load the Credentials
            this.aries = this.getAriesInstance()
            await this.getCredentials()

            window.$webCredentialHandler = this.$webCredentialHandler
            this.username = this.getCurrentUser().username
        },
        methods: {
            ...mapGetters('aries', {getAriesInstance: 'getInstance'}),
            ...mapGetters(['getCurrentUser']),
            ...mapActions('aries', {initAries: 'init'}),
            getCredentials: async function () {
                try {
                    let resp = await this.aries.verifiable.getCredentials()
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

                await fetchCredentials(this.aries)
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
                aries: null,
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


