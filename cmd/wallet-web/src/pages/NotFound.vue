<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import useBreakpoints from '@/plugins/breakpoints';
import Header from '@/components/Header/Header.vue';
import Footer from '@/components/Footer/Footer.vue';
import StyledButton from '@/components/StyledButton/StyledButton.vue';
import AppLink from '@/components/AppLink/AppLink.vue';
import IconError from '@/components/icons/IconError.vue';

const router = useRouter();
const breakpoints = useBreakpoints();
const loading = ref(false);

function handleClick() {
  loading.value = true;
  router.push({ name: 'vaults' });
}
</script>

<template>
  <div
    class="
      flex-grow
      justify-start
      items-center
      mx-auto
      max-w-7xl
      h-screen
      max-h-screen
      shadow-main-container
      flex flex-col
    "
  >
    <Header :has-custom-gradient="true" :show-menu-dropdown="false">
      <template #gradientContainer>
        <div class="absolute h-15 bg-gradient-full oval" />
      </template>
    </Header>
    <div
      class="
        overflow-hidden
        relative
        z-10
        flex-grow
        justify-start
        items-center
        px-6
        pt-32
        w-full
        h-full
        flex flex-col
        bg-neutrals-softWhite
      "
    >
      <IconError />
      <span class="mt-6 text-3xl text-neutrals-dark">Page not found ⸱ 404</span>
      <span class="mx-2 mt-2 text-base text-center text-neutrals-medium">
        This page cannot be found, or you have navigated to an invalid URL.
      </span>
      <StyledButton type="btn-primary" :loading="loading" class="relative mt-6">
        <AppLink :to="{ name: 'vaults' }" @click="handleClick">Go Home</AppLink>
      </StyledButton>
    </div>
    <Footer
      v-if="!(breakpoints.xs || breakpoints.sm)"
      class="sticky bottom-0 z-20 border-t border-neutrals-thistle bg-neutrals-magnolia"
    />
  </div>
</template>

<style scoped>
.oval {
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  filter: blur(50px);
  width: 15.625rem; /* 250px */
  top: 2.0625rem; /* 33px */
}
</style>
