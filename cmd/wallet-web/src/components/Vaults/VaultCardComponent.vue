<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div
    v-if="type === 'regular'"
    class="vaultContainer relative flex h-auto flex-row rounded-xl border border-neutrals-dark border-opacity-20 bg-neutrals-white px-6 md:h-40 md:flex-col xl:w-64"
  >
    <AppLinkComponent
      :id="id"
      :to="{ name: 'credentials', params: { vaultId } }"
      @click="handleClick"
    >
      <div class="flex flex-row pt-5 md:flex-col">
        <div
          :class="[`flex justify-center items-center w-12 h-12 bg-gradient-${color} rounded-full`]"
        >
          <img class="h-5 w-6" src="@/assets/img/vaults.svg" alt="Vault Icon" />
        </div>
        <div class="px-3 pb-4 md:px-0">
          <span class="block pt-4 text-lg font-bold text-neutrals-dark"> {{ name }}</span>
          <span class="block text-sm font-bold text-neutrals-medium">
            {{ numOfCreds }}
          </span>
        </div>
      </div>
    </AppLinkComponent>
    <div class="absolute right-4 pt-4">
      <slot />
    </div>
  </div>
  <div v-else-if="type === 'addNew'" class="relative">
    <div
      class="flex h-24 w-full flex-col items-center justify-center rounded-xl border border-neutrals-dark border-opacity-10 bg-neutrals-moist px-6 pt-4 md:h-40 md:pt-5 xl:w-64"
    >
      <div class="flex h-8 w-8 items-center justify-center rounded-full bg-neutrals-white">
        <img
          class="h-5 w-6 text-primary-purple"
          src="@/assets/img/icons-sm--plus-icon.svg"
          alt="Add Icon"
        />
      </div>
      <span class="block pt-2 pb-4 text-base font-bold text-neutrals-dark">
        {{ t('Vaults.AddModal.addVault') }}</span
      >
    </div>
    <button
      :id="id"
      class="absolute top-0 left-0 z-0 h-full w-full cursor-pointer"
      @click="showAddVault = !showAddVault"
    />
    <AddVaultComponent :show="showAddVault" :existing-names="existingNames" @close="handleClose" />
  </div>
</template>

<script>
import { ref } from 'vue';
import AddVaultComponent from '@/components/Vaults/AddVaultModalComponent.vue';
import AppLinkComponent from '@/components/AppLink/AppLinkComponent.vue';
import { useI18n } from 'vue-i18n';
import { mapActions } from 'vuex';

export default {
  name: 'VaultCardComponent',
  components: {
    AddVaultComponent,
    AppLinkComponent,
  },
  props: {
    id: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'purple',
    },
    name: {
      type: String,
      required: true,
    },
    numOfCreds: {
      type: String,
      default: null,
    },
    type: {
      type: String,
      default: 'regular',
    },
    vaultId: {
      type: String,
      default: null,
    },
    existingNames: {
      type: Array,
      default: null,
    },
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  data() {
    return {
      showAddVault: ref(false),
    };
  },
  methods: {
    ...mapActions(['updateSelectedVaultId']),
    handleClick: function () {
      this.updateSelectedVaultId(this.vaultId);
      this.$router.push({ name: 'credentials', params: { vaultId: this.vaultId } });
    },
    handleClose: function () {
      this.showAddVault = false;
    },
  },
};
</script>
<style scoped>
.vaultContainer {
  box-shadow: 0px 4px 20px 0px rgba(25, 12, 33, 0.1);
}
</style>
