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

    export default {
        beforeCreate: async function () {
            this.registrar = new RegisterWallet(this.$polyfill, this.$webCredentialHandler, await this.$arieslib,
                this.$trustblocAgent, await this.$trustblocStartupOpts)
        },
        methods: {
            logout: async function () {
                await this.registrar.unregister()
                this.$store.dispatch('resetUser')
                this.$router.push("/login");
            }
        }
    }

</script>

