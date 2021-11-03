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
    <ul class="p-4 text-base text-start">
      <slot />
    </ul>
  </div>
</template>

<script>
export default {
  name: 'FlyoutMenuList',
};
</script>
