/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div style="display: flex;">
        <md-ripple>
            <md-button :md-ripple="false" v-on:click="registerDevice" class="md-primary"
                       style="background-color: #4f424c !important; font-size: 16px; " v-if="register !== 'success'">Register Device
            </md-button>
        </md-ripple>
        <div v-if="register === 'success'">
            <i class="md-title" style="font-size: 20px; color: white !important;">
                <md-icon style="color: green; margin: 10px;">check_circle</md-icon>
                Device successfully registered. </i>
        </div>
        <div v-if="register === 'failure'">
            <i class="md-title" style="font-size: 20px; color: white !important;">
                <md-icon style="color: red; margin: 10px;">report</md-icon>
                Device failed to register. </i>
        </div>
    </div>
</template>

<script>
    import {DeviceRegister} from "./wallet"
    import {mapActions, mapGetters} from 'vuex'
    export default {
        created: async function () {
            this.deviceRegister = new DeviceRegister(this.getAgentOpts())
        },
        data() {
            return {
                register: 'none',
            }
        },
        methods: {
            ...mapActions({logoutUser: 'logout'}),
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            ...mapGetters(['getAgentOpts']),
            registerDevice: async function () {
                this.register = await this.deviceRegister.register();
            },
        },
    }
</script>
