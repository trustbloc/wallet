/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>

  <div>
    <md-table v-model="users" :table-header-color="tableHeaderColor">
      <md-table-row slot="md-table-row"   style="text-align:left">
        <badger-accordion>
          <badger-accordion-item>
            <template slot="header"><h5><md-icon>keyboard_arrow_down</md-icon>{{name}}</h5></template>
            <template slot="content">
              <md-tabs class="md-info md-ripple" md-alignment="left" >
                <!-- Display the json of stored credential -->
                <md-tab id="tab-home" md-label="Source">
                  <md-content>
                  <vue-json-pretty
                          :data="this.credential"
                         >
                  </vue-json-pretty>
                  </md-content>
                </md-tab>
                <!--Todo add preview of credential tab - document view of credentials -->
              </md-tabs>
            </template>
          </badger-accordion-item>
        </badger-accordion>
      </md-table-row>

    </md-table>
  </div>

</template>

<script>
  import {BadgerAccordion, BadgerAccordionItem} from 'vue-badger-accordion';
  import VueJsonPretty from 'vue-json-pretty';

  export default {
    name: "simple-table",
    components: {
      BadgerAccordion,
      BadgerAccordionItem,
        VueJsonPretty,
    },
    props: {
      name: {
        type: String,
      },
      credential: {
        type: Object,
      }
    },
    data() {
      return {
        selected: [],
        second: false
      }
    }
  }
</script>

<style>
  .md-ripple {
    margin-top: 10px;
  }

</style>
