<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <md-table :table-header-color="tableHeaderColor">
      <md-table-row slot="md-table-row" style="text-align: left">
        <badger-accordion>
          <badger-accordion-item>
            <template slot="header">
              <span>
                <md-icon style="color: #0b97c4; float: left" class="md-size-2x">{{
                  headerIcon
                }}</md-icon>
                <div>
                  <md-header style="font-size: 16px; padding: 0 16px">{{ name }}</md-header>
                  <md-subheader>{{ description }}</md-subheader>
                </div>
              </span>
            </template>
            <md-divider></md-divider>
            <template slot="content">
              <!-- Display the json of stored credential -->
              <md-content>
                <vue-json-pretty :data="data"> </vue-json-pretty>
              </md-content>
            </template>
          </badger-accordion-item>
        </badger-accordion>
      </md-table-row>
    </md-table>
  </div>
</template>

<script>
import { BadgerAccordion, BadgerAccordionItem } from 'vue-badger-accordion';
import VueJsonPretty from 'vue-json-pretty';

export default {
  name: 'SimpleTable',
  components: {
    BadgerAccordion,
    BadgerAccordionItem,
    VueJsonPretty,
  },
  props: {
    name: {
      type: String,
    },
    description: {
      type: String,
    },
    data: {
      type: Object,
    },
    tableHeaderColor: {
      type: String,
      default: '',
    },
    headerIcon: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      selected: [],
      second: false,
    };
  },
};
</script>
