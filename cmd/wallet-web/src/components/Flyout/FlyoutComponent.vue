<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div class="relative">
    <slot name="button" :toggle-flyout-menu="toggleFlyoutMenu" :set-show-tooltip="setShowTooltip" />
    <ToolTipComponent v-if="showTooltip && !showFlyoutMenu" :tool-tip-label="toolTipLabel" />
    <slot v-if="showFlyoutMenu" name="menu" :toggle-flyout-menu="toggleFlyoutMenu" />
  </div>
</template>

<script>
import ToolTipComponent from '@/components/ToolTip/ToolTipComponent.vue';

export default {
  name: 'FlyoutComponent',
  components: {
    ToolTipComponent,
  },
  props: {
    toolTipLabel: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      showTooltip: false,
      showFlyoutMenu: false,
    };
  },
  mounted() {
    document.addEventListener('click', this.close);
  },
  beforeUnmount() {
    document.removeEventListener('click', this.close);
  },
  methods: {
    toggleFlyoutMenu() {
      this.showFlyoutMenu = !this.showFlyoutMenu;
      this.showTooltip = false;
    },
    close(e) {
      if (!this.$el.contains(e.target)) this.showFlyoutMenu = false;
    },
    setShowTooltip(value) {
      this.showTooltip = value;
    },
  },
};
</script>
