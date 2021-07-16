/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>directions_car</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> Name: {{item.credentialSubject.family_name || 'N/A' }},
                {{item.credentialSubject.given_name || 'N/A' }} </p>
            <p> Date of Birth: {{item.credentialSubject.birthdate || "N/A" }} </p>
            <p> Driving Privileges: {{item.credentialSubject.driving_privileges || "N/A"}} </p>
            <p> Issuing Authority: {{item.credentialSubject.issuing_authority || "N/A" }} </p>
            <p> Address: {{item.credentialSubject.resident_address || "N/A" }},
                {{item.credentialSubject.resident_city || 'N/A' }},
                {{item.credentialSubject.resident_postal_code || 'N/A' }}</p>

        </div>
    </div>
</template>

<script>
    export default {
        name: "drivers-license",
        props: {
            item: Object,
        },
    };
</script>
