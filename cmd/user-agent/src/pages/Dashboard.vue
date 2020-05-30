/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="content">
    <div class="md-layout">

      <div class="md-layout-item md-medium-size-50 md-xsmall-size-100 md-size-30">
        <stats-card data-background-color="green">
          <template slot="header">
            <md-icon>how_to_reg</md-icon>
          </template>

          <template slot="content">
            <div style="padding-top: 5px;float: left">
                <h3 class="title">Login to register wallet</h3>
            </div>

            <div style="padding-top: 104px">
              <component v-bind:is="component"></component>
            </div>
          </template>
        </stats-card>
      </div>
      <div class="md-layout-item md-medium-size-50 md-xsmall-size-100 md-size-30">
        <stats-card data-background-color="orange">
          <template slot="header">
            <md-icon>flip_to_back</md-icon>
          </template>

          <template slot="content">
            <sidebar-link to="/DIDManagement">
              <h3 class="title">DID Management</h3>
            </sidebar-link>

          </template>
        </stats-card>

          <stats-card data-background-color="purple">
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

    </div>
    <!-- Stored credentials-->
    <div class="content">
      <div class="md-layout">
        <div class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-50">
          <md-card>
            <md-card-header data-background-color="green">
              <h4 class="title"> <md-icon>content_paste</md-icon> My Stored Credentials</h4>
            </md-card-header>
            <md-card-content>
              <simple-table v-for="vc in verifiableCredential" v-bind:key=vc.name
                            v-bind:name="vc.name"
                            v-bind:data="vc.credential">
              </simple-table>
            </md-card-content>
          </md-card>
        </div>
        <!-- Stored presentation -->
        <div class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-50">
          <md-card>
            <md-card-header data-background-color="orange">
              <h4 class="title"> <md-icon>content_paste</md-icon> My Stored Presentations</h4>
            </md-card-header>
            <md-card-content>
              <simple-table v-for="vp in verifiablePresentation" v-bind:key=vp.name
                            v-bind:name="vp.name"
                            v-bind:data="vp.presentation">
              </simple-table>
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
  import LoginForm from "./chapi/Login.vue";

  let vcData = [];
  let vpData = [];

  async function fetchCredentials() {
    // Get the VC data
    for (let i = 0; i < vcData.length; i++) {
      try {
        let resp = await window.$aries.verifiable.getCredential({
          id: vcData[i].id
        })
        console.log('get vc ' + vcData[i].id)

        vcData[i].credential = JSON.parse(resp.verifiableCredential)
      } catch (e) {
        console.error('get vc failed : errMsg=' + e)
      }
    }
  }

  async function fetchPresentations() {
    // Get the VP data
    for (let i = 0; i < vpData.length; i++) {
      try {
        let resp = await window.$aries.verifiable.getPresentation({
          id: `${vpData[i].id}`
        })
        vpData[i].presentation = resp.verifiablePresentation
      } catch (e) {
        console.error('get vp failed : errMsg=' + e + " ID " + vpData[i].id)
      }
    }
  }

  export default {
    components: {
      StatsCard,
      SimpleTable,
    },
    beforeCreate: async function () {
      // Load the Credentials
      let aries = await this.$arieslib
      window.$webCredentialHandler = this.$webCredentialHandler
      window.$aries = aries
      await this.getCredentials(aries)
      await this.getPresentations(aries)
      this.component = LoginForm
    },
    methods: {
      getCredentials: async function (aries) {
        try {
          let resp = await aries.verifiable.getCredentials()
          if (!resp.result) {
            console.log('no credentials exists')
            return
          }

          vcData = resp.result
          if (resp.result.length === 0) {
            console.log('no credentials exists')
          }

        } catch (e) {
          console.error('get credentials failed : errMsg=' + e)
        }

        await fetchCredentials()
        this.verifiableCredential = vcData
      },
      getPresentations: async function (aries) {
        try {
          let resp = await aries.verifiable.getPresentations()
          if (!resp.result) {
            console.log('no presentation exists')
            return
          }

          vpData = resp.result
          if (resp.result.length === 0) {
            console.log('no presentation exists')
          }
        } catch (e) {
          console.error('get presentation failed : errMsg=' + e)
        }

        await fetchPresentations()
        this.verifiablePresentation = vpData
      }
    },
    data() {
      return {
        verifiableCredential: [],
        verifiablePresentation: [],
        component: null,
      }
    }
  }
</script>

<style>
  .title{
    text-transform: capitalize;
  }
  .md-content {
    overflow: auto;
    padding: 1px;
    font-size: 6px;
    line-height: 16px;
  }

</style>


