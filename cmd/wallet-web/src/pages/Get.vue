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
import DIDAuthForm from './DIDAuth.vue';
import DIDConnectForm from './DIDConnect.vue';
import PresentationDefQueryForm from './PresentationDefQuery.vue';
import MultipleQueryForm from './MultipleQuery.vue';
import WACIShareForm from './WACIShare.vue';
import { extractQueryTypes, WACIPolyfillHandler, CHAPIEventHandler } from './mixins';

const QUERY_FORMS = [
  {
    id: 'PresentationDefinitionQuery',
    component: PresentationDefQueryForm,
    match: (types) => ['PresentationExchange', 'DIDConnect'].every((elem) => types.includes(elem)),
  },
  {
    id: 'DIDAuth',
    component: DIDAuthForm,
    match: (types) => types.length == 1 && types.includes('DIDAuth'),
  },
  {
    id: 'DIDConnect',
    component: DIDConnectForm,
    match: (types) => ['DIDConnect'].every((elem) => types.includes(elem)),
  },
  {
    id: 'WACIShare',
    component: WACIShareForm,
    match: (types) => types.includes('WACIShare'),
    protocolHandler: WACIPolyfillHandler,
  },
  // default: MultipleQueryForm
];

function findForm(credEvent) {
  let { query } = credEvent.credentialRequestOptions.web.VerifiablePresentation;
  query = Array.isArray(query) ? query : [query];
  const types = extractQueryTypes(query);

  // TODO user filter instead of forloop
  for (let form of QUERY_FORMS) {
    if (form.match(types)) {
      return {
        component: form.component,
        protocolHandler: form.protocolHandler
          ? new form.protocolHandler(credEvent)
          : new CHAPIEventHandler(credEvent),
      };
    }
  }

  console.debug(
    'no matching query type handler found, switching to default MultipleQueryForm & CHAPI event handler'
  );

  // default form & protocol handler
  return {
    component: MultipleQueryForm,
    protocolHandler: new CHAPIEventHandler(credEvent),
  };
}

export default {
  data() {
    return {
      component: '',
      protocolHandler: null,
    };
  },
  beforeCreate: async function () {
    this.credentialEvent = await this.$webCredentialHandler.receiveCredentialEvent();
    if (!this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation) {
      console.log("invalid web credential type, expected 'VerifiablePresentation'");
      return;
    }

    const { component, protocolHandler } = findForm(this.credentialEvent);

    this.component = component;
    this.protocolHandler = protocolHandler;
  },
  computed: {
    dynamo() {
      return this.component;
    },
  },
};
</script>

<style lang="css">
.md-menu-content {
  width: auto;
  max-width: 100% !important;
}

.md-menu-content .md-ripple > span {
  position: relative !important;
}

.md-table-head {
  background-color: #00bcd4 !important;
  color: white !important;
  font-size: 10px !important;
  font-weight: 500 !important;
}

.md-table-head-label .md-sortable {
  padding-left: 10px !important;
}

.md-empty-state-icon {
  width: 60px !important;
  min-width: 60px !important;
  height: 50px !important;
  font-size: 50px !important;
  margin: 0;
}

.md-empty-state-label {
  font-size: 18px !important;
  font-weight: 500 !important;
  line-height: 20px !important;
}

.viewport {
  font-size: 16px !important;
}

.center-span {
  text-align: center !important;
}
</style>
