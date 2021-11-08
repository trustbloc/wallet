<template>
  <div class="h-11">
    <!--Todo make flyout icon and text configurable -->
    <button
      v-if="type === 'default'"
      id="default"
      class="inline-flex items-center py-2 px-3 w-screen md:w-auto"
    >
      <div class="flex-none w-6 h-6">
        <img src="@/assets/img/icons-sm--vault-icon.svg" />
      </div>
      <div class="flex-grow pl-2 text-left">
        <span class="text-sm font-bold text-neutrals-dark">{{ i18n.allVaultLabel }}</span>
      </div>
      <div class="flex-none w-6 h-6">
        <img src="@/assets/img/icons-sm--chevron-down-icon.svg" />
      </div>
    </button>
    <button
      v-if="type === 'outline'"
      id="outline"
      class="
        w-11
        h-11
        bg-neutrals-white
        rounded-lg
        focus:border-neutrals-chatelle
        border-opacity-10
        focus:ring-primary-purple focus:ring-opacity-70 focus:ring-2
        focus-within:ring-offset-2
        border border-neutrals-chatelle
        hover:border-neutrals-mountainMist-light
      "
      @focus="showTooltip = !showTooltip"
    >
      <div class="flex-none p-2">
        <!-- TODO: Issue-816 Implement svg color change on hover -->
        <img
          id="flyoutMenuId"
          src="@/assets/img/more-icon.svg"
          @click="toggleFlyoutMenuList"
          @mouseover="showTooltip = !showTooltip"
        />
      </div>
      <div v-if="showTooltip" id="tooltip">
        <tool-tip :tool-tip-label="i18n.toolTipLabel"></tool-tip>
      </div>
    </button>
    <div v-if="showFlyoutMenuList" id="flyoutMenuList" class="relative">
      <flyout-menu-list :credential-id="credentialId" />
    </div>
  </div>
</template>

<script>
import ToolTip from '@/components/ToolTip/ToolTip.vue';
import FlyoutMenuList from '@/components/FlyoutMenu/FlyoutMenuList.vue';

export default {
  name: 'FlyoutMenu',
  components: {
    ToolTip,
    FlyoutMenuList,
  },
  props: {
    type: {
      type: String,
      default: 'default',
    },
    credentialId: {
      type: String,
    },
  },
  data() {
    return {
      toolTipLabel: {
        type: String,
        default: 'default',
      },
      showTooltip: false,
      showFlyoutMenuList: false,
    };
  },
  computed: {
    i18n() {
      return this.$t('CredentialDetails');
    },
  },
  mounted() {
    document.addEventListener('click', this.close);
  },
  beforeDestroy() {
    document.removeEventListener('click', this.close);
  },
  methods: {
    toggleFlyoutMenuList() {
      this.showFlyoutMenuList = !this.showFlyoutMenuList;
      this.showTooltip = false;
    },
    close(e) {
      if (!this.$el.contains(e.target)) {
        this.showFlyoutMenuList = false;
      }
    },
  },
};
</script>
