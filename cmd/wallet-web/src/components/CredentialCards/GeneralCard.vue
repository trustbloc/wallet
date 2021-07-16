/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>perm_identity</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> {{ item.description || 'Verifiable Credential'}}</p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "general-card",
        props: {
            item: Object,
        },
    };
</script>
