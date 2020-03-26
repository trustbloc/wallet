/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="content">
    <div class="md-layout">

      <div
        class="md-layout-item md-medium-size-50 md-xsmall-size-100 md-size-25">
        <stats-card data-background-color="green">
          <template slot="header">
            <md-icon>store</md-icon>
          </template>

          <template slot="content">
            <sidebar-link to="/RegisterWallet">
              <h3 class="title">Register Wallet</h3>
            </sidebar-link>
          </template>
        </stats-card>
      </div>
      <div
        class="md-layout-item md-medium-size-50 md-xsmall-size-100 md-size-25">
        <stats-card data-background-color="orange">
          <template slot="header">
            <md-icon>border_outer</md-icon>
          </template>

          <template slot="content">
            <sidebar-link to="/MyVC">
              <h3 class="title">Generate Presentation</h3>
            </sidebar-link>

          </template>
        </stats-card>
      </div>
      <div
              class="md-layout-item md-medium-size-50 md-xsmall-size-100 md-size-25">
        <stats-card data-background-color="blue">
          <template slot="header">
            <md-icon>account_balance</md-icon>
          </template>

          <template slot="content">
            <label>For testing</label>
            <sidebar-link to="/StoreVC">
              <h3 class="title">Store Credentials</h3>
            </sidebar-link>
          </template>
        </stats-card>
      </div>
      <div
              class="md-layout-item md-medium-size-50 md-xsmall-size-100 md-size-25">
        <stats-card data-background-color="purple">
          <template slot="header">
            <md-icon>description</md-icon>
          </template>

          <template slot="content">
            <label>For testing</label>
            <sidebar-link to="/GetVC">
              <h3 class="title">Get Credentials</h3>
            </sidebar-link>
          </template>
        </stats-card>
      </div>
    </div>
    <div class="content">
      <div class="md-layout">
        <div
                class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100"
        >
          <md-card>
            <md-card-header data-background-color="green">
              <h4 class="title"> <md-icon>content_paste</md-icon> My Stored Credentials</h4>
            </md-card-header>
            <md-card-content>
              <simple-table table-header-color="green"></simple-table>
            </md-card-content>
          </md-card>
        </div>
      </div>
    </div>
  </div>

</template>

<script>
import {StatsCard} from "@/components";
import {SimpleTable} from "@/components";

export default {
  components: {
    StatsCard,
    SimpleTable,
  },
  data() {
    return {
      }
    },
  };
</script>

<style>
  .title{
    text-transform: capitalize;
  }
</style>
