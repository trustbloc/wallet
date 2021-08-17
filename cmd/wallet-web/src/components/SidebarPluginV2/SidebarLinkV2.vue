<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <li
    class="
      justify-start
      items-center
      w-full
      h-16
      font-bold
      hover:bg-gradient-to-r hover:opacity-100 hover:from-neutrals-black
      opacity-60
      bar
      flex flex-row
    "
  >
    <router-link class="flex flex-row justify-start items-center px-10" v-bind="$attrs">
      <img v-if="icon" :src="require(`@/assets/img/${icon}`)" class="w-8 h-8" />
      <div v-else class="w-8 h-8" />
      <span class="ml-4 text-lg text-neutrals-white">{{ heading }}</span>
    </router-link>
  </li>
</template>

<script>
export default {
  name: 'SidebarLinkV2',
  props: {
    heading: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      default: '',
    },
  },
};
</script>

<style scoped>
.bar:hover:before {
  position: absolute;
  content: "''";
  display: block;
  background-color: theme('colors.primary.purple.hashita');
  height: theme('spacing.16');
  width: 4px;
}
</style>
