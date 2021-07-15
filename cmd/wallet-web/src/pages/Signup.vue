/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="bg-neutrals-softWhite flex  min-w-screen min-h-screen items-center justify-center px-5 py-5
              bg-scroll lg:bg-onboarding xl:bg-onboarding xs:bg-onboarding sm:bg-onboarding 2xl:bg-onboarding md:bg-onboarding bg-no-repeat">
        <div class="bg-gradient-dark rounded-xl overflow-hidden text-neutrals-softWhite text-xl md:text-3xl
    2xl:h-full 2xl:w-xl xl:h-full xl:w-xl sm:h-full sm:w-xl md:h-full md:w-xl w-full -mt-40 sm:-mt-56 xs:-mt-56 lg:mt-0 2xl:mt-0 xl:mt-0 md:mt-10">
            <!--Trustbloc Intro div  -->
            <div class="grid 2xl:grid-cols-2 lg:grid-cols-2 md:grid-cols-2 sm:grid-cols-1 xs:grid-cols-1
           bg-flare bg-no-repeat w-auto h-3/5 divide-x divide-neutrals-medium divide-opacity-25">
                <div class="col-span-1 hidden xl:block 2xl:block lg:block md:block py-14 px-8">
                    <ul class="list-none md:list-disc text-sm px-4 py-4">
                        <div class="flex items-center">
                            <img class="h-2 w-8 mr-2 mt-4 md-4 my-2 mx-4" :src="logoUrl">
                            <h2 class="font-bold text-4xl lg:text-6xl">TrustBloc</h2>
                        </div>
                        <ul class="w-full rounded-lg mt-2 mb-3 font-normal">
                            <li class="pl-3 pr-4 py-3 text-sm">
                                <div class="flex-1 flex items-center overflow-y-auto max-w-full">
                                    <img class="flex w-8 h-full" src="@/assets/img/onboarding-icon-1.svg"/>
                                    <span class="text-sm px-6 text-neutrals-softWhite align-middle">Keep your digital identity safe</span>
                                </div>
                            </li>
                            <li class="pl-3 pr-4 py-3 text-sm">
                                <div class="flex-1 flex items-center overflow-y-auto max-w-full">
                                    <img class="flex w-8 h-full" src="@/assets/img/onboarding-icon-2.svg"/>
                                    <span class="text-sm px-6 text-neutrals-softWhite align-middle">Store digital IDs, certifications, and
                    more—all in one secure wallet</span>
                                </div>
                            </li>
                            <li class="pl-3 pr-4 py-3 text-sm">
                                <div class="flex-1 flex items-center overflow-y-auto max-w-full">
                                    <img class="flex w-8 h-full" src="@/assets/img/onboarding-icon-3.svg"/>
                                    <span class="text-sm px-6 text-neutrals-softWhite align-middle">Verify your identity in person or online</span>
                                </div>
                            </li>
                        </ul>
                    </ul>
                </div>
                <!--Trustbloc Sign-up provider div -->
                <div class="col-span-1 md:block sm:object-none sm:object-center md:py-18 lg:px-4 xl:py-14 xl:px-4 lg:py-14">
                    <div class="flex justify-center items-center lg:hidden md:hidden xl:hidden xs:block sm:block ">
                        <img class="h-4 w-8 mr-1 mt-4 md-4 my-2" :src="logoUrl">
                        <h2 class="font-bold text-8xl lg:text-2xl md:text-4xl sm:text-6xl xs:text-8xl">TrustBloc</h2>
                    </div>
                    <div class="text-center items-center mb-10">
                        <h1 class="font-bold text-3xl 2xl:text-4xl lg:text-2xl md:text-4xl sm:text-6xl xs:text-8xl text-neutrals-softWhite">
                            Sign up. It's free!
                        </h1>
                    </div>
                    <div class="flex 2xl:m-8 xl:mx-8 lg:mx-4 md:mx-4 mx-4 align-middle content-center">
                        <div class="w-full px-8 mb-8 2xl:m-8 xl:m-8 lg:m-8 md:m-4 sm:mx-2 sm:mb-8 md:mb-8">
                            <button class="2xl:flex 2xl:flex-wrap lg:flex lg:flex-wrap md:flex md:flex-wrap text-xs 2xl:text-xs xl:text-xs lg:text-sm
                md:text-xs sm:text-xl text-center content-center w-full py-2 mb-4 px-6
                        font-bold bg-neutrals-softWhite text-neutrals-dark rounded-md"
                                    v-for="(provider, index) in providers" :key="index"
                                    v-on:click="beginOIDCLogin(provider.id)">
                                <img class="object-contain inline-block  max-h-4 max-w-2 mr-2"
                                     :src="provider.logoURL"/>
                                {{ provider.signUpText }}
                            </button>
                        </div>

                    </div>
                </div>
                <!-- Foooter -->
                <ContentFooter></ContentFooter>
            </div>
        </div>
    </div>
