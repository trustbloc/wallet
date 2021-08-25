<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="flex relative flex-col justify-between w-80 h-screen shadow min-w-80 bg-gradient-dark"
  >
    <img class="absolute -top-1 left-0" src="@/assets/img/navbar-flare.png" />
    <div class="flex z-10 flex-col justify-start items-start pb-8 h-full">
      <div class="flex justify-start items-center px-10">
        <Logo class="mt-13" />
      </div>
      <div class="flex flex-col flex-grow justify-start mt-8 w-full">
        <ul class="">
          <slot />
        </ul>
        <div class="my-5 mx-10 h-px opacity-20 bg-neutrals-white" />
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
import Logo from '@/components/Logo/Logo.vue';
import Signout from '@/components/Signout/Signout.vue';
import LocaleSwitcher from '@/components/LocaleSwitcher/LocaleSwitcher.vue';

export default {
  name: 'Navbar',
  components: {
    Signout,
    Logo,
    LocaleSwitcher,
  },
  data() {
    return {
      date: new Date().getFullYear(),
    };
  },
  computed: {
    i18n() {
      return this.$t('Footer');
    },
  },
};
</script>

<!-- TODO: remove after sass styles are deleted -->
<style>
.locale-switcher > span {
  color: #fff;
}
</style>
