<template>
  <div
    class="
      overflow-hidden
      absolute
      right-0
      z-20
      mt-2
      w-56
      max-w-max
      text-neutrals-medium
      bg-neutrals-white
      rounded-lg
      shadow-xl
    "
  >
    <div class="p-4 text-base">
      <button id="renameCredential" class="block pb-2 font-bold">
        {{ t('CredentialDetails.renameCredential') }}
      </button>
      <button id="moveCredential" class="block pb-2 font-bold">
        {{ t('CredentialDetails.moveCredential') }}
      </button>
      <button
        id="deleteCredential"
        class="block font-bold text-primary-vampire"
        @click="toggleDeleteCredentialModal()"
      >
        {{ t('CredentialDetails.deleteCredential') }}
      </button>
    </div>
    <!-- todo move to components folder-->
    <div
      v-if="showDeleteCredentialModal"
      class="
        flex
        overflow-y-auto
        fixed
        inset-0
        z-50
        justify-center
        items-center
        bg-neutrals-dark bg-opacity-50
      "
    >
      <div class="relative mx-6 lg:mx-auto max-w-6xl bg-neutrals-white rounded-2xl modal-width">
        <!--content-->
        <div class="flex relative flex-col items-center px-8 pt-10 w-full">
          <div class="flex justify-center items-center w-15 h-15 bg-primary-valencia rounded-full">
            <svg width="32" height="32" xmlns="http://www.w3.org/2000/svg">
              <g transform="translate(3 6)" fill="none" fill-rule="evenodd">
                <rect stroke="#ffffff" stroke-width="2" x="1" y="1" width="24" height="19" rx="4" />
                <ellipse fill="#ffffff" cx="8" cy="13" rx="4" ry="2" />
                <circle fill="#ffffff" cx="8" cy="8" r="2" />
                <path
                  stroke="#ffffff"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M16 8h5M16 13h2"
                />
              </g>
            </svg>
          </div>
          <span class="pt-5 pb-3 text-lg font-bold text-neutrals-dark">
            {{ t('CredentialDetails.deleteCredential') }}?
          </span>
          <div class="relative flex-auto">
            <p class="pb-12 text-base text-center lg:text-left text-neutrals-medium">
              {{ t('CredentialDetails.deleteCredentialConfirmMessage') }}
            </p>
          </div>
        </div>
        <!--footer-->
        <div
          class="
            md:flex-row
            lg:flex-row
            gap-4
            justify-start
            md:justify-between
            lg:justify-between
            items-center
            px-5
            md:px-8
            lg:px-8
            pt-4
            pb-5
            text-center
            md:text-left
            lg:text-left
            bg-neutrals-magnolia
            rounded-b-2xl
            flex flex-col
            border-t border-0 border-neutrals-lilacSoft
          "
        >
          <button
            class="w-full md:w-auto lg:w-auto btn-outline"
            type="button"
            @click="toggleDeleteCredentialModal()"
          >
            {{ t('CredentialDetails.deleteButtonCancel') }}
          </button>
          <button
            id="deleteButton"
            class="order-first md:order-last lg:order-last w-full md:w-auto lg:w-auto btn-danger"
            type="button"
            @click="deleteCredential(credentialId)"
          >
            {{ t('CredentialDetails.deleteButtonLabel') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { decode } from 'js-base64';
import { useI18n } from 'vue-i18n';

export default {
  name: 'FlyoutMenuList',
  props: {
    credentialId: {
      type: String,
      required: true,
    },
  },
  setup() {
    const { t, locale } = useI18n();
    return { t, locale };
  },
  data() {
    return {
      showDeleteCredentialModal: false,
      agent: null,
    };
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    async deleteCredential(credentialID) {
      const { user, token } = this.getCurrentUser().profile;
      const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
      const credID = decode(credentialID);
      try {
        await credentialManager.remove(token, credID);
        this.$router.push({ name: 'dashboard' });
      } catch (e) {
        console.error('failed to remove credential:', e);
      }
    },
    toggleDeleteCredentialModal() {
      this.showDeleteCredentialModal = !this.showDeleteCredentialModal;
    },
  },
};
</script>

<style scoped>
.modal-width {
  width: 28rem;
  box-shadow: 0px 12px 48px 0px rgba(25, 12, 33, 0.2);
}
</style>
