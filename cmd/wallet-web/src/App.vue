<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="font-sans">
    <div v-if="!$root.loaded" class="loader">
      <Spinner />
      <span>Loading Agent...</span>
    </div>
    <router-view v-else />
  </div>
</template>

<script>
import Spinner from '@/components/Spinner/Spinner.vue';
import EventBus from '@/EventBus';

export default {
  components: { Spinner },
  data: () => ({
    isLoading: true,
  }),
  mounted() {
    EventBus.$on('i18n-load-start', () => (this.isLoading = true));
    EventBus.$on('i18n-load-complete', () => (this.isLoading = false));
  },
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
