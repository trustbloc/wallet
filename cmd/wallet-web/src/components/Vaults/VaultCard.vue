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
    <app-link :to="{ name: 'credentials', params: { vaultId } }" @click="handleClick">
      <div class="flex flex-row md:flex-col pt-5">
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
      </div>
    </app-link>
    <div class="absolute right-4 pt-4">
      <slot />
    </div>
  </div>
  <div v-else-if="type === 'addNew'" class="relative">
    <div
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
    >
      <div class="flex justify-center items-center w-8 h-8 bg-neutrals-white rounded-full">
        <img
          class="w-6 h-5 text-primary-purple"
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
      @click="showAddVault = !showAddVault"
      class="absolute top-0 left-0 z-0 w-full h-full cursor-pointer"
    />
    <add-vault :show="showAddVault" :existing-names="existingNames" />
  </div>
</template>

<script>
import { ref } from 'vue';
import AddVault from '@/components/Vaults/AddVaultModal.vue';
import AppLink from '@/components/AppLink/AppLink.vue';
import { useI18n } from 'vue-i18n';
import { mapActions } from 'vuex';

export default {
  name: 'VaultCard',
  components: {
    AddVault,
    AppLink,
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
  },
};
</script>
<style scoped>
.vaultContainer {
  box-shadow: 0px 4px 20px 0px rgba(25, 12, 33, 0.1);
}
</style>
