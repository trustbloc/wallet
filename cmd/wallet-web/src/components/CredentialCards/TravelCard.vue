/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>flight</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> Name: {{item.credentialSubject.name || "N/A" }} </p>
            <p> Gender: {{item.credentialSubject.sex || "N/A" }} </p>
            <p> Country: {{item.credentialSubject.country || "N/A" }} </p>
            <p> Travel Card ID: {{item.credentialSubject.travelcardid || "N/A" }} </p>
            <p> Date of Birth: {{item.credentialSubject.dob || "N/A" }} </p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "travel-card",
        props: {
            item: Object,
        },
    };
</script>
