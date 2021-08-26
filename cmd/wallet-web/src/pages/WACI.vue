<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <component :is="component"></component>
  </div>
</template>

<script>
import WACIShareForm from './WACIShare.vue';
import { WACIRedirectHandler } from './mixins';

function findForm(credEvent) {
  // for now, has implemented only one form and handler
  return {
    component: WACIShareForm,
    protocolHandler: new WACIRedirectHandler(credEvent),
  };
}

export default {
  data() {
    return {
      component: null,
      protocolHandler: null,
    };
  },
  created: function () {
    // TODO read oob & redirect URL from query string
    this.credentialEvent = {};

    const { component, protocolHandler } = findForm(this.credentialEvent);

    this.component = component;
    this.protocolHandler = protocolHandler;
  },
};
</script>
