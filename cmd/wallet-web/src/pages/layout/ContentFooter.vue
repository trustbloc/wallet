/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="w-full mx-6 px-15 py-6 text-center">
    <div class="text-neutrals-medium mx-auto md:block hidden text-center">
       © {{ date }} TrustBloc ･ {{ i18n.privacyPolicy }} ･ {{ i18n.terms }}  ･
      <select v-model="$i18n.locale" class="py-4" style="background-color: transparent">
        <option value="en" class="text-neutrals-dark block text-sm" id="english">English</option>
        <option value="fr" class="text-neutrals-medium block text-sm" id="french">Français</option>
      </select>
    </div>
    <div class="text-neutrals-medium md:hidden block">
      <ul class="list-inside">
        <li><span class="whitespace-nowrap">{{ i18n.privacyPolicy }}</span> ･ <span class="whitespace-nowrap">{{ i18n.terms }}</span></li>
        <li><span class="whitespace-nowrap">© {{ date }} TrustBloc</span> ･
          <select class="py-4" style="background-color: transparent">
            <option class="text-neutrals-dark block text-sm">English</option>
            <option class="text-neutrals-medium block text-sm">Français</option>
          </select>
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
export default {
  name: "content-footer",
  computed: {
		i18n() {
			return this.$t("ContentFooter")
		}
	},
  data() {
    return {
      date: new Date().getFullYear(),
    };
  },
};
</script>
