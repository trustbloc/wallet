/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import DashboardLayout from "@/pages/Layout/DashboardLayout.vue";

import Dashboard from "@/pages/Dashboard.vue";
import Login from "@/pages/chapi/Login.vue";
import Logout from "@/pages/chapi/Logout.vue";
import TablePresentation from "@/pages/TablePresentation.vue";
import StoreInWallet from "@/pages/chapi/Store.vue";
import GetFromWallet from "@/pages/chapi/Get.vue";
import WalletWorker from "@/pages/chapi/Worker.vue";
import WebWallet from "@/pages/WebWallet.vue";
import DIDManagement from "@/pages/DIDManagement.vue";
import Connections from "@/pages/Connections.vue";
import Relationships from "@/pages/Relationships.vue";
import IssueCredential from "@/pages/IssueCredential.vue";
import PresentProof from "@/pages/PresentProof.vue";
import NotFound from '@/pages/PageNotFound'
import BlockNoAuth from "@/pages/chapi/BlockNoAuth.vue";

const routes = [
    {
        path: "/",
        component: DashboardLayout,
        // name: "dashboard",
        redirect: "dashboard",
        children: [
            {
                path: "dashboard",
                name: "dashboard",
                component: Dashboard,
                meta: {requiresAuth: true}
            },
            {
                path: "logout",
                name: "logout",
                component: Logout
            },
            {
                path: "MyVC",
                name: "my-credential",
                component: TablePresentation
            },
            {
                path: "WebWallet",
                name: "web-wallet",
                component: WebWallet
            },
            {
                path: "DIDManagement",
                name: "did-management",
                component: DIDManagement,
                meta: {requiresAuth: true}
            },
            {
                path: "connections",
                name: "connections",
                component: Connections,
                meta: {requiresAuth: true}
            },
            {
                path: "relationships",
                name: "relationships",
                component: Relationships,
                meta: {requiresAuth: true}
            },
            {
                path: "issue-credential",
                name: "issue-credential",
                component: IssueCredential,
                meta: {requiresAuth: true}
            },
            {
                path: "present-proof",
                name: "present-proof",
                component: PresentProof,
                meta: {requiresAuth: true}
            }
        ]
    },
    {
        path: "/login",
        name: "login",
        component: Login
    },
    {
        path: '*',
        component: NotFound
    },
    {
        path: "/StoreInWallet",
        name: "chapi-store",
        component: StoreInWallet,
        meta: {blockNoAuth: true}
    },
    {
        path: "/GetFromWallet",
        name: "chapi-get",
        component: GetFromWallet,
        meta: {blockNoAuth: true}
    },
    {
        path: "/worker",
        name: "chapi-worker",
        component: WalletWorker
    },
    {
        path: '/needauth',
        name: "block-no-auth",
        component: BlockNoAuth
    },
];


export default routes;
