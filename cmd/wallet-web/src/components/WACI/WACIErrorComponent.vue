<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex h-full w-full grow flex-col items-center justify-center">
    <img src="@/assets/img/icons-error.svg" />
    <span class="mt-6 mb-3 text-center text-xl font-bold text-neutrals-dark">{{
      t('CHAPI.Share.Error.heading')
    }}</span>
    <span class="w-full text-center text-lg text-neutrals-medium">{{
      t('CHAPI.Share.Error.body')
    }}</span>
    <StyledButtonComponent
      v-bind="$attrs"
      id="share-credentials-ok-btn"
      type="btn-primary"
      class="mt-6"
    >
      {{ t('CHAPI.Share.Error.tryAgain') }}
    </StyledButtonComponent>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';

export default {
  name: 'WACIErrorComponent',
  components: {
    StyledButtonComponent,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
};
</script>
