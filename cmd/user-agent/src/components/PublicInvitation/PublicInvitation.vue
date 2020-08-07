<template>
  <md-card class="md-card-plain">
    <md-card-header data-background-color="green">
      <h4 class="title">{{ title }}</h4>
      <md-button :disabled="!isMediatorRegistered" v-on:click="copyInvitationToClipboard"
                 class="md-icon-button md-dense md-raised md-info right refresh-connections">
        <md-icon>content_copy</md-icon>
      </md-button>
    </md-card-header>
    <md-card-content style="background-color: white;">
      <div style="display: flow-root;text-align: center;">
        <div style="line-break: anywhere;">{{ content }}</div>
      </div>
      <input type="hidden" id="created-invitation" :value="content">
    </md-card-content>
  </md-card>
</template>

<script>
import {mapGetters, mapActions} from 'vuex'

export default {
  name: "public-invitation",
  props: {
    title: {
      type: String,
      default: 'Make New Friends By Sharing This Invitation!'
    },
  },
  methods: {
    ...mapActions(['createInvitation']),
    copyInvitationToClipboard: function () {
      let inv = document.querySelector('#created-invitation')
      inv.setAttribute('type', 'text')
      inv.select()

      document.execCommand('copy');

      inv.setAttribute('type', 'hidden')
      window.getSelection().removeAllRanges()
    },
    async generatePublicInvitation() {
      if (!this.isMediatorRegistered) {
        this.content = 'Invitation can\'t be generated without a mediator. Please, set up a mediator first.'

        return
      }
      let res = await this.createInvitation(this.agentDefaultLabel)
      // encodes invitation to base64 string
      this.content = window.btoa(JSON.stringify(res))
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
</style>