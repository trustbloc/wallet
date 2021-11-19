<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex justify-center items-start w-screen h-screen">
    <div
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
                <img :src="getCrendentialIcon(credential.icon)" />
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
                  @selected="vaultSelected"
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
  CHAPIEventHandler,
  getCredentialDisplayData,
  getCredentialType,
  getCrendentialIcon,
  isVPType,
} from './mixins';
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import CustomSelect from '@/components/CustomSelect/CustomSelect';
import { useI18n } from 'vue-i18n';

export default {
  components: { CustomSelect },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      processedCredentials: [],
      errors: [],
      vaults: [],
      selectedVault: '',
      credentialDisplayData: '',
    };
  },
  created: async function () {
    // Load the Credentials
    this.credentialEvent = new CHAPIEventHandler(
      await this.$webCredentialHandler.receiveCredentialEvent()
    );
    const { dataType, data } = this.credentialEvent.getEventData();

    if (!isVPType(dataType)) {
      this.errors.push(`unknown credential data type '${dataType}'`);
      return;
    }

    const { user, token } = this.getCurrentUser().profile;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });

    this.credentialDisplayData = await this.getCredentialManifestData();
    // prepare cards
    this.prepareCards(data);
    this.presentation = data;
    const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    this.fetchAllVaults(token, collectionManager);
  },
  methods: {
    ...mapGetters(['getCurrentUser', 'getStaticAssetsUrl', 'getCredentialManifestData']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    toggleDetails(credential) {
      credential.showDetails = !credential.showDetails;
    },
    prepareCards: function (data) {
      // If only one credential then show details by default
      if (data.verifiableCredential.length === 1) {
        const manifest = this.getManifest(data.verifiableCredential[0]);
        const credential = this.getCredentialDisplayData(data.verifiableCredential[0], manifest);
        this.processedCredentials.push({ ...credential, showDetails: true });
      } else {
        data.verifiableCredential.map((vc) => {
          const manifest = this.getManifest(vc);
          const credential = this.getCredentialDisplayData(vc, manifest);
          this.processedCredentials.push({ ...credential, showDetails: false });
        });
      }
    },
    store: function () {
      this.errors.length = 0;
      const { token } = this.getCurrentUser().profile;

      this.credentialManager
        .save(
          token,
          {
            presentation: this.presentation,
          },
          { collection: this.selectedVault }
        )
        .then(() => {
          this.credentialEvent.done();
        })
        .catch((e) => {
          console.error(e);
          this.errors.push(`failed to save credential`);
        });
    },
    vaultSelected: function (e) {
      this.selectedVault = e;
    },
    cancel: function () {
      this.credentialEvent.cancel();
    },
    getCredentialType: function (vc) {
      return getCredentialType(vc.type);
    },
    getCrendentialIcon: function (icon) {
      return getCrendentialIcon(this.getStaticAssetsUrl(), icon);
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
