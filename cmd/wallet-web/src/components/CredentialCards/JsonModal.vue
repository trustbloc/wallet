/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <div>
    <md-dialog :md-active.sync="showDialog" style="max-width: 75% !important">
      <md-dialog-title style="background-color: #00bcd4">
        {{ credDisplayName(item) }}
      </md-dialog-title>
      <md-content class="md-dialog-container">
        <vue-json-pretty :data="item"> </vue-json-pretty>
      </md-content>

      <md-dialog-actions style="float: right; padding: 20px">
        <md-button class="md-primary" @click="showDialog = false"> Close </md-button>
      </md-dialog-actions>
    </md-dialog>
    <span class="infoIcon" @click.stop="showDialog = true"><md-icon>info</md-icon></span>
  </div>
</template>

<script>
import { getCredentialType } from '@/pages/mixins';
import VueJsonPretty from 'vue-json-pretty';
export default {
  name: 'DialogCustom',
  components: {
    VueJsonPretty,
  },
  props: {
    item: Object,
  },
  data: () => ({
    showDialog: false,
  }),
  methods: {
    credDisplayName: function (vc) {
      return vc.name ? vc.name : getCredentialType(vc.type);
    },
  },
};
</script>

<style scoped>
.infoIcon {
  position: absolute;
  right: 0;
  top: 0;
  padding: 10px 15px;
  opacity: 0.4;
  transition: all 0.5s ease;
}
.infoIcon:hover {
  opacity: 1;
}
.md-dialog-container {
  max-width: 1000px !important;
  max-height: 100% !important;
  overflow-y: scroll !important;
  overflow-x: scroll !important;
  word-wrap: break-word !important;
}
</style>
