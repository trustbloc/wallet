/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
	<div
		class="flex flex-col justify-between items-center min-w-screen min-h-screen
              bg-scroll bg-no-repeat bg-neutrals-softWhite bg-onboarding-sm md:bg-onboarding px-6"
	>
		<div class="flex flex-grow flex-col justify-center items-center">
			<!--Trustbloc Sign-up provider div -->
			<div
				class="flex flex-col justify-start items-center h-auto w-full sm:w-screen max-w-xl bg-gradient-dark rounded-xl overflow-hidden text-xl mx-6 px-6 mt-20 md:text-3xl"
			>
				<!-- TODO: add href url to the root component once it is implemented -->
				<Logo class="py-12" href="" />
				<div class="text-center items-center mb-10 md:mb-8">
					<span class="text-2xl md:text-4xl font-bold text-neutrals-white">
						Sign in to your account
					</span>
				</div>
				<div class="flex justify-center content-center w-full py-24 min-h-xl">
					<Spinner v-if="loading" />
					<button
						v-else
						v-for="(provider, index) in providers"
						:key="index"
						class="w-full h-11 max-w-xs flex flex-wrap items-center text-sm font-bold text-neutrals-dark py-2 px-4 mb-4
						bg-neutrals-softWhite rounded-md"
						@click="beginOIDCLogin(provider.id)"
					>
						<img class="object-contain inline-block max-h-6 mr-2" :src="provider.logoURL" />
						<span class="flex flex-wrap">{{ provider.signInText }}</span>
					</button>
				</div>
				<div class="text-center py-10 md:py-12">
					<p class="text-base font-normal text-neutrals-softWhite">
						Don't have an account?
						<a class="text-primary-blue whitespace-nowrap" href="/Signup">Sign up</a>
					</p>
				</div>
			</div>
		</div>
		<!-- Foooter -->
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
	components: {
		ContentFooter,
		Logo,
		Spinner,
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
		beginOIDCLogin: async function(providerID) {
			let leftPosition, topPosition, width, height;
			// Dimensions as per defined in ux design.
			width = 700;
			height = 770;
			leftPosition = window.screen.width / 2 - (width / 2 + 10);
			//Allow for title and status bars.
			topPosition = window.screen.height / 2 - (height / 2 + 50);
			//Open the pop up window.
			window.open(
				this.serverURL() + '/oidc/login?provider=' + providerID,
				'_blank',
				'status=no,height=' +
					height +
					',width=' +
					width +
					',resizable=yes,left=' +
					leftPosition +
					',top=' +
					topPosition +
					',screenX=' +
					leftPosition +
					',screenY=' +
					topPosition +
					',toolbar=no,menubar=no,scrollbars=no,location=no,directories=no'
			);
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
			this.$router.push(this.redirect);
		},
	},
};
</script>
