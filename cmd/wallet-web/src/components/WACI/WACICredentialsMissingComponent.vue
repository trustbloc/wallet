<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-center py-16 px-5">
    <img src="@/assets/img/icons-error.svg" />
    <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
      t('CHAPI.Share.CredentialsMissing.heading')
    }}</span>
    <span class="text-lg text-center text-neutrals-medium">{{
      t('CHAPI.Share.CredentialsMissing.body')
    }}</span>
    <StyledButtonComponent
      v-bind="$attrs"
      id="share-credentials-ok-btn"
      type="btn-outline"
      class="mt-6"
    >
      {{ t('CHAPI.Share.CredentialsMissing.ok') }}
    </StyledButtonComponent>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';

export default {
  name: 'WACICredentialsMissingComponent',
  components: {
    StyledButtonComponent,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
};
</script>
