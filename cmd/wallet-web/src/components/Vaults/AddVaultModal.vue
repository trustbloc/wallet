<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <modal :show="showModal" :show-close-button="true">
    <template #content>
      <!--content-->
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
          {{ t('Vaults.addVault') }}
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
          @update="receivedVaultName($event)"
        />
      </div>
    </template>
    <!--button container-->
    <template #actionButton>
      <styled-button
        class="order-first md:order-last lg:order-last"
        type="primary"
        @click="addVault"
      >
        {{ t('Vaults.add') }}
      </styled-button>
    </template>
  </modal>
</template>

<script>
import Modal from '@/components/Modal/Modal.vue';
import InputField from '@/components/InputField/InputField';
import { CollectionManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import { useI18n } from 'vue-i18n';
import { vaultsMutations } from '@/pages/Vaults.vue';
import StyledButton from '@/components/StyledButton/StyledButton';
import { ref, watch } from 'vue/dist/vue.esm-bundler';

const props = {
  show: {
    type: Boolean,
    default: false,
  },
};
export default {
  name: 'AddVault',
  components: { StyledButton, InputField, Modal },
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

    return {
      t,
      showModal,
    };
  },
  data() {
    return {
      vaultName: '',
    };
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    receivedVaultName(data) {
      this.vaultName = data;
    },
    async addVault() {
      // Integration with collections API.
      const { user, token } = this.getCurrentUser().profile;
      const collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
      const name = this.vaultName;
      try {
        const id = await collectionManager.create(token, { name });
        if (id) {
          vaultsMutations.setVaultsOutdated(true);
          this.showModal = false;
        }
        // TODO: add an error state to display in the UI
      } catch (e) {
        console.error('Error creating a new vault:', e);
      }
    },
  },
};
</script>
