<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="input-container">
    <input
      :id="'input-' + label"
      v-bind="$attrs"
      :disabled="!characterCount"
      @input="$emit('update', $event.target.value)"
      @keyup.enter="characterCount"
    />
    <label :for="'input-' + label" class="input-label">{{ label }}</label>
    <span class="input-helper">{{ helperMessage }}</span>
    <div class="fader" />
    <!-- TODO: use inline svg instead once https://github.com/trustbloc/edge-agent/issues/816 is fixed -->
    <img src="@/assets/img/danger-icon.svg" />
    <span class="input-wordlimit">{{ characterCount }}</span>
  </div>
</template>

<script>
export default {
  name: 'InputField',
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'update',
  },
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
  emits: ['update'],
  computed: {
    characterCount() {
      return this.value
        ? this.value.length + '/' + this.$attrs['maxlength']
        : 0 + '/' + this.$attrs['maxlength'];
    },
  },
};
</script>
