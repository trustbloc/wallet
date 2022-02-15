<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <modal :show="showModal" @close="handleClose">
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
          {{ t('Vaults.DeleteModal.deleteVault') }}?
        </span>
        <div class="relative flex-auto">
          <p class="pb-12 text-base text-center text-neutrals-medium">
            {{ t('Vaults.DeleteModal.deleteVaultConfirmMessage') }}
          </p>
        </div>
      </div>
    </template>
    <template #actionButton>
      <styled-button
        class="order-first md:order-last w-full md:w-auto"
        type="danger"
        :loading="loading"
        @click="deleteVault(vaultId)"
      >
        {{ t('Vaults.DeleteModal.deleteVaultButton') }}
      </styled-button>
    </template>
  </modal>
</template>

<script>
import { computed, ref, watch } from 'vue';
import { useStore } from 'vuex';
import { CollectionManager } from '@trustbloc/wallet-sdk';
import { useI18n } from 'vue-i18n';
import { vaultsMutations } from '@/pages/Vaults.vue';
import Modal from '@/components/Modal/Modal.vue';
import StyledButton from '@/components/StyledButton/StyledButton';

export default {
  name: 'DeleteVaultModal',
  components: {
    StyledButton,
    Modal,
  },
  props: {
    target: {
      type: String,
      default: 'body',
    },
    vaultId: {
      type: String,
      required: true,
    },
    show: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['close'],
  setup(props) {
    const store = useStore();
    const agentInstance = computed(() => store.getters['agent/getInstance']);
    const currentUser = computed(() => store.getters.getCurrentUser);
    const { t } = useI18n();
    const showModal = ref(false);
    watch(
      () => props.show,
      (show) => {
        showModal.value = show;
      }
    );

    return { agentInstance, currentUser, showModal, t };
  },
  data() {
    return {
      loading: false,
    };
  },
  methods: {
    async deleteVault(vaultId) {
      this.loading = true;
      const { user, token } = this.currentUser.profile;
      const collectionManager = new CollectionManager({ agent: this.agentInstance, user });
      try {
        await collectionManager.remove(token, vaultId);
        vaultsMutations.setVaultsOutdated(true);
        this.showModal = false;
        this.loading = false;
        this.$emit('close');
      } catch (e) {
        console.error('Error removing a vault:', e);
      }
    },
    handleClose() {
      this.showModal = false;
      this.$emit('close');
    },
  },
};
</script>
