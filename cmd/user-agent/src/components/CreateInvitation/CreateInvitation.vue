/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <md-card class="md-card-plain">
    <md-dialog :md-active.sync="dialog">
      <md-dialog-title>{{ dialogTitle }}</md-dialog-title>
      <div class="content">
        <div class="md-layout">
          <div class="md-layout-item">
            <pre>{{ dialogContent }}</pre>
          </div>
        </div>
      </div>
      <md-dialog-actions>
        <copy-button @click.native="dialog = false" :content="dialogContent"/>
      </md-dialog-actions>
    </md-dialog>

    <md-card-header data-background-color="green">
      <h4 class="title">{{ title }}</h4>
    </md-card-header>
    <md-card-content style="background-color: white;">
      <md-field>
        <label>Alias</label>
        <md-input v-model="alias" required></md-input>
      </md-field>
      <div style="display: flow-root">
        <span class="error" v-if="error">{{ error }}</span>
        <md-button class="md-button md-info md-square right"
                   id='createInvitation'
                   v-on:click="submit">
          <b>Create</b>
        </md-button>
      </div>
    </md-card-content>
  </md-card>
</template>

<script>
import {mapGetters, mapActions} from 'vuex'
import CopyButton from "../CopyButton/CopyButton";

export default {
  name: "create-invitation",
  props: {
    title: {
      type: String,
      default: 'Create Invitation'
    },
  },
  components: {CopyButton},
  methods: {
    ...mapActions(['createInvitation']),
    async submit() {
      this.error = ""

      if (!this.isMediatorRegistered) {
        this.error = "Please register a router first!"
        return
      }

      if (this.alias.trim().length === 0) {
        this.error = "Please fill in the field!"
        return
      }

      try {
        let res = await this.createInvitation(this.alias.trim())
        this.showDialog("Invitation created!", JSON.stringify(res, null, 2))
      } catch (e) {
        this.error = e.message
      }
    },
    showDialog(title, content) {
      this.dialogTitle = title
      this.dialogContent = content
      this.dialog = true
    },
  },
  computed: mapGetters(['isMediatorRegistered']),
  data: () => ({
    dialog: false,
    dialogTitle: "",
    dialogContent: "",
    alias: '',
    error: '',
  })
}
</script>

<style scoped>
.title {
  display: -webkit-inline-box;
}

.right {
  float: right;
}

.error {
  color: red;
}
</style>