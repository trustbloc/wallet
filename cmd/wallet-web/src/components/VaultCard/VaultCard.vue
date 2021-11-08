<template>
  <div>
    <div
      v-if="type === 'regular'"
      class="
        md:flex-col
        px-6
        w-full
        md:w-64
        h-20
        md:h-auto
        bg-neutrals-white
        rounded-xl
        flex flex-row
      "
    >
      <div class="pt-5">
        <div
          :class="[`flex justify-center items-center w-12 h-12 bg-gradient-${color} rounded-full`]"
        >
          <img class="w-6 h-5" src="@/assets/img/vaults.svg" alt="Vault Icon" />
        </div>
      </div>
      <div>
        <div class="px-3 md:px-0 pb-4">
          <span class="block pt-4 text-lg font-bold text-neutrals-dark"> {{ name }}</span>
          <span class="block text-sm text-neutrals-medium">
            {{ numOfCreds }}
          </span>
        </div>
      </div>
    </div>
    <div
      v-if="type === 'addNew'"
      class="
        px-6
        w-full
        md:w-64
        h-24
        md:h-full
        bg-neutrals-moist
        rounded-xl
        border border-neutrals-dark border-opacity-10
        flex flex-col
      "
    >
      <div class="flex justify-center pt-5 md:pt-10">
        <div class="flex justify-center items-center w-8 h-8 bg-neutrals-white rounded-full">
          <img
            class="w-6 h-5 text-primary-purple"
            src="@/assets/img/icons-sm--plus-icon.svg"
            alt="Add Icon"
          />
        </div>
      </div>
      <div class="flex justify-center">
        <span class="block pt-2 pb-4 text-base font-bold text-neutrals-dark">
          {{ $t('Vaults.addVault') }}</span
        >
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'VaultCard',
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
      required: true,
    },
    type: {
      type: String,
      default: 'regular',
    },
  },
};
</script>
