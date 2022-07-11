<template>
  <div class="w-full">
    <ul
      class="flex w-full flex-row items-center justify-start bg-gradient-to-r from-haiti to-jagger md:rounded"
    >
      <li
        v-for="(tab, index) in tabs"
        :key="index"
        :class="[
          'w-full p-4 h-min-full rounded cursor-pointer',
          index === activeTab && 'bg-neutrals-mobster bg-opacity-40',
        ]"
      >
        <label
          :for="`DIDManagement-tab-${index}`"
          :class="[
            'font-bold cursor-pointer opacity-60 hover:opacity-100 focus:opacity-100 flex justify-center flex-wrap lg:text-lg md:text-base xs:text-sm text-center text-neutrals-white',
            index === activeTab && 'opacity-100',
          ]"
          v-text="tab"
        />
        <input
          :id="`DIDManagement-tab-${index}`"
          v-model="activeTab"
          type="radio"
          :name="`DIDManagement--tab`"
          :value="index"
          class="hidden"
        />
      </li>
    </ul>
    <template v-for="(tab, index) in tabs">
      <div v-if="index === activeTab" :key="index" class="grow">
        <slot :name="`tabPanel-${index}`" />
      </div>
    </template>
  </div>
</template>

<script>
export default {
  props: {
    tabs: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
};
</script>
