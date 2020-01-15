/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div id="app">
    <nav>
      <router-link :to="{ name: 'createInvitation' }">create invitation</router-link>
    </nav>
    <nav>
      <router-link :to="{ name: 'storeVC' }">store vc</router-link>
    </nav>
    <router-view></router-view>
  </div>
</template>

<script>
  export default {
    name: "App",
  };
</script>
