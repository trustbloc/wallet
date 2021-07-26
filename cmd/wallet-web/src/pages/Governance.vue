/* Copyright SecureKey Technologies Inc. All Rights Reserved. SPDX-License-Identifier: Apache-2.0 */

<template>
  <div v-if="govnVC">
    <div style="margin-top: 20px" class="center-header">
      <md-button
        v-if="trusted"
        class="md-icon-button green-icon-button"
        @click="displayTrustConsent"
      >
        <md-icon class="black-icon">verified </md-icon>
      </md-button>

      <md-badge v-else md-content="1" style="margin: 0px" class="md-primary" md-dense>
        <md-button class="md-icon-button green-icon-button" @click="displayTrustConsent">
          <md-icon class="black-icon">verified </md-icon>
        </md-button>
      </md-badge>

      <span class="clickable-title" @click.prevent="displayTrustConsent">
        {{ requestOrigin }}
      </span>
    </div>

    <div :hidden="hideTrustConsent">
      <md-card>
        <md-card-header style="background-color: #eaeaea; margin: 10px">
          <md-card-header-text>
            <div class="md-title">
              <md-icon style="color: #00cc66" class="md-size-2x">verified_user </md-icon>
              {{ govnVC.credentialSubject.name }}
            </div>
            <div class="md-subhead">{{ govnVC.credentialSubject.description }}</div>
          </md-card-header-text>

          <md-card-media>
            <img :src="govnLogoSrc" @error="imageError = true" />
          </md-card-media>
        </md-card-header>

        <md-subheader>
          <span v-if="issuer" style="margin-top: 10px"
            >As per <b style="font-size: 16px">{{ govnVC.credentialSubject.name }}</b> rules,
            issuers should prove they're accredited by {{ govnVC.credentialSubject.name }} before
            issuing credentials.This app automatically challenged the current issuer,
            <b style="font-size: 14px; color: #0e9a00">{{ requestOrigin }}</b
            >, to do so, and received proof that it meets the requirement.</span
          >
          <span v-else style="margin-top: 10px"
            >As per <b style="font-size: 16px">{{ govnVC.credentialSubject.name }}</b> rules,
            verifiers should prove they're accredited by {{ govnVC.credentialSubject.name }} before
            asking for your credentials.This app automatically challenged the current verifier,
            <b style="font-size: 14px; color: #0e9a00">{{ requestOrigin }}</b
            >, to do so, and received proof that it meets the requirement.</span
          >
        </md-subheader>

        <md-card-actions>
          <md-button
            v-if="!trusted"
            class="md-raised"
            style="background-color: #0e9a00 !important"
            @click="trustRequester"
            >Acknowledge
          </md-button>

          <md-button
            v-if="!trusted"
            style="background-color: #fb4934 !important"
            @click="doNotTrustRequester"
            >Cancel
          </md-button>

          <md-button
            v-if="trusted"
            class="md-raised"
            style="background-color: #0b97c4 !important"
            @click="trustRequester"
          >
            <md-icon>verified_user</md-icon>
            {{ issuer ? 'This issuer is trusted' : 'This verifier is trusted' }}
          </md-button>
        </md-card-actions>
      </md-card>
    </div>
  </div>

  <div v-else>
    <div style="margin-top: 20px" class="center-header">
      <md-button
        v-if="trusted"
        class="md-icon-button green-icon-button"
        @click="displayTrustConsent"
      >
        <md-icon class="black-icon">verified</md-icon>
      </md-button>

      <md-badge v-else md-content="1" style="margin: 0px" md-dense>
        <md-button
          class="md-icon-button"
          style="background-color: #dce775 !important; margin: 0px"
          @click="displayTrustConsent"
        >
          <md-icon class="black-icon">notifications</md-icon>
        </md-button>
      </md-badge>

      <span class="clickable-title" @click.prevent="displayTrustConsent">
        {{ requestOrigin }}
      </span>
    </div>

    <div :hidden="hideTrustConsent">
      <md-card>
        <md-card-header style="background-color: #eaeaea; margin: 10px">
          <md-card-header-text>
            <div class="md-title">
              <md-icon style="color: #d73a49" class="md-size-3x">error_outline </md-icon>
              Not trusted
            </div>
            <div class="md-subhead">
              This {{ issuer ? 'issuer' : 'verifier' }} may not be trustworthy
            </div>
          </md-card-header-text>

          <md-card-media>
            <img :src="govnLogoSrc" @error="imageError = true" />
          </md-card-media>
        </md-card-header>

        <md-subheader>
          <span v-if="issuer" style="margin-top: 10px"
            >Issuers should prove they're accredited by trusted governance authority before issuing
            credentials. The current issuer,
            <b style="font-size: 14px; color: #d73a49">{{ requestOrigin }}</b
            >, has not proven they are accredited by a trusted governance authority.</span
          >
          <span v-else style="margin-top: 10px"
            >Verifiers should prove they're accredited by trusted governance authority before asking
            for your credentials. The current verifier,
            <b style="font-size: 14px; color: #d73a49">{{ requestOrigin }}</b
            >, has not proven they are accredited by a trusted governance authority.</span
          >
        </md-subheader>

        <md-card-actions>
          <md-button
            v-if="!trusted"
            class="md-raised"
            style="background-color: #0e9a00 !important"
            @click="trustRequester"
            >Proceed Anyways
          </md-button>

          <md-button style="background-color: #fb4934 !important" @click="doNotTrustRequester">
            <md-icon>pan_tool</md-icon>
            {{ issuer ? 'Do not trust this issuer' : 'Do not trust this verifier' }}
          </md-button>
        </md-card-actions>
      </md-card>
    </div>
  </div>
</template>
<script>
let govnDefaultLogo;

try {
  govnDefaultLogo = require('@/assets/img/govn.png');
} catch (e) {
  console.warn('unable to load default governance logo !');
}

export default {
  props: {
    govnVC: null,
    requestOrigin: null,
    issuer: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      imageError: false,
      // TODO should communicate store/services to avoid re-displaying trust consent
      trusted: false,
      showTrustConsent: false,
    };
  },
  computed: {
    hideTrustConsent() {
      return !this.showTrustConsent;
    },
    govnLogoSrc() {
      return !this.govnVC || this.imageError ? govnDefaultLogo : this.govnVC.credentialSubject.logo;
    },
  },
  methods: {
    displayTrustConsent: function () {
      console.log('this,.issuer', this.issuer);
      this.showTrustConsent = !this.showTrustConsent;
    },
    trustRequester: function () {
      this.showTrustConsent = false;
      this.trusted = true;
      //TODO should communicate store/services to avoid re-displaying trust consent
    },
    doNotTrustRequester: function () {
      this.showTrustConsent = false;
      this.trusted = false;
      //TODO should communicate store/services to avoid re-displaying trust consent
    },
  },
};
</script>

<style lang="css" scoped>
.green-icon-button {
  background-color: #00cc66 !important;
  margin: 0px;
}

.md-button:focus,
.md-button:active,
.md-button:hover,
.md-button.md-default:focus,
.md-button.md-default:active,
.md-button.md-default:hover {
  background-color: #00cc66 !important;
}

.black-icon {
  color: black !important;
  height: 25px !important;
}

.clickable-title {
  margin-left: 10px;
  font-size: 15px;
  font-weight: 400;
  font-style: italic;
  font-family: monospace;
}

.center-header {
  text-align: center !important;
}
</style>
