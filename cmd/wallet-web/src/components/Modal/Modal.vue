<template>
  <Modal :is-open="!!component" :title="title" @onClose="handleModalClose">
    <component :is="component" v-bind="props" @onClose="handleClose" />
  </Modal>
</template>

<script>
import { ModalBus } from '@/EventBus';

import Modal from './ModalLayout';

export default {
  components: { Modal },
  data() {
    return {
      component: null,
      title: '',
      props: null,
      closeOnClick: true,
    };
  },
  created() {
    ModalBus.$on('open', ({ component, title = '', props = null, closeOnClick = true }) => {
      this.component = component;
      this.title = title;
      this.props = props;
      this.closeOnClick = closeOnClick;
    });
    document.addEventListener('keyup', this.handleKeyup);
  },
  beforeDestroy() {
    document.removeEventListener('keyup', this.handleKeyup);
  },
  methods: {
    handleModalClose() {
      if (!this.closeOnClick) return;
      this.handleClose();
    },
    handleClose() {
      this.component = null;
    },
    handleKeyup(e) {
      // esc code
      if (e.keyCode === 27) this.handleClose();
    },
  },
};
</script>
