<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <modal :show="showModal" :show-close-button="true">
    <template #content>
      <div
        class="
          flex
          justify-between
          items-center
          py-4
          px-8
          w-full
          border-b-2 border-neutrals-thistle
        "
      >
        <span class="text-lg font-bold text-neutrals-dark">
          {{ t('Vaults.renameVault') }}
        </span>
      </div>
      <div class="flex items-center px-8 pt-10 w-full">
        <input-field
          v-model="vaultName"
          :helper-message="t('Vaults.addHelperMessage')"
          :label="t('Vaults.addlabel')"
          :placeholder="t('Vaults.placeholderLabel')"
          :value="vaultName"
          type="text"
          maxlength="42"
          @input="updateVaultName($event)"
        />
      </div>
    </template>
    <template #actionButton>
      <styled-button
        class="order-first md:order-last w-full md:w-auto"
        type="primary"
        :loading="loading"
        @click="renameVault(vaultId)"
      >
        {{ t('Vaults.rename') }}
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
import InputField from '@/components/InputField/InputField';
import StyledButton from '@/components/StyledButton/StyledButton';

export default {
  name: 'RenameVaultModal',
  components: { StyledButton, InputField, Modal },
  props: {
    show: {
      type: Boolean,
      default: false,
    },
    vaultId: {
      type: String,
      required: true,
    },
    selectedVaultName: {
      type: String,
      required: true,
    },
  },
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
      vaultName: '',
      loading: false,
    };
  },
  methods: {
    updateVaultName(name) {
      this.vaultName = name;
    },
    async renameVault(vaultId) {
      this.loading = true;
      if (this.selectedVaultName === this.vaultName) {
        this.loading = false;
        this.showModal = false;
      }
      const { user, token } = this.currentUser.profile;
      const collectionManager = new CollectionManager({ agent: this.agentInstance, user });
      try {
        await collectionManager.update(token, vaultId, { name: this.vaultName });
        vaultsMutations.setVaultsOutdated(true);
        this.showModal = false;
        this.loading = false;
        // TODO: add an error state to display in the UI
      } catch (e) {
        console.error('Error while renaming vault:', e);
      }
    },
  },
};
</script>
