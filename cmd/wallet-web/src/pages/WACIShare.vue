<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex overflow-hidden relative justify-start items-start w-full h-full">
    <!-- Loading State -->
    <div v-if="loading" class="flex flex-grow justify-center items-center w-full h-full">
      <Spinner />
    </div>
    <!-- Sharing State -->
    <div
      v-else-if="sharing"
      class="flex flex-col flex-grow justify-center items-center w-full h-full"
    >
      <Spinner />
      <span class="mt-8 text-base md:text-lg text-neutrals-dark">{{
        t('CHAPI.Share.sharingCredential')
      }}</span>
    </div>
    <!-- Error State -->
    <div
      v-else-if="errors.length"
      class="flex flex-col flex-grow justify-center items-center w-full h-full"
    >
      <div class="flex flex-col flex-grow justify-center items-center w-full h-full">
        <img src="@/assets/img/icons-error.svg" />
        <span class="mt-6 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
          t('CHAPI.Share.Error.heading')
        }}</span>
        <span class="w-full text-lg text-center text-neutrals-medium">{{
          t('CHAPI.Share.Error.body')
        }}</span>
        <styled-button id="share-credentials-ok-btn" type="primary" class="mt-6" @click="cancel">
          {{ t('CHAPI.Share.Error.tryAgain') }}
        </styled-button>
      </div>
    </div>
    <!-- Credentials Missing State -->
    <div
      v-else-if="showCredentialsMissing"
      class="flex flex-col flex-grow justify-center items-center w-full h-full"
    >
      <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
        <img src="@/assets/img/icons-error.svg" />
        <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
          t('CHAPI.Share.CredentialsMissing.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          t('CHAPI.Share.CredentialsMissing.body')
        }}</span>
        <styled-button id="share-credentials-ok-btn" type="outline" class="mt-6" @click="cancel">
          {{ t('CHAPI.Share.CredentialsMissing.ok') }}
        </styled-button>
      </div>
    </div>
    <!-- Main State -->
    <div
      v-else
      class="flex overflow-hidden flex-col flex-grow justify-between items-center w-full h-full"
    >
      <div class="overflow-auto w-full flex justify-center">
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
          <span class="mb-6 text-3xl font-bold">{{ t('CHAPI.Share.shareCredential') }}</span>
          <div class="flex flex-row justify-start items-start mb-4 w-full">
            <div class="flex-none w-12 h-12 border-opacity-10">
              <!-- todo issue-1055 Read meta data from external urls -->
              <img src="@/assets/img/generic-issuer-icon.svg" />
            </div>
            <div class="flex flex-col pl-3">
              <span
                class="flex-1 mb-1 text-sm font-bold text-left text-neutrals-dark overflow-ellipsis"
              >
                <!-- todo issue-1055 Read meta data from external urls -->
                Requestor
              </span>
              <div class="flex flex-row justify-center items-center">
                <img src="@/assets/img/small-lock-icon.svg" />
                <span class="flex-1 pl-1 text-xs text-left text-neutrals-medium overflow-ellipsis">
                  {{ requestOrigin }}
                </span>
              </div>
            </div>
          </div>

          <span class="text-sm text-neutrals-dark">{{
            t('CHAPI.Share.headline', processedCredentials.length, { issuer: 'Requestor' })
          }}</span>

          <!-- Single Credential Preview (with details) -->
          <credential-preview
            v-if="processedCredentials.length === 1"
            :credential="processedCredentials[0]"
          />
          <!-- List of Credential Banners (Links to Details for each) -->
          <!-- TODO: issue-1391 -->
          <!-- <ul v-else-if="processedCredentials.length" class="space-y-5 w-full">
              <li v-for="(credential, index) in processedCredentials" :key="index">
                <credential-banner
                  :id="credential.id"
                  :brand-color="credential.brandColor"
                  :icon="credential.icon"
                  :title="credential.title"
                />
              </li>
            </ul> -->
        </div>
      </div>

      <div
        class="
          sticky
          bottom-0
          z-20
          justify-center
          items-center
          pt-4
          pr-5
          pb-4
          pl-5
          w-full
          bg-neutrals-magnolia
          flex flex-row
          border-t border-neutrals-thistle
        "
      >
        <div class="flex flex-row flex-grow justify-between items-center w-full max-w-3xl">
          <styled-button id="cancelBtn" type="outline" @click="cancel">
            {{ t('CHAPI.Share.decline') }}
          </styled-button>
          <styled-button id="share-credentials" type="primary" @click="share">
            {{ t('CHAPI.Share.share') }}
          </styled-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { ErrorCodes, toRaw } from 'vue';
import { mapGetters } from 'vuex';
import { useI18n } from 'vue-i18n';
import { DIDComm } from '@trustbloc/wallet-sdk';
import { getCredentialType, getCredentialDisplayData, getCredentialIcon } from '@/utils/mixins';
import Spinner from '@/components/Spinner/Spinner.vue';
import StyledButton from '@/components/StyledButton/StyledButton.vue';
// import CredentialBanner from '@/components/WACI/CredentialBanner.vue';
import CredentialPreview from '@/components/WACI/CredentialPreview.vue';

export default {
  components: {
    Spinner,
    StyledButton,
    // CredentialBanner,
    CredentialPreview,
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
    };
  },
  computed: {
    showCredentialsMissing() {
      return this.processedCredentials.length === 0;
    },
  },
  created: async function () {
    this.loading = true;
    this.protocolHandler = this.$parent.protocolHandler;
    const invitation = toRaw(this.protocolHandler.message());
    const { user, token } = this.getCurrentUser().profile;
    this.credentialDisplayData = await this.getCredentialManifestData();

    //initiate credential share flow.
    this.didcomm = new DIDComm({ agent: this.getAgentInstance(), user });
    try {
      const { threadID, presentations } = await this.didcomm.initiateCredentialShare(
        token,
        invitation,
        { userAnyRouterConnection: true }
      );

      this.threadID = threadID;
      this.presentations = presentations;
    } catch (e) {
      if (e.includes('12009')) {
        this.showCredentialsMissing = true;
      } else {
        this.errors.push('Error initiating credential share');
      }
      console.error('initiating credential share failed,', e);
      this.loading = false;
      return;
    }

    this.prepareRecords(this.presentations);
    this.requestOrigin = this.protocolHandler.requestor();
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getCredentialManifestData', 'getStaticAssetsUrl']),
    getCredentialIcon: function (icon) {
      return getCredentialIcon(this.getStaticAssetsUrl(), icon);
    },
    prepareRecords: function (presentations) {
      try {
        const credentials = presentations.reduce(
          (acc, val) => acc.concat(val.verifiableCredential),
          []
        );
        credentials.map((credential) => {
          const manifest = this.getManifest(credential);
          const processedCredential = this.getCredentialDisplayData(credential, manifest);
          // TODO: issue1410 - add logic to retrieve the list of vaults in which the credential is stored
          const vaultName = 'Unavailable';
          this.processedCredentials.push({ ...processedCredential, vaultName });
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

      let { status, url } = ack;

      // TODO check if status="FAIL", then should redirect to generic error screen, it means WACI flow didn't succeed

      this.protocolHandler.done(url);

      this.sharing = false;
    },
    cancel() {
      this.protocolHandler.cancel();
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
  },
};
</script>
