/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="sidebar fixed inset-0 flex z-40">
      <!-- Sidebar starts -->
      <div class="md:w-1/4 xl:w-1/5 2xl:w-1/6 h-screen absolute sm:relative gradient shadow md:h-auto flex-col justify-between hidden sm:flex">
        <div class="h-full flex flex-col justify-start align-left px-12 pb-8">
          <div class="flex justify-start items-center">
            <Logo class="mt-12" href="" />
          </div>
          <div class="sidebar-wrapper flex-grow flex flex-col justify-between mt-8">
            <slot name="content"></slot>
                <slot>
                  <ul class="mt-12">
                    <li class="">
                        <sidebar-link
                          v-for="(link, index) in sidebarLinks"
                          :key="link.name + index"
                          :to="link.path"
                          :link="link">
                         </sidebar-link>
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
  import SidebarLink from "./SidebarLink.vue";
  import Logo from '@/components/Logo/Logo.vue';

    export default {
      components: {
        SidebarLink,
        Logo,
      },
      props: {
        sidebarLinks: {
          type: Array,
          default: () => []
        },
        autoClose: {
          type: Boolean,
          default: true
        }
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
  background: linear-gradient(to bottom ,#14061D, #261131,#0c0116,#13113F,#1A0C22);
}
</style>
