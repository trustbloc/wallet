<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <a
    tabindex="0"
    class="cursor-pointer"
    @click="handleLocaleSwitch(newLocale)"
    @keyup.enter="handleLocaleSwitch(newLocale)"
  >
    {{ newLocale.name }}
  </a>
</template>

<script>
import supportedLocales from '@/config/supportedLocales';
import store from '@/store';
import { updateI18nLocale } from '@/plugins/i18n';

export default {
  name: 'LocaleSwitcherComponent',
  computed: {
    newLocale: function () {
      return store.getters.getLocale.id === 'en'
        ? supportedLocales.find((loc) => loc.id === 'fr')
        : supportedLocales.find((loc) => loc.id === 'en');
    },
  },
  methods: {
    async handleLocaleSwitch(newLocale) {
      if (this.$i18n.locale !== newLocale || store.getters.getLocale.id !== newLocale) {
        await updateI18nLocale(newLocale.id);
        store.dispatch('setLocale', newLocale);
        this.$router.replace({
          name: this.$router.currentRoute._value.name,
          params: { ...this.$router.currentRoute._value.params, locale: newLocale.base },
        });
      }
    },
  },
};
</script>
