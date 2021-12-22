<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="input-container">
    <input
      :id="'input-' + label"
      ref="input"
      v-bind="$attrs"
      :name="name"
      :disabled="!characterCount"
      @input="onInput"
      @keyup.enter="characterCount"
      @click="active = true"
    />
    <label :for="'input-' + label" class="input-label">{{ label }}</label>
    <span v-if="valid || (!submitted && !value.length)" class="input-helper">{{
      helperMessage
    }}</span>
    <span v-else class="text-sm font-bold text-primary-vampire">{{ errorMessage }}</span>
    <div class="fader" />
    <!-- TODO: use inline svg instead once https://github.com/trustbloc/edge-agent/issues/816 is fixed -->
    <img
      v-if="(value.length && !valid) || (submitted && !valid)"
      class="danger-icon"
      src="@/assets/img/danger-icon.svg"
    />
    <img
      v-else-if="submitted && valid"
      class="checkmark-icon"
      src="@/assets/img/icons-checkmark.svg"
    />
    <span v-else class="input-word-limit">{{ characterCount }}</span>
  </div>
</template>

<script>
import { computed, onMounted, toRefs, watchEffect, ref } from 'vue';

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
      required: true,
    },
    value: {
      type: String,
      required: true,
    },
    helperMessage: {
      type: String,
      default: '',
    },
    errorMessage: {
      type: String,
      default: '',
    },
    submitted: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['input', 'update'],
  setup(props, { attrs, emit }) {
    const { value, label, submitted } = toRefs(props);
    const characterCount = computed(() =>
      value ? value.value.length + '/' + attrs['maxlength'] : 0 + '/' + attrs['maxlength']
    );
    const name = computed(() => label.value.toLowerCase());
    const input = ref(null);
    const pattern = new RegExp(attrs['pattern']);
    const valid = computed(() => {
      return pattern.test(value.value);
    });

    onMounted(() => {
      watchEffect(() => {
        if (submitted.value) {
          input.value.focus();
        }
      });
    });

    return { characterCount, input, name, pattern, valid };
  },
  data() {
    return {
      active: false,
    };
  },
  methods: {
    onInput: function (event) {
      this.$emit('input', {
        name: event.target.value,
        valid: this.pattern.test(event.target.value),
      });
    },
  },
};
</script>
