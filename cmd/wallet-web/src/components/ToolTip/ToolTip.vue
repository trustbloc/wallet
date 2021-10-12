<template>
  <div class="relative top-6 right-32 w-40">
    <div class="py-1 px-2 text-neutrals-white bg-neutrals-dark rounded-lg align-center">
      <span class="absolute triangle"></span>
      {{ toolTipLabel }}
    </div>
  </div>
</template>

<script>
export default {
  name: 'Tooltip.vue',
  props: {
    toolTipLabel: {
      type: String,
      default: 'default',
    },
  },
};
</script>

<style scoped>
.triangle {
  border-width: 0 6px 9px;
  border-bottom-color: #190c21;
  top: -9px;
  left: 132px;
  transform: translate-x(-50%);
}
</style>
