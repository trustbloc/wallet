<template>
<div>
  <p>Redirecting...</p>
</div>
</template>

<script>
import {mapGetters} from "vuex";

export default {
  data () {
    return {
    }
  },
  props: {
    provider: {
      type: String,
      default: null
    }
  },
  computed:{
    getState() {
      return this.getLoginHandleSetup();
    },
  },
  name: "LoginRedirect",
   created() {
      this.window.height = window.innerHeight
},
  mounted() {
   console.log("wow in different tab", this.getLoginHandleSetup(), this.getState)
   if (this.getState !== "success") {
     window.location.href= this.serverURL() + '/oidc/login?provider=' + this.provider
   }
  },
  methods:{
    ...mapGetters(['serverURL']),
    ...mapGetters(['getLoginHandleSetup']),
  }
}
</script>

<style scoped>
.popup {
  position: fixed;
  opacity: 0.9;
  height: 100vh;
  width: 100vw;
  overflow: visible;
  z-index: 100;
  top: 0;
  left: 0;
}

</style>