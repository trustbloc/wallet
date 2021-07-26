<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="gradient">
    <div class="sticky top-0 z-40">
      <div class="flex justify-between items-center px-6 w-full h-20">
        <!-- left navbar -->
        <div class="flex"></div>
        <!-- right navbar -->
        <div class="flex relative items-center">
          <div>
            <i class="px-2 text-xl md-title text-white">
              <md-icon style="color: white; margin: 10px" class="px-2">account_circle</md-icon>
              {{ $store.getters.getCurrentUser ? $store.getters.getCurrentUser.username : '' }}</i
            >
          </div>
          <img
            id="dashboardSettings"
            src="@/assets/img/settings.png"
            style="color: white"
            class="w-6 h-6 rounded-full shadow-lg"
            @click="dropDownOpen = !dropDownOpen"
          />
        </div>
      </div>
      <!-- dropdown menu -->
      <div
        class="
          absolute
          right-0
          bottom-10
          mr-6
          w-48
          rounded-b-lg
          shadow-xl
          bg-gray-100
          border border-t-0
          text-gray-700
        "
        :class="dropDownOpen ? '' : 'hidden'"
      >
        <a href="/DIDManagement" class="block py-2 px-4 hover:bg-gray-200">
          <md-icon class="px-4 mr-2 w-8 h-8 hover:bg-gray-300 bg-indigo-100">flip_to_back</md-icon
          >DID Management</a
        >
      </div>
      <!-- dropdown menu end -->
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      dropDownOpen: false,
    };
  },
  methods: {
    toggleSidebar() {
      this.$sidebar.displaySidebar(!this.$sidebar.showSidebar);
    },
  },
};
</script>
<style scoped>
.gradient {
  background: linear-gradient(to bottom, #000000, #0c0116, #14061d, #261131);
}
</style>
