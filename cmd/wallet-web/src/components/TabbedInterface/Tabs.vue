<template>
  <div class="w-full">
    <ul
      class="
        justify-start
        items-center
        w-full
        bg-gradient-to-r
        from-haiti
        to-jagger
        md:rounded
        flex flex-row
      "
    >
      <li
        v-for="(tab, index) in tabs"
        :key="tab.title"
        :class="[
          'w-full p-4 h-min-full rounded cursor-pointer',
          index === selectedIndex && 'bg-neutrals-mobster bg-opacity-40',
        ]"
        @click="selectTab(index)"
      >
        <span
          :class="[
            'font-bold opacity-60 hover:opacity-100 focus:opacity-100 flex justify-center flex-wrap lg:text-lg md:text-base xs:text-sm text-center text-neutrals-white',
            index === selectedIndex && 'opacity-100',
          ]"
          >{{ tab.title }}</span
        >
      </li>
    </ul>
    <slot></slot>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedIndex: 0,
      tabs: [],
    };
  },
  created() {
    this.tabs = this.$children;
  },
  mounted() {
    this.selectTab(0);
  },
  methods: {
    selectTab(i) {
      this.selectedIndex = i;
      this.tabs.forEach((tab, index) => {
        tab.isActive = index === i;
      });
    },
  },
};
</script>
