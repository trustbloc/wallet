/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div>
        <md-dialog :md-active.sync="showDialog" style="max-width: 75% !important;">
            <md-dialog-title
                    style="text-align: center; margin-top: 20px;">
                {{credDisplayName(item)}}
            </md-dialog-title>
            <md-content style="width:90%; margin:auto;">
                <vue-json-pretty
                        :data="item">
                </vue-json-pretty>
            </md-content>

            <md-dialog-actions style="float: right; padding: 20px;">
                <md-button
                        class="md-primary"
                        @click="showDialog = false"
                >
                    Close
                </md-button>
            </md-dialog-actions>
        </md-dialog>
        <span  @click.stop="showDialog = true" class="infoIcon"><md-icon>info</md-icon></span>
    </div>
</template>

<script>
    import { getCredentialType} from "@/pages/chapi/wallet";
    import VueJsonPretty from 'vue-json-pretty';
    export default {
        name: "DialogCustom",
        components: {
            VueJsonPretty,
        },
        props: {
            item: Object,
        },
        data: () => ({
            showDialog: false,
        }),
        methods: {
            credDisplayName: function (vc) {
                return vc.name ? vc.name : getCredentialType(vc.type)
            },
        },
    };
</script>

<style scoped>
    .infoIcon {
        position: absolute;
        right: 0;
        top: 0;
        padding: 10px 15px;
        opacity: .4;
        transition: all 0.5s ease;
    }
    .infoIcon:hover {
        opacity: 1;
    }
    .md-dialog-container {
        width: 100% !important;
    }
</style>
