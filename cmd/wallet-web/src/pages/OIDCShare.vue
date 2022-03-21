<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    v-if="!showMainState"
    class="flex flex-col flex-grow justify-center items-center w-full h-full"
  >
    <!-- Loading State -->
    <WACI-loading v-if="loading" />

    <!-- Sharing State -->
    <WACI-loading
      v-else-if="sharing"
      :message="t('CHAPI.Share.sharingCredential', processedCredentials.length)"
    />

    <!-- Error State -->
    <WACI-error v-else-if="errors.length" @click="cancel" />

    <!-- Credentials Missing State -->
    <WACI-credentials-missing v-else-if="showCredentialsMissing" @click="cancel" />

    <!-- Success State -->
    <WACI-success
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
  <div
    v-else
    class="flex overflow-hidden flex-col flex-grow justify-between items-center w-full h-full"
  >
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
        <span class="mb-6 text-3xl font-bold">{{
          t('CHAPI.Share.shareCredential', processedCredentials.length)
        }}</span>
        <div class="flex flex-row justify-start items-start mb-4 w-full">
          <div class="flex-none w-12 h-12 border-opacity-10">
            <!-- todo issue-1055 Read meta data from external urls -->
            <img src="@/assets/img/generic-issuer-icon.svg" />
          </div>
          <div class="flex flex-col pl-3">
            <span class="flex-1 mb-1 text-sm font-bold text-left text-neutrals-dark text-ellipsis">
              <!-- todo issue-1055 Read meta data from external urls -->
              Requestor
            </span>
            <div class="flex flex-row justify-center items-center">
              <img src="@/assets/img/small-lock-icon.svg" />
              <span class="flex-1 pl-1 text-xs text-left text-neutrals-medium text-ellipsis">
                {{ requestOrigin }}
              </span>
            </div>
          </div>
        </div>

        <span class="text-sm text-neutrals-dark">{{
          t('CHAPI.Share.headline', { issuer: 'Requestor' }, processedCredentials.length)
        }}</span>

        <!-- Single Credential Overview (with details) -->
        <credential-overview
          v-if="processedCredentials.length === 1"
          class="my-5 waci-share-credential-overview-root"
          :credential="processedCredentials[0]"
        >
          <template #bannerBottomContainer>
            <div
              class="
                absolute
                justify-start
                items-start
                px-5
                pt-13
                pb-3
                w-full
                bg-neutrals-white
                rounded-b-xl
                flex flex-row
                waci-share-credential-overview-vault
              "
            >
              <span class="flex text-sm font-bold text-neutrals-dark">
                {{ t('CredentialDetails.Banner.vault') }}
              </span>
              <span class="flex ml-3 text-sm text-neutrals-medium">
                {{ processedCredentials[0].vaultName }}
              </span>
            </div>
          </template>
          <template #credentialDetails>
            <credential-details-table
              :heading="t('WACI.Share.whatIsShared')"
              :credential="processedCredentials[0]"
            />
          </template>
        </credential-overview>

        <!-- List of Credential Banners (Links to Details for each) -->
        <ul v-else-if="processedCredentials.length > 1" class="mt-6 space-y-5 w-full">
          <li v-for="(credential, index) in processedCredentials" :key="index">
            <credential-banner
              :id="credential.id"
              :brand-color="credential.brandColor"
              :icon="credential.icon"
              :title="credential.title"
              @click="handleOverviewClick(credential.id)"
            />
          </li>
        </ul>
      </div>
    </div>

    <WACI-action-buttons-container>
      <template #leftButton>
        <styled-button id="cancelBtn" type="outline" @click="cancel">
          {{ t('CHAPI.Share.decline') }}
        </styled-button>
      </template>
      <template #rightButton>
        <styled-button id="share-credentials" type="primary" @click="share">
          {{ t('CHAPI.Share.share') }}
        </styled-button>
      </template>
    </WACI-action-buttons-container>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';
