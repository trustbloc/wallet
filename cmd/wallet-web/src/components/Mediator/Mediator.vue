/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <md-card class="md-card-plain">
    <md-card-header data-background-color="green">
      <h4 class="title">{{ title }}</h4>
    </md-card-header>
    <md-card-content style="background-color: white;">
      <md-field>
        <label>Mediator URL</label>
        <md-input placeholder="http://router.example.com" v-model="URL" required></md-input>
      </md-field>
      <div style="display: flow-root">
        <span class="error" v-if="error">{{ error }}</span>
        <md-button v-bind:disabled="disabled" class="md-button md-info md-square right" v-on:click="register">
          <b>Register</b>
        </md-button>
      </div>
    </md-card-content>
    <md-card-content v-for="connection in getMediatorConnections" style="background-color: white;" :key="connection">
      <div class="md-layout router">
        <div class="md-layout-item md-layout" style="padding-right: 0px">
          <div class="md-layout-item router-done">
            <div>Mediator is registered {{ connection }}</div>
          </div>
          <div class="md-layout-item md-size-25" style="padding-right: 0px">
            <md-button class="md-button md-danger md-square right" id='routerUnregister'
                       v-on:click="unregisteredMediator(connection)">
              <b>Unregister</b>
            </md-button>
          </div>
        </div>
      </div>
    </md-card-content>
  </md-card>
</template>

<script>
import {mapGetters, mapActions} from 'vuex'

export default {
  name: "mediator",
  props: {
    title: {
      type: String,
      default: 'Mediator'
    },
  },
  methods: {
    ...mapActions(['unregisteredMediator', 'registeredMediator']),
    async register() {
      this.error = ""

      let URL = this.URL.trim().replace(/\/$/, "");
      if (URL.length === 0) {
        this.error = "Please fill in the field!"
        return
      }

      this.disabled = true

      try {
        await this.registeredMediator(URL)
      } catch (e) {
        this.error = e.message
      }

      this.disabled = false
    },
  },
  computed: mapGetters(['getMediatorConnections']),
  data: () => ({
    URL: '',
    error: '',
    disabled: false,
  })
}
</script>

<style scoped>
.router-done {
  display: flex;
  align-items: center;
}

.right {
  float: right;
}

.error {
  color: red;
}
</style>