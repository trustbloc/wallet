/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <md-ripple>
            <md-button v-on:click="logout" class="md-primary"
                       style="background-color: #4f424c !important; font-size: 16px">Logout
            </md-button>
        </md-ripple>
    </div>
</template>

<script>

    import {RegisterWallet} from "./wallet"
    import {mapActions, mapGetters} from 'vuex'

    export default {
        created: async function () {
            this.registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, this.getAriesInstance(),
                this.$trustblocAgent, this.getTrustblocOpts())
        },
        methods: {
            ...mapActions({logoutUser: 'logout'}),
            ...mapGetters('aries', {getAriesInstance: 'getInstance'}),
            ...mapGetters(['getTrustblocOpts']),
            logout: async function () {
                await this.registrar.uninstallHandlers()
                await this.logoutUser()
                this.$router.push("/login");
            }
        }
    }

</script>

