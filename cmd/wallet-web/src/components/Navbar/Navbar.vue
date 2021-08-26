<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="
      flex flex-grow
      md:flex-grow-0
      relative
      flex-col
      justify-between
      w-screen
      md:w-80 md:h-screen
      shadow
      min-w-80
      md:bg-gradient-dark
    "
  >
    <div class="hidden md:block absolute w-full opacity-40 bg-gradient-full oval" />
    <div class="flex md:z-10 flex-col justify-start items-start pb-8 h-full">
      <div class="hidden md:flex justify-start items-center px-10">
        <Logo class="mt-13" />
      </div>
      <div class="flex flex-col flex-grow justify-start mt-8 w-full">
        <ul>
          <slot />
        </ul>
        <div class="my-5 mx-10 h-px opacity-20 bg-neutrals-white" />
        <signout />
      </div>
    </div>
    <div class="flex md:z-10 flex-col items-start px-10 text-sm text-neutrals-white">
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

<style>
.oval {
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  filter: blur(50px);
  width: 22.625rem; /* 362px */
  height: 8.25rem; /* 132px */
  top: -5rem; /* 80px */
}
/* TODO: remove after sass styles are deleted */
.locale-switcher > span {
  color: #fff;
}
</style>
