/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="input-container">
        <input :id="'input-' + label" v-bind="$attrs" v-on="$listeners" :value="value" @input="$emit('update', $event.target.value)">
        <label :for="'input-' + label" class="input-label">{{ label }}</label>
        <span class="input-helper">{{ helperMessage }}</span>
        <div class="fader">
            <!-- TODO: use inline svg instead once https://github.com/trustbloc/edge-agent/issues/816 is fixed -->
            <img src="@/assets/img/danger-icon.svg" />
        </div>
    </div>
</template>

<script>
export default {
    inheritAttrs: false,
    name: "input-field",
    props: {
        label: {
            type: String,
            default: '',
        },
        helperMessage: {
            type: String,
            default: '',
        },
        value: {
            type: String,
            default: '',
        },
    },
    model: {
        prop: "value",
        event: "update"
    },
    methods: {},
}
</script>
