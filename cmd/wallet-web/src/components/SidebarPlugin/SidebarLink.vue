<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <li class="md-list-item">
    <router-link
      class="md-list-item-router md-list-item-container md-button-clean"
      v-bind="$attrs"
      @click="hideSidebar"
    >
      <div class="md-list-item-content md-ripple">
        <slot>
          <md-icon>{{ link.icon }}</md-icon>
          <p>{{ link.name }}</p>
        </slot>
      </div>
    </router-link>
  </li>
</template>
<script>
export default {
  name: 'SidebarLink',
  inject: {
    autoClose: {
      default: true,
    },
  },
  props: {
    link: {
      type: [String, Object],
      default: () => {
        return {
          name: '',
          path: '',
          icon: '',
        };
      },
    },
    tag: {
      type: String,
      default: 'router-link',
    },
  },
  methods: {
    hideSidebar() {
      if (this.autoClose && this.$sidebar && this.$sidebar.showSidebar === true) {
        this.$sidebar.displaySidebar(false);
      }
    },
  },
};
</script>
<style></style>
