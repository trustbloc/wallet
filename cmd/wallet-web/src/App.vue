/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="font-sans">
    <div class="loader" v-if="!$root.loaded">
      <md-progress-spinner :md-diameter="100" :md-stroke="10" md-mode="indeterminate"></md-progress-spinner>
      <div>Loading Agent...</div>
    </div>
    <router-view v-if="$root.loaded"></router-view>
  </div>

</template>

<script>
import { setDocumentLang, setDocumentTitle } from "@/utils/i18n/document";
export default {
  mounted() {
    this.$watch(
      "$i18n.locale",
      (newLocale, oldLocale) => {
        if (newLocale === oldLocale) {
          return
        }
        setDocumentLang(newLocale)
        setDocumentTitle(this.$t("App.title"))
      },
      { immediate: true }
    )
  }
};
</script>

<style scoped>
.loader {
  width: 100%;
  height: 100%;
  position: absolute;
  left: 50%;
  top: 30%;
  margin-left: -4em;
}
</style>
