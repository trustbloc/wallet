<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="py-6 px-15 mx-6 w-full text-center">
    <div class="hidden md:block mx-auto text-center text-neutrals-medium">
      © {{ date }} TrustBloc ･ {{ i18n.privacyPolicy }} ･ {{ i18n.terms }} ･
      <LocaleSwitcher />
    </div>
    <div class="block md:hidden text-neutrals-medium">
      <ul class="list-inside">
        <li>
          <span class="whitespace-nowrap">{{ i18n.privacyPolicy }}</span> ･
          <span class="whitespace-nowrap">{{ i18n.terms }}</span>
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

export default {
  name: 'ContentFooter',
  components: {
    LocaleSwitcher,
  },
  data() {
    return {
      date: new Date().getFullYear(),
    };
  },
  computed: {
    i18n() {
      return this.$t('ContentFooter');
    },
  },
};
</script>
