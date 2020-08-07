<template>
  <md-card class="md-card-plain">
    <md-card-header data-background-color="green">
      <h4 class="title">{{ title }}</h4>
    </md-card-header>
    <md-card-content style="background-color: white;">
      <md-field>
        <label>Invitation ({{ type }})</label>
        <md-textarea v-model="invitation" required></md-textarea>
      </md-field>
      <div style="display: flow-root">
        <span class="error" v-if="error">{{ error }}</span>
        <span class="success" v-if="success">{{ success }}</span>
        <md-button class="md-button md-info md-square right" v-on:click="submit">
          <b>Receive</b>
        </md-button>
      </div>
    </md-card-content>
  </md-card>
</template>

<script>
import {mapGetters, mapActions} from 'vuex'

export default {
  name: "receive-invitation",
  props: {
    title: {
      type: String,
      default: 'Receive Invitation'
    },
    type: {
      type: String,
      default: 'json'
    },
  },
  methods: {
    ...mapActions(['acceptInvitation', 'createInvitation']),
    async submit() {
      this.error = ""

      if (!this.isMediatorRegistered) {
        this.error = "Please register a mediator first!"
        return
      }

      if (this.invitation.trim().length === 0) {
        this.error = "Please fill in the field!"
        return
      }

      let invitation = this.invitation.trim();
      if (this.type === "base64") {
        invitation = window.atob(this.invitation.trim())
      }

      try {
        // accepts invitation thought out-of-band protocol
        let res = await this.acceptInvitation({
          my_label: this.agentDefaultLabel,
          invitation: JSON.parse(invitation),
        })

        // shows connection id that has been received
        this.success = `Your connection ID is ${res['connection_id']}`
        this.invitation = ''
      } catch (e) {
        this.error = e.message
      }
    },
  },
  computed: mapGetters(['isMediatorRegistered', 'agentDefaultLabel']),
  data: () => ({
    invitation: '',
    error: '',
    success: '',
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

.success {
  color: green;
}

.error {
  color: red;
}

</style>