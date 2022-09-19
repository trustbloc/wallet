<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <!-- Loading State -->
  <div v-if="loading" class="flex h-screen w-screen items-start justify-center">
    <div
      class="flex h-80 w-full max-w-md items-center justify-center border-neutrals-black bg-gray-light md:border md:border-t-0"
    >
      <SpinnerIcon />
    </div>
  </div>
  <!-- Sharing State -->
  <div v-else-if="sharing" class="flex h-screen w-screen items-start justify-center">
    <div
      class="flex h-80 w-full max-w-md flex-col items-center justify-center border-neutrals-black bg-gray-light md:border md:border-t-0"
    >
      <SpinnerIcon />
      <span class="mt-8 text-base text-neutrals-dark">{{
        t('CHAPI.Share.sharingCredential')
      }}</span>
    </div>
  </div>
  <!-- Error State -->
  <div
    v-else-if="!showCredentialsMissing && errors.length"
    class="flex h-screen w-screen items-start justify-center"
  >
    <div
      class="flex h-auto w-full max-w-md flex-col items-center justify-center border-neutrals-black bg-gray-light md:border md:border-t-0"
    >
      <div class="flex flex-col items-center justify-start py-16 px-5">
        <img src="@/assets/img/icons-error.svg" />
        <span class="mt-5 mb-3 text-center text-xl font-bold text-neutrals-dark">{{
          t('CHAPI.Share.Error.heading')
        }}</span>
        <span class="text-center text-lg text-neutrals-medium">{{
          t('CHAPI.Share.Error.body')
        }}</span>
      </div>
      <div
        class="flex w-full flex-row items-center justify-center border-t border-neutrals-thistle bg-neutrals-magnolia py-4 px-5"
      >
        <button id="share-credentials-ok-btn" class="btn-primary" @click="cancel">
          {{ t('CHAPI.Share.Error.tryAgain') }}
        </button>
      </div>
    </div>
  </div>
  <!-- Credentials Missing State -->
  <div v-else-if="showCredentialsMissing" class="flex h-screen w-screen items-start justify-center">
    <div
      class="flex h-auto w-full max-w-md flex-col items-center justify-center border-neutrals-black bg-gray-light md:border md:border-t-0"
    >
      <div class="flex flex-col items-center justify-start py-16 px-5">
        <img src="@/assets/img/icons-error.svg" />
        <span class="mt-5 mb-3 text-center text-xl font-bold text-neutrals-dark">{{
          t('CHAPI.Share.CredentialsMissing.heading')
        }}</span>
        <span class="text-center text-lg text-neutrals-medium">{{
          t('CHAPI.Share.CredentialsMissing.body')
        }}</span>
      </div>
      <div
        class="flex w-full flex-row items-center justify-center border-t border-neutrals-thistle bg-neutrals-magnolia py-4 px-5"
      >
        <button id="share-credentials-ok-btn" class="btn-outline" @click="cancel">
          {{ t('CHAPI.Share.CredentialsMissing.ok') }}
        </button>
      </div>
    </div>
  </div>
  <!-- Main State -->
  <div
    v-else
    class="flex h-screen max-h-screen w-screen items-start justify-center overflow-scroll"
  >
    <div class="w-full max-w-md border-neutrals-black bg-gray-light md:border md:border-t-0">
      <div class="p-5">
        <!-- Heading -->
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

        <span class="text-neutrals-dark">{{
          t('CHAPI.Share.headline', processedCredentials.length, { issuer: 'Requestor' })
        }}</span>

        <!-- Credentials Preview -->
        <div
          v-if="processedCredentials.length"
          class="my-6 flex w-full flex-col items-center justify-start"
        >
          <ul class="w-full space-y-5">
            <li v-for="(credential, index) in processedCredentials" :key="index">
              <!-- Credential Preview -->
              <button
                class="group credentialPreviewContainer inline-flex h-20 w-full items-center rounded-xl border p-5 text-sm font-bold focus-within:ring-2 focus-within:ring-offset-2 md:h-24 md:text-base"
                :class="
                  credential.styles.background.color !== '#fff'
                    ? `border-neutrals-black border-opacity-10 notWhiteCredentialPreview`
                    : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`
                "
                :style="focusStyleColor(credential.styles.background.color)"
                @click="toggleDetails(credential)"
              >
                <div class="h-12 w-12 flex-none border-opacity-10">
                  <img :src="getCredentialIconFunction(credential)" />
                </div>
                <div class="grow p-4">
                  <span
                    :class="[
                      `text-sm md:text-base font-bold text-left text-ellipsis`,
                      credential.styles.background.color !== '#fff'
                        ? `text-neutrals-white`
                        : `text-neutrals-dark`,
                    ]"
                  >
                    {{ credential.title }}
                  </span>
                </div>
              </button>
              <!-- Credential Details -->
              <div
                v-if="credential.showDetails"
                class="mt-5 flex w-full flex-col items-start justify-start md:mt-6"
              >
                <span class="py-3 text-base font-bold text-neutrals-dark">What's being shared</span>

                <!-- TODO: move this to reusable components -->
                <table class="w-full border-t border-neutrals-chatelle">
                  <tr
                    v-for="(property, key) of credential.properties"
                    :key="key"
                    class="border-b border-dotted border-neutrals-thistle"
                  >
                    <td v-if="property.value" class="py-4 pr-6 pl-3 text-neutrals-medium">
                      {{ property.label }}
                    </td>
                    <td
                      v-if="property.schema.contentMediaType != 'image/png' && property.value"
                      class="break-words py-4 pr-6 pl-3 text-neutrals-dark"
                    >
                      {{ property.value }}
                    </td>
                    <td
                      v-if="property.schema.contentMediaType === 'image/png' && property.value"
                      class="break-words py-4 pr-6 pl-3 text-neutrals-dark"
                    >
                      <img :src="property.value" class="h-20 w-20" />
                    </td>
                  </tr>
                </table>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <div
        class="sticky bottom-0 flex w-full flex-row items-center justify-between border-t border-neutrals-thistle bg-neutrals-magnolia py-4 px-5"
      >
        <button id="cancelBtn" class="btn-outline" @click="cancel">
          {{ t('CHAPI.Share.decline') }}
        </button>
        <button id="share-credentials" class="btn-primary" @click="share">
          {{ t('CHAPI.Share.share') }}
        </button>
      </div>
    </div>
  </div>
</template>
<script>
import { toRaw } from 'vue';
import { CredentialManager } from '@trustbloc/wallet-sdk';
import {
  normalizeQuery,
  getCredentialIcon,
  prepareCredentialManifest,
  resolveManifest,
} from '@/mixins';
import { mapGetters } from 'vuex';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';
import { useI18n } from 'vue-i18n';

export default {
  components: {
    SpinnerIcon,
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
    const query = normalizeQuery(toRaw(this.protocolHandler.getEventData().query));
    const { user, token } = this.getCurrentUser().profile;
    this.credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });

    try {
      const { results } = await this.credentialManager.query(token, query);
      this.presentation = results;

      const manifest = prepareCredentialManifest(
        this.presentation[0],
        this.getCredentialManifests(),
        this.protocolHandler.requestor()
      );
      this.processedCredentials = await resolveManifest(
        this.credentialManager,
        this.getCredentialManifests(),
        this.token,
        {
          manifest,
          response: this.presentation[0],
        }
      );
    } catch (e) {
      this.errors.push(e);
      console.error('get credentials failed,', e);
      this.loading = false;
    }

    // TODO: governance VC check
    this.requestOrigin = this.protocolHandler.requestor();
    this.loading = false;
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser', 'getCredentialManifests', 'getStaticAssetsUrl']),
    getCredentialIconFunction: function (credential) {
      return credential?.styles?.thumbnail?.uri?.includes('https://')
        ? credential?.styles?.thumbnail?.uri
        : getCredentialIcon(this.getStaticAssetsUrl(), credential?.styles?.thumbnail?.uri);
    },
    toggleDetails(credential) {
      credential.showDetails = !credential.showDetails;
    },
    focusStyleColor(color) {
      return {
        'background-color': color,
        '--focus-color': color,
      };
    },
    async share() {
      this.sharing = true;
      const { profile, preference } = this.getCurrentUser();
      const { controller, proofType, verificationMethod } = preference;
      const { domain, challenge } = toRaw(this.protocolHandler.getEventData());

      const _present = async (presentation) => {
        return (
          await this.credentialManager.present(
            profile.token,
            { presentation },
            {
              controller,
              proofType,
              domain,
              challenge,
              verificationMethod,
            }
          )
        ).presentation;
      };
      try {
        // TODO: match expected format
        const results = await Promise.all(this.presentation.map(_present));
        // typically single presentation, but some verifier queries might produce multiple presentation.
        if (results.length === 1) {
          this.protocolHandler.present(results[0]);
        } else {
          this.protocolHandler.present(results);
        }
      } catch (e) {
        this.errors.push(e);
        console.error('share credentials failed,', e);
        this.sharing = false;
      }

      this.sharing = false;
    },
    cancel() {
      this.protocolHandler.cancel();
    },
  },
};
</script>

<style scoped>
.notWhiteCredentialPreview:focus {
  outline: 2px solid var(--focus-color);
  outline-offset: 2px;
}
</style>
