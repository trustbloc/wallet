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
import * as trustblocAgent from "@trustbloc/trustbloc-agent"
import MaterialDashboard from "./material-dashboard";
import {mapActions, mapGetters} from "vuex";

Vue.config.productionTip = false

Vue.prototype.$polyfill = polyfill
Vue.prototype.$webCredentialHandler = webCredentialHandler
Vue.prototype.$trustblocAgent = trustblocAgent

// configure router
const router = new VueRouter({
    mode: 'history',
    routes, // short for routes: routes
    linkExactActiveClass: "nav-item active"
});

router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.requiresAuth)) {
        if (store.getters.getCurrentUser) {
            next();
        } else if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
            next()
        } else {
            next({
                name: "login",
                params: {redirect: to.name},
            });
        }
    } else if (to.matched.some(record => record.meta.blockNoAuth)) {
        if (store.dispatch('loadUser') && store.getters.getCurrentUser) {
            next()
        } else {
            next({
                name: "block-no-auth",
                params: {loginURL: "login"},
            });
        }
    } else {
        next();
    }
})

Vue.use(VueRouter);
Vue.use(MaterialDashboard);

new Vue({
    store,
    el: "#app",
    data: () => ({
        loaded: false,
    }),
    methods: {
        ...mapActions(['initOpts', 'loadUser']),
        ...mapActions('aries', {initAries: 'init'}),
        ...mapGetters('aries', {isAriesInitialized: 'isInitialized'}),
        ...mapGetters(['getTrustblocOpts']),
    },
    mounted: async function () {
        // load opts
        await this.initOpts()

        // load user if already logged in
        this.loadUser()

        // load aries if user already logged in and aries not initialized (scenario: page refresh)
        if (store.getters.getCurrentUser && !this.isAriesInitialized()) {
            await this.initAries()
        }

        // removes spinner
        this.loaded = true
    },
    render: h => h(App),
    router
});

