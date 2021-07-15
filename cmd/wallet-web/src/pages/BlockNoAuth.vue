/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div style="margin: 5%">
        <md-card md-with-hover>
            <md-card-header data-background-color="red">
                <div class="md-title">Access Denied</div>
            </md-card-header>

            <md-card-content style="font-size: 16px">
                    You don't have authorization to access this page,
                please login to your wallet <a v-bind:href="loginURL" target="_blank" style="font-weight: 600; font-style: italic">here</a> and try again.
            </md-card-content>

        </md-card>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                loginURL: this.$route.params['loginURL'] ? this.$route.params['loginURL'] : '/',
            };
        }
    }

</script>

