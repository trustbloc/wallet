/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>local_gas_station</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> Address: {{ item.credentialSubject.address.address|| 'N/A'}}</p>
            <p> Category: {{item.credentialSubject.category || 'N/A' }} </p>
            <p> HS Code: {{item.credentialSubject.hscode || 'N/A' }} </p>
            <p> Volume: {{item.credentialSubject.volume || 'N/A' }} </p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "crude-product-card",
        props: {
            item: Object,
        },
    };
</script>
