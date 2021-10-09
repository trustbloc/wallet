<template>
  <div
    class="
      overflow-hidden
      absolute
      right-0
      z-20
      mt-2
      w-56
      max-w-max
      text-neutrals-medium
      bg-neutrals-white
      rounded-lg
      shadow-xl
      bg-white
    "
  >
    <div class="p-4 text-base text-start">
      <button class="block pb-2 font-bold">{{ i18n.renameCredential }}</button>
      <button class="block pb-2 font-bold">{{ i18n.moveCredential }}</button>
      <button class="block font-bold text-primary-vampire">{{ i18n.deleteCredential }}</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FlyoutMenuList',
  computed: {
    i18n() {
      return this.$t('CredentialDetails');
    },
  },
};
</script>

<style scoped></style>
