/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>house</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> Name: {{item.credentialSubject.familyName || 'Smith' }},
                {{item.credentialSubject.givenName || 'John' }} </p>
            <p> Date of Birth: {{item.credentialSubject.birthDate || 'N/A' }} </p>
            <p> Country of Birth: {{item.credentialSubject.birthCountry || 'N/A' }} </p>
            <p> Resident Since: {{item.credentialSubject.residentSince || 'N/A' }} </p>
            <p> LPR Number: {{item.credentialSubject.lprNumber || 'N/A' }} </p>
            <p> LPR Category: {{item.credentialSubject.lprCategory || 'N/A' }}</p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "permanent-resident-card",
        props: {
            item: Object,
        },
    };
</script>
