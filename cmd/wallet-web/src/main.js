/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import store from './store'
import App from './App.vue';
import VueRouter from "vue-router";
import routes from "./router/index";
import * as polyfill from "credential-handler-polyfill";
import * as webCredentialHandler from "web-credential-handler";
import {mapActions, mapGetters} from "vuex";
import VueCookies from 'vue-cookies'
import '@/assets/css/tailwind.css'
import SideBar from "@/components/SidebarPlugin";
import VueMaterial from "vue-material";
import "vue-material/dist/vue-material.css";
import "./assets/scss/material-dashboard.scss";
import i18n from './i18n'

Vue.config.productionTip = false

Vue.prototype.$polyfill = polyfill
Vue.prototype.$webCredentialHandler = webCredentialHandler

// configure router
const router = new VueRouter({
    mode: 'history',
    routes, // short for routes: routes
    linkExactActiveClass: "nav-item active"
});

router.beforeEach((to, from, next) => {
    store.dispatch('agent/flushStore')
    if (to.matched.some(record => record.meta.requiresAuth)) {
        if (store.getters.getCurrentUser) {
            next();
        } else if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
            next()
        } else {
            next({
                name: "signup",
                params: {redirect: to.name},
            });
        }
    } else if (to.matched.some(record => record.meta.blockNoAuth)) {
        if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
            next()
        } else {
            console.log("other options ")
            next({
                name: "block-no-auth",
                params: {loginURL: "signup"},
            });
        }
    } else {
        next();
    }
})

Vue.use(VueRouter);
Vue.use(SideBar);
Vue.use(VueCookies);
Vue.$cookies.config('7d');
Vue.use(VueMaterial);

new Vue({
    store,
    el: "#app",

    data: () => ({
        loaded: false,
    }),

    methods: {
        ...mapActions(['initOpts', 'loadUser']),
        ...mapActions('agent', {initAgent: 'init'}),
        ...mapGetters('agent', {isAgentInitialized: 'isInitialized'}),
        ...mapGetters(['getAgentOpts']),
    },

    mounted: async function () {
        // load opts
        await this.initOpts()

        // load user if already logged in
        this.loadUser()

        // load agent if user already logged in and agent not initialized (scenario: page refresh)
        if (store.getters.getCurrentUser && !this.isAgentInitialized()) {
            await this.initAgent()
        }

        // removes spinner
        this.loaded = true
    },

    render: h => h(App),
    i18n,
    router
});

window.onbeforeunload = function(){
    if (store) {
        store.dispatch('agent/flushStore')
    }
}


