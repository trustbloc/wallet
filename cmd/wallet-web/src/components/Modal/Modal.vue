<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <teleport v-if="showModal" :to="target">
    <!--Modal Content-->
    <div
      class="
        flex
        overflow-y-auto
        fixed
        inset-0
        z-50
        justify-center
        items-center
        bg-neutrals-dark bg-opacity-50
      "
      @click.self="close"
    >
      <div class="relative flex-grow mx-6 lg:mx-auto max-w-lg bg-neutrals-white rounded-2xl">
        <button v-if="showCloseButton" class="absolute right-0 pt-3 pr-3 w-10 h-10">
          <!-- TODO: use inline svg instead once https://github.com/trustbloc/edge-agent/issues/816 is fixed -->
          <img
            class="w-6 h-6 cursor-pointer"
            src="@/assets/img/Icons-sm--close-icon.svg"
            alt="Close Icon"
            @click="showModal = false"
          />
        </button>
        <slot name="content" />
        <!-- Buttons Container -->
        <div
          class="
            md:flex-row
            gap-4
            justify-start
            md:justify-between
            items-center
            px-5
            md:px-8
            pt-4
            pb-5
            text-center
            bg-neutrals-magnolia
            rounded-b-2xl
            flex flex-col
            modal-footer
            border-t border-0 border-neutrals-lilacSoft
          "
        >
          <styled-button type="outline" class="w-full md:w-auto" @click="showModal = false">
            {{ t('Modal.cancel') }}
          </styled-button>
          <slot name="actionButton" />
        </div>
      </div>
    </div>
  </teleport>
</template>

<script>
import { useI18n } from 'vue-i18n';
import StyledButton from '@/components/StyledButton/StyledButton';
import { ref, watch } from 'vue';

export default {
  name: 'Modal',
  components: {
    StyledButton,
  },
  props: {
    target: {
      type: String,
      default: 'body',
    },
    show: {
      type: Boolean,
      default: false,
    },
    showCloseButton: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const showModal = ref(false);
    watch(
      () => props.show,
      (show) => {
        showModal.value = show;
      }
    );

    return {
      t,
      showModal,
    };
  },
  methods: {
    close() {
      if (this.showCloseButton) this.showModal = false;
    },
  },
};
</script>
<style scoped>
.modal-footer {
  box-shadow: inset 0px 1px 0px 0px rgb(219, 215, 220);
}
</style>
