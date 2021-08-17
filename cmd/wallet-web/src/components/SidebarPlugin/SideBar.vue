<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex fixed inset-0 z-40 sidebar">
    <!-- Sidebar starts -->
    <div
      class="
        hidden
        sm:flex
        absolute
        sm:relative
        flex-col
        justify-between
        md:w-1/4
        xl:w-1/5
        2xl:w-1/6
        h-screen
        md:h-auto
        shadow
        gradient
      "
    >
      <div class="flex flex-col justify-start px-12 pb-8 h-full align-left">
        <div class="flex justify-start items-center">
          <Logo class="mt-12" href="" />
        </div>
        <div class="flex flex-col flex-grow justify-between mt-8 sidebar-wrapper">
          <slot name="content" />
          <slot>
            <ul class="mt-12">
              <li class="">
                <sidebar-link
                  v-for="(link, index) in sidebarLinks"
                  :key="link.name + index"
                  :to="link.path"
                  :link="link"
                />
              </li>
            </ul>
          </slot>
        </div>
      </div>
      <!-- Sidebar ends -->
    </div>
  </div>
</template>
<script>
import SidebarLink from '@/components/SidebarPlugin/SidebarLink.vue';
import Logo from '@/components/Logo/Logo.vue';

export default {
  name: 'SideBar',
  components: {
    SidebarLink,
    Logo,
  },
  props: {
    sidebarLinks: {
      type: Array,
      default: () => [],
    },
    autoClose: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      moved: true,
    };
  },
  methods: {},
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
.gradient {
  background: linear-gradient(to bottom, #14061d, #261131, #0c0116, #13113f, #1a0c22);
}
</style>
