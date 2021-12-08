<template>
  <div
    class="
      inline-flex
      absolute
      top-14
      right-0
      z-20
      justify-center
      py-1
      px-5
      bg-neutrals-dark
      rounded-lg
    "
  >
    <div
      class="
        absolute
        z-10
        bg-neutrals-dark
        rounded-sm
        transform
        rotate-45
        border-8 border-neutrals-dark
        triangle
      "
    />
    <span class="z-20 text-sm text-center text-neutrals-white whitespace-nowrap">{{
      toolTipLabel
    }}</span>
  </div>
</template>

<script>
export default {
  name: 'ToolTip',
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
  top: -4px;
  right: 20px;
}
</style>
