/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div id="app">
    <ul>
      <li><router-link to="/" exact>Homepage</router-link></li>
      <li><router-link to="/RegisterWallet">Register Wallet</router-link></li>
      <li><router-link to="/MyVC">My Credentials</router-link></li>
    </ul>
    <router-view></router-view>
  </div>
</template>

<script>
  export default {
    name: "App",
  };
</script>
