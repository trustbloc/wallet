/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="cardBody">
        <div class="cardDetailsL">
            <md-icon>person</md-icon>
        </div>
        <div class="cardDetailsR">
            <p> Name: {{item.credentialSubject.name || "N/A" }} </p>
            <p> Email: {{item.credentialSubject.email || "N/A" }} </p>
            <p> Student ID: {{item.credentialSubject.studentid || "N/A" }} </p>
            <p> University: {{item.credentialSubject.university || "N/A" }} </p>
            <p> Semester: {{item.credentialSubject.semester || "N/A" }} </p>
        </div>
    </div>
</template>

<script>
    export default {
        name: "student-card",
        props: {
            item: Object,
        },
    };
</script>
