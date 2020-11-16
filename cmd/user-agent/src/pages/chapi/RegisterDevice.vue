/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <md-ripple>
            <md-button :md-ripple="false" v-on:click="registerDevice" class="md-primary"
                       style="background-color: #4f424c !important; font-size: 16px; ">Register Device
            </md-button>
        </md-ripple>
    </div>
</template>

<script>
    import {DeviceRegister} from "./wallet"
    import {mapActions, mapGetters} from 'vuex'
    export default {
        created: async function () {
            this.deviceRegister = new DeviceRegister(this.getAgentOpts())
        },
        methods: {
            ...mapActions({logoutUser: 'logout'}),
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            ...mapGetters(['getAgentOpts']),
            registerDevice: async function () {
                await this.deviceRegister.register()
            }
        },
        mounted: function () {
            if (!window.PublicKeyCredential) {
                alert("Error: this browser does not support WebAuthn");
                return;
            }
        }
    }
</script>
