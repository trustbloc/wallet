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
    import DIDConnectForm from "./DIDConnect.vue";
    import GetCredentialsForm from "./GetCredentials.vue";
    import PresentationDefQueryForm from "./PresentationDefQuery.vue";
    import jp from 'jsonpath';

    const types = [
        {id:'PresentationDefinitionQuery', component:PresentationDefQueryForm},
        {id:'DIDAuth', component:DIDAuthForm},
        {id:'DIDConnect', component:DIDConnectForm}
    ]

    function getComponent(credEvent){
        for (let i in types) {
            let type = types[i]

            let found = jp.query(credEvent, `$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="${type.id}")]`);
            if (found.length  > 0){
                return type.component
            }

            found = jp.query(credEvent, `$..credentialRequestOptions.web.VerifiablePresentation.query`);
            if (found.length  > 0 && found[0].type == type.id){
                return type.component
            }
        }

        // default form
        return GetCredentialsForm
    }

    export default {
        beforeCreate: async function () {
            this.credentialEvent = await this.$webCredentialHandler.receiveCredentialEvent();
            if (!this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation) {
                console.log("invalid web credential type, expected 'VerifiablePresentation'")
                return
            }

            console.log("incoming web credential event", this.credentialEvent)
            this.component = getComponent(this.credentialEvent)
        },
        data() {
            return {
                component: ""
            };
        }
    }
</script>

<style lang="css">
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

    .viewport {
        font-size: 16px !important;
    }

    .center-span {
        text-align: center !important;
    }

</style>
