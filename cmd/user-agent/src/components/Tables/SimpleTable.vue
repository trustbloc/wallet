/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>

  <div>
    <md-table :table-header-color="tableHeaderColor">
      <md-table-row slot="md-table-row"   style="text-align:left">
        <badger-accordion>
          <badger-accordion-item>
            <template slot="header"><p><md-icon>keyboard_arrow_down</md-icon>{{name}}</p></template>
            <template slot="content">
              <!-- Display the json of stored credential -->
              <md-content>
                <vue-json-pretty
                        :data="this.data"
                >
                </vue-json-pretty>
              </md-content>
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
      data: {
        type: Object,
      },
      tableHeaderColor: {
        type: String,
        default: ""
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

