/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { createApp } from 'vue';
import * as polyfill from 'credential-handler-polyfill';
import * as webCredentialHandler from 'web-credential-handler';
import i18n from '@/plugins/i18n';
import store from '@/store';
import router from '@/router';
import '@/assets/css/tailwind.css';
import App from '@/App.vue';
import ToastNotification from '@/components/ToastNotification/ToastNotification.vue';

const app = createApp(App);

app.use(router);
app.use(store);
app.use(i18n);
app.provide('polyfill', polyfill);
app.provide('webCredentialHandler', webCredentialHandler);
app.component('ToastNotification', ToastNotification);

app.mount('#app');
