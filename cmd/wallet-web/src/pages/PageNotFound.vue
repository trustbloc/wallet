<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div>
    <h1>Not Found</h1>
    <p>
      Ooops, we couldn't find what you are looking for. Why don't you
      <router-link :to="{ name: 'dashboard' }"> go home instead? </router-link>
    </p>
  </div>
</template>

<script>
export default {};
</script>

<style scoped>
h1 {
  font-size: 100px;
}

p {
  font-size: 50px;
  font-weight: 100;
}

div {
  text-align: center;
}
</style>
