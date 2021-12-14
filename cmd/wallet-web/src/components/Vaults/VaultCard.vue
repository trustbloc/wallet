<template>
  <div
    v-if="type === 'regular'"
    class="
      relative
      md:flex-col
      px-6
      xl:w-64
      h-auto
      md:h-40
      bg-neutrals-white
      rounded-xl
      border border-neutrals-dark border-opacity-20
      vaultContainer
      flex flex-row
    "
  >
    <router-link
      :to="{ name: 'credentials', params: { vaultId } }"
      class="flex flex-row md:flex-col pt-5"
    >
      <div
        :class="[`flex justify-center items-center w-12 h-12 bg-gradient-${color} rounded-full`]"
      >
        <img class="w-6 h-5" src="@/assets/img/vaults.svg" alt="Vault Icon" />
      </div>
      <div class="px-3 md:px-0 pb-4">
        <span class="block pt-4 text-lg font-bold text-neutrals-dark"> {{ name }}</span>
        <span class="block text-sm font-bold text-neutrals-medium">
          {{ numOfCreds }}
        </span>
      </div>
    </router-link>
    <div class="absolute right-4 pt-4">
      <slot />
    </div>
  </div>
  <div v-else-if="type === 'addNew'">
    <button
      class="
        justify-center
        items-center
        px-6
        pt-4
        md:pt-5
        w-full
        xl:w-64
        h-24
        md:h-40
        bg-neutrals-moist
        rounded-xl
        border border-neutrals-dark border-opacity-10
        flex flex-col
      "
      @click="showAddVault = !showAddVault"
    >
      <div class="flex justify-center items-center w-8 h-8 bg-neutrals-white rounded-full">
        <img
          class="w-6 h-5 text-primary-purple"
          src="@/assets/img/icons-sm--plus-icon.svg"
          alt="Add Icon"
        />
      </div>
      <span class="block pt-2 pb-4 text-base font-bold text-neutrals-dark">
        {{ t('Vaults.addVault') }}</span
      >
    </button>
    <add-vault :show="showAddVault" />
  </div>
</template>

<script>
import { ref } from 'vue';
import AddVault from '@/components/Vaults/AddVaultModal.vue';
import { useI18n } from 'vue-i18n';

export default {
  name: 'VaultCard',
  components: {
    AddVault,
  },
  props: {
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
};
</script>
<style scoped>
.vaultContainer {
  box-shadow: 0px 4px 20px 0px rgba(25, 12, 33, 0.1);
}
</style>
