<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
    <img src="@/assets/img/icons-error.svg" />
    <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
      t('CHAPI.Share.CredentialsMissing.heading')
    }}</span>
    <span class="text-lg text-center text-neutrals-medium">{{
      t('CHAPI.Share.CredentialsMissing.body')
    }}</span>
    <styled-button v-bind="$attrs" id="share-credentials-ok-btn" type="btn-outline" class="mt-6">
      {{ t('CHAPI.Share.CredentialsMissing.ok') }}
    </styled-button>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';
import StyledButton from '@/components/StyledButton/StyledButton.vue';

export default {
  name: 'WACICredentialsMissing',
  components: {
    StyledButton,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
};
</script>
