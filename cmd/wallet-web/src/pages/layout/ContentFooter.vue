/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="fixed bottom-0 mx-auto px-18 mb-4 w-full sm:px-20 xs:px-20 lg:mx-0 xl:mx-0 2xl:mx-0 text-center">
    <div class="text-neutrals-medium mx-auto 2xl:block xl:block lg:block md:block hidden  text-center">
       © {{ date }} TrustBloc ･ Privacy Policy ･ Terms of Service  ･
      <select class="py-4" style="background-color: transparent">
        <option class="text-neutrals-dark block  text-sm"  id="english">English</option>
        <option class="text-neutrals-medium block  text-sm"  id="french">Français</option>
      </select>
    </div>
    <div class="text-neutrals-medium relative 2xl:hidden xl:hidden lg:hidden md:hidden block">
      <ul class="list-inside">
        <li> Privacy Policy ･ Terms of Service</li>
        <li>© {{ date }} TrustBloc ･
          <select class="py-4" style="background-color: transparent">
            <option class="text-neutrals-dark block  text-sm">English</option>
            <option class="text-neutrals-medium block  text-sm">Français</option>
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
