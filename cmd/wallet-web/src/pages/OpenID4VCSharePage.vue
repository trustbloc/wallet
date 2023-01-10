<!--
 * Copyright Avast Software. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { computed, onMounted, ref, toRaw } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { useI18n } from 'vue-i18n';
import { parseJWTVC } from '@/mixins';
import { CollectionManager, CredentialManager, DIDManager, OpenID4VP } from '@trustbloc/wallet-sdk';
import { OpenID4VCMutations } from '@/layouts/OpenID4VCLayout.vue';
import { OpenID4VCShareLayoutMutations } from '@/layouts/OpenID4VCShareLayout.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import CredentialBannerComponent from '@/components/WACI/CredentialBannerComponent.vue';
import CredentialOverviewComponent from '@/components/WACI/CredentialOverviewComponent.vue';
import WACIActionButtonsContainerComponent from '@/components/WACI/WACIActionButtonsContainerComponent.vue';
import WACICredentialsMissingComponent from '@/components/WACI/WACICredentialsMissingComponent.vue';
import WACIErrorComponent from '@/components/WACI/WACIErrorComponent.vue';
import WACILoadingComponent from '@/components/WACI/WACILoadingComponent.vue';
import WACISuccessComponent from '@/components/WACI/WACISuccessComponent.vue';
import CredentialDetailsTableComponent from '@/components/WACI/CredentialDetailsTableComponent.vue';
import OIDCShareOverviewPage from '@/pages/OIDCShareOverviewPage.vue';

// Hooks
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const store = useStore();

// Local Variables
const errors = ref([]);
const requestOrigin = ref();
const loading = ref(true);
const noCredentialFound = ref(false);
const sharing = ref(false);
const sharedSuccessfully = ref(false);
const processedCredentials = ref([]);
const token = ref(null);
const showMainState = computed(
  () =>
    !loading.value &&
    !errors.value.length &&
    !sharing.value &&
    !noCredentialFound.value &&
    !sharedSuccessfully.value
);
const credentialManager = ref(null);
const collectionManager = ref(null);
const didManager = ref(null);
const openID4VP = ref(null);
const presentation = ref(null);

// Store Getters
const currentUser = computed(() => store.getters['getCurrentUser']);
const agentInstance = computed(() => store.getters['agent/getInstance']);

// Methods
async function prepareRecords(presentations) {
  try {
    const credentials = presentations.reduce(
      (acc, val) => acc.concat(val.verifiableCredential),
      []
    );
    await credentials.map(async (credential) => {
      const cred = parseJWTVC(credential);
      const { id, collection, name, issuanceDate, resolved } =
        await credentialManager.value.getCredentialMetadata(token.value, cred.id);
      const {
        content: { name: vaultName },
      } = await collectionManager.value.get(token.value, collection);
      processedCredentials.value.push({
        id,
        name: name || resolved[0].title,
        issuanceDate,
        ...resolved[0],
        vaultName,
      });
    });
  } catch (e) {
    errors.value.push('No credentials found matching requested criteria.');
    console.error('get credentials failed,', e);
    loading.value = false;
  }
}
async function share() {
  sharing.value = true;
  const { contents } = await didManager.value.getAllDIDs(token.value);
  // Selecting the first verification method to be used for now
  // since selection logic is not defined in the current openid4vp spec
  const kid = Object.values(contents)[0].didDocument.verificationMethod[0].id;
  await openID4VP.value
    .submitOIDCPresentation({
      authToken: token.value,
      kid,
      presentation: toRaw(presentation.value),
      expiry: Math.floor(Date.now() / 1000 + 60 * 10),
    })
    .then(() => {
      sharedSuccessfully.value = true;
    })
    .catch((e) => {
      console.error(e);
      errors.value.push('presentation submission failed:', e);
    });
  sharing.value = false;
}
function finish() {
  router.push('/credentials');
}
function cancel() {
  window.location = window.location.origin;
}
function handleOverviewClick(id) {
  OpenID4VCMutations.setSelectedCredentialId(id);
  OpenID4VCShareLayoutMutations.setComponent(OIDCShareOverviewPage);
}

onMounted(async () => {
  const { profile } = currentUser.value;
  const { user } = profile;
  token.value = profile.token;

  credentialManager.value = new CredentialManager({ agent: agentInstance.value, user });
  collectionManager.value = new CollectionManager({ agent: agentInstance.value, user });
  didManager.value = new DIDManager({ agent: agentInstance.value, user });
  openID4VP.value = new OpenID4VP({ agent: agentInstance.value, user });

  presentation.value = await openID4VP.value.initiateOIDCPresentation({
    authToken: token.value,
    url: route.query.url,
  });

  await prepareRecords(presentation.value);
  OpenID4VCMutations.setProcessedCredentials(processedCredentials.value);
  // TODO: get Requestor
  loading.value = false;
});
</script>

<template>
  <div v-if="!showMainState" class="flex h-full w-full grow flex-col items-center justify-center">
    <!-- Loading State -->
    <WACILoadingComponent v-if="loading" />

    <!-- Sharing State -->
    <WACILoadingComponent
      v-else-if="sharing"
      :message="t('CHAPI.Share.sharingCredential', processedCredentials.length)"
    />

    <!-- Error State -->
    <WACIErrorComponent v-else-if="errors.length" @click="cancel" />

    <!-- Credentials Missing State -->
    <WACICredentialsMissingComponent v-else-if="noCredentialFound" @click="cancel" />

    <!-- Success State -->
    <WACISuccessComponent
      v-else-if="sharedSuccessfully"
      id="share-credentials-ok-btn"
      :heading="t('WACI.Share.success', processedCredentials.length)"
      :message="
        t(
          'WACI.Share.message',
          {
            subject:
              processedCredentials.length === 1
                ? processedCredentials[0].title
                : processedCredentials.length,
            // TODO: issue-1055 Read meta data from external urls
            requestor: 'Requestor',
          },
          processedCredentials.length
        )
      "
      :button-label="t('WACI.ok')"
      @click="finish"
    />
  </div>

  <!-- Main State -->
  <div v-else class="flex h-full w-full grow flex-col items-center justify-between overflow-hidden">
    <div class="flex w-full justify-center overflow-auto">
      <div
        class="flex h-full w-full max-w-3xl grow flex-col items-start justify-start py-8 px-5 md:px-0"
      >
        <span class="mb-6 text-3xl font-bold">{{
          t('CHAPI.Share.shareCredential', processedCredentials.length)
        }}</span>
        <div class="mb-4 flex w-full flex-row items-start justify-start">
          <div class="h-12 w-12 flex-none border-opacity-10">
            <!-- todo issue-1055 Read meta data from external urls -->
            <img src="@/assets/img/generic-issuer-icon.svg" />
          </div>
          <div class="flex flex-col pl-3">
            <span class="mb-1 flex-1 text-ellipsis text-left text-sm font-bold text-neutrals-dark">
              <!-- todo issue-1055 Read meta data from external urls -->
              Requestor
            </span>
            <div class="flex flex-row items-center justify-center">
              <img src="@/assets/img/small-lock-icon.svg" />
              <span class="flex-1 text-ellipsis pl-1 text-left text-xs text-neutrals-medium">
                {{ requestOrigin }}
              </span>
            </div>
          </div>
        </div>

        <span class="text-sm text-neutrals-dark">{{
          t('CHAPI.Share.headline', { issuer: 'Requestor' }, processedCredentials.length)
        }}</span>

        <!-- Single Credential Overview (with details) -->
        <CredentialOverviewComponent
          v-if="processedCredentials.length === 1"
          class="waci-share-credential-overview-root my-5"
          :credential="processedCredentials[0]"
        >
          <template #bannerBottomContainer>
            <div
              class="waci-share-credential-overview-vault absolute flex w-full flex-row items-start justify-start rounded-b-xl bg-neutrals-white px-5 pt-13 pb-3"
            >
              <span class="flex text-sm font-bold text-neutrals-dark">
                {{ t('CredentialDetails.Banner.vault') }}
              </span>
              <span class="ml-3 flex text-sm text-neutrals-medium">
                {{ processedCredentials[0].vaultName }}
              </span>
            </div>
          </template>
          <template #credentialDetails>
            <CredentialDetailsTableComponent
              :heading="t('WACI.Share.whatIsShared')"
              :credential="processedCredentials[0]"
            />
          </template>
        </CredentialOverviewComponent>

        <!-- List of Credential Banners (Links to Details for each) -->
        <ul v-else-if="processedCredentials.length > 1" class="mt-6 w-full space-y-5">
          <li v-for="(credential, index) in processedCredentials" :key="index">
            <CredentialBannerComponent
              :id="credential.id"
              :styles="credential.styles"
              :title="credential.name"
              @click="handleOverviewClick(credential.id)"
            />
          </li>
        </ul>
      </div>
    </div>

    <WACIActionButtonsContainerComponent>
      <template #leftButton>
        <StyledButtonComponent id="cancelBtn" type="btn-outline" @click="cancel">
          {{ t('CHAPI.Share.decline') }}
        </StyledButtonComponent>
      </template>
      <template #rightButton>
        <StyledButtonComponent id="share-credentials" type="btn-primary" @click="share">
          {{ t('CHAPI.Share.share') }}
        </StyledButtonComponent>
      </template>
    </WACIActionButtonsContainerComponent>
  </div>
</template>

<style>
.waci-share-credential-overview-root {
  padding-bottom: 2.5rem;
}
.waci-share-credential-overview-vault {
  top: 2.5rem;
  left: 0;
  /* TODO: replace with tailwind shadow once defined in config */
  box-shadow: 0px 2px 12px 0px rgba(25, 12, 33, 0.1);
}
</style>
