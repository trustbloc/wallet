<script setup>
import { computed, onMounted, ref } from 'vue';
import { useStore } from 'vuex';
import SpinnerIcon from '@/components/icons/SpinnerIcon.vue';

const props = defineProps({
  providerID: {
    type: String,
    default: null,
  },
});

// Local Variables
const userInfo = ref();

// Hooks
const store = useStore();

// Store Getters
const serverUrl = computed(() => store.getters['serverURL']);

// Store Actions
const updateUserOnboard = () => store.dispatch('updateUserOnboard');
const updateLoginSuspended = () => store.dispatch('updateLoginSuspended');

// Methods
async function checkIfUserLoggedIn() {
  userInfo.value = await fetch(serverUrl.value + '/oidc/userinfo', {
    method: 'GET',
    credentials: 'include',
  });
}
function onclose() {
  updateLoginSuspended();
}

onMounted(async () => {
  await checkIfUserLoggedIn();
  if (userInfo.value.ok) {
    updateUserOnboard();
    window.top.close();
  } else {
    window.location.href = serverUrl.value + '/oidc/login?provider=' + props.providerID;
    window.addEventListener('beforeunload', onclose);
  }
});
</script>

<template>
  <div class="flex flex-col justify-center items-center w-screen h-screen">
    <SpinnerIcon />
    <span class="mt-8 text-base text-neutrals-dark md:text-lg">Redirecting . . .</span>
  </div>
</template>
