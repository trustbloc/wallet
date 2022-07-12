<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div v-if="!showMainState" class="flex flex-col grow justify-center items-center w-full h-full">
    <!-- Loading State -->
    <WACILoadingComponent v-if="loading" />

    <!-- Saving State -->
    <WACILoadingComponent v-else-if="saving" :message="t('WACI.Issue.savingCredential')" />

    <!-- Error State -->
    <WACIErrorComponent v-else-if="errors.length" @click="cancel" />

    <!-- Success State -->
    <WACISuccessComponent
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
  <div v-else class="flex overflow-hidden flex-col grow justify-between items-center w-full h-full">
    <div class="flex overflow-auto justify-center w-full">
      <div
        class="flex flex-col grow justify-start items-start py-8 px-5 w-full max-w-3xl h-full md:px-0"
      >
        <span class="mb-6 text-3xl font-bold">{{ t('WACI.Issue.saveCredential') }}</span>

        <CredentialOverviewComponent
          v-if="processedCredentials.length === 1"
          :credential="processedCredentials[0]"
        >
          <template #bannerBottomContainer>
            <div
              class="flex flex-col grow justify-start items-start px-4 mt-5 w-full bg-neutrals-lilacSoft rounded-t-lg border-b border-neutrals-dark"
            >
              <label for="select-key" class="mb-1 text-sm font-bold text-neutrals-dark">{{
                t('Vaults.selectVault')
              }}</label>
              <CustomSelectComponent
                id="waci-issue-select-vault"
                :options="vaults"
                default="Default Vault"
                @selected="setSelectedVault"
              />
            </div>
          </template>
          <template #credentialDetails>
            <CredentialDetailsTableComponent
              :heading="t('WACI.Issue.verifiedInformation')"
              :credential="processedCredentials[0]"
              class="mt-8"
            />
          </template>
        </CredentialOverviewComponent>
      </div>
    </div>

    <WACIActionButtonsContainerComponent>
      <template #leftButton>
        <StyledButtonComponent id="cancelBtn" type="btn-outline" @click="cancel">
          {{ t('CHAPI.Share.decline') }}
        </StyledButtonComponent>
      </template>
      <template #rightButton>
        <StyledButtonComponent id="storeVCBtn" type="btn-primary" @click="save">
          {{ t('WACI.Issue.save') }}
        </StyledButtonComponent>
      </template>
    </WACIActionButtonsContainerComponent>
  </div>
</template>

<script>
import { toRaw } from 'vue';
import { mapGetters } from 'vuex';
import { CollectionManager, CredentialManager, DIDComm } from '@trustbloc/wallet-sdk';
import { useI18n } from 'vue-i18n';
import { WACIStore } from '@/layouts/WACILayout.vue';
import CustomSelectComponent from '@/components/CustomSelect/CustomSelectComponent.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import CredentialOverviewComponent from '@/components/WACI/CredentialOverviewComponent.vue';
import WACIActionButtonsContainerComponent from '@/components/WACI/WACIActionButtonsContainerComponent.vue';
import WACIErrorComponent from '@/components/WACI/WACIErrorComponent.vue';
import WACILoadingComponent from '@/components/WACI/WACILoadingComponent.vue';
import WACISuccessComponent from '@/components/WACI/WACISuccessComponent.vue';
import CredentialDetailsTableComponent from '@/components/WACI/CredentialDetailsTableComponent.vue';

export default {
  components: {
    CredentialDetailsTableComponent,
    CredentialOverviewComponent,
    CustomSelectComponent,
    StyledButtonComponent,
    WACIActionButtonsContainerComponent,
    WACIErrorComponent,
    WACILoadingComponent,
    WACISuccessComponent,
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
    this.token = token;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
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
    this.setSelectedVault(await defVault);

    // prepare cards
    await this.prepareCards();

    this.loading = false;
  },
  methods: {
    ...mapGetters(['getCurrentUser']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    prepareCards: async function () {
      const { fulfillment, manifest } = this.interactionData;
      this.processedCredentials = await this.credentialManager.resolveManifest(this.token, {
        manifest,
        fulfillment,
      });
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
      this.saving = false;
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
