<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex relative flex-col justify-between w-80 h-screen shadow bg-gradient-dark">
    <img class="absolute -top-1 left-0" src="@/assets/img/sidebar-flare.png" />
    <div class="flex z-10 flex-col justify-start items-start pb-8 h-full">
      <div class="flex justify-start items-center px-10">
        <Logo class="mt-13" />
      </div>
      <div class="flex flex-col flex-grow justify-start mt-8 w-full">
        <slot>
          <ul class="">
            <sidebar-link
              v-for="(link, index) in sidebarLinks"
              :key="link.name + index"
              :to="link.path"
              :link="link"
            />
          </ul>
        </slot>
        <div class="my-5 mx-10 h-px opacity-20 bg-neutrals-white"></div>
        <signout />
      </div>
    </div>
    <div class="flex z-10 flex-col items-start px-10 text-sm text-neutrals-white">
      <span
        tabindex="0"
        class="
          mb-2
          focus:rounded focus:ring-1
          opacity-60
          hover:opacity-100
          focus:opacity-100
          cursor-pointer
          ring-primary-blue
          underline-white
        "
        >{{ i18n.privacyPolicy }}</span
      >
      <span
        tabindex="0"
        class="
          focus:rounded focus:ring-1
          opacity-60
          hover:opacity-100
          focus:opacity-100
          cursor-pointer
          ring-primary-blue
          underline-white
        "
        >{{ i18n.terms }}</span
      >
      <div class="flex flex-row justify-start items-center my-6">
        <span
          tabindex="0"
          class="
            whitespace-nowrap
            focus:rounded focus:ring-1
            opacity-60
            hover:opacity-100
            focus:opacity-100
            cursor-pointer
            ring-primary-blue
            underline-white
          "
          >© {{ date }} TrustBloc</span
        >
        <span class="px-2 opacity-60">･</span>
        <!-- TODO: remove locale-switcher class after sass styles are deleted -->
        <locale-switcher
          class="
            focus:rounded focus:ring-1
            opacity-60
            hover:opacity-100
            focus:opacity-100
            text-neutrals-white
            underline-white
            ring-primary-blue
            locale-switcher
          "
        />
      </div>
    </div>
  </div>
</template>
<script>
import SidebarLink from './SidebarLink.vue';
import Logo from '@/components/Logo/Logo.vue';
import Signout from '@/components/Signout/Signout.vue';
import LocaleSwitcher from '@/components/LocaleSwitcher/LocaleSwitcher.vue';

export default {
  name: 'Sidebar',
  components: {
    SidebarLink,
    Signout,
    Logo,
    LocaleSwitcher,
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
      date: new Date().getFullYear(),
    };
  },
  computed: {
    i18n() {
      return this.$t('ContentFooter');
    },
  },
  methods: {},
};
</script>

<!-- TODO: remove after sass styles are deleted -->
<style>
.locale-switcher > span {
  color: #fff;
}
</style>