</template>

<script>
    import {CHAPIHandler, RegisterWallet} from "./mixins"
    import {DeviceLogin} from "@trustbloc/wallet-sdk"
    import ContentFooter from "@/pages/layout/ContentFooter.vue"
    import {mapActions, mapGetters} from 'vuex'
    import axios from 'axios';

    export default {
        created: async function () {
            await this.fetchProviders()
            //TODO: issue-601 Implement cookie logic with information from the backend.
            this.deviceLogin = new DeviceLogin(this.getAgentOpts()['edge-agent-server']);

            let redirect = this.$route.params['redirect']
            this.redirect = redirect ? {name: redirect} : `${__webpack_public_path__}`

            // if session found, then no need to login
            this.loadUser()
            if (this.getCurrentUser()) {
                this.handleSuccess()
                return
            }

            // call server to get user info and process login
            await this.loadOIDCUser()

            // register or let user inside wallet
            if (this.getCurrentUser()) {
                await this.finishOIDCLogin()
                this.handleSuccess()
                return
            }

            if (this.$cookies.isKey('registerSuccess')) {
                this.registered = true;
            }
        },
        components: {
            ContentFooter,
        },
        data() {
            return {
                providers: [],
                hubOauthProvider: this.hubURL(),
                statusMsg: '',
                loading: true,
                registered: false,
                messageNum: 0,
                messages: [
                    "Getting Started..",
                    "This could take a minute.",
                    "Please do not refresh the page.",
                    "Please wait...",
                ],
                logoUrl: this.getLogoUrl(),
            };
        },
        methods: {
            ...mapActions({
                loadUser: 'loadUser',
                loadOIDCUser: 'loadOIDCUser',
                startUserSetup: 'startUserSetup',
                completeUserSetup: 'completeUserSetup',
                refreshUserPreference: 'refreshUserPreference'
            }),
            ...mapGetters(['getCurrentUser', 'getAgentOpts', 'serverURL', 'hubURL', 'getStaticAssetsUrl']),
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            beginOIDCLogin: async function (providerID) {
                window.location.href = this.serverURL() + "/oidc/login?provider=" + providerID
            },

            // Fetching the providers from hub-auth
            fetchProviders: async function () {
                await axios.get(this.hubURL() + '/oauth2/providers')
                    .then(response => {
                        this.providers = response.data.authProviders
                    })
            },

            // Get logo url based on docker configuration
            getLogoUrl: function() {
              let staticAssetsUrl = this.getStaticAssetsUrl()
              if (staticAssetsUrl) {
                return this.logoUrl = `${staticAssetsUrl}/images/logo.svg`
              }
              return this.logoUrl = `${require('@/assets/img/logo.svg')}`
            },

            sleep(ms) {
                return new Promise(resolve => setTimeout(resolve, ms));
            },
            finishOIDCLogin: async function () {
                let user = this.getCurrentUser()
                this.registerUser(user)

                // all credential handlers registration should happen here, ex: CHAPI, WACI etc
                let chapi = new CHAPIHandler(this.$polyfill, this.$webCredentialHandler, this.getAgentOpts().credentialMediatorURL)
                await chapi.install(this.getCurrentUser().username)
            },
            async registerUser(user) {
                if (!user.preference) {
                    this.startUserSetup()

                    let user = this.getCurrentUser()

                    // first time login, register this user
                    await new RegisterWallet(this.getAgentInstance(), this.getAgentOpts()).register({
                        name: user.username,
                        user: user.profile.user,
                        token: user.profile.token
                    }, this.completeUserSetup)
                }

                this.refreshUserPreference()
            },
            handleSuccess() {
                this.$router.push(this.redirect);
            },
            loginDevice: async function () {
                await this.deviceLogin.login();
            },

        }
    }

</script>
