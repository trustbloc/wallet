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
      shadow
      min-w-80
      md:w-80 md:h-auto md:min-h-screen md:bg-gradient-dark
      flex flex-grow
    "
  >
    <div class="hidden md:block absolute opacity-40 bg-gradient-full oval" />
    <div class="flex md:z-10 flex-col justify-start items-start pb-8 h-full">
      <div class="hidden md:flex justify-start items-center px-10">
        <Logo class="cursor-pointer mt-13" @click="$router.push({ name: 'vaults' })" />
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
          focus:outline-none
          ring-primary-blue
          underline-white
        "
        >{{ t('Footer.privacyPolicy') }}</span
      >
      <span
        tabindex="0"
        class="
          focus:rounded focus:ring-1
          opacity-60
          hover:opacity-100
          focus:opacity-100
          cursor-pointer
          focus:outline-none
          ring-primary-blue
          underline-white
        "
        >{{ t('Footer.terms') }}</span
      >
      <div class="flex flex-row justify-start items-center my-6">
        <span
          tabindex="0"
          class="
            whitespace-nowrap
            focus:rounded focus:ring-1
            opacity-60
            hover:opacity-100
            cursor-pointer
            focus:opacity-100 focus:outline-none
            ring-primary-blue
            underline-white
          "
          >© {{ date }} TrustBloc</span
        >
        <span class="px-2 opacity-60">･</span>
        <locale-switcher
          class="
            focus:rounded focus:ring-1
            opacity-60
            hover:opacity-100
            focus:opacity-100 focus:outline-none
            text-neutrals-white
            ring-primary-blue
            underline-white
          "
        />
      </div>
    </div>
  </div>
</template>
<script>
import { watchEffect } from 'vue';
import { navbarMutations } from '@/components/Navbar';
import Logo from '@/components/Logo/Logo.vue';
import Signout from '@/components/Signout/Signout.vue';
import LocaleSwitcher from '@/components/LocaleSwitcher/LocaleSwitcher.vue';
import { useI18n } from 'vue-i18n';

export default {
  name: 'Navbar',
  components: {
    Signout,
    Logo,
    LocaleSwitcher,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      date: new Date().getFullYear(),
    };
  },
  created: function () {
    watchEffect(() => {
      this.setCurrentPage(this.$route.name);
    });
  },
  methods: {
    setCurrentPage(newPage) {
      navbarMutations.setCurrentPage(`navbar-link-${newPage}`);
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
