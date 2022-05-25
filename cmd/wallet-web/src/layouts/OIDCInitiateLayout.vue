<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <keep-alive> </keep-alive>
</template>

<script>
import uuid from 'uuid/v4';
import { readOpenIDConfiguration, sendCredentialAuthorizeRequest } from '@/mixins';
import Cookies from 'js-cookie';

export default {
  name: 'OIDCInitiateLayout',
  created: async function () {
    {
      const opState = this.$route.query.op_state || uuid();
      const { issuer, credential_type, manifest_id } = this.$route.query;
      console.log('what is issuer', issuer, this.$route.query);
      const configuration = await readOpenIDConfiguration(issuer);
      Cookies.set(
        opState,
        JSON.stringify({
          issuer,
          credentialTypes: Array.isArray(credential_type) ? credential_type : [credential_type],
          manifestID: manifest_id,
        })
      );
      sendCredentialAuthorizeRequest(
        configuration,
        this.$route.query,
        `${location.protocol}//${location.host}/oidc/save`,
        opState
      );
    }
  },
};
</script>
