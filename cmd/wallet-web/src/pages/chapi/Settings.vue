/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <!-- right navbar -->
  <div class="flex items-center relative">
    <md-icon  style="color: white; margin: 20px;" @click="dropDownOpen = !dropDownOpen">settings</md-icon>
    <!-- dropdown menu -->
    <div class="absolute bg-gray-100 border border-t-0 shadow-xl text-gray-700 rounded-b-lg w-48 bottom-10 right-0 mr-6" :class="dropDownOpen ? '' : 'hidden'">
      <a href="#" class="block px-4 py-2 hover:bg-gray-200">Account</a>
      <a href="#" class="block px-4 py-2 hover:bg-gray-200">Settings</a>
      <a href="#" class="block px-4 py-2 hover:bg-gray-200">Logout</a>
    </div>
    <!-- dropdown menu end -->
  </div>
</template>

<script>
export default {
  name: "Settings",
  data() {
    return {
      dropDownOpen: false
    }
  },
}
</script>

<style scoped>

</style>
