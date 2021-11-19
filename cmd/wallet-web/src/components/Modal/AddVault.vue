<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <!-- TODO:Issue-1194 Move this to resuable modal component once vue3 resuable modal component change gets implemented -->
  <div
    class="
      flex
      overflow-y-auto
      fixed
      inset-0
      z-50
      justify-center
      items-center
      bg-neutrals-dark bg-opacity-50
    "
  >
    <div
      class="
        relative
        mx-6
        lg:mx-auto
        max-w-6xl
        bg-neutrals-white
        rounded-2xl
        opacity-100
        modal-width
      "
    >
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
        <div>
          <!-- TODO: use inline svg instead once https://github.com/trustbloc/edge-agent/issues/816 is fixed -->
          <img
            class="w-6 h-6 cursor-pointer"
            src="@/assets/img/Icons-sm--close-icon.svg"
            alt="Close Icon"
            @click="closeModal"
          />
        </div>
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
          border-t border-0 border-neutrals-thistle
        "
      >
        <button class="w-full md:w-auto lg:w-auto btn-outline" type="button" @click="closeModal">
          {{ t('Vaults.cancel') }}
        </button>
        <button
          id="deleteButton"
          class="order-first md:order-last lg:order-last w-full md:w-auto lg:w-auto btn-primary"
          type="button"
          @click="addVault"
        >
          {{ t('Vaults.add') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import InputField from '@/components/InputField/InputField';
import { CollectionManager } from '@trustbloc/wallet-sdk';
import { mapGetters } from 'vuex';
import { useI18n } from 'vue-i18n';

export default {
  name: 'AddVault',
  components: { InputField },
  emits: ['close'],
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      vaultName: '',
    };
  },
  methods: {
    ...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    ...mapGetters(['getCurrentUser']),
    closeModal() {
      this.$emit('close');
    },
    receivedVaultName(data) {
      this.vaultName = data;
    },
    addVault() {
      // Integration with collections API.
      let { user, token } = this.getCurrentUser().profile;
      this.username = this.getCurrentUser().username;
      let collectionManager = new CollectionManager({ agent: this.getAgentInstance(), user });
      const name = this.vaultName;
      const id = collectionManager.create(token, { name });
      if (id) {
        this.closeModal();
      }
    },
  },
};
</script>
<style scoped>
.modal-width {
  width: 32rem;
}
</style>
