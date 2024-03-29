<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <ModalComponent :show="showModal" :show-close-button="true" @close="handleClose">
    <template #errorToast>
      <ToastNotificationComponent
        v-if="systemError"
        :title="t('CredentialDetails.RenameModal.errorToast.title')"
        :description="t('CredentialDetails.RenameModal.errorToast.description')"
        type="error"
      />
    </template>
    <template #content>
      <div
        class="flex w-full items-center justify-between border-b-2 border-neutrals-thistle py-4 px-8"
      >
        <span class="text-lg font-bold text-neutrals-dark">
          {{ t('CredentialDetails.RenameModal.renameCredential') }}
        </span>
      </div>
      <div class="flex w-full items-center px-8 pt-10">
        <InputFieldComponent
          :v-model="credentialName"
          type="text"
          :label="t('CredentialDetails.RenameModal.label')"
          :value="credentialName"
          :helper-message="t('CredentialDetails.helperMessage')"
          :error-message="errorMessage"
          :placeholder="t('CredentialDetails.RenameModal.placeholderLabel')"
          :pattern="pattern"
          required
          :empty-error="t('CredentialDetails.emptyError')"
          minlength="1"
          maxlength="52"
          :submitted="submitted"
          autocomplete="off"
          @input="updateCredentialName"
        />
      </div>
    </template>
    <template #actionButton>
      <StyledButtonComponent
        id="rename-credential-button"
        class="order-first w-full md:order-last md:w-auto"
        type="btn-primary"
        :loading="loading"
        @click="renameCredential(credentialId, vaultName)"
      >
        {{ t('CredentialDetails.RenameModal.rename') }}
      </StyledButtonComponent>
    </template>
  </ModalComponent>
</template>

<script>
import { computed, ref, watch } from 'vue';
import { useStore } from 'vuex';
import { CollectionManager, CredentialManager } from '@trustbloc/wallet-sdk';
import { useI18n } from 'vue-i18n';
import ModalComponent from '@/components/Modal/ModalComponent.vue';
import InputFieldComponent from '@/components/InputField/InputFieldComponent';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent';
import ToastNotificationComponent from '@/components/ToastNotification/ToastNotificationComponent';
import { credentialMutations } from '@/pages/CredentialDetailsPage';

export default {
  name: 'RenameCredentialModalComponent',
  components: {
    StyledButtonComponent,
    InputFieldComponent,
    ModalComponent,
    ToastNotificationComponent,
  },
  props: {
    show: {
      type: Boolean,
      default: false,
    },
    credentialId: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    vaultName: {
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
      credentialName: '',
      loading: false,
      submitted: false,
      systemError: false,
    };
  },
  computed: {
    errorMessage() {
      if (this.submitted && this.nameEmpty) {
        return this.t('CredentialDetails.emptyError');
      } else if (!this.nameValid) {
        return this.t('CredentialDetails.patternError');
      } else return '';
    },
    // Returns true if the length of the trimmed name is 0
    nameEmpty() {
      return this.credentialName.trim().replace(/ +/g, '\\s+').length === 0;
    },
    // Returns true if the name only contains letters, numbers or whitespace characters
    nameValid() {
      const re = new RegExp('^[a-zA-Z\\d]+(?:[a-zA-Z\\d\\s]+)*$');
      return re.test(this.credentialName.trim().replace(/  +/g, ' '));
    },
    pattern() {
      return `(^[a-zA-Z\\d]+(?:[a-zA-Z\\d\\s]+)*$)`;
    },
  },
  watch: {
    showModal: function () {
      if (!this.showModal.value) {
        this.submitted = false;
        this.credentialName = this.name;
      }
    },
  },
  methods: {
    updateCredentialName({ name }) {
      this.credentialName = name;
      this.submitted = false;
    },
    async renameCredential(credentialId, vaultName) {
      this.submitted = true;
      if (this.credentialName.length && this.nameValid) {
        this.loading = true;
        const { user, token } = this.currentUser.profile;
        const credentialManager = new CredentialManager({ agent: this.agentInstance, user });
        const collectionManager = new CollectionManager({ agent: this.agentInstance, user });
        const { contents } = await collectionManager.getAll(token);
        this.vaults = Object.values(contents).map((vault) => vault);
        const vaultID = this.vaults.find((vault) => vault.name === vaultName).id;
        try {
          await credentialManager.updateCredentialMetadata(token, credentialId, {
            name: this.credentialName,
            collection: vaultID,
          });
          credentialMutations.setCredentialOutdated(true);
          this.showModal = false;
          this.loading = false;
          this.submitted = false;
          this.$emit('close');
        } catch (e) {
          this.systemError = true;
          this.loading = false;
          console.error('Error while renaming credential:', e);
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
