/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <md-toolbar md-elevation="0" class="md-transparent" :style="topNavStyle">
        <div class="md-toolbar-row">
            <div class="md-toolbar-section-end">

                <div>
                    <i class="md-title" style="font-size: 20px; color: white !important;">
                        <md-icon style="color: white; margin: 10px;">person_outline</md-icon>
                        logged in as {{$store.getters.getCurrentUser ? $store.getters.getCurrentUser.username : ''}}</i>

                </div>
                <logout/>
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
    </md-toolbar>
</template>

<script>

    import Logout from "@/pages/chapi/Logout.vue";

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
            Logout,
        },
        computed: {
            topNavStyle() {
                return {
                    backgroundImage: `url(${require("@/assets/img/topbar-2.jpg")})`,
                    'background-size': 'cover',
                    height: '70px'
                };
            }
        }
    };
</script>

<style lang="css"></style>
