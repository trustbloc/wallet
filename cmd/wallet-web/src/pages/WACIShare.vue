<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

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
    <!-- Sharing State -->
    <div
      v-else-if="sharing"
      class="
        justify-center
        items-center
        w-full
        max-w-md
        h-80
        bg-gray-light
        rounded-lg
        md:border
        border-neutrals-chatelle
        flex flex-col
      "
    >
      <Spinner />
      <span class="mt-8 text-base text-neutrals-dark">{{
        $t('CHAPI.Share.sharingCredential')
      }}</span>
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
          $t('CHAPI.Share.Error.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          $t('CHAPI.Share.Error.body')
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
          {{ $t('CHAPI.Share.Error.tryAgain') }}
        </button>
      </div>
    </div>
    <!-- Credentials Missing State -->
    <div
      v-else-if="showCredentialsMissing"
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
          $t('CHAPI.Share.CredentialsMissing.heading')
        }}</span>
        <span class="text-lg text-center text-neutrals-medium">{{
          $t('CHAPI.Share.CredentialsMissing.body')
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
          {{ $t('CHAPI.Share.CredentialsMissing.ok') }}
        </button>
      </div>
    </div>
    <!-- Main State -->
    <div v-else class="w-full max-w-md bg-gray-light rounded-lg md:border border-neutrals-chatelle">
      <div class="p-5">
        <!-- Heading -->
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

        <span class="text-neutrals-dark">{{
          $tc('CHAPI.Share.headline', processedCredentials.length, { issuer: 'Requestor' })
        }}</span>

        <!-- Credentials Preview -->
        <div
          v-if="processedCredentials.length"
          class="flex flex-col justify-start items-center mt-6 mb-6 w-full"
        >
          <ul class="space-y-5 w-full">
            <li v-for="(credential, index) in processedCredentials" :key="index">
              <!-- Credential Preview -->
              <button
                :class="[
                  `group inline-flex items-center rounded-xl p-5 text-sm md:text-base font-bold border w-full h-20 md:h-24 focus-within:ring-2 focus-within:ring-offset-2 credentialPreviewContainer`,
                  credential.brandColor.length
                    ? `bg-gradient-${credential.brandColor} border-neutrals-chatelle border-opacity-10 focus-within:ring-primary-${credential.brandColor}`
                    : `bg-neutrals-white border-neutrals-thistle hover:border-neutrals-chatelle focus-within:ring-neutrals-victorianPewter`,
                ]"
                @click="toggleDetails(credential)"
              >
                <div class="flex-none w-12 h-12 border-opacity-10">
                  <img :src="require(`@/assets/img/${credential.icon}`)" />
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
                    v-for="(property, index) of credential.properties"
                    :key="index"
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
      </div>

      <div
        class="
          justify-between
          items-center
          pt-4
          pr-5
          pb-4
          pl-5
          w-full
          bg-neutrals-magnolia
          rounded-b-lg
          flex flex-row
          border-t border-neutrals-thistle
        "
      >
        <button id="cancelBtn" class="btn-outline" @click="cancel">
          {{ $t('CHAPI.Share.decline') }}
        </button>
        <button id="share-credentials" class="btn-primary" @click="share">
          {{ $t('CHAPI.Share.share') }}
        </button>
      </div>
    </div>
  </div>
</template>
<script>
import { DIDComm } from '@trustbloc/wallet-sdk';
import { wait, getCredentialType, getCredentialDisplayData } from './mixins';
import { mapGetters } from 'vuex';
import Spinner from '@/components/Spinner/Spinner.vue';
import credentialDisplayData from '@/config/credentialDisplayData';

export default {
  components: {
    Spinner,
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
    const invitation = this.protocolHandler.message();
    const { user, token } = this.getCurrentUser().profile;

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
      this.errors.push('Error initiating credential share');
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
    ...mapGetters(['getCurrentUser']),
    toggleDetails(credential) {
      credential.showDetails = !credential.showDetails;
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
          this.processedCredentials.push({ ...processedCredential, showDetails: false });
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

      try {
        await this.didcomm.completeCredentialShare(
          profile.token,
          this.threadID,
          this.presentations,
          {
            controller,
            proofType,
            verificationMethod,
          },
          true
        );
      } catch (e) {
        this.errors.push(e);
        console.error('share credentials failed,', e);
        this.sharing = false;
        return;
      }

      //TODO this delay to be removed once we have WACI ack feature available for present proof.
      await wait(2000);

      this.protocolHandler.done();
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
      return credentialDisplayData[currentCredentialType] || credentialDisplayData.fallback;
    },
  },
};
</script>
