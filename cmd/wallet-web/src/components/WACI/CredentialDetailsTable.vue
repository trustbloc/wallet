<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-start w-full">
    <span v-if="heading.length" class="py-4 pl-3 text-xl font-bold text-neutrals-dark">{{
      heading
    }}</span>
    <table class="w-full border-t border-neutrals-chatelle">
      <tr
        v-for="(property, key) of credential?.properties"
        :key="key"
        class="border-b border-neutrals-thistle border-dotted"
      >
        <td class="py-4 pr-6 pl-3 w-2/5 whitespace-nowrap text-neutrals-medium">
          {{ property?.label }}
        </td>
        <td
          v-if="property?.type === 'image' || property?.schema?.format?.includes('image')"
          class="py-4 pr-6 pl-3 w-3/5 break-words text-neutrals-dark"
        >
          <img :src="property?.value" class="w-20 h-20" />
        </td>
        <td v-else class="py-4 pr-6 pl-3 w-3/5 break-words text-neutrals-dark">
          {{ property?.value }}
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';

export default {
  name: 'CredentialDetailsTable',
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
  setup() {
    const { t } = useI18n();
    return { t };
  },
};
</script>
