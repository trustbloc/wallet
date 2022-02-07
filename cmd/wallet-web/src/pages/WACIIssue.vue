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

    <!-- Saving State -->
    <WACI-loading v-else-if="saving" :message="t('WACI.Issue.savingCredential')" />

    <!-- Error State -->
    <WACI-error v-else-if="errors.length" @click="cancel" />

    <!-- Success State -->
    <WACI-success
      v-else-if="savedSuccessfully"
      id="issue-credentials-ok-btn"
      :heading="t('WACI.Issue.success', processedCredentials.length)"
      :message="
        t('WACI.Issue.message', {
          subject: processedCredentials[0].title,
        })
      "
      :button-label="successButtonLabel"
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
        <span class="mb-6 text-3xl font-bold">{{ t('WACI.Issue.saveCredential') }}</span>

        <credential-overview
          v-if="processedCredentials.length === 1"
          :credential="processedCredentials[0]"
        >
          <template #bannerBottomContainer>
            <div
              class="
                justify-start
                items-start
                px-4
                mt-5
                w-full
                bg-neutrals-lilacSoft
                rounded-t-lg
                flex flex-col flex-grow
                border-b border-neutrals-dark
              "
            >
              <label for="select-key" class="mb-1 text-sm font-bold text-neutrals-dark">{{
                t('Vaults.selectVault')
              }}</label>
              <custom-select
                id="waci-issue-select-vault"
                :options="vaults"
                default="Default Vault"
                @selected="setSelectedVault"
              />
            </div>
          </template>
          <template #credentialDetails>
            <credential-details-table
              :heading="t('WACI.Issue.verifiedInformation')"
              :credential="processedCredentials[0]"
              class="mt-8"
            />
          </template>
        </credential-overview>
      </div>
    </div>

    <WACI-action-buttons-container>
      <template #leftButton>
        <styled-button id="cancelBtn" type="outline" @click="cancel">
          {{ t('CHAPI.Share.decline') }}
        </styled-button>
      </template>
      <template #rightButton>
        <styled-button id="storeVCBtn" type="primary" @click="save">
          {{ t('WACI.Issue.save') }}
        </styled-button>
      </template>
    </WACI-action-buttons-container>
  </div>
</template>

<script>
import { toRaw } from 'vue';
import { mapGetters } from 'vuex';
import { CollectionManager, DIDComm } from '@trustbloc/wallet-sdk';
import { useI18n } from 'vue-i18n';
import { getCredentialDisplayData, getCredentialType, getCredentialIcon } from '@/mixins';
import { WACIStore } from '@/layouts/WACI.vue';
import CustomSelect from '@/components/CustomSelect/CustomSelect.vue';
import StyledButton from '@/components/StyledButton/StyledButton.vue';
import CredentialOverview from '@/components/WACI/CredentialOverview.vue';
import WACIActionButtonsContainer from '@/components/WACI/WACIActionButtonsContainer.vue';
import WACIError from '@/components/WACI/WACIError.vue';
import WACILoading from '@/components/WACI/WACILoading.vue';
import WACISuccess from '@/components/WACI/WACISuccess.vue';
import CredentialDetailsTable from '@/components/WACI/CredentialDetailsTable.vue';

