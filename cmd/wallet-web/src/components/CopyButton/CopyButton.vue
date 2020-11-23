/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <md-button :disabled="!content" v-on:click="copyToClipboard" class="md-icon-button md-dense md-raised md-info">
    <md-icon>content_copy</md-icon>
    <textarea style="display: none" id="content" :value="content"/>
  </md-button>
</template>

<script>
export default {
  name: "copy-button",
  props: {
    content: {
      type: String,
      default: ''
    },
  },
  methods: {
    copyToClipboard: function () {
      let payload = document.querySelector('#content')
      payload.setAttribute('style', 'display: unset')
      payload.select()

      document.execCommand('copy');
      payload.setAttribute('style', 'display: none')
      window.getSelection().removeAllRanges()
    },
  }
}
</script>