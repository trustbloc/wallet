<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <!-- Loading state -->
  <div v-if="loading" class="flex justify-center items-start w-screen h-screen">
    <div
      class="
        flex
        justify-center
        items-center
        w-full
        max-w-md
        h-80
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
      "
    >
      <Spinner />
    </div>
  </div>
  <!-- Sharing State -->
  <div v-else-if="sharing" class="flex justify-center items-start w-screen h-screen">
    <div
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-80
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
        flex flex-col
      "
    >
      <Spinner />
      <span class="mt-8 text-base text-neutrals-dark">{{
        t('CHAPI.Share.sharingCredential')
      }}</span>
    </div>
  </div>
  <!-- Error State -->
  <div
    v-else-if="!showCredentialsMissing && errors.length"
    class="flex justify-center items-start w-screen h-screen"
  >
    <div
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-auto
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
        flex flex-col
      "
    >
      <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
        <img src="@/assets/img/icons-error.svg" />
        <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
          t('CHAPI.Share.Error.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          t('CHAPI.Share.Error.body')
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
          {{ t('CHAPI.Share.Error.tryAgain') }}
        </button>
      </div>
    </div>
  </div>
  <!-- Credentials Missing State -->
  <div v-else-if="showCredentialsMissing" class="flex justify-center items-start w-screen h-screen">
    <div
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-auto
        bg-gray-light
        md:border md:border-t-0
        border-neutrals-black
        flex flex-col
      "
    >
      <div class="flex flex-col justify-start items-center pt-16 pr-5 pb-16 pl-5">
        <img src="@/assets/img/icons-error.svg" />
        <span class="mt-5 mb-3 text-xl font-bold text-center text-neutrals-dark">{{
          t('CHAPI.Share.CredentialsMissing.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          t('CHAPI.Share.CredentialsMissing.body')
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
        <button id="share-credentials-ok-btn" class="btn-outline" @click="cancel">
          {{ t('CHAPI.Share.CredentialsMissing.ok') }}
        </button>
      </div>
    </div>
  </div>
  <!-- Main State -->
  <div
    v-else
    class="flex overflow-scroll justify-center items-start w-screen h-screen max-h-screen"
  >
    <div class="w-full max-w-md bg-gray-light md:border md:border-t-0 border-neutrals-black">
      <div class="p-5">
        <!-- Heading -->
        <div class="flex flex-row justify-start items-start mb-4 w-full">
          <div class="flex-none w-12 h-12 border-opacity-10">
            <!-- TODO: issue-1055 Read meta data from external urls -->
            <img src="@/assets/img/generic-issuer-icon.svg" />
          </div>
          <div class="flex flex-col pl-3">
            <span
              class="flex-1 mb-1 text-sm font-bold text-left text-neutrals-dark overflow-ellipsis"
            >
              <!-- TODO: issue-1055 Read meta data from external urls -->
              Verifier
            </span>
            <div class="flex flex-row justify-center items-center">
              <img src="@/assets/img/small-lock-icon.svg" />
              <span class="flex-1 pl-1 text-xs text-left text-neutrals-medium overflow-ellipsis">
                {{ requestOrigin }}
              </span>
            </div>
          </div>
        </div>

        <span class="text-neutrals-dark">{{
          t('CHAPI.Share.headline', credsFound.length, { issuer: 'Verifier' })
        }}</span>
        <div
          v-if="credsFound.length"
          class="flex flex-col justify-start items-center mt-6 mb-6 w-full"
        >
          <ul class="space-y-5 w-full">
            <li v-for="(credential, index) in credsFound" :key="index">
              <!-- Credential Preview -->
              <button
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
              </button>
              <!-- Credential Details -->
              <div
                v-if="credential.showDetails"
                class="flex flex-col justify-start items-start mt-5 md:mt-6 w-full"
              >
                <span class="py-3 text-base font-bold text-neutrals-dark">What's being shared</span>

                <!-- TODO: move this to reusable components -->
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

        <div v-if="credsFound.length && issuersFound.length" style="margin: 30px"></div>

        <!-- Issuers Container -->
        <div v-if="issuersFound.length">
          <p id="result-header-2">
            Found {{ issuersFound.length }} issuer{{ issuersFound.length > 1 ? 's' : '' }} who can
            issue credentials matching above criteria in your wallet:
          </p>

          <ul>
            <li
              v-for="(issuer, key) in issuersFound"
              :id="'issuer-' + key"
              :key="key"
              class="
                group
                items-center
                p-5
                mb-5
                w-full
                h-20
                md:h-24
                rounded-xl
                border
                focus-within:ring-2 focus-within:ring-offset-2
                flex flex-col
              "
            >
              <span class="text-lg font-bold">{{ issuer.name }}</span>
              <span class="text-sm md:text-base">{{ issuer.description }}</span>
            </li>
          </ul>
        </div>
      </div>

      <!-- Bottom Buttons Container -->
      <div
        class="
          sticky
          bottom-0
          justify-between
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
        <button id="cancelBtn" class="btn-outline" @click="cancel">
          {{ t('CHAPI.Share.decline') }}
        </button>
        <button id="share-credentials" class="btn-primary" @click="createPresentation">
          {{ t('CHAPI.Share.share') }}
        </button>
      </div>
    </div>
  </div>
</template>
<script>
import {
  filterCredentialsByType,
  getCredentialType,
  getCredentialDisplayData,
  getCredentialIcon,
  WalletGetByQuery,
} from '@/mixins';
import { mapGetters } from 'vuex';
import Spinner from '@/components/Spinner/Spinner.vue';
import { useI18n } from 'vue-i18n';

const manifestCredType = 'IssuerManifestCredential';

export default {
  components: { Spinner },
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
      credsFound: [],
      issuersFound: [],
      credentialDisplayData: '',
    };
  },
  computed: {
    showCredentialsMissing() {
      return this.credsFound.length === 0;
    },
  },
  created: async function () {
    const { user, token } = this.getCurrentUser().profile;
    this.wallet = new WalletGetByQuery(
      this.getAgentInstance(),
      this.$parent.protocolHandler,
      this.getAgentOpts(),
      user
    );
    // make sure mediator is connected
    await this.wallet.connectMediator();
    this.credentialDisplayData = await this.getCredentialManifestData();

    this.requestOrigin = this.$parent.protocolHandler.requestor();

    try {
      this.presentation = await this.wallet.getPresentationSubmission(token);
    } catch (e) {
      this.errors.push(e);
      console.error('get credentials failed,', e);
      this.loading = false;
      return;
    }
    const credentials = this.presentation.verifiableCredential;
    const credsFound = filterCredentialsByType(credentials, [manifestCredType]);
    credsFound.map((credential) => {
      const manifest = this.getManifest(credential);
      const processedCredential = this.getCredentialDisplayData(credential, manifest);
      this.credsFound.push({ ...processedCredential, showDetails: false });
    });
    this.issuersFound = filterCredentialsByType(credentials, [manifestCredType], true);
    this.loading = false;
  },
  methods: {
    ...mapGetters([
      'getCurrentUser',
      'getAgentOpts',
      'getCredentialManifestData',
      'getStaticAssetsUrl',
    ]),
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    toggleDetails(credential) {
      credential.showDetails = !credential.showDetails;
    },
    createPresentation: async function () {
      this.sharing = true;
      try {
        await this.wallet.createAndSendPresentation(this.getCurrentUser(), this.presentation);
      } catch (e) {
        console.error(e);
        this.errors.push('share credentials failed,', e);
        this.sharing = false;
      }
      this.sharing = false;
    },
    cancel: function () {
      this.wallet.cancel();
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
  },
};
</script>
