<template>
  <div>
    <div class="flex flex-row justify-between items-center mb-8 w-full">
      <div class="flex flex-grow">
        <h3 class="text-neutrals-dark">{{ i18n.credentialDetails }}</h3>
      </div>
      <div
        class="
          inline-flex
          items-center
          bg-neutrals-white
          rounded-lg
          border border-neutrals-chatelle
        "
      >
        <flyout-menu type="outline" />
      </div>
    </div>
    <div class="flex relative z-10 flex-row justify-start items-center w-full h-24">
      <div class="flex-none w-12 h-12 border-opacity-10">
        <img src="@/assets/img/credential--generic-icon.svg" />
      </div>
      <!-- Todo: Integrate the credential heading with real data-->
      <div class="flex-grow p-4 text-left text-neutrals-dark overflow-ellipsis"></div>
    </div>
    <div class="absolute items-center -mt-5 w-full bg-neutrals-white rounded-xl">
      <table class="mt-11 mb-6 ml-6 detailHeadingData">
        <thead>
          <tr class="text-left text-neutrals-dark">
            <th>Added on</th>
            <th>Expires on</th>
            <th>Last used</th>
            <th>Vault</th>
          </tr>
        </thead>
        <tbody class="text-left text-neutals-medium">
          <tr>
            <!-- Todo: Integrate the credential heading data with real data-->
            <td></td>
          </tr>
        </tbody>
      </table>
      <div class="mt-12">
        <p class="py-5 text-xl font-bold text-neutrals-dark">Verified Information</p>
        <table class="w-full border-t-2 border-neutrals-chatelle">
          <tr class="border-b-2 border-neutrals-thistle border-dotted">
            <!-- Todo: Integrate the credential data label and real data-->
            <td class="py-4 pr-6 pl-3 text-neutrals-medium tdlabel"></td>
            <td class="py-4 pr-6 pl-3 text-neutrals-dark tddata"></td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';

export default {
  name: 'CredentialDetail',
  components: {
    FlyoutMenu,
  },
  computed: {
    i18n() {
      return this.$t('Credentials');
    },
  },
};
</script>
<style scoped>
.detailSubHeadingContainer {
  height: 7.3rem; /* 117px */
  box-shadow: 0px 2px 12px 0px rgba(25, 12, 33, 0.1);
  position: absolute;
}
.detailHeadingData {
  width: 42.6rem; /* 682px */
  height: 3.06rem; /* 49px */
}
.tdlabel {
  width: 21rem;
}
.tddata {
  width: 31rem;
}
</style>
