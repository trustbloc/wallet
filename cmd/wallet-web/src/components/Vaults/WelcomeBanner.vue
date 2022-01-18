<template>
  <div
    class="
      justify-start
      p-6
      md:py-8 md:px-10 md:rounded-xl md:border
      border-b
      flex flex-col
      border-opacity-25 border-primary-purple
    "
  >
    <div class="flex flex-grow justify-between items-start mb-8 w-full">
      <div class="flex flex-col justify-start items-start pr-6">
        <span class="text-lg font-bold text-neutrals-dark">{{
          t('Vaults.WelcomeBanner.heading')
        }}</span>
        <span class="text-base text-neutrals-medium">
          {{ t('Vaults.WelcomeBanner.subheading') }}
        </span>
      </div>
      <button v-bind="$attrs">
        <img src="@/assets/img/Icons-sm--close-icon.svg" class="w-6 max-w-none h-6" />
      </button>
    </div>
    <div class="grid grid-flow-row lg:grid-flow-col gap-x-12 gap-y-6 h-full">
      <slot />
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';

export default {
  name: 'WelcomeBanner',
  setup() {
    const { t } = useI18n();
    return { t };
  },
};
</script>
