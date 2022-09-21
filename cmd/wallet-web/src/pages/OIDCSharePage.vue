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
import { mapGetters } from 'vuex';
import { useI18n } from 'vue-i18n';
import { decode, encode } from 'js-base64';
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { OIDCMutations } from '@/layouts/OIDCLayout.vue';
import { OIDCShareLayoutMutations } from '@/layouts/OIDCShareLayout.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import CredentialBannerComponent from '@/components/WACI/CredentialBannerComponent.vue';
import CredentialOverviewComponent from '@/components/WACI/CredentialOverviewComponent.vue';
import WACIActionButtonsContainerComponent from '@/components/WACI/WACIActionButtonsContainerComponent.vue';
import WACICredentialsMissingComponent from '@/components/WACI/WACICredentialsMissingComponent.vue';
import WACIErrorComponent from '@/components/WACI/WACIErrorComponent.vue';
import WACILoadingComponent from '@/components/WACI/WACILoadingComponent.vue';
import CredentialDetailsTableComponent from '@/components/WACI/CredentialDetailsTableComponent.vue';
import OIDCShareOverviewPage from '@/pages/OIDCShareOverviewPage.vue';

const isBase64Param = (param) => {
  if (!param) {
    return false;
  }
  try {
    return btoa(atob(param)) === param;
  } catch (error) {
    return false;
  }
};

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
      noCredentialFound: false,
      sharing: false,
      processedCredentials: [],
      token: null,
    };
  },
  computed: {
    showMainState() {
      return !this.loading && !this.errors.length && !this.sharing && !this.noCredentialFound;
    },
  },
  created: async function () {
    this.loading = true;
    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    this.collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    const extractClaimsFromQuery = (claims) => {
      let decodedClaims;

      if (isBase64Param(claims)) {
        decodedClaims = JSON.parse(decode(claims));
      } else {
        try {
          decodedClaims = JSON.parse(claims);
        } catch (error) {
          decodedClaims = JSON.parse(JSON.stringify(claims));
        }
      }

      return decodedClaims;
    };
    this.claims = extractClaimsFromQuery(this.$route.query.claims);

    //initiate credential share flow.
    try {
      const { results } = await this.credentialManager.query(this.token, [
        {
          type: 'PresentationExchange',
          credentialQuery: [this.claims.vp_token.presentation_definition],
        },
      ]);
      this.presentations = results;
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

    this.prepareRecords(this.presentations);
    OIDCMutations.setProcessedCredentials(this.processedCredentials);
    // TODO: get Requestor
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
      this.generateIdToken();
      this.generateVPToken();
    },
    share() {
      this.sharing = true;
      let ack;

      try {
        ack = {
          status: 'OK',
          url: `${this.$route.query.redirect_uri}?state=${this.$route.query.state}&id_token=${this.idToken}&vp_token=${this.vpToken}`,
        };
      } catch (e) {
        this.errors.push(e);
        console.error('share credentials failed,', e);
        this.sharing = false;
        return;
      }

      const { status, url } = ack;
      this.redirectUrl = url;

      if (status === 'OK') this.finish();

      this.sharing = false;
    },
    finish() {
      window.location.href = this.redirectUrl;
    },
    cancel() {
      window.location = window.location.origin;
    },
    generateIdToken() {
      const header = JSON.stringify({
        alg: 'none',
      });
      const payload = JSON.stringify({
        iss: this.$route.query.client_id,
        sub: this.$route.query.client_id,
        aud: this.$route.query.client_id,
        iat: Date.now(),
        exp: Date.now(),
      });
      const signature = JSON.stringify({});
      const encodedHeader = encode(header).slice(0, -1);
      const encodedPayload = encode(payload).slice(0, -1);
      const encodedSignature = encode(signature).slice(0, -1);
      this.idToken = `${encodedHeader}.${encodedPayload}.${encodedSignature}`;
    },
    async generateVPToken() {
      const { controller } = this.getCurrentUser().preference;
      try {
        const { presentation } = await this.credentialManager.present(
          this.token,
          { rawCredentials: this.presentations[0].verifiableCredential },
          { controller }
        );
        console.log('presentation', presentation);
        this.vpToken = encodeURIComponent(JSON.stringify(presentation));
        console.log('vpToken', this.vpToken);
      } catch (e) {
        console.error('error sharing a credential:', e);
      }
    },
    handleOverviewClick: function (id) {
      OIDCMutations.setSelectedCredentialId(id);
      OIDCShareLayoutMutations.setComponent(OIDCShareOverviewPage);
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
