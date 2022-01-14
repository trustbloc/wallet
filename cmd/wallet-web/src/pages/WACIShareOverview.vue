<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex overflow-hidden relative justify-start items-start w-full h-full">
    <div class="flex overflow-hidden flex-col flex-grow justify-between items-center w-full h-full">
      <div class="flex overflow-auto justify-center w-full">
        <div
          class="
            flex-grow
            justify-start
            items-start
            pt-8
            pr-5
            pb-8
            pl-5
            md:pr-0 md:pl-0
            w-full
            max-w-3xl
            h-full
            flex flex-col
          "
        >
          <span class="mb-6 text-3xl font-bold">{{ t('CHAPI.Share.shareCredential', 2) }}</span>
          <credential-overview
            class="mb-5 waci-share-credential-overview-root"
            :credential="credential"
          >
            <template #bannerBottomContainer>
              <div
                class="
                  absolute
                  justify-start
                  items-start
                  px-5
                  pb-3
                  w-full
                  rounded-b-xl
                  pt-13
                  bg-neutrals-white
                  flex flex-row
                  waci-share-credential-overview-vault
                "
              >
                <span class="flex text-sm font-bold text-neutrals-dark">
                  {{ t('CredentialDetails.Banner.vault') }}
                </span>
                <span class="flex ml-3 text-sm text-neutrals-medium">
                  {{ credential?.vaultName }}
                </span>
              </div>
            </template>
            <template #credentialDetails>
              <credential-details-table
                :heading="t('WACI.Share.whatIsShared')"
                :credential="credential"
              />
            </template>
          </credential-overview>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { WACIGetters, WACIStore } from '@/layouts/WACI.vue';
import CredentialOverview from '@/components/WACI/CredentialOverview.vue';
import CredentialDetailsTable from '@/components/WACI/CredentialDetailsTable.vue';

export default {
  components: {
    CredentialOverview,
    CredentialDetailsTable,
  },
  setup() {
    const id = ref(WACIStore.selectedCredentialId);
    const credential = ref(WACIGetters.getProcessedCredentialById(id.value));
    watch(
      () => WACIStore.selectedCredentialId,
      () => {
        id.value = WACIStore.selectedCredentialId;
        credential.value = WACIGetters.getProcessedCredentialById(id.value);
      }
    );
    const { t } = useI18n();
    return { credential, t };
  },
};
</script>
