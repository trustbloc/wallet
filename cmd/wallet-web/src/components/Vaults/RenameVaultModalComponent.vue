<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <ModalComponent :show="showModal" :show-close-button="true" @close="handleClose">
    <template #content>
      <div
        class="flex w-full items-center justify-between border-b-2 border-neutrals-thistle py-4 px-8"
      >
        <span class="text-lg font-bold text-neutrals-dark">
          {{ t('Vaults.RenameModal.renameVault') }}
        </span>
      </div>
      <div class="flex w-full items-center px-8 pt-10">
        <InputFieldComponent
          v-model="vaultName"
          type="text"
          :label="t('Vaults.label')"
          :value="vaultName"
          :helper-message="t('Vaults.helperMessage')"
          :error-message="errorMessage"
          :placeholder="t('Vaults.placeholderLabel')"
          :pattern="pattern"
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
      <StyledButtonComponent
        class="order-first w-full md:order-last md:w-auto"
        type="btn-primary"
        :loading="loading"
        @click="renameVault(vaultId)"
      >
        {{ t('Vaults.RenameModal.rename') }}
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
import InputFieldComponent from '@/components/InputField/InputFieldComponent';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent';

export default {
  name: 'RenameVaultModalComponent',
  components: { StyledButtonComponent, InputFieldComponent, ModalComponent },
  props: {
    show: {
      type: Boolean,
      default: false,
    },
    existingNames: {
      type: Array,
      required: true,
    },
    vaultId: {
      type: String,
      required: true,
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
      vaultName: '',
      loading: false,
      submitted: false,
    };
  },
  computed: {
    errorMessage() {
      if (this.submitted && this.nameEmpty) {
        return this.t('Vaults.emptyError');
      } else if (!this.nameValid) {
        return this.t('Vaults.patternError');
      } else if (this.nameAlreadyExists) {
        return this.t('Vaults.alreadyExistsError');
      } else return '';
    },
    // Returns true if the length of the trimmed name is 0
    nameEmpty() {
      return this.vaultName.trim().replace(/ +/g, '\\s+').length === 0;
    },
    // Returns true if the name only contains letters, numbers or whitespace characters
    nameValid() {
      const re = new RegExp('^[a-zA-Z\\d]+(?:[a-zA-Z\\d\\s]+)*$');
      return re.test(this.vaultName.trim().replace(/  +/g, ' '));
    },
    // Returns true if the trimmed name matches one of the existing vaults' names
    nameAlreadyExists() {
      let existingNamesRegex = `(?=^${this.existingNames[0].replace(/ +/g, '\\s+')}$)`;
      this.existingNames.map((name) => {
        existingNamesRegex += `|(?=^${name.replace(/ +/g, '\\s+')}$)`;
      });
      const re = new RegExp(existingNamesRegex);
      return re.test(this.vaultName.trim().replace(/  +/g, ' '));
    },
    // Returns a string regex to pass down to the InputField's pattern attribute
    pattern() {
      let existingNamesRegex = '';
      this.existingNames.map((name) => {
        existingNamesRegex += `(?!^${name.replace(/ +/g, '\\s+')}$)`;
      });
      return `${existingNamesRegex}(^[a-zA-Z\\d]+(?:[a-zA-Z\\d\\s]+)*$)`;
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
    updateVaultName({ name }) {
      this.vaultName = name;
      this.submitted = false;
    },
    async renameVault(vaultId) {
      this.submitted = true;
      if (this.vaultName.length && this.nameValid && !this.nameAlreadyExists) {
        this.loading = true;
        const { user, token } = this.currentUser.profile;
        const collectionManager = new CollectionManager({ agent: this.agentInstance, user });
        try {
          await collectionManager.update(token, vaultId, { name: this.vaultName });
          vaultsMutations.setVaultsOutdated(true);
          this.showModal = false;
          this.loading = false;
          this.submitted = false;
          this.$emit('close');
        } catch (e) {
          console.error('Error while renaming vault:', e);
        }
      }
    },
    handleClose() {
      this.showModal = false;
      this.$emit('close');
    },
  },
};
</script>
