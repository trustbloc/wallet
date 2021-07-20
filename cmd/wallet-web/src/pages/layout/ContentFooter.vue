/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="w-full mx-6 px-15 py-6 text-center">
    <div class="text-neutrals-medium mx-auto md:block hidden text-center">
       © {{ date }} TrustBloc ･ Privacy Policy ･ Terms of Service  ･
      <select class="py-4" style="background-color: transparent">
        <option class="text-neutrals-dark block text-sm" id="english">English</option>
        <option class="text-neutrals-medium block text-sm" id="french">Français</option>
      </select>
    </div>
    <div class="text-neutrals-medium md:hidden block">
      <ul class="list-inside">
        <li><span class="whitespace-nowrap">Privacy Policy</span> ･ <span class="whitespace-nowrap">Terms of Service</span></li>
        <li><span class="whitespace-nowrap">© {{ date }} TrustBloc</span> ･
          <select class="py-4" style="background-color: transparent">
            <option class="text-neutrals-dark block text-sm">English</option>
            <option class="text-neutrals-medium block text-sm">Français</option>
          </select>
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
export default {
  name: "content-footer",
  data() {
    return {
      date: new Date().getFullYear(),
    };
  },
};
</script>
