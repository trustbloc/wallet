<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <md-card class="md-card-plain associated-presentations">
    <md-card-header data-background-color="green">
      <h4 class="title">{{ title }}</h4>
      <div class="title-btn-right">
        <md-button class="md-icon-button md-dense md-raised md-info" @click="getPresentations">
          <md-icon>cached</md-icon>
        </md-button>
      </div>
    </md-card-header>
    <md-card-content class="white">
      <div v-if="presentations.length === 0" class="text-center">No connections</div>
      <md-content class="md-content-connections md-scrollbar">
        <md-list class="md-triple-line">
          <md-list-item v-for="presentation in presentations" :key="presentation.id">
            <div class="md-list-item-text">
              <span v-if="presentation.name">
                <span style="color: red">{{
                  presentation.name ? presentation.name : presentation.id
                }}</span>
                presentation was exchanged by
                <span style="color: green">{{ presentation.label }}</span>
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
  name: 'AssociatedPresentation',
  props: {
    title: {
      type: String,
      default: 'Associated Presentations',
    },
    presentations: {
      type: Array,
    },
  },
  methods: {
    ...mapActions(['getPresentations']),
  },
  mounted() {
    this.getPresentations();
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
