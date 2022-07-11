<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex w-full flex-col items-start justify-start">
    <span v-if="heading.length" class="py-4 pl-3 text-xl font-bold text-neutrals-dark">{{
      heading
    }}</span>
    <table class="w-full border-t border-neutrals-chatelle">
      <tr
        v-for="(property, key) of credential?.properties"
        :key="key"
        class="border-b border-dotted border-neutrals-thistle"
      >
        <td class="w-2/5 whitespace-nowrap py-4 pr-6 pl-3 text-neutrals-medium">
          {{ property?.label }}
        </td>
        <td
          v-if="property?.type === 'image' || property?.schema?.format?.includes('image')"
          class="w-3/5 break-words py-4 pr-6 pl-3 text-neutrals-dark"
        >
          <img :src="property?.value" class="h-20 w-20" />
        </td>
        <td v-else class="w-3/5 break-words py-4 pr-6 pl-3 text-neutrals-dark">
          {{ property?.value }}
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
export default {
  name: 'CredentialDetailsTableComponent',
  props: {
    heading: {
      type: String,
      default: '',
    },
    credential: {
      type: Object,
      required: true,
    },
  },
};
</script>
