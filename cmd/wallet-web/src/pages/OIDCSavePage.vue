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
  <div
    v-else
    class="flex overflow-hidden flex-col flex-grow justify-between items-center w-full h-full"
  >
    <div class="flex overflow-auto justify-center w-full">
      <div
        class="flex flex-col flex-grow justify-start items-start pt-8 pr-5 md:pr-0 pb-8 pl-5 md:pl-0 w-full max-w-3xl h-full"
      >
        <span class="mb-6 text-3xl font-bold">{{ t('WACI.Issue.saveCredential') }}</span>

        <div
          v-for="(credential, index) in processedCredentials"
          :key="index"
          class="flex flex-col justify-start w-full max-w-3xl"
        >
          <CredentialOverviewComponent :credential="credential">
            <template #bannerBottomContainer>
              <div
                class="flex flex-col flex-grow justify-start items-start px-4 mt-5 w-full bg-neutrals-lilacSoft rounded-t-lg border-b border-neutrals-dark"
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
import { mapGetters } from 'vuex';
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { useI18n } from 'vue-i18n';
import CustomSelectComponent from '@/components/CustomSelect/CustomSelectComponent.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import CredentialOverviewComponent from '@/components/WACI/CredentialOverviewComponent.vue';
import WACIActionButtonsContainerComponent from '@/components/WACI/WACIActionButtonsContainerComponent.vue';
import WACIErrorComponent from '@/components/WACI/WACIErrorComponent.vue';
import WACILoadingComponent from '@/components/WACI/WACILoadingComponent.vue';
import WACISuccessComponent from '@/components/WACI/WACISuccessComponent.vue';
import CredentialDetailsTableComponent from '@/components/WACI/CredentialDetailsTableComponent.vue';
import { readOpenIDConfiguration, requestCredential, requestToken } from '@/mixins';
import Cookies from 'js-cookie';
import jp from 'jsonpath';

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

    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    const defVault = this.fetchAllVaults(token, collectionManager);
    // Await here to allow these operations to run asynchronously along with the didcomm request
    this.setSelectedVault(await defVault);

    const txState = Cookies.get(this.$route.query.state);
    if (!txState) {
      //TODO put error handling for invalid state
      throw 'invalid state, session expired!';
    }

    const { issuer, credentialTypes } = JSON.parse(txState);
    const configuration = await readOpenIDConfiguration(issuer);

    const { access_token, token_type } = await requestToken(configuration.token_endpoint, {
      redirect_uri: `${location.protocol}//${location.host}/oidc/save`,
      code: this.$route.query.code,
    });

    let processedCredentials = [];

    this.saveData = await Promise.all(
      credentialTypes.map(async (credentialType) => {
        const { credential } = await requestCredential(configuration.credential_endpoint, {
          access_token,
          token_type,
          credentialType,
        });

        const { processed, descriptorID, manifest } = await this.prepareCards(
          credential,
          configuration.credential_manifests,
          credentialType
        );

        processedCredentials.push(...processed);

        return {
          credential,
          manifest,
          descriptorID,
        };
      })
    );

    this.processedCredentials = processedCredentials;
    this.loading = false;
  },
  methods: {
    ...mapGetters(['getCurrentUser']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    save: function () {
      this.errors.length = 0;
      this.saving = true;

      const { profile, preference } = this.getCurrentUser();
      const { controller, proofType, verificationMethod } = preference;

      this.saveData.forEach(({ credential, manifest, descriptorID }) => {
        this.credentialManager.save(
          profile.token,
          { credentials: [credential] },
          {
            manifest,
            descriptorMap: [
              {
                id: descriptorID,
                format: 'ldp_vc',
                path: '$[0]',
              },
            ],
            collection: this.selectedVault,
          }
        );
      });

      this.savedSuccessfully = true;
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
      this.$router.push('/credentials');
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
    prepareCards: async function (credential, manifests, type) {
      let manifest;
      let descriptorID;
      for (const m of manifests) {
        const match = jp.query(m, `$.output_descriptors[?(@.schema=="${type}")].id`);
        if (match.length > 0) {
          manifest = m;
          descriptorID = match[0];
          break;
        }
      }

      if (!manifest) {
        throw 'unable to find matching manifest'; // TODO handle this error, Issue #1531
      }

      const processed = await this.credentialManager.resolveManifest(this.token, {
        credential,
        manifest,
        descriptorID,
      });

      return { processed, descriptorID, manifest };
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
