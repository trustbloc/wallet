/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div
    class="sidebar"
    :data-color="sidebarItemColor"
    :data-image="sidebarBackgroundImage"
    :style="sidebarStyle"
  >
    <div class="logo">
      <a href="#" class="simple-text logo-mini">
        <div class="logo-img">
          <img :src="imgLogo" alt="" />
        </div>
      </a>

      <a
        href="https://github.com/trustbloc/edge-agent"
        class="simple-text"
      >
        {{ title }}
      </a>
    </div>
    <div class="sidebar-wrapper">
      <slot name="content"></slot>
      <md-list class="nav">
        <!--By default vue-router adds an active class to each route link. This way the links are colored when clicked-->
        <slot>
          <sidebar-link
            v-for="(link, index) in sidebarLinks"
            :key="link.name + index"
            :to="link.path"
            :link="link"
          >
          </sidebar-link>
        </slot>
      </md-list>
      <div class="dev-mode">
        <div>
          <md-checkbox :value="!isDevMode" @change="updateDevMode(!isDevMode)">Developer Mode
          </md-checkbox>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import SidebarLink from "./SidebarLink.vue";
import {mapGetters, mapMutations} from 'vuex'

export default {
  components: {
    SidebarLink
  },
  methods: mapMutations(['updateDevMode']),
  props: {
    title: {
      type: String,
      default: "User Agent Wallet"
    },
    sidebarBackgroundImage: {
      type: String,
      default: require("@/assets/img/sidebar-2.jpg")
    },
    imgLogo: {
      type: String,
      default: require("@/assets/img/logo.png")
    },
    sidebarItemColor: {
      type: String,
      default: "green",
      validator: value => {
        let acceptedValues = ["", "purple", "blue", "green", "orange", "red"];
        return acceptedValues.indexOf(value) !== -1;
      }
    },
    sidebarLinks: {
      type: Array,
      default: () => []
    },
    autoClose: {
      type: Boolean,
      default: true
    }
  },
  provide() {
    return {
      autoClose: this.autoClose
    };
  },
  computed: {
    ...mapGetters(['isDevMode']),
    sidebarStyle() {
      return {
        backgroundImage: `url(${this.sidebarBackgroundImage})`
      };
    }
  }
};
</script>
<style>
@media screen and (min-width: 995px) {
  .nav-mobile-menu {
    display: none;
  }
}
</style>

<style scoped>
.dev-mode {
  position: absolute;
  bottom: 0;
  width: 100%
}

.dev-mode > div {
  text-align: center;
}
</style>
