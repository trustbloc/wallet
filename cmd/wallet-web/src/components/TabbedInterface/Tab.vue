<template>
  <div v-show="isActive" class="px-4 md:px-0">
    <slot></slot>
  </div>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: 'Tab',
    },
  },
  data() {
    return {
      isActive: true,
    };
  },
};
</script>
