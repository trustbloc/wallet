<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <div class="h-11">
    <!--Todo make flyout icon and text configurable -->
    <button
      :id="id"
      v-if="type === 'default'"
      class="inline-flex items-center py-2 px-3 w-screen md:w-auto"
    >
      <div class="flex-none w-6 h-6">
        <img src="@/assets/img/icons-sm--vault-icon.svg" />
      </div>
      <div class="flex-grow pl-2 text-left">
        <span class="text-sm font-bold text-neutrals-dark truncate">{{
          t('CredentialDetails.allVaultLabel')
        }}</span>
      </div>
      <div class="flex flex-none justify-end w-6 h-6">
        <img src="@/assets/img/icons-sm--chevron-down-icon.svg" />
      </div>
    </button>
    <button
      :id="id"
      v-if="type === 'outline'"
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
      @click="toggleFlyoutMenuList"
      @focus="showTooltip = false"
    >
      <!-- TODO: Issue-816 Implement svg color change on hover -->
      <img
        alt="flyout menu icon"
        class="p-2"
        src="@/assets/img/more-icon.svg"
        @mouseover="showTooltip = true"
        @mouseout="showTooltip = false"
      />
      <tool-tip v-if="showTooltip" :tool-tip-label="t('CredentialDetails.toolTipLabel')" />
    </button>
    <button
      :id="id"
      v-if="type === 'vault'"
      class="
        w-8
        h-8
        bg-neutrals-white
        hover:bg-neutrals-softWhite
        focus:bg-neutrals-mischka
        rounded-full
      "
    >
      <!-- TODO: Issue-816 Implement svg color change on hover -->
      <img class="p-2" src="@/assets/img/more-icon.svg" @click="toggleFlyoutMenuList" />
    </button>
    <div v-if="showFlyoutMenuList" class="relative" @click.prevent="toggleFlyoutMenuList">
      <slot />
    </div>
  </div>
</template>

<script>
import ToolTip from '@/components/ToolTip/ToolTip.vue';
import { useI18n } from 'vue-i18n';

export default {
  name: 'FlyoutMenu',
  components: {
    ToolTip,
  },
  props: {
    type: {
      type: String,
      default: 'default',
    },
    id: {
      type: String,
      required: true,
    },
  },
  setup() {
    const { t, locale } = useI18n();
    return { t, locale };
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
  mounted() {
    document.addEventListener('click', this.close);
  },
  beforeUnmount() {
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