import { useI18n } from 'vue-i18n';
import { getCredentialType, getCredentialDisplayData, getCredentialIcon } from '@/mixins';
import { decode, encode } from 'js-base64';
import { CredentialManager } from '@trustbloc-cicd/wallet-sdk';
import { OIDCMutations } from '@/layouts/OIDC.vue';
import { OIDCShareLayoutMutations } from '@/layouts/OIDCShareLayout.vue';
import StyledButton from '@/components/StyledButton/StyledButton.vue';
import CredentialBanner from '@/components/WACI/CredentialBanner.vue';
import CredentialOverview from '@/components/WACI/CredentialOverview.vue';
import WACIActionButtonsContainer from '@/components/WACI/WACIActionButtonsContainer.vue';
import WACICredentialsMissing from '@/components/WACI/WACICredentialsMissing.vue';
import WACIError from '@/components/WACI/WACIError.vue';
import WACILoading from '@/components/WACI/WACILoading.vue';
import WACISuccess from '@/components/WACI/WACISuccess.vue';
import CredentialDetailsTable from '@/components/WACI/CredentialDetailsTable.vue';
import OIDCShareOverview from '@/pages/OIDCShareOverview.vue';
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
    CredentialBanner,
    CredentialDetailsTable,
    CredentialOverview,
    StyledButton,
    WACIActionButtonsContainer,
    WACICredentialsMissing,
    WACIError,
    WACILoading,
    WACISuccess,
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
      credentialDisplayData: {},
      sharedSuccessfully: false,
      token: null,
    };
  },
  computed: {
    showCredentialsMissing() {
      return this.processedCredentials.length === 0;
    },
    showMainState() {
      return (
        !this.loading &&
        !this.errors.length &&
        !this.sharing &&
        !this.showCredentialsMissing &&
        !this.sharedSuccessfully
      );
    },
  },
  created: async function () {
    this.loading = true;
    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.credentialDisplayData = await this.getCredentialManifestData();
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
      const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
      const { results } = await credentialManager.query(this.token, [
        {
          type: 'PresentationExchange',
          credentialQuery: [this.claims.vp_token.presentation_definition],
        },
      ]);
      this.presentations = results;
      this.generateIdToken();
    } catch (e) {
      if (!e.message.includes('12009')) {
        this.errors.push('Error initiating credential share');
      }
      console.error('initiating credential share failed,', e);
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
    ...mapGetters(['getCurrentUser', 'getCredentialManifestData', 'getStaticAssetsUrl']),
    getCredentialIcon: function (icon) {
      return getCredentialIcon(this.getStaticAssetsUrl(), icon);
    },
    prepareRecords: function (presentations) {
      const credentials = presentations.reduce(
        (acc, val) => acc.concat(val.verifiableCredential),
        []
      );
      try {
        credentials.map((credential) => {
          const manifest = this.getManifest(credential);
          const processedCredential = this.getCredentialDisplayData(credential, manifest);
          const vaultName = 'Unavailable';
          this.processedCredentials.push({ ...processedCredential, vaultName });
        });
      } catch (e) {
        this.errors.push('No credentials found matching requested criteria.');
        console.error('get credentials failed,', e);
        this.loading = false;
      }
    },
    share() {
      this.sharing = true;
      const { profile, preference } = this.getCurrentUser();
      const { controller, proofType, verificationMethod } = preference;
      const VPToken = JSON.stringify(this.claims.vp_token);
      let ack;

      try {
        // Using "try" because eventually we will be making an AJAX call here
        ack = {
          status: 'OK',
          url: `${this.$route.query.redirect_uri}?state=${this.$route.query.state}&id_token=${this.idToken}&vp_token=${VPToken}`,
        };
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
      window.location.href = this.redirectUrl;
    },
    getCredentialType: function (vc) {
      return getCredentialType(vc.type);
    },
    getCredentialDisplayData: function (vc, manifestCredential) {
      return getCredentialDisplayData(vc, manifestCredential);
    },
    getManifest: function (credential) {
      const currentCredentialType = this.getCredentialType(credential);
      return (
        this.credentialDisplayData[currentCredentialType] || this.credentialDisplayData.fallback
      );
    },
    generateIdToken: function () {
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
    handleOverviewClick: function (id) {
      OIDCMutations.setSelectedCredentialId(id);
      OIDCShareLayoutMutations.setComponent(OIDCShareOverview);
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
