/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>bar_chart</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> Certificate Number: {{ item.credentialSubject.cmtr.certificateNumber || 'N/A'}}</p>
            <p> Company Name: {{item.credentialSubject.cmtr.companyName || "N/A" }} </p>
            <p> Phone Number: {{item.credentialSubject.cmtr.companyPhoneNumber || "N/A" }}</p>
            <p> Website: {{item.credentialSubject.cmtr.companyWebsite || "N/A" }} </p>
            <p> Authorizing Party: {{item.credentialSubject.cmtr.authorizingPartyName || 'N/A' }} </p>
            <p> Authorizing Date: {{item.credentialSubject.cmtr.authorizingPartyDate || 'N/A' }} </p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "mill-test-card",
        props: {
            item: Object,
        },
    };
</script>
