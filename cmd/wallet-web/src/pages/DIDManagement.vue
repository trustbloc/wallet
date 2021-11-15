<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-start w-full">
    <tabs :tabs="tabs">
      <template #tabPanel-0><preferences v-model:allDIDs="allDIDs" /></template>
      <template #tabPanel-1><create-orb v-model:allDIDs="allDIDs" /></template>
      <template #tabPanel-2><import-did v-model:allDIDs="allDIDs" /></template>
    </tabs>
  </div>
</template>

<script>
import CreateOrb from '@/components/DIDManagement/CreateOrb.vue';
import ImportDid from '@/components/DIDManagement/ImportDid.vue';
import Preferences from '@/components/DIDManagement/Preferences.vue';
import Tabs from '@/components/TabbedInterface/Tabs.vue';

export default {
  name: 'DIDManagement',
  components: {
    CreateOrb,
    ImportDid,
    Preferences,
    Tabs,
  },
  data: function () {
    return {
      allDIDs: [],
      tabs: [
        'Digital Identity Preference',
        'Create ORB Digital Identity',
        'Import Any Digital Identity',
      ],
    };
  },
};
</script>