export default {
  components: {
    CredentialDetailsTable,
    CredentialOverview,
    CustomSelect,
    StyledButton,
    WACIActionButtonsContainer,
    WACIError,
    WACILoading,
    WACISuccess,
  },
  setup() {
    const { t, locale } = useI18n();
    return { t, locale };
  },
  data() {
    return {
      loading: true,
      saving: false,
      errors: [],
      vaults: [],
      selectedVault: '',
      credentialDisplayData: '',
      processedCredentials: [],
      savedSuccessfully: false,
    };
  },
  computed: {
    showMainState() {
      return !this.loading && !this.errors.length && !this.saving && !this.savedSuccessfully;
    },
    successButtonLabel() {
      return this.redirectUrl ? this.t('WACI.ok') : this.t('WACI.Issue.viewCredential');
    },
  },
  created: async function () {
    this.loading = true;
    this.protocolHandler = WACIStore.protocolHandler;
    const invitation = toRaw(this.protocolHandler.message());
    const { user, token } = this.getCurrentUser().profile;
    const credManifest = this.getCredentialManifestData();
    const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    const defVault = this.fetchAllVaults(token, collectionManager);

    this.didcomm = new DIDComm({ agent: this.getAgentInstance(), user });
    try {
      const { threadID, presentations, fulfillment, manifest, domain, challenge, error } =
        await this.didcomm.initiateCredentialIssuance(token, invitation, {
          userAnyRouterConnection: true,
        });

      // business error
      if (error) {
        const { status, url, code } = error;
        if (url) {
          this.protocolHandler.done(url);
        }

        this.handleError(error);
        return;
      }

      // TODO: [Issue#1336] - read manifest, presentations, normalized, comment, fields to enhance UI
      this.interactionData = { threadID, fulfillment, domain, challenge, error, manifest };
    } catch (e) {
      this.handleError(e);
      return;
    }

    // Await here to allow these operations to run asynchronously along with the didcomm request
    this.credentialDisplayData = await credManifest;
    this.setSelectedVault(await defVault);

    // prepare cards
    this.prepareCards();

    this.loading = false;
  },
  methods: {
    ...mapGetters(['getCurrentUser', 'getStaticAssetsUrl', 'getCredentialManifestData']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    prepareCards: function () {
      const { fulfillment } = this.interactionData;
      // If only one credential then show details by default
      if (fulfillment.verifiableCredential.length === 1) {
        const manifest = this.getManifest(fulfillment.verifiableCredential[0]);
        const credential = this.getCredentialDisplayData(
          fulfillment.verifiableCredential[0],
          manifest
        );
        this.processedCredentials.push({ ...credential });
      } else {
        fulfillment.verifiableCredential.map((vc) => {
          const manifest = this.getManifest(vc);
          const credential = this.getCredentialDisplayData(vc, manifest);
          this.processedCredentials.push({ ...credential });
        });
      }
    },
    save: async function () {
      this.errors.length = 0;
      this.saving = true;

      const { profile, preference } = this.getCurrentUser();
      const { controller, proofType, verificationMethod } = preference;

      const { threadID, manifest, presentations, domain, challenge } = this.interactionData;

      let ack;
      try {
        ack = await this.didcomm.completeCredentialIssuance(
          profile.token,
          threadID,
          presentations && presentations.length > 0 ? presentations[0] : null,
          manifest,
          {
            controller,
            proofType,
            verificationMethod,
            domain,
            challenge,
          },
          {
            waitForDone: true,
            autoAccept: true,
            collectionID: this.selectedVault,
          }
        );
      } catch (e) {
        this.handleError(e);
        return;
      }

      console.debug('acknowledgment msg', ack);
      const { status, url } = ack;
      this.redirectUrl = url;
      // TODO: [Issue#1336] - check if status="FAIL", then should redirect to generic error screen, it means WACI flow didn't succeed
      if (status === 'OK') this.savedSuccessfully = true;
      this.saving = false;
    },
    handleError: function (e) {
      console.error('failed to perform credential issuance flow', e);
      this.errors.push(e);
      this.loading = false;
    },
    setSelectedVault: function (e) {
      this.selectedVault = e;
    },
    finish() {
      this.protocolHandler.done(
        this.redirectUrl ? this.redirectUrl : `${window.location.origin}/${this.locale}/credentials`
      );
    },
    cancel: function () {
      this.protocolHandler.cancel();
    },
    getCredentialType: function (vc) {
      return getCredentialType(vc.type);
    },
    getCredentialIcon: function (icon) {
      return getCredentialIcon(this.getStaticAssetsUrl(), icon);
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
    fetchAllVaults: async function (token, collectionManager) {
      const { contents } = await collectionManager.getAll(token);
      this.vaults = Object.values(contents).map((vault) => vault);
      // Default vault is selected vault by default, it is created on wallet setup and must be only one.
      const defaultVaultId = this.vaults.find((vault) => vault.name === 'Default Vault').id;
      return defaultVaultId;
    },
  },
};
</script>

<style scoped>
.chapi-container {
  width: 28rem;
}

.credentialPreviewContainer:not(:focus-within):hover {
  box-shadow: 0px 4px 12px 0px rgba(25, 12, 33, 0.1);
}

.footerContainer {
  box-shadow: inset 0px 1px 0px 0px rgb(219, 215, 220);
}
</style>
