<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <router-link
    :to="{ name: 'credential-details', params: { id } }"
    class="inline-flex items-center credentialCard"
  >
    <div class="flex-none w-12 h-12 border-opacity-10">
      <img :src="require(`@/assets/img/${icon}`)" />
    </div>
    <div class="flex-grow p-4">
      <span class="text-sm md:text-base font-bold text-left text-neutrals-dark overflow-ellipsis">
        {{ title }}
      </span>
    </div>
    <div class="flex-none credentialArrowContainer">
      <div class="p-1">
        <img src="@/assets/img/credential--arrow-right-icon.svg" />
      </div>
    </div>
  </router-link>
</template>

<script>
export default {
  name: 'CredentialPreview',
  props: {
    id: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
  },
};
</script>
