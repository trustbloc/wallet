<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="h-screen">
    <component :is="component"></component>
  </div>
</template>
<script>
import { inject, provide, computed, reactive } from 'vue';
import DIDAuthPage from '@/pages/DIDAuthPage.vue';
import DIDConnectPage from '@/pages/DIDConnectPage.vue';
import CHAPISharePage from '@/pages/CHAPISharePage.vue';
import MultipleQueryPage from '@/pages/MultipleQueryPage.vue';
import WACISharePage from '@/pages/WACISharePage.vue';
import { extractQueryTypes, WACIPolyfillHandler, CHAPIEventHandler } from '@/mixins';

const QUERY_FORMS = [
  {
    id: 'PresentationDefinitionQuery',
    component: CHAPISharePage,
    match: (types) => ['PresentationExchange', 'DIDConnect'].every((elem) => types.includes(elem)),
  },
  {
    id: 'DIDAuth',
    component: DIDAuthPage,
    match: (types) => types.length == 1 && types.includes('DIDAuth'),
  },
  {
    id: 'DIDConnect',
    component: DIDConnectPage,
    match: (types) => ['DIDConnect'].every((elem) => types.includes(elem)),
  },
  {
    id: 'WACIShare',
    component: WACISharePage,
    match: (types) => types.includes('WACIShare'),
    protocolHandler: WACIPolyfillHandler,
  },
  // default: MultipleQueryForm
];

function findForm(credEvent) {
  let { query } = credEvent.credentialRequestOptions.web.VerifiablePresentation;
  query = Array.isArray(query) ? query : [query];

  const types = extractQueryTypes(query);
  const found = QUERY_FORMS.filter((form) => form.match(types));

  if (found.length > 0) {
    console.log('yas chapi share in findform');
    phStore.protocolHandlerr = found[0].protocolHandler
      ? new found[0].protocolHandler(credEvent)
      : new CHAPIEventHandler(credEvent);
    return {
      component: found[0].component,
      protocolHandler: found[0].protocolHandler
        ? new found[0].protocolHandler(credEvent)
        : new CHAPIEventHandler(credEvent),
    };
  }

  console.debug(
    'no matching query type handler found, switching to default MultipleQueryForm & CHAPI event handler'
  );

  // default form & protocol handler
  return {
    component: MultipleQueryPage,
    protocolHandler: new CHAPIEventHandler(credEvent),
  };
}
export const phStore = reactive({
  protocolHandlerr: null,
});

export default {
  provide() {
    return {
      protocolHandler: this.protocolHandler,
    };
  },
  setup() {
    const webCredentialHandler = inject('webCredentialHandler');
    console.log('when is this called');
    return { webCredentialHandler };
  },
  data() {
    return {
      component: '',
      protocolHandler: null,
    };
  },
  computed: {
    dynamo() {
      return this.component;
    },
  },
  provide: {
    protocolHandler: computed(() => this.protocolHandler),
  },
  // provide() {
  //   return {
  //     protocolHandler: this.protocolHandler,
  //   };
  // },
  beforeCreate: async function () {
    this.credentialEvent = await this.webCredentialHandler.receiveCredentialEvent();
    if (!this.credentialEvent.credentialRequestOptions.web.VerifiablePresentation) {
      console.log("invalid web credential type, expected 'VerifiablePresentation'");
      return;
    }

    const { component, protocolHandler } = findForm(this.credentialEvent);
    console.log('before create in get layout');
    this.component = component;
    this.protocolHandler = protocolHandler;
    console.log('protocolhandler set here');
    //this.provide('protocolHandler', this.protocolHandler);
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
