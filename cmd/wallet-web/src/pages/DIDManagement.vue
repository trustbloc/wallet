<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-start w-full">
    <tabs>
      <tab title="Digital Identity Preference"><preferences :all-d-i-ds.sync="allDIDs" /></tab>
      <tab title="Create ORB Digital Identity"><create-orb :all-d-i-ds.sync="allDIDs" /></tab>
      <tab title="Import Any Digital Identity"><import-did :all-d-i-ds.sync="allDIDs" /></tab>
    </tabs>
  </div>
</template>

<script>
import CreateOrb from '@/components/DIDManagement/CreateOrb.vue';
import ImportDid from '@/components/DIDManagement/ImportDid.vue';
import Preferences from '@/components/DIDManagement/Preferences.vue';
import Tab from '@/components/TabbedInterface/Tab.vue';
import Tabs from '@/components/TabbedInterface/Tabs.vue';

export default {
  name: 'DIDManagement',
  components: {
    CreateOrb,
    ImportDid,
    Preferences,
    Tab,
    Tabs,
  },
  data: function () {
    return {
      allDIDs: [],
    };
  },
};
</script>
