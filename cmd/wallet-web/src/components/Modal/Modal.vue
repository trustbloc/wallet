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
    >
      <div class="flex-grow mx-6 lg:mx-auto max-w-lg bg-neutrals-white rounded-2xl">
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
          <styled-button
            type="outline"
            class="w-full md:w-auto lg:w-auto"
            @click="showModal = false"
          >
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
};
</script>
<style scoped>
.modal-footer {
  box-shadow: inset 0px 1px 0px 0px rgb(219, 215, 220);
}
</style>
