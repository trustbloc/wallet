<template>
  <div
    v-show="isOpen"
    class="flex overflow-y-auto fixed inset-0 z-50 justify-center items-center bg-gradient-harold"
    :class="{ open: isOpen }"
    @click="$emit('onClose')"
  >
    <div
      class="relative mx-6 lg:mx-auto max-w-6xl bg-neutrals-white rounded-2xl modal-width"
      :class="{ open: isOpen }"
      @click.stop
    >
      <div v-if="title" class="pt-5 pb-3 text-lg font-bold text-neutrals-dark">{{ title }}</div>
      <div>
        <slot />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Modal',
  props: {
    isOpen: {
      type: Boolean,
      default: false,
    },
    title: {
      type: String,
      default: 'Credential',
    },
  },
};
</script>

<style scoped>
.modal-width {
  width: 28rem;
  box-shadow: 0px 12px 48px 0px rgba(25, 12, 33, 0.2);
}
</style>
