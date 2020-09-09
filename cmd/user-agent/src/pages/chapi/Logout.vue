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
    import {mapActions} from 'vuex'

    export default {
        beforeCreate: async function () {
            this.registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, await this.$arieslib,
                this.$trustblocAgent, await this.$trustblocStartupOpts)
        },
        methods: {
            ...mapActions({logoutUser: 'logout'}),
            logout: async function () {
                await this.registrar.uninstallHandlers()
                this.logoutUser()
                this.$router.push("/login");
            }
        }
    }

</script>

