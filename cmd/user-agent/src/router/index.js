/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import DashboardLayout from "@/pages/Layout/DashboardLayout.vue";

import Dashboard from "@/pages/Dashboard.vue";
import Login from "@/pages/chapi/Login.vue";
import Logout from "@/pages/chapi/Logout.vue";
import TableList from "@/pages/TableList.vue";
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
const routes = [
    {
        path: "/",
        component: DashboardLayout,
        name: "dashboard",
        redirect: "dashboard",
        children: [
            {
                path: "dashboard",
                component: Dashboard,
                meta: { requiresAuth: true }
            },
            {
                path: "login",
                component: Login
            },
            {
                path: "logout",
                component: Logout
            },
            {
                path: "ViewVC",
                name: "View Wallet",
                component: TableList
            },
            {
                path: "MyVC",
                name: "Generate Presentation",
                component: TablePresentation
            },
            {
                path: "WebWallet",
                name: "Web Wallet Demo",
                component: WebWallet
            },
            {
                path: "DIDManagement",
                name: "DID Management",
                component: DIDManagement
            },
            {
                path: "connections",
                name: "Connections",
                component: Connections
            },
            {
                path: "relationships",
                name: "Relationships",
                component: Relationships
            }, {
                path: "issue-credential",
                name: "Issue Credential",
                component: IssueCredential
            }, {
                path: "present-proof",
                name: "Present Proof",
                component: PresentProof
            }
        ]
    },
    {
        path: '*',
        name: 'NotFound',
        component: NotFound
    },
    {
        path: "/StoreInWallet",
        component: StoreInWallet
    },
    {
        path: "/GetFromWallet",
        component: GetFromWallet
    },
    {
        path: "/worker",
        component: WalletWorker
    }
];



export default routes;
