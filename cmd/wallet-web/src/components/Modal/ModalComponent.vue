<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->
<template>
  <teleport v-if="showModal" :to="target">
    <!--Modal Content-->
    <div
      class="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-neutrals-dark bg-opacity-50"
      @click.self="close"
    >
      <div class="fixed inset-0 flex items-center justify-center overflow-y-auto">
        <slot name="errorToast" />
      </div>
      <div class="relative grow mx-6 max-w-lg bg-neutrals-white rounded-2xl lg:mx-auto">
        <button v-if="showCloseButton" class="absolute right-0 pt-3 pr-3 w-10 h-10" @click="close">
          <!-- TODO: use inline svg instead once https://github.com/trustbloc/wallet/issues/816 is fixed -->
          <img
            class="h-6 w-6 cursor-pointer"
            src="@/assets/img/Icons-sm--close-icon.svg"
            alt="Close Icon"
          />
        </button>
        <slot name="content" />
        <!-- Buttons Container -->
        <div
          class="flex flex-col gap-4 justify-start items-center px-5 pt-4 pb-5 text-center bg-neutrals-magnolia rounded-b-2xl border-0 border-t border-neutrals-lilacSoft md:flex-row md:justify-between md:px-8 modal-footer"
        >
          <StyledButtonComponent type="btn-outline" class="w-full md:w-auto" @click="cancel">
            {{ t('Modal.cancel') }}
          </StyledButtonComponent>
          <slot name="actionButton" />
        </div>
      </div>
    </div>
  </teleport>
</template>

<script>
import { useI18n } from 'vue-i18n';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent';
import { ref, watch } from 'vue';

export default {
  name: 'ModalComponent',
  components: {
    StyledButtonComponent,
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
  emits: ['close'],
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
      if (this.showCloseButton) {
        this.showModal = false;
        this.$emit('close');
      }
    },
    cancel() {
      this.showModal = false;
      this.$emit('close');
    },
  },
};
</script>
<style scoped>
.modal-footer {
  box-shadow: inset 0px 1px 0px 0px rgb(219, 215, 220);
}
</style>
