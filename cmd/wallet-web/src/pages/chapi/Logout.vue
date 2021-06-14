/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div>
    <button
      type="button"
      v-on:click="logout"
      class="border-2 border-solid border-red-50 md-button logout-button"
    >
      Log Out
    </button>
  </div>
</template>

<script>
  import { RegisterWallet } from "./wallet";
  import { mapActions, mapGetters } from "vuex";

  export default {
    created: async function() {
      this.registrar = new RegisterWallet(
        this.$polyfill,
        this.$webCredentialHandler,
        this.getAgentInstance(),
        this.getAgentOpts()
      );
    },
    methods: {
      ...mapActions({ logoutUser: "logout" }),
      ...mapGetters("agent", { getAgentInstance: "getInstance" }),
      ...mapGetters(["getAgentOpts"]),
      logout: async function() {
        await this.registrar.uninstallHandlers();
        await this.logoutUser();
        this.$router.push("/login");
      },
    },
  };
</script>

<style scoped>
  .logout-button {
    width: 240px;
    height: 48px;
    font-size: 16px; /* remove once global styles for new designs are defined */
    font-weight: bold;
    border: 2px solid;
    border-image-slice: 1;
    border-image-source: linear-gradient(
      to left,
      #3f5fd3,
      #743ad5,
      #d53a9d,
      #cd3a67
    );
  }
  /*--Remove this once vue-material css is removed */
  .md-button {
    text-transform: none !important;
    background: transparent !important;
    font-size: large;
    font-family: sans-serif;
  }
  .md-button:hover {
    background: transparent !important;
  }
  .md-button:active:after {
    background: transparent !important;
  }
</style>
