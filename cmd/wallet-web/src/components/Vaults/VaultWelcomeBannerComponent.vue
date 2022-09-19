<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div
    class="flex flex-col justify-start border-b border-primary-purple border-opacity-25 p-6 md:rounded-xl md:border md:py-8 md:px-10"
  >
    <div class="mb-8 flex w-full grow items-start justify-between">
      <div class="flex flex-col items-start justify-start pr-6">
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
    <div class="grid h-full grid-flow-row gap-x-12 gap-y-6 lg:grid-flow-col">
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
