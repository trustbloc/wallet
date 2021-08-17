<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <a class="cursor-pointer" @click="handleLocaleSwitch(newLocale)">
    <!-- TODO: remove span and text color after sass styles are removed -->
    <span class="text-neutrals-medium">{{ newLocale.name }}</span>
  </a>
</template>

<script>
import supportedLocales from '@/config/supportedLocales';
import store from '@/store';
import { updateI18nLocale } from '@/plugins/i18n';

export default {
  name: 'LocaleSwitcher',
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
          name: this.$router.history.current.name,
          params: { ...this.$router.history.current.params, locale: newLocale.base },
        });
      }
    },
  },
};
</script>
