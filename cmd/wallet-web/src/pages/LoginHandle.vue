<template>
  <div>
    <p>Redirecting...</p>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'LoginRedirect',
  props: {
    provider: {
      type: String,
      default: null,
    },
  },
  data() {
    return {
      userinfo: null,
    };
  },
  created: async function () {
    await this.checkIfUserLoggedIn();
    if (this.userInfo.ok) {
      this.updateUserOnboard();
      window.top.close();
    } else {
      window.location.href = this.serverURL() + '/oidc/login?provider=' + this.provider;
    }
  },
  methods: {
    ...mapActions({
      updateUserOnboard: 'updateUserOnboard',
    }),
    ...mapGetters(['serverURL']),
    checkIfUserLoggedIn: async function () {
      this.userInfo = await fetch(this.serverURL() + '/oidc/userinfo', {
        method: 'GET',
        credentials: 'include',
      });
    },
  },
};
</script>
