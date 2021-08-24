<template>
  <div>
    <!--Todo make flyout icon and text configurable -->
    <button class="inline-flex py-2 px-3 w-screen md:w-auto">
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
  </div>
</template>

<script>
export default {
  name: 'FlyoutMenu',
};
</script>
