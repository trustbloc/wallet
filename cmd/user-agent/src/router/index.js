/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import DashboardLayout from "@/pages/Layout/DashboardLayout.vue";

import Dashboard from "@/pages/Dashboard.vue";
import TableList from "@/pages/TableList.vue";
import TablePresentation from "@/pages/TablePresentation.vue";
import RegisterWallet from "@/pages/RegisterWallet.vue";
import StoreCredentials from "@/pages/StoreVC.vue";
import GetCredentials from "@/pages/GetVC.vue";
import CreateDID from "@/views/CreateDID.vue";
import Interop from "@/pages/Interop.vue";
const routes = [
    {
        path: "/",
        component: DashboardLayout,
        name: "dashboard",
        redirect: "dashboard",
        children: [
            {
                path: "dashboard",
                name: "Welcome to User Agent",
                component: Dashboard
            },
            {
                path: "ViewVC",
                name: "View Wallet",
                component: TableList
            },
            {
                path: "RegisterWallet",
                name: "Register Wallet",
                component: RegisterWallet
            },
            {
                path: "MyVC",
                name: "Generate Presentation",
                component: TablePresentation
            },
            {
                path: "StoreVC",
                name: "Store Credentials",
                component: StoreCredentials
            },
            {
                path: "GetVC",
                name: "Get Credentials",
                component: GetCredentials
            },
            {
                path: "CreateDID",
                name: "Create DID",
                component: CreateDID
            },
            {
                path: "Interop",
                name: "Interop",
                component: Interop
            }
        ]
    }

];
export default routes;
