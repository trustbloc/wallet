/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <component v-bind:is="component"></component>
    </div>
</template>
<script>

    import DIDAuthForm from "./DIDAuth.vue";
    import GetCredentialsForm from "./GetCredentials.vue";

    export default {
        beforeCreate: async function () {
            this.$polyfill.loadOnce()
            this.credentialEvent = await this.$webCredentialHandler.receiveCredentialEvent();
            if (!this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation) {
                console.log("invalid web credential type, expected 'VerifiablePresentation'")
                return
            }

            console.log("incoming web credential event", this.credentialEvent)

            this.verifiable = this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation
            this.query = {}
            if (this.verifiable.query) {
                // supporting only one query for now
                this.query = Array.isArray(this.verifiable.query) ? this.verifiable.query[0] : this.verifiable.query;
            }

            if (this.query.type === 'DIDAuth') {
                this.component = DIDAuthForm
            } else {
                this.component = GetCredentialsForm
            }
        },
        data() {
            return {
                component: ""
            };
        }
    }
</script>

<style lang="scss">
    .md-menu-content {
        width: auto;
        max-width: 100% !important;
    }

    .md-menu-content .md-ripple > span {
        position: relative !important;
    }

    .md-table-head {
        background-color: #00bcd4 !important;
        color: white !important;
        font-size: 10px !important;
        font-weight: 500 !important;
    }

    .md-table-head-label .md-sortable {
        padding-left: 10px !important;
    }

    .md-empty-state-icon {
        width: 60px !important;
        min-width: 60px !important;
        height: 50px !important;
        font-size: 50px !important;
        margin: 0;
    }

    .md-empty-state-label {
        font-size: 18px !important;
        font-weight: 500 !important;
        line-height: 20px !important;
    }

</style>