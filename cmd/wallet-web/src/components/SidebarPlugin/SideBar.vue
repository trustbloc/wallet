/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="fixed inset-0 flex z-40">
      <!-- Sidebar starts -->
      <div class="w-2/12  absolute sm:relative gradient shadow md:h-auto flex-col justify-between hidden sm:flex">
        <div class="px-16 py-8">
          <div class="flex items-center">
            <img class="h-10 w-10 mr-2 mt-2" src="@/assets/img/logo.png" alt="">
            <h1 class="font-semibold text-white text-2xl lg:text-4xl tracking-tight">TrustBloc</h1>
          </div>
          <div class="sidebar-wrapper">
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

    export default {
      components: {
        SidebarLink
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
