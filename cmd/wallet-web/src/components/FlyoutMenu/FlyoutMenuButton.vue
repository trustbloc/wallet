<template>
  <div>
    <button
      :id="id"
      :class="[`block font-bold block pb-2 text-${color}`]"
      @click="$emit('click', $event)"
    >
      {{ text }}
      <slot />
    </button>
  </div>
</template>

<script>
export default {
  name: 'FlyoutMenuButton',
  props: {
    text: {
      type: String,
      required: true,
    },
    color: {
      type: String,
      default: 'neutrals-medium',
    },
    id: {
      type: String,
      required: true,
    },
  },
  computed: {
    i18n() {
      return this.$t('CredentialDetails');
    },
  },
};
</script>
