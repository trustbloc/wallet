<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div
    class="
      flex-grow
      justify-start
      items-center
      mx-auto
      max-w-7xl
      h-screen
      max-h-screen
      flex flex-col
      shadow-main-container
    "
  >
    <Header :has-custom-gradient="true">
      <template #gradientContainer>
        <div class="absolute h-15 bg-gradient-full oval" />
      </template>
    </Header>
    <component :is="component" class="z-10 flex-grow h-full bg-neutrals-softWhite"></component>
    <Footer class="sticky bottom-0 z-20 border-t border-neutrals-thistle bg-neutrals-magnolia" />
  </div>
</template>

<script>
import { extractOOBGoalCode } from '@trustbloc/wallet-sdk';
import { WACIRedirectHandler } from '@/utils/mixins';
import { decode } from 'js-base64';
import WACIShareForm from '@/pages/WACIShare.vue';
import WACIIssuanceForm from '@/pages/WACIIssue.vue';
import Header from '@/components/Header/Header.vue';
import Footer from '@/components/Footer/Footer.vue';

function findForm(query) {
  if (!query.oob) {
    // TODO [Issue#1325] should be redirected to standard error screen.
    throw 'access denied, oob invitation missing';
  }

  const invitation = JSON.parse(decode(query.oob));

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
  components: {
    Header,
    Footer,
  },
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
<style>
.oval {
  left: 50%;
  transform: translateX(-50%);
  border-radius: 50%;
  filter: blur(50px);
  width: 15.625rem; /* 250px */
  top: 2.0625rem; /* 33px */
}
</style>
