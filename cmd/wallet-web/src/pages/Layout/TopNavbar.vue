/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="gradient">
      <div class="sticky top-0 z-40">
        <div class="w-full h-20 px-6 flex items-center justify-between">
          <!-- left navbar -->
          <div class="flex">
          </div>
          <!-- right navbar -->
          <div class="flex items-center relative">
            <div>
              <i class="md-title text-xl px-2 text-white">
                <md-icon  style="color: white; margin: 10px;" class="px-2">account_circle</md-icon>
                {{$store.getters.getCurrentUser ? $store.getters.getCurrentUser.username : ''}}</i>
            </div>
            <img id="dashboardSettings" src="@/assets/img/settings.png"  style="color:white;"
                 class="w-6 h-6 rounded-full shadow-lg" @click="dropDownOpen = !dropDownOpen">
          </div>
        </div>
        <!-- dropdown menu -->
        <div class="absolute bg-gray-100 border border-t-0 shadow-xl text-gray-700 rounded-b-lg w-48 bottom-10 right-0 mr-6" :class="dropDownOpen ? '' : 'hidden'">
          <a href="/DIDManagement" class="block px-4 py-2 hover:bg-gray-200">
            <md-icon class="h-8 w-8 mr-2 px-4 hover:bg-gray-300 bg-indigo-100">flip_to_back</md-icon>DID Management</a>
        </div>
        <!-- dropdown menu end -->
      </div>
    </div>
</template>

<script>
    export default {
        methods: {
            toggleSidebar() {
                this.$sidebar.displaySidebar(!this.$sidebar.showSidebar);
            }
        },
      data() {
        return {
          dropDownOpen: false
        }
      },
    };
</script>
<style scoped>
.gradient {
  background: linear-gradient(to bottom ,#000000, #0c0116,#14061D, #261131);
}
</style>
