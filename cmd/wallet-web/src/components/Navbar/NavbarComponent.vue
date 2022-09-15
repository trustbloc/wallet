<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="relative flex w-screen min-w-80 grow flex-col justify-between shadow md:h-auto md:min-h-screen md:w-80 md:grow-0 md:bg-gradient-dark"
  >
    <div class="oval absolute hidden bg-gradient-full opacity-40 md:block" />
    <div class="flex h-full flex-col items-start justify-start pb-8 md:z-10">
      <div class="hidden items-center justify-start px-10 md:flex">
        <LogoComponent class="mt-13 cursor-pointer" @click="$router.push({ name: 'vaults' })" />
      </div>
      <div class="mt-8 flex w-full grow flex-col justify-start">
        <ul>
          <slot />
        </ul>
        <div class="my-5 mx-10 h-px bg-neutrals-white opacity-20" />
        <SignoutComponent />
      </div>
    </div>
    <div class="flex flex-col items-start px-10 text-sm text-neutrals-white md:z-10">
      <span
        tabindex="0"
        class="underline-white mb-2 cursor-pointer opacity-60 ring-primary-blue hover:opacity-100 focus:rounded focus:opacity-100 focus:outline-none focus:ring-1"
        >{{ t('Footer.privacyPolicy') }}</span
      >
      <span
        tabindex="0"
        class="underline-white cursor-pointer opacity-60 ring-primary-blue hover:opacity-100 focus:rounded focus:opacity-100 focus:outline-none focus:ring-1"
        >{{ t('Footer.terms') }}</span
      >
      <div class="my-6 flex flex-row items-center justify-start">
        <span
          tabindex="0"
          class="underline-white cursor-pointer whitespace-nowrap opacity-60 ring-primary-blue hover:opacity-100 focus:rounded focus:opacity-100 focus:outline-none focus:ring-1"
          >© {{ date }} TrustBloc</span
        >
        <span class="px-2 opacity-60">･</span>
        <LocaleSwitcherComponent
          class="underline-white text-neutrals-white opacity-60 ring-primary-blue hover:opacity-100 focus:rounded focus:opacity-100 focus:outline-none focus:ring-1"
        />
      </div>
    </div>
  </div>
</template>
<script>
import { watchEffect } from 'vue';
import { navbarMutations } from '@/components/Navbar';
import LogoComponent from '@/components/Logo/LogoComponent.vue';
import SignoutComponent from '@/components/Signout/SignoutComponent.vue';
import LocaleSwitcherComponent from '@/components/LocaleSwitcher/LocaleSwitcherComponent.vue';
import { useI18n } from 'vue-i18n';

export default {
  name: 'NavbarComponent',
  components: {
    SignoutComponent,
    LogoComponent,
    LocaleSwitcherComponent,
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
