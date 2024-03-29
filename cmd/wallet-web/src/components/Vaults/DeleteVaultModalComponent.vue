<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <ModalComponent :show="showModal" @close="handleClose">
    <template #content>
      <div class="relative flex w-full flex-col items-center px-8 pt-10">
        <div class="flex h-15 w-15 items-center justify-center rounded-full bg-primary-valencia">
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
          <p class="pb-12 text-center text-base text-neutrals-medium">
            {{ t('Vaults.DeleteModal.deleteVaultConfirmMessage') }}
          </p>
        </div>
      </div>
    </template>
    <template #actionButton>
      <StyledButtonComponent
        class="order-first w-full md:order-last md:w-auto"
        type="btn-danger"
        :loading="loading"
        @click="deleteVault(vaultId)"
      >
        {{ t('Vaults.DeleteModal.deleteVaultButton') }}
      </StyledButtonComponent>
    </template>
  </ModalComponent>
</template>

<script>
import { computed, ref, watch } from 'vue';
import { useStore } from 'vuex';
import { CollectionManager } from '@trustbloc/wallet-sdk';
import { useI18n } from 'vue-i18n';
import { vaultsMutations } from '@/pages/VaultsPage.vue';
import ModalComponent from '@/components/Modal/ModalComponent.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent';

export default {
  name: 'DeleteVaultModalComponent',
  components: {
    StyledButtonComponent,
    ModalComponent,
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
