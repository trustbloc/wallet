<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<!-- TODO: issue #1336, temporary WACI issuance page, to be standardized.-->

<template>
  <div class="flex justify-center items-start w-full">
    <!-- Loading State -->
    <div
      v-if="loading"
      class="
        flex
        justify-center
        items-center
        w-full
        max-w-md
        h-80
        bg-gray-light
        rounded-lg
        md:border
        border-neutrals-chatelle
      "
    >
      <Spinner />
    </div>

    <!-- Error State -->
    <div
      v-else-if="errors.length"
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-auto
        bg-gray-light
        rounded-lg
        md:border
        border-neutrals-chatelle
        flex flex-col
      "
    >
      <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
        <img src="@/assets/img/icons-error.svg" />
        <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
          t('CHAPI.Issue.Error.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          t('CHAPI.Issue.Error.body')
        }}</span>
      </div>
      <div
        class="
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
        <button id="share-credentials-ok-btn" class="btn-primary" @click="cancel">
          {{ t('CHAPI.Issue.Error.tryAgain') }}
        </button>
      </div>
    </div>

    <!-- Main State -->
    <div
      v-else
      class="
        overflow-scroll
        pt-5
        max-h-screen
        bg-neutrals-softWhite
        rounded-b
        border border-neutrals-black
        chapi-container
      "
    >
      <span class="px-5 text-xl font-bold text-neutrals-dark">Save credential</span>
      <div v-if="processedCredentials.length" class="flex flex-col justify-center px-5">
        <ul class="grid grid-cols-1 gap-4 my-8">
          <li
            v-for="(credential, index) in processedCredentials"
            :key="index"
            @click="toggleDetails(index)"
          >
            <div
              :class="[
                `group inline-flex items-center rounded-xl p-5 text-sm md:text-base font-bold border w-full h-20 md:h-24 focus-within:ring-2 focus-within:ring-offset-2 credentialPreviewContainer`,
                credential.brandColor.length
                  ? `bg-gradient-${credential.brandColor} border-neutrals-black border-opacity-10 focus-within:ring-primary-${credential.brandColor}`
                  : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`,
              ]"
              @click="toggleDetails(credential)"
            >
              <div class="flex-none w-12 h-12 border-opacity-10">
                <img :src="getCredentialIcon(credential.icon)" />
              </div>
              <div class="flex-grow p-4">
                <span
                  :class="[
                    `text-sm md:text-base font-bold text-left overflow-ellipsis`,
                    credential.brandColor.length ? `text-neutrals-white` : `text-neutrals-dark`,
                  ]"
                >
                  {{ credential.title }}
                </span>
              </div>
            </div>
            <div
              v-if="credential.showDetails"
              class="flex flex-col justify-start items-start mt-5 md:mt-6 w-full"
            >
              <!-- TODO: populate with dynamic vault list -->
              <div
                class="
                  justify-start
                  items-start
                  px-4
                  mb-8
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
                  id="valutOptions"
                  :options="vaults"
                  default="Default Vault"
                  @selected="setSelectedVault"
                ></custom-select>
              </div>

              <span class="py-3 text-base font-bold text-neutrals-dark">Verified Information</span>

              <!-- TODO: move this to a reusable component -->
              <table class="w-full border-t border-neutrals-chatelle">
                <tr
                  v-for="(property, key) of credential.properties"
                  :key="key"
                  class="border-b border-neutrals-thistle border-dotted"
                >
                  <td class="py-4 pr-6 pl-3 text-neutrals-medium">{{ property.label }}</td>
                  <td
                    v-if="property.type != 'image'"
                    class="py-4 pr-6 pl-3 text-neutrals-dark break-words"
                  >
                    {{ property.value }}
                  </td>
                  <td
                    v-if="property.type === 'image'"
                    class="py-4 pr-6 pl-3 text-neutrals-dark break-words"
                  >
                    <img :src="property.value" class="w-20 h-20" />
                  </td>
                </tr>
              </table>
            </div>
          </li>
        </ul>
      </div>
      <div v-if="errors.length">
        <b>Please correct the following error(s):</b>
        <ul>
          <!-- TODO: implement error as per designs -->
          <li v-for="error in errors" :key="error" class="text-sm text-primary-valencia">
            {{ error }}
          </li>
        </ul>
      </div>
      <div
        class="
          flex
          sticky
          bottom-0
          justify-between
          p-5
          w-full
          h-auto
          bg-neutrals-magnolia
          footerContainer
        "
      >
        <button id="cancelBtn" class="btn-outline" @click="cancel">Decline</button>
        <button id="storeVCBtn" class="btn-primary" @click="store">Save</button>
      </div>
    </div>
  </div>
</template>

<script>
import {
  getCredentialDisplayData,
  getCredentialType,
  getCredentialIcon,
  isVPType,
} from '@/utils/mixins';
import { toRaw } from 'vue';
import { mapGetters } from 'vuex';
import { CollectionManager, DIDComm } from '@trustbloc/wallet-sdk';
import CustomSelect from '@/components/CustomSelect/CustomSelect.vue';
import Spinner from '@/components/Spinner/Spinner.vue';
import { useI18n } from 'vue-i18n';

export default {
  components: { CustomSelect, Spinner },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      loading: true,
      processedCredentials: [],
      errors: [],
      vaults: [],
      selectedVault: '',
      credentialDisplayData: '',
    };
  },
  created: async function () {
    this.loading = true;
    this.protocolHandler = this.$parent.protocolHandler;
    const invitation = toRaw(this.protocolHandler.message());
    const { user, token } = this.getCurrentUser().profile;
    const credManifest = this.getCredentialManifestData();

    // initiate WACI credential issuance flow.
    this.didcomm = new DIDComm({ agent: this.getAgentInstance(), user });
    try {
      const { threadID, presentations, fulfillment, domain, challenge, error } =
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

      // TODO [Issue#1336] read manifest, presentations, normalized, comment, fields to enhance UI
      this.interactionData = { threadID, fulfillment, domain, challenge, error };
    } catch (e) {
      this.handleError(e);
      return;
    }

    this.credentialDisplayData = await credManifest;

    // prepare cards
    this.prepareCards();

    // get all vaults
    const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    this.fetchAllVaults(token, collectionManager);

    this.loading = false;
  },
  methods: {
    ...mapGetters(['getCurrentUser', 'getStaticAssetsUrl', 'getCredentialManifestData']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    toggleDetails(credential) {
      credential.showDetails = !credential.showDetails;
    },
    prepareCards: function () {
      const { fulfillment } = this.interactionData;
      // If only one credential then show details by default
      if (fulfillment.verifiableCredential.length === 1) {
        const manifest = this.getManifest(fulfillment.verifiableCredential[0]);
        const credential = this.getCredentialDisplayData(
          fulfillment.verifiableCredential[0],
          manifest
        );
        this.processedCredentials.push({ ...credential, showDetails: true });
      } else {
        fulfillment.verifiableCredential.map((vc) => {
          const manifest = this.getManifest(vc);
          const credential = this.getCredentialDisplayData(vc, manifest);
          this.processedCredentials.push({ ...credential, showDetails: false });
        });
      }
    },
    store: async function () {
      this.errors.length = 0;
      this.loading = true;

      const { profile, preference } = this.getCurrentUser();
      const { controller, proofType, verificationMethod } = preference;

      const { threadID, presentations, domain, challenge } = this.interactionData;

      let ack;
      try {
        ack = await this.didcomm.completeCredentialIssuance(
          profile.token,
          threadID,
          presentations && presentations.length > 0 ? presentations[0] : null,
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
      let { status, url } = ack;

      this.loading = false;
      // TODO [Issue#1336]  check if status="FAIL", then should redirect to generic error screen, it means WACI flow didn't succeed
      this.protocolHandler.done(url);
    },
    handleError: function (e) {
      console.error('failed to perform credential issuance flow', e);
      this.errors.push(e);
      this.loading = false;
    },
    setSelectedVault: function (e) {
      this.selectedVault = e;
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
      this.setSelectedVault(defaultVaultId);
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
