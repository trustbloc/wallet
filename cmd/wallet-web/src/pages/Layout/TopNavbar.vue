/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="footer-gradient">
        <div class="md-toolbar-row">
            <div class="px-6 md-toolbar-section-end">
                <div>
                    <i class="md-title text-xl text-white">
                      <md-icon  style="color: white; margin: 10px;">account_circle</md-icon>
                        {{$store.getters.getCurrentUser ? $store.getters.getCurrentUser.username : ''}}</i>
                </div>
                <settings/>
                <md-button
                        class="md-just-icon md-simple md-toolbar-toggle"
                        :class="{ toggled: $sidebar.showSidebar }"
                        @click="toggleSidebar">
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </md-button>
            </div>
        </div>
    </div>
</template>

<script>
    import Settings from "@/pages/chapi/Settings";
    export default {
        data() {
            return {}
        },
        methods: {
            toggleSidebar() {
                this.$sidebar.displaySidebar(!this.$sidebar.showSidebar);
            }
        },
        components: {
          Settings,
        },
        computed: {
        }
    };
</script>
<style scoped>
.footer-gradient {
  background: linear-gradient(to bottom ,#000000, #0c0116,#14061D, #261131);
}
</style>
