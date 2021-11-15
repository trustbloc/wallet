<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="py-6 px-15 mx-6 w-full text-center">
    <div class="hidden md:block mx-auto text-center text-neutrals-medium">
      © {{ date }} TrustBloc ･ {{ t('Footer.privacyPolicy') }} ･ {{ t('Footer.terms') }} ･
      <LocaleSwitcher />
    </div>
    <div class="block md:hidden text-neutrals-medium">
      <ul class="list-inside">
        <li>
          <span class="whitespace-nowrap">{{ t('Footer.privacyPolicy') }}</span> ･
          <span class="whitespace-nowrap">{{ t('Footer.terms') }}</span>
        </li>
        <li>
          <span class="whitespace-nowrap">© {{ date }} TrustBloc</span> ･
          <LocaleSwitcher />
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
import LocaleSwitcher from '@/components/LocaleSwitcher/LocaleSwitcher.vue';
import { useI18n } from 'vue-i18n';

export default {
  name: 'Footer',
  components: {
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
};
</script>
