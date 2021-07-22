/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
	<div
		class="flex flex-col justify-between items-center min-w-screen min-h-screen bg-scroll bg-no-repeat bg-neutrals-softWhite bg-onboarding-sm md:bg-onboarding px-6"
	>
		<div class="flex flex-grow flex-col justify-center items-center">
			<div
				class="h-auto md:max-w-4xl bg-gradient-dark rounded-xl overflow-hidden text-xl mt-20 md:text-3xl"
			>
				<!--Trustbloc Intro div  -->
				<div
					class="grid md:grid-cols-2 grid-cols-1 bg-no-repeat bg-flare w-full h-full divide-x divide-neutrals-medium divide-opacity-25 md:px-20"
				>
					<div class="col-span-1 hidden md:block py-24 pr-16">
						<Logo class="mb-12" href="" />

						<div class="flex-1 flex items-center overflow-y-auto max-w-full mb-8">
							<img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-1.svg" />
							<span class="text-base pl-5 text-neutrals-white align-middle">
								{{ i18n.leftContainer.span1 }}
							</span>
						</div>

						<div class="flex-1 flex items-center overflow-y-auto max-w-full mb-8">
							<img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-2.svg" />
							<span class="text-base pl-5 text-neutrals-white align-middle">
								{{ i18n.leftContainer.span2 }}
							</span>
						</div>

						<div class="flex-1 flex items-center overflow-y-auto max-w-full">
							<img class="flex w-10 h-10" src="@/assets/img/onboarding-icon-3.svg" />
							<span class="text-base pl-5 text-neutrals-white align-middle">
								{{ i18n.leftContainer.span3 }}
							</span>
						</div>
					</div>
					<!--Trustbloc Sign-up provider div -->
					<div class="col-span-1 md:block object-none object-center">
						<div class="px-6 md:pt-16 md:pb-12 md:pl-16 md:pr-0">
							<Logo class="md:hidden justify-center mt-12 my-2" href="" />
							<div class="text-center items-center mb-10">
								<h1 class="text-2xl md:text-4xl font-bold text-neutrals-white">
									{{ i18n.heading }}
								</h1>
							</div>
							<div class="flex justify-center content-center w-full py-24 min-h-xl">
								<Spinner v-if="loading" />
								<button
									v-else
									v-for="(provider, index) in providers"
									:key="index"
									class="w-full h-11 max-w-xs flex flex-wrap items-center text-sm font-bold text-neutrals-dark py-2 px-4 mb-4
                                    bg-neutrals-softWhite rounded-md"
									@click="beginOIDCLogin(provider.id)">
									<img class="object-contain inline-block max-h-6 mr-2" :src="provider.logoURL" />
									<span class="flex flex-wrap" id="signUpText">{{ provider.signUpText }}</span>
								</button>
							</div>
							<div class="text-center py-10 md:pb-0 md:pt-12">
								<p class="text-base font-normal text-neutrals-white">
									{{ i18n.redirect }}
									<a class="underline-blue text-primary-blue whitespace-nowrap" href="/Signin">{{ i18n.signin }}</a>
								</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<ContentFooter />
	</div>
</template>

<script>
import { CHAPIHandler, RegisterWallet } from './mixins';
import { DeviceLogin } from '@trustbloc/wallet-sdk';
import ContentFooter from '@/pages/layout/ContentFooter.vue';
import Logo from '@/components/Logo/Logo.vue';
import Spinner from '@/components/Spinner/Spinner.vue';
import { mapActions, mapGetters } from 'vuex';
import axios from 'axios';

export default {
	created: async function() {
		await this.fetchProviders();
		//TODO: issue-601 Implement cookie logic with information from the backend.
		this.deviceLogin = new DeviceLogin(this.getAgentOpts()['edge-agent-server']);

		let redirect = this.$route.params['redirect'];
		this.redirect = redirect ? { name: redirect } : `${__webpack_public_path__}`;

		// if session found, then no need to login
		this.loadUser();
		if (this.getCurrentUser()) {
			this.handleSuccess();
			return;
		}

		// call server to get user info and process login
		await this.loadOIDCUser();

		// register or let user inside wallet
		if (this.getCurrentUser()) {
			await this.finishOIDCLogin();
			this.handleSuccess();
			return;
		}

		if (this.$cookies.isKey('registerSuccess')) {
			this.registered = true;
		}
		this.loading = false;
	},
	components: {
		ContentFooter,
		Logo,
		Spinner,
	},
	computed: {
		i18n() {
			return this.$t("Signup")
		}
	},
	data() {
		return {
			providers: [],
			hubOauthProvider: this.hubURL(),
			statusMsg: '',
			loading: true,
			registered: false,
		};
	},
	methods: {
		...mapActions({
			loadUser: 'loadUser',
			loadOIDCUser: 'loadOIDCUser',
			startUserSetup: 'startUserSetup',
			completeUserSetup: 'completeUserSetup',
			refreshUserPreference: 'refreshUserPreference',
		}),
		...mapGetters(['getCurrentUser', 'getAgentOpts', 'serverURL', 'hubURL']),
		...mapGetters('agent', { getAgentInstance: 'getInstance' }),
    beginOIDCLogin: async function (providerID) {
      window.location.href = this.serverURL() + '/oidc/login?provider=' + providerID
    },
		// Fetching the providers from hub-auth
		fetchProviders: async function() {
			await axios.get(this.hubURL() + '/oauth2/providers').then((response) => {
				this.providers = response.data.authProviders;
			});
		},
		finishOIDCLogin: async function() {
			let user = this.getCurrentUser();
			this.registerUser(user);

			// all credential handlers registration should happen here, ex: CHAPI, WACI etc
			let chapi = new CHAPIHandler(
				this.$polyfill,
				this.$webCredentialHandler,
				this.getAgentOpts().credentialMediatorURL
			);
			await chapi.install(this.getCurrentUser().username);
		},
		async registerUser(user) {
      if (!user.preference) {
				this.startUserSetup();

				let user = this.getCurrentUser();
				// first time login, register this user
				await new RegisterWallet(this.getAgentInstance(), this.getAgentOpts()).register(
					{
						name: user.username,
						user: user.profile.user,
						token: user.profile.token,
					},
					this.completeUserSetup
				);
			}

			this.refreshUserPreference();
		},
		handleSuccess() {
			const route = this.$router.resolve({ name: this.redirect });
			window.open(route.href, '_top');
		},
	},
};
</script>
