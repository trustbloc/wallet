/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>

  <div>
    <md-table v-model="users" :table-header-color="tableHeaderColor">
      <md-table-row slot="md-table-row"  slot-scope="{ item }"  style="text-align:left">
        <badger-accordion>
          <badger-accordion-item>
            <template slot="header"><h5><md-icon>keyboard_arrow_down</md-icon>{{item.name}}</h5></template>
            <template slot="content">{{item.credential}}</template>
          </badger-accordion-item>
        </badger-accordion>
      </md-table-row>

    </md-table>
  </div>

</template>
<!--TODO : Fix the hard coded data with the actual data issue-89-->
<script>
  import {BadgerAccordion, BadgerAccordionItem} from 'vue-badger-accordion'

  export default {
    name: "simple-table",
    components: {
      BadgerAccordion,
      BadgerAccordionItem
    },
    props: {
      tableHeaderColor: {
        type: String,
        default: ""
      }
    },
    data() {
      return {
        selected: [],
        users: [
          {
            name: "Foo Student Card",
            id: "1",
            credential: "{\n" +
                    "    \"@context\": [\n" +
                    "        \"https://www.w3.org/2018/credentials/v1\"\n" +
                    "    ],\n" +
                    "    \"credentialSchema\": [],\n" +
                    "    \"credentialStatus\": {\n" +
                    "        \"id\": \"http://issuer.vcs.example.com:8070/status/1\",\n" +
                    "        \"type\": \"CredentialStatusList2017\"\n" +
                    "    },\n" +
                    "    \"credentialSubject\": {\n" +
                    "        \"email\": \"foo@bar.com\",\n" +
                    "        \"id\": \"Foo\",\n" +
                    "        \"name\": \"Foo\",\n" +
                    "        \"semester\": \"3\",\n" +
                    "        \"studentid\": \"1234568\",\n" +
                    "        \"university\": \"Faber College\"\n" +
                    "    },\n" +
                    "    \"id\": \"http://vc-issuer-1.com/c3172aff-4744-4e6d-a7d2-41e9055738ba\",\n" +
                    "    \"issuanceDate\": \"2020-03-24T22:13:39.009553847Z\",\n" +
                    "    \"issuer\": {\n" +
                    "        \"id\": \"did:trustbloc:testnet.trustbloc.local:EiBqgRZF5uDni1hIOfl-HqEX5GH4tNtLAEcpH1VlSHg_Qg==\",\n" +
                    "        \"name\": \"vc-issuer-1\"\n" +
                    "    },\n" +
                    "    \"proof\": {\n" +
                    "        \"created\": \"2020-03-24T22:13:39Z\",\n" +
                    "        \"jws\": \"eyJhbGciOiJFZDI1NTE5U2lnbmF0dXJlMjAxOCIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..5a0VBXSE9wYqsjxL9XRoS-8afFM150uTUr-SZNFuERV48FoxDbSzzPGnFLkItJPqC5TtSXZjS2z_FTMc2RfjAg\",\n" +
                    "        \"type\": \"Ed25519Signature2018\",\n" +
                    "        \"verificationMethod\": \"did:trustbloc:testnet.trustbloc.local:EiBqgRZF5uDni1hIOfl-HqEX5GH4tNtLAEcpH1VlSHg_Qg==#key-1\"\n" +
                    "    },\n" +
                    "    \"type\": [\n" +
                    "        \"VerifiableCredential\",\n" +
                    "        \"StudentCard\"\n" +
                    "    ]\n" +
                    "}",
          },
          {
            name: "Foo PR Card",
           id: "2",
            credential: "{\n" +
                    "    \"@context\": [\n" +
                    "        \"https://www.w3.org/2018/credentials/v1\"\n" +
                    "    ],\n" +
                    "    \"credentialSchema\": [],\n" +
                    "    \"credentialStatus\": {\n" +
                    "        \"id\": \"http://issuer.vcs.example.com:8070/status/1\",\n" +
                    "        \"type\": \"CredentialStatusList2017\"\n" +
                    "    },\n" +
                    "    \"credentialSubject\": {\n" +
                    "        \"email\": \"foo@bar.com\",\n" +
                    "        \"id\": \"Foo\",\n" +
                    "        \"name\": \"Foo\",\n" +
                    "        \"semester\": \"3\",\n" +
                    "        \"studentid\": \"1234568\",\n" +
                    "        \"university\": \"Faber College\"\n" +
                    "    },\n" +
                    "    \"id\": \"http://vc-issuer-1.com/c3172aff-4744-4e6d-a7d2-41e9055738ba\",\n" +
                    "    \"issuanceDate\": \"2020-03-24T22:13:39.009553847Z\",\n" +
                    "    \"issuer\": {\n" +
                    "        \"id\": \"did:trustbloc:testnet.trustbloc.local:EiBqgRZF5uDni1hIOfl-HqEX5GH4tNtLAEcpH1VlSHg_Qg==\",\n" +
                    "        \"name\": \"vc-issuer-1\"\n" +
                    "    },\n" +
                    "    \"proof\": {\n" +
                    "        \"created\": \"2020-03-24T22:13:39Z\",\n" +
                    "        \"jws\": \"eyJhbGciOiJFZDI1NTE5U2lnbmF0dXJlMjAxOCIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..5a0VBXSE9wYqsjxL9XRoS-8afFM150uTUr-SZNFuERV48FoxDbSzzPGnFLkItJPqC5TtSXZjS2z_FTMc2RfjAg\",\n" +
                    "        \"type\": \"Ed25519Signature2018\",\n" +
                    "        \"verificationMethod\": \"did:trustbloc:testnet.trustbloc.local:EiBqgRZF5uDni1hIOfl-HqEX5GH4tNtLAEcpH1VlSHg_Qg==#key-1\"\n" +
                    "    },\n" +
                    "    \"type\": [\n" +
                    "        \"VerifiableCredential\",\n" +
                    "        \"StudentCard\"\n" +
                    "    ]\n" +
                    "}",
          }

        ]
      };
    }
  };
</script>
