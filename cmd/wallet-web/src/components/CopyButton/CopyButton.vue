<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <md-button
    :disabled="!content"
    class="md-icon-button md-dense md-raised md-info"
    @click="copyToClipboard"
  >
    <md-icon>content_copy</md-icon>
    <textarea id="content" style="display: none" :value="content" />
  </md-button>
</template>

<script>
export default {
  name: 'CopyButton',
  props: {
    content: {
      type: String,
      default: '',
    },
  },
  methods: {
    copyToClipboard: function () {
      let payload = document.querySelector('#content');
      payload.setAttribute('style', 'display: unset');
      payload.select();

      document.execCommand('copy');
      payload.setAttribute('style', 'display: none');
      window.getSelection().removeAllRanges();
    },
  },
};
</script>
