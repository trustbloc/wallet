<template>
  <div>
    <!--Todo make flyout icon and text configurable -->
    <button v-if="type === 'default'" class="inline-flex items-center py-2 px-3 w-screen md:w-auto">
      <div class="flex-none w-6 h-6">
        <img src="@/assets/img/icons-sm--vault-icon.svg" />
      </div>
      <div class="flex-grow pl-2 text-left">
        <span class="text-sm font-bold text-neutrals-dark">All Vaults</span>
      </div>
      <div class="flex-none w-6 h-6">
        <img src="@/assets/img/icons-sm--chevron-down-icon.svg" />
      </div>
    </button>
    <button v-if="type === 'outline'" class="py-2 px-3 w-11 h-11">
      <div class="flex-none">
        <img src="@/assets/img/more-icon.svg" />
      </div>
    </button>
  </div>
</template>

<script>
export default {
  name: 'FlyoutMenu',
  props: {
    type: {
      type: String,
      default: 'default',
    },
  },
};
</script>
