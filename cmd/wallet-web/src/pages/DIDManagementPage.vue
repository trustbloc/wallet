<!--
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="flex flex-col justify-start items-start w-full">
    <TabsComponent :tabs="tabs">
      <template #tabPanel-0><PreferencesComponent v-model:allDIDs="allDIDs" /></template>
      <template #tabPanel-1><CreateOrbComponent v-model:allDIDs="allDIDs" /></template>
      <template #tabPanel-2><ImportDIDComponent v-model:allDIDs="allDIDs" /></template>
    </TabsComponent>
  </div>
</template>

<script>
import CreateOrbComponent from '@/components/DIDManagement/CreateOrbComponent.vue';
import ImportDIDComponent from '@/components/DIDManagement/ImportDIDComponent.vue';
import PreferencesComponent from '@/components/DIDManagement/DIDPreferencesComponent.vue';
import TabsComponent from '@/components/TabbedInterface/TabsComponent.vue';

export default {
  name: 'DIDManagementPage',
  components: {
    CreateOrbComponent,
    ImportDIDComponent,
    PreferencesComponent,
    TabsComponent,
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
