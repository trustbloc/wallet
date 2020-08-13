/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <md-card class="md-card-plain">
    <md-card-header data-background-color="green">
      <h4 class="title">{{ title }}</h4>
      <div class="title-btn-right">
        <copy-button :content="content"/>
      </div>
    </md-card-header>
    <md-card-content class="white">
      <div class="text-center">
        <div style="line-break: anywhere;" v-if="content">{{ content }}</div>
        <div class="error" v-if="error">{{ error }}</div>
      </div>
      <input type="hidden" id="created-invitation" :value="content">
    </md-card-content>
  </md-card>
</template>

<script>
import {mapGetters, mapActions} from 'vuex'
import CopyButton from "../CopyButton/CopyButton";

export default {
  name: "public-invitation",
  props: {
    title: {
      type: String,
      default: 'Make New Friends By Sharing This Invitation!'
    },
  },
  components: {CopyButton},
  methods: {
    ...mapActions(['createInvitation']),
    async generatePublicInvitation() {
      if (!this.isMediatorRegistered) {
        this.error = 'Invitation can\'t be generated without a mediator. Please, set up a mediator first.'

        return
      }

      try {
        let res = await this.createInvitation(this.agentDefaultLabel)
        // encodes invitation to base64 string
        this.content = window.btoa(JSON.stringify(res))
        this.error = ''
      } catch (e) {
        console.error(e)
        this.error = 'Something went wrong :('
      }
    },
  },
  beforeMount() {
    this.generatePublicInvitation()
  },
  watch: {
    isMediatorRegistered() {
      this.generatePublicInvitation()
    }
  },
  computed: mapGetters(['isMediatorRegistered', 'agentDefaultLabel']),
  data: () => ({
    content: '',
    error: '',
  })
}
</script>
<style scoped>
.white {
  background-color: white;
}

.error {
  line-break: anywhere;
  color: red;
}

.title {
  width: calc(100% - 42px);
  display: -webkit-inline-box;
}

.title-btn-right {
  right: 20px;
  position: absolute;
  display: -webkit-inline-box;
}
</style>