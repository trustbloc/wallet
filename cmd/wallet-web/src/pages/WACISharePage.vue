<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

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
<script>
import { toRaw } from 'vue';
import { mapGetters } from 'vuex';
import { useI18n } from 'vue-i18n';
import { CollectionManager, CredentialManager, DIDComm } from '@trustbloc/wallet-sdk';
import { WACIMutations, WACIStore } from '@/layouts/WACILayout.vue';
import { WACIShareLayoutMutations } from '@/layouts/WACIShareLayout.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import CredentialBannerComponent from '@/components/WACI/CredentialBannerComponent.vue';
import CredentialOverviewComponent from '@/components/WACI/CredentialOverviewComponent.vue';
import WACIActionButtonsContainerComponent from '@/components/WACI/WACIActionButtonsContainerComponent.vue';
import WACICredentialsMissingComponent from '@/components/WACI/WACICredentialsMissingComponent.vue';
import WACIErrorComponent from '@/components/WACI/WACIErrorComponent.vue';
import WACILoadingComponent from '@/components/WACI/WACILoadingComponent.vue';
import WACISuccessComponent from '@/components/WACI/WACISuccessComponent.vue';
import CredentialDetailsTableComponent from '@/components/WACI/CredentialDetailsTableComponent.vue';
import WACIShareOverviewPage from '@/pages/WACIShareOverviewPage.vue';

export default {
  components: {
    CredentialBannerComponent,
    CredentialDetailsTableComponent,
    CredentialOverviewComponent,
    StyledButtonComponent,
    WACIActionButtonsContainerComponent,
    WACICredentialsMissingComponent,
    WACIErrorComponent,
    WACILoadingComponent,
    WACISuccessComponent,
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      errors: [],
      requestOrigin: '',
      loading: true,
      sharing: false,
      processedCredentials: [],
      sharedSuccessfully: false,
      noCredentialFound: false,
    };
  },
  computed: {
    showMainState() {
      return (
        !this.loading &&
        !this.errors.length &&
        !this.sharing &&
        !this.noCredentialFound &&
        !this.sharedSuccessfully
      );
    },
  },
  created: async function () {
    this.loading = true;
    this.protocolHandler = WACIStore.protocolHandler;
    const invitation = toRaw(this.protocolHandler.message());
    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    //initiate credential share flow.
    this.didcomm = new DIDComm({ agent: this.getAgentInstance(), user });
    try {
      const { threadID, presentations } = await this.didcomm.initiateCredentialShare(
        token,
        invitation,
        {
          userAnyRouterConnection: true,
        }
      );
      this.threadID = threadID;
      this.presentations = presentations;
    } catch (e) {
      if (!e.message.includes('12009')) {
        this.errors.push('Error initiating credential share');
      }
      console.error('initiating credential share failed,', e);
      // Error code 12009 is for no result found message
      this.noCredentialFound = true;
      this.loading = false;
      return;
    }

    await this.prepareRecords(this.presentations);
    this.requestOrigin = this.protocolHandler.requestor();
    WACIMutations.setProcessedCredentials(this.processedCredentials);
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    prepareRecords: async function (presentations) {
      try {
        const credentials = presentations.reduce(
          (acc, val) => acc.concat(val.verifiableCredential),
          []
        );
        await credentials.map(async (credential) => {
          const { id, collection, name, issuanceDate, resolved } =
            await this.credentialManager.getCredentialMetadata(this.token, credential.id);
          const {
            content: { name: vaultName },
          } = await this.collectionManager.get(this.token, collection);
          this.processedCredentials.push({
            id,
            name: name || resolved[0].title,
            issuanceDate,
            ...resolved[0],
            vaultName,
          });
        });
      } catch (e) {
        this.errors.push('No credentials found matching requested criteria.');
        console.error('get credentials failed,', e);
        this.loading = false;
      }
    },
    async share() {
      this.sharing = true;
      const { profile, preference } = this.getCurrentUser();
      const { controller, proofType, verificationMethod } = preference;

      let ack;

      try {
        ack = await this.didcomm.completeCredentialShare(
          profile.token,
          this.threadID,
          this.presentations,
          {
            controller,
            proofType,
            verificationMethod,
          },
          { waitForDone: true }
        );
      } catch (e) {
        this.errors.push(e);
        console.error('share credentials failed,', e);
        this.sharing = false;
        return;
      }

      const { status, url } = ack;
      this.redirectUrl = url;
      // TODO check if status="FAIL", then should redirect to generic error screen, it means WACI flow didn't succeed
      if (status === 'OK') this.sharedSuccessfully = true;

      this.sharing = false;
    },
    finish() {
      this.protocolHandler.done(this.redirectUrl ? this.redirectUrl : window.location.origin);
    },
    cancel() {
      this.protocolHandler.cancel();
    },
    handleOverviewClick: function (id) {
      WACIMutations.setSelectedCredentialId(id);
      WACIShareLayoutMutations.setComponent(WACIShareOverviewPage);
    },
  },
};
</script>
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
