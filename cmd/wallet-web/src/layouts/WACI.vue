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
import WACIShareForm from '@/pages/WACIShare.vue';
import WACIIssuanceForm from '@/pages/WACIIssue.vue';
import { WACIRedirectHandler } from '@/utils/mixins';
import { extractOOBGoalCode } from '@trustbloc/wallet-sdk';
const { Base64 } = require('js-base64');

function findForm(query) {
  if (!query.oob) {
    // TODO [Issue#1325] should be redirected to standard error screen.
    throw 'access denied, oob invitation missing';
  }

  const invitation = JSON.parse(Base64.decode(query.oob));

  switch (extractOOBGoalCode(invitation)) {
    case 'streamlined-vc':
      return {
        component: WACIIssuanceForm,
        protocolHandler: new WACIRedirectHandler(invitation),
      };
    case 'streamlined-vp':
      return {
        component: WACIShareForm,
        protocolHandler: new WACIRedirectHandler(invitation),
      };
    default:
      // TODO [Issue#1326] should throw error if goal-code is missing once other trustbloc components starts sending right goal codes.
      return {
        component: WACIShareForm,
        protocolHandler: new WACIRedirectHandler(invitation),
      };
  }
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
