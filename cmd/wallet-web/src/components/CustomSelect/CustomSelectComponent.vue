<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="custom-select" :tabindex="tabindex" @blur="open = false">
    <div class="custom-select inline-flex" :class="{ open: open }" @click="open = !open">
      <div class="selected grow">
        {{ selected }}
      </div>
      <div class="flex h-6 w-6 flex-none justify-end">
        <img src="@/assets/img/icons-sm--chevron-down-icon.svg" />
      </div>
    </div>
    <div class="items" :class="{ selectHide: !open }">
      <div
        v-for="(option, i) of options"
        :key="i"
        @click="
          selected = option.name;
          open = false;
          $emit('selected', option.id);
        "
      >
        {{ option.name }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    options: {
      type: Array,
      required: true,
    },
    default: {
      type: String,
      required: false,
      default: null,
    },
  },
  emits: ['selected'],
  data() {
    return {
      selected: this.default,
      open: false,
      tabindex: 0,
    };
  },
};
</script>

<style scoped>
.custom-select {
  position: relative;
  width: 100%;
  text-align: left;
  outline: none;
}

.custom-select .selected {
  background-color: theme('colors.neutrals.bonjour');
  color: theme('colors.neutrals.dark');
  cursor: pointer;
}

.custom-select .selected:after {
  position: absolute;
}

.custom-select .items {
  color: theme('colors.neutrals.dark');
  top: theme('spacing.1');
  overflow: hidden;
  position: absolute;
  background-color: theme('colors.neutrals.white');
  border-radius: theme('spacing.2');
  padding-top: theme('spacing.2');
  padding-bottom: theme('spacing.2');
  box-shadow: 0px 4px 20px 0px rgba(25, 12, 33, 0.2);
  width: 100%;
  z-index: 1;
  font-size: theme('fontSize.base');
  font-weight: theme('fontWeight.bold');
}

.custom-select .items div {
  color: theme('colors.neutrals.medium');
  padding-top: theme('spacing.2');
  padding-bottom: theme('spacing.2');
  padding-left: theme('spacing.4');
  padding-right: theme('spacing.4');
  cursor: pointer;
}

.custom-select .items div:hover {
  background-color: theme('colors.neutrals.softWhite');
  color: theme('colors.neutrals.dark');
  font-weight: theme('fontWeight.bold');
}

.selectHide {
  display: none;
}
</style>
