<template>
  <div>
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
      <span class="pt-5 pb-3 text-lg font-bold text-neutrals-dark align-center">
        {{ i18n.deleteCredential }}?
      </span>
      <div class="relative flex-auto">
        <p class="pb-12 text-base text-center text-neutrals-medium lg:text-start">
          {{ i18n.deleteCredentialConfirmMessage }}
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
        bg-neutrals-magnolia
        rounded-b-2xl
        flex flex-col
        lg:text-start
        md:text-start
        border-t border-0 border-neutrals-lilacSoft
      "
    >
      <button
        class="w-full md:w-auto lg:w-auto btn-outline"
        type="button"
        @click="$emit('onClose')"
      >
        {{ i18n.deleteButtonCancel }}
      </button>
      <button
        id="deleteButton"
        class="order-first md:order-last lg:order-last w-full md:w-auto lg:w-auto btn-danger"
        type="button"
        @click="deleteCredential(credentialId)"
      >
        {{ i18n.deleteButtonLabel }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { CredentialManager } from '@trustbloc-cicd/wallet-sdk';
import base64url from 'base64url';

export default {
  name: 'DeleteCredential',
  props: {
    credentialId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      agent: null,
    };
  },
  computed: {
    i18n() {
      return this.$t('CredentialDetails');
    },
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    async deleteCredential(credentialID) {
      const { user, token } = this.getCurrentUser().profile;
      const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
      const credID = base64url.decode(credentialID);
      try {
        const resp = await credentialManager.remove(token, credID);
        this.$router.push({ name: 'dashboard' });
      } catch (e) {
        console.error('failed to remove credential:', e);
      }
    },
  },
};
</script>
