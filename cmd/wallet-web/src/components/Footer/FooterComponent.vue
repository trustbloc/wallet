<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="mx-6 w-full py-6 px-15 text-center">
    <div class="mx-auto hidden text-center text-neutrals-medium md:block">
      © {{ date }} TrustBloc ･ {{ t('Footer.privacyPolicy') }} ･ {{ t('Footer.terms') }} ･
      <LocaleSwitcherComponent />
    </div>
    <div class="block text-neutrals-medium md:hidden">
      <ul class="list-inside">
        <li>
          <span class="whitespace-nowrap">{{ t('Footer.privacyPolicy') }}</span> ･
          <span class="whitespace-nowrap">{{ t('Footer.terms') }}</span>
        </li>
        <li>
          <span class="whitespace-nowrap">© {{ date }} TrustBloc</span> ･
          <LocaleSwitcherComponent />
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
import LocaleSwitcherComponent from '@/components/LocaleSwitcher/LocaleSwitcherComponent.vue';
import { useI18n } from 'vue-i18n';

export default {
  name: 'FooterComponent',
  components: {
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
};
</script>
