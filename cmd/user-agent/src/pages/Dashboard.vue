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
            <md-icon>store</md-icon>
          </template>

          <template slot="content">
            <sidebar-link to="/RegisterWallet">
              <h3 class="title">Register Wallet</h3>
            </sidebar-link>
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
      </div>
      <div class="md-layout-item md-medium-size-50 md-xsmall-size-100 md-size-30">
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

  let vcData = [];
  let vpData = [];

  async function fetchCredentials() {
    // Get the VC data
    for (let i = 0; i < vcData.length; i++) {
      await window.$aries.verifiable.getCredential({
        id: vcData[i].id
      }).then(resp => {
        console.log('get vc ' + vcData[i].id)

        vcData[i].credential= JSON.parse(resp.verifiableCredential)
              }
      ).catch(err =>
              console.log('get vc failed : errMsg=' + err)
      )
    }
  }

  async function fetchPresentations() {
    // Get the VP data
    for (let i = 0; i < vpData.length; i++) {
      await window.$aries.verifiable.getPresentation({
       id: `${vpData[i].id}`
     }).then(resp => {
               vpData[i].presentation = resp.verifiablePresentation
             }
     ).catch(err =>
      console.log('get vp failed : errMsg=' + err + " ID " +  vpData[i].id)
     )
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
    },
    methods: {
      getCredentials: async function (aries) {
        await aries.verifiable.getCredentials()
                .then(resp => {
                          vcData = resp.result
                          if (vcData && vcData.length == 0) {
                            console.log('no credentials exists')
                          }
                        }
                ).catch(err => {
                          console.log('get credentials failed : errMsg=' + err)
                        }
                )

        await fetchCredentials()
        this.verifiableCredential = vcData
      },
      getPresentations: async function(aries){
        await aries.verifiable.getPresentations()
                .then(resp => {
                          vpData = resp.result
                          if (vpData && vpData.length == 0) {
                            console.log('no presentation exists')
                          }
                        }
                ).catch(err => {
                          console.log('get presentation failed : errMsg=' + err)
                        }
                )

        await fetchPresentations()
        this.verifiablePresentation = vpData
      }
    },
    data() {
      return {
        verifiableCredential: [],
        verifiablePresentation: []
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


