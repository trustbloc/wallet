<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div
    class="flex flex-col justify-start p-6 border-b border-primary-purple border-opacity-25 md:py-8 md:px-10 md:rounded-xl md:border"
  >
    <div class="flex grow justify-between items-start mb-8 w-full">
      <div class="flex flex-col justify-start items-start pr-6">
        <span class="text-lg font-bold text-neutrals-dark">{{
          t('Vaults.WelcomeBanner.heading')
        }}</span>
        <span class="text-base text-neutrals-medium">
          {{ t('Vaults.WelcomeBanner.subheading') }}
        </span>
      </div>
      <button v-bind="$attrs">
        <img src="@/assets/img/Icons-sm--close-icon.svg" class="h-6 w-6 max-w-none" />
      </button>
    </div>
    <div class="grid grid-flow-row gap-x-12 gap-y-6 h-full lg:grid-flow-col">
      <slot />
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';

export default {
  name: 'VaultWelcomeBannerComponent',
  setup() {
    const { t } = useI18n();
    return { t };
  },
};
</script>
