<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <modal  v-if="showModal">
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
            <p class="pb-12 text-base text-center text-neutrals-medium">
              {{ t('CredentialDetails.deleteCredentialConfirmMessage') }}
            </p>
          </div>
        </div>
        <!-- Buttons Container -->
        <div
          class="
            md:flex-row
            gap-4
            justify-start
            md:justify-between
            items-center
            px-5
            md:px-8
            pt-4
            pb-5
            text-center
            bg-neutrals-magnolia
            rounded-b-2xl
            flex flex-col
            modal-footer
            border-t border-0 border-neutrals-lilacSoft
          "
        >
          <styled-button class="btn-outline" type="outline" @click="closeModal">
            {{ t('App.cancel') }}
          </styled-button>
          <styled-button
            id="deleteButton"
            class="order-first md:order-last lg:order-last"
            type="danger"
            @click="deleteCredential(credentialId)"
          >
            {{ t('CredentialDetails.deleteButtonLabel') }}
          </styled-button>
        </div>
  </modal>
</template>

<script>
import Modal from '@/components/Modal/Modal.vue';
import {ref, watch} from 'vue';
import {mapGetters} from 'vuex';
import {CredentialManager} from '@trustbloc-cicd/wallet-sdk';
import {decode} from 'js-base64';
import {useI18n} from 'vue-i18n';
import StyledButton from "@/components/StyledButton/StyledButton";

const props = {
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
};
export default {
  name: 'DeleteCredentialModal',
  components: {
    StyledButton,
    Modal,
  },
  props,
  setup(props) {
    const { t } = useI18n();
    const showModal = ref(false);
    watch(
      () => props.show,
      (show) => {
        showModal.value = show;
      }
    );

    function closeModal() {
      showModal.value = false;
    }

    return {
      t,
      showModal,
      closeModal,
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
    async deleteCredential(credentialID) {
      const { user, token } = this.getCurrentUser().profile;
      const credentialManager = new CredentialManager({ agent: this.getAgentInstance(), user });
      const credID = decode(credentialID);
      try {
        const resp = await credentialManager.remove(token, credID);
        this.$router.push({ name: 'credentials' });
      } catch (e) {
        console.error('failed to remove credential:', e);
      }
    },
  },
};
</script>
<style scoped>
.modal-width {
  width: 28rem;
}

.modal-footer {
  box-shadow: inset 0px 1px 0px 0px rgb(219, 215, 220);
}
</style>
