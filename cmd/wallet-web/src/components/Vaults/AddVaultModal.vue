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
          {{ t('Vaults.AddModal.addVault') }}
        </span>
      </div>
      <div class="flex items-center px-8 pt-10 w-full">
        <input-field
          v-model="vaultName"
          type="text"
          :label="t('Vaults.label')"
          :value="vaultName"
          :helper-message="t('Vaults.helperMessage')"
          :error-message="errorMessage"
          :placeholder="t('Vaults.placeholderLabel')"
          pattern="^[a-zA-Z\d]+$"
          required
          :empty-error="t('Vaults.emptyError')"
          minlength="1"
          maxlength="42"
          :submitted="submitted"
          autocomplete="off"
          @input="updateVaultName"
        />
      </div>
    </template>
    <template #actionButton>
      <styled-button
        class="order-first md:order-last w-full md:w-auto"
        type="primary"
        :loading="loading"
        @click="addVault"
      >
        {{ t('Vaults.AddModal.add') }}
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
  name: 'AddVault',
  components: { StyledButton, InputField, Modal },
  props: {
    show: {
      type: Boolean,
      default: false,
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
      submitted: false,
      nameValid: false,
    };
  },
  computed: {
    errorMessage() {
      if (this.submitted && !this.vaultName.length) {
        return this.t('Vaults.emptyError');
      } else if (!this.nameValid) {
        return this.t('Vaults.patternError');
        // TODO: issue-1359 - check if name already exists
      } else if (false) {
        return this.t('Vaults.alreadyExistsError');
      } else return '';
    },
  },
  watch: {
    showModal: function () {
      if (!this.showModal.value) {
        this.submitted = false;
        this.vaultName = '';
      }
    },
  },
  methods: {
    updateVaultName({ name, valid }) {
      this.vaultName = name;
      this.nameValid = valid;
      this.submitted = false;
    },
    async addVault() {
      this.submitted = true;
      if (this.vaultName.length && this.nameValid) {
        this.loading = true;
        const { user, token } = this.currentUser.profile;
        const collectionManager = new CollectionManager({ agent: this.agentInstance, user });
        // await new Promise((resolve) => setTimeout(resolve, 100000));
        try {
          const id = await collectionManager.create(token, { name: this.vaultName });
          if (id) {
            vaultsMutations.setVaultsOutdated(true);
            this.showModal = false;
            this.loading = false;
            this.submitted = false;
          }
          // TODO: add an error state to display in the UI
        } catch (e) {
          console.error('Error creating a new vault:', e);
        }
      }
    },
  },
};
</script>
