<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <ModalComponent :show="showModal">
    <template #content>
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
          {{ t('CredentialDetails.DeleteModal.deleteCredential') }}?
        </span>
        <div class="relative flex-auto">
          <p class="pb-12 text-base text-center text-neutrals-medium">
            {{ t('CredentialDetails.DeleteModal.deleteCredentialConfirmMessage') }}
          </p>
        </div>
      </div>
    </template>
    <template #actionButton>
      <StyledButtonComponent
        id="delete-credential-button"
        class="order-first md:order-last w-full md:w-auto"
        type="btn-danger"
        @click="deleteCredential()"
      >
        {{ t('CredentialDetails.DeleteModal.deleteButtonLabel') }}
      </StyledButtonComponent>
    </template>
  </ModalComponent>
</template>

<script>
import ModalComponent from '@/components/Modal/ModalComponent.vue';
import { ref, watch } from 'vue';
import { mapGetters } from 'vuex';
import { CredentialManager } from '@trustbloc/wallet-sdk';
import { useI18n } from 'vue-i18n';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent';

export default {
  name: 'DeleteCredentialModalComponent',
  components: {
    StyledButtonComponent,
    ModalComponent,
  },
  props: {
    target: {
      type: String,
      default: 'body',
    },
    credentialId: {
      type: String,
      required: true,
    },
    show: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const showModal = ref(false);
    watch(
      () => props.show,
      (show) => {
        showModal.value = show;
      }
    );

    return {
      props,
      showModal,
      t,
    };
  },
  data() {
    return {
      agent: null,
    };
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    async deleteCredential() {
      const { user, token } = this.getCurrentUser().profile;
      const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
      try {
        await credentialManager.remove(token, this.props.credentialId);
        this.$router.push({ name: 'credentials' });
      } catch (e) {
        console.error('failed to remove credential:', e);
      }
    },
  },
};
</script>
