/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>school</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> Name: {{ item.credentialSubject.name || 'John Smith'}}</p>
            <p> Degree Type: {{item.credentialSubject.degree.type || 'N/A' }} </p>
            <p> Degree Name: {{item.credentialSubject.degree.name || item.credentialSubject.degree.degree || 'N/A' }} </p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "university-card",
        props: {
            item: Object,
        },
    };
</script>
