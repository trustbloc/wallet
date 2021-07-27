<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <md-card class="md-card-plain associated-credentials">
    <md-card-header data-background-color="green">
      <h4 class="title">{{ title }}</h4>
      <div class="title-btn-right">
        <md-button class="md-icon-button md-dense md-raised md-info" @click="getCredentials">
          <md-icon>cached</md-icon>
        </md-button>
      </div>
    </md-card-header>
    <md-card-content class="white">
      <div v-if="credentials.length === 0" class="text-center">No connections</div>
      <md-content class="md-content-connections md-scrollbar">
        <md-list class="md-triple-line">
          <md-list-item v-for="credential in credentials" :key="credential.id">
            <div class="md-list-item-text">
              <span v-if="credential.name">
                <span style="color: red">{{
                  credential.name ? credential.name : credential.id
                }}</span>
                credential was issued by
                <span style="color: green">{{ credential.label }}</span>
              </span>
            </div>
          </md-list-item>
        </md-list>
      </md-content>
    </md-card-content>
  </md-card>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'AssociatedCredentials',
  props: {
    title: {
      type: String,
      default: 'Associated credentials',
    },
    credentials: {
      type: Array,
    },
  },
  methods: {
    ...mapActions(['getCredentials']),
  },
  mounted() {
    this.getCredentials();
  },
};
</script>

<style scoped>
.white {
  background-color: white;
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

.md-list {
  width: 100%;
  display: inline-block;
  vertical-align: center;
}

.title {
  display: -webkit-inline-box;
}

.md-content-connections {
  width: 100%;
  max-height: 500px;
}
</style>
