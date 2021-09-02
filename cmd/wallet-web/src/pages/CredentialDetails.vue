<template>
  <div class="flex flex-col justify-start items-start py-6 px-3">
    <div class="flex flex-row justify-between items-center mb-8 w-full">
      <div class="flex flex-grow">
        <!-- TODO: remove m-0 and font-bold once vue-material is scrapped -->
        <h3 class="m-0 font-bold text-neutrals-dark">{{ i18n.heading }}</h3>
      </div>
      <div
        class="inline-flex items-center rounded-lg border border-neutrals-chatelle bg-transparent"
      >
        <flyout-menu type="outline" />
      </div>
    </div>
    <!-- TODO: pass actual credential data to banner -->
    <!-- <banner :credential="credential" /> -->
    <banner />
    <!-- List of Credential Details -->
    <div class="flex flex-col justify-start items-start mt-12 w-full">
      <span class="mb-5 text-xl font-bold text-neutrals-dark">{{ i18n.verifiedInformation }}</span>
      <table class="w-full border-t border-neutrals-chatelle">
        <tr class="border-b border-neutrals-thistle border-dotted">
          <!-- Todo: Integrate the credential data label and real data-->
          <td class="py-4 pr-6 pl-3 text-neutrals-medium"></td>
          <td class="py-4 pr-6 pl-3 text-neutrals-dark"></td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import Banner from '@/components/CredentialDetails/Banner';
import FlyoutMenu from '@/components/FlyoutMenu/FlyoutMenu';

export default {
  name: 'CredentialDetails',
  components: {
    Banner,
    FlyoutMenu,
  },
  computed: {
    i18n() {
      return this.$t('CredentialDetails');
    },
  },
};
</script>
