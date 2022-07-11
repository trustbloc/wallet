<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex overflow-hidden relative justify-start items-start w-full h-full">
    <div class="flex overflow-hidden flex-col grow justify-between items-center w-full h-full">
      <div class="flex overflow-auto justify-center w-full">
        <div
          class="flex flex-col grow justify-start items-start py-8 px-5 w-full max-w-3xl h-full md:px-0"
        >
          <span class="mb-6 text-3xl font-bold">{{ t('CHAPI.Share.shareCredential', 2) }}</span>
          <CredentialOverviewComponent
            class="waci-share-credential-overview-root mb-5"
            :credential="credential"
          >
            <template #bannerBottomContainer>
              <div
                class="waci-share-credential-overview-vault absolute flex w-full flex-row items-start justify-start rounded-b-xl bg-neutrals-white px-5 pt-13 pb-3"
              >
                <span class="flex text-sm font-bold text-neutrals-dark">
                  {{ t('CredentialDetails.Banner.vault') }}
                </span>
                <span class="ml-3 flex text-sm text-neutrals-medium">
                  {{ credential?.vaultName }}
                </span>
              </div>
            </template>
            <template #credentialDetails>
              <CredentialDetailsTableComponent
                :heading="t('WACI.Share.whatIsShared')"
                :credential="credential"
              />
            </template>
          </CredentialOverviewComponent>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { OIDCGetters, OIDCStore } from '@/layouts/OIDCLayout.vue';
import CredentialOverviewComponent from '@/components/WACI/CredentialOverviewComponent.vue'; // TODO: Change the route to something more generic
import CredentialDetailsTableComponent from '@/components/WACI/CredentialDetailsTableComponent.vue'; // TODO: Change the route to something more generic

export default {
  components: {
    CredentialOverviewComponent,
    CredentialDetailsTableComponent,
  },
  setup() {
    const id = ref(OIDCStore.selectedCredentialId);
    const credential = ref(OIDCGetters.getProcessedCredentialById(id.value));
    watch(
      () => OIDCStore.selectedCredentialId,
      () => {
        id.value = OIDCStore.selectedCredentialId;
        credential.value = OIDCGetters.getProcessedCredentialById(id.value);
      }
    );
    const { t } = useI18n();
    return { credential, t };
  },
};
</script>
