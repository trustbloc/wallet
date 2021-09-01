<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="
      relative
      flex-col
      md:flex-grow-0
      justify-between
      w-screen
      min-w-80
      md:w-80 md:h-screen md:bg-gradient-dark
      shadow
      flex flex-grow
    "
  >
    <div class="hidden md:block absolute bg-gradient-full opacity-40 oval" />
    <div class="flex md:z-10 flex-col justify-start items-start pb-8 h-full">
      <div class="hidden md:flex justify-start items-center px-10">
        <Logo class="mt-13" />
      </div>
      <div class="flex flex-col flex-grow justify-start mt-8 w-full">
        <ul>
          <slot />
        </ul>
        <div class="my-5 mx-10 h-px bg-neutrals-white opacity-20" />
        <signout />
      </div>
    </div>
    <div class="flex md:z-10 flex-col items-start px-10 text-sm text-neutrals-white">
      <span
        tabindex="0"
        class="
          mb-2
          focus:rounded focus:ring-1
          ring-primary-blue
          opacity-60
          hover:opacity-100
          focus:opacity-100
          cursor-pointer
          underline-white
        "
        >{{ i18n.privacyPolicy }}</span
      >
      <span
        tabindex="0"
        class="
          focus:rounded focus:ring-1
          ring-primary-blue
          opacity-60
          hover:opacity-100
          focus:opacity-100
          cursor-pointer
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
            ring-primary-blue
            opacity-60
            hover:opacity-100
            focus:opacity-100
            cursor-pointer
            underline-white
          "
          >© {{ date }} TrustBloc</span
        >
        <span class="px-2 opacity-60">･</span>
        <!-- TODO: remove locale-switcher class after sass styles are deleted -->
        <locale-switcher
          class="
            text-neutrals-white
            focus:rounded focus:ring-1
            ring-primary-blue
            opacity-60
            hover:opacity-100
            focus:opacity-100
            underline-white
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

<style scoped>
.oval {
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  filter: blur(50px);
  width: 22.625rem; /* 362px */
  height: 8.25rem; /* 132px */
  top: -5rem; /* 80px */
}
</style>

<!-- TODO: remove after sass styles are deleted -->
<style>
.locale-switcher > span {
  color: #fff;
}
</style>
