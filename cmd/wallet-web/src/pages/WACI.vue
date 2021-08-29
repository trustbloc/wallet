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

function findForm(query) {
  // for now, implemented only one form and handler
  let protocolHandler;

  if (query.redirect) {
    protocolHandler = new WACIRedirectHandler(query.oob, query.redirect);
  } else {
    throw 'unable to find protocol handler.';
  }

  return {
    component: WACIShareForm,
    protocolHandler,
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
    const { component, protocolHandler } = findForm(this.$route.query);

    this.component = component;
    this.protocolHandler = protocolHandler;
  },
};
</script>
