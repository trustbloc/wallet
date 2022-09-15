<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex h-screen w-screen items-start justify-center">
    <div
      class="chapi-container max-h-screen overflow-scroll rounded-b border border-neutrals-black bg-neutrals-softWhite pt-5"
    >
      <span class="px-5 text-xl font-bold text-neutrals-dark">Save credential</span>
      <div v-if="processedCredentials.length" class="flex flex-col justify-center px-5">
        <ul class="my-8 grid grid-cols-1 gap-4">
          <li v-for="(credential, index) in processedCredentials" :key="index">
            <div
              class="group credentialPreviewContainer inline-flex h-20 w-full items-center rounded-xl border p-5 text-sm font-bold md:h-24 md:text-base"
              :class="
                credential.styles.background.color !== '#fff'
                  ? `border-neutrals-black border-opacity-10`
                  : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle`
              "
              :style="`background-color: ${credential.styles.background.color}`"
              @click="toggleDetails(credential)"
            >
              <div class="h-12 w-12 flex-none border-opacity-10">
                <img :src="getCredentialIconSrc(credential)" />
              </div>
              <div class="grow p-4">
                <span
                  class="text-ellipsis text-left text-sm font-bold md:text-base"
                  :style="`color: ${credential.styles.text.color}`"
                >
                  {{ credential.title }}
                </span>
              </div>
            </div>
            <div
              v-if="credential.showDetails"
              class="mt-5 flex w-full flex-col items-start justify-start md:mt-6"
            >
              <!-- TODO: populate with dynamic vault list -->
              <div
                class="mb-8 flex w-full grow flex-col items-start justify-start rounded-t-lg border-b border-neutrals-dark bg-neutrals-lilacSoft px-4"
              >
                <label for="select-key" class="mb-1 text-sm font-bold text-neutrals-dark">{{
                  t('Vaults.selectVault')
                }}</label>
                <CustomSelectComponent
                  id="valutOptions"
                  :options="vaults"
                  default="Default Vault"
                  @selected="setSelectedVault"
                ></CustomSelectComponent>
              </div>

              <span class="py-3 text-base font-bold text-neutrals-dark">Verified Information</span>

              <!-- TODO: move this to a reusable component -->
              <table class="w-full border-t border-neutrals-chatelle">
                <tr
                  v-for="(property, key) of credential.properties"
                  :key="key"
                  class="border-b border-dotted border-neutrals-thistle"
                >
                  <td class="py-4 pr-6 pl-3 text-neutrals-medium">{{ property.label }}</td>
                  <td
                    v-if="property.schema.contentMediaType === 'image/png'"
                    class="break-words py-4 pr-6 pl-3 text-neutrals-dark"
                  >
                    <img :src="property.value" class="h-20 w-20" />
                  </td>
                  <td v-else class="break-words py-4 pr-6 pl-3 text-neutrals-dark">
                    {{ property.value }}
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
        class="footerContainer sticky bottom-0 flex h-auto w-full justify-between bg-neutrals-magnolia p-5"
      >
        <button id="cancelBtn" class="btn-outline" @click="cancel">Decline</button>
        <button id="storeVCBtn" class="btn-primary" @click="store">Save</button>
      </div>
    </div>
  </div>
</template>

<script>
import { useStore } from 'vuex';
import * as webCredentialHandler from 'web-credential-handler';
import {
  CHAPIEventHandler,
  getCredentialIcon,
  isVPType,
  prepareCredentialManifest,
  resolveManifest,
} from '@/mixins';
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import CustomSelectComponent from '@/components/CustomSelect/CustomSelectComponent.vue';
import { useI18n } from 'vue-i18n';

export default {
  components: { CustomSelectComponent },
  setup() {
    const store = useStore();
    const { t } = useI18n();
    const getStaticAssetsUrl = () => store.getters.getStaticAssetsUrl;
    const getCredentialIconSrc = (credential) => {
      return credential?.styles?.thumbnail?.uri?.includes('https://')
        ? credential?.styles?.thumbnail?.uri
        : getCredentialIcon(getStaticAssetsUrl(), credential?.styles?.thumbnail?.uri);
    };

    return { getCredentialIconSrc, t, webCredentialHandler };
  },
  data() {
    return {
      processedCredentials: [],
      errors: [],
      vaults: [],
      selectedVault: '',
      loading: true,
    };
  },
  created: async function () {
    this.loading = true;
    // Load the Credentials
    this.credentialEvent = new CHAPIEventHandler(
      await this.webCredentialHandler.receiveCredentialEvent()
    );
    const { dataType, data } = this.credentialEvent.getEventData();

    if (!isVPType(dataType)) {
      this.errors.push(`unknown credential data type '${dataType}'`);
      return;
    }

    const { user, token } = this.getCurrentUser().profile;
    this.token = token;
    this.presentation = data;
    this.manifest = prepareCredentialManifest(
      this.presentation,
      this.getCredentialManifests(),
      this.credentialEvent.requestor()
    );
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
    const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
    await this.fetchCredentials();
    await this.fetchVaults(token, collectionManager);
    this.loading = false;
  },
  methods: {
    ...mapGetters(['getCurrentUser', 'getCredentialManifests']),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    toggleDetails(credential) {
      credential.showDetails = !credential.showDetails;
    },
    fetchCredentials: async function () {
      this.processedCredentials = await resolveManifest(
        this.credentialManager,
        this.getCredentialManifests(),
        this.token,
        {
          manifest: this.manifest,
          response: this.presentation,
        }
      );
      if (this.processedCredentials.length === 1) this.processedCredentials[0].showDetails = true;
    },
    fetchVaults: async function (token, collectionManager) {
      const { contents } = await collectionManager.getAll(token);
      this.vaults = Object.values(contents).map((vault) => vault);
      // Default vault is selected vault by default, it is created on wallet setup and must be only one.
      const defaultVaultId = this.vaults.find((vault) => vault.name === 'Default Vault').id;
      this.setSelectedVault(defaultVaultId);
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
          { collection: this.selectedVault, manifest: this.manifest }
        )
        .then(() => {
          this.credentialEvent.done();
        })
        .catch((e) => {
          console.error(e);
          this.errors.push(`failed to save credential`);
        });
    },
    setSelectedVault: function (e) {
      this.selectedVault = e;
    },
    cancel: function () {
      this.credentialEvent.cancel();
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
