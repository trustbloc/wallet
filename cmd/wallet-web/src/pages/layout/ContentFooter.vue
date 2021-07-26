/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <div class="py-6 px-15 mx-6 w-full text-center">
    <div class="hidden md:block mx-auto text-center text-neutrals-medium">
      © {{ date }} TrustBloc ･ {{ i18n.privacyPolicy }} ･ {{ i18n.terms }} ･
      <select v-model="$i18n.locale" class="py-4" style="background-color: transparent">
        <option id="english" value="en" class="block text-sm text-neutrals-dark">English</option>
        <option id="french" value="fr" class="block text-sm text-neutrals-medium">Français</option>
      </select>
    </div>
    <div class="block md:hidden text-neutrals-medium">
      <ul class="list-inside">
        <li>
          <span class="whitespace-nowrap">{{ i18n.privacyPolicy }}</span> ･
          <span class="whitespace-nowrap">{{ i18n.terms }}</span>
        </li>
        <li>
          <span class="whitespace-nowrap">© {{ date }} TrustBloc</span> ･
          <select class="py-4" style="background-color: transparent">
            <option class="block text-sm text-neutrals-dark">English</option>
            <option class="block text-sm text-neutrals-medium">Français</option>
          </select>
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
export default {
  name: 'ContentFooter',
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
