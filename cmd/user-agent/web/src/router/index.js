/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

export default new Router({
    mode: 'history',
    routes: [
        {
            name: "homepage",
            path: "/",
            component: () => import("@/views/Homepage")
        },
        {
            name: "createInvitation",
            path: "/CreateInvitation",
            component: () => import("@/views/CreateInvitation")
        },
        {
            name: "storeVC",
            path: "/StoreVC",
            component: () => import("@/views/StoreVC")
        },
        {
            name: "registerWallet",
            path: "/RegisterWallet",
            component: () => import("@/views/RegisterWallet")
        },
        {
            name: "worker",
            path: "/Worker",
            component: () => import("@/views/Worker")
        }
    ]
});
