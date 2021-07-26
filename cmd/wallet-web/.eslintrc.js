module.exports = {
	root: true,
	parser: 'vue-eslint-parser',
	parserOptions: {
		vueFeatures: {
			filter: true,
			interpolationAsNonHTML: false,
		},
	},
	extends: [
		// TODO: Switch to the following after upgrading to Vue.js 3.x
		// 'plugin:vue/vue3-recommended'
		'plugin:vue/recommended',
		'plugin:tailwindcss/recommended',
		'plugin:eslint-comments/recommended',
		'plugin:i18n-json/recommended',
		'plugin:prettier/recommended',
	],
	plugins: ['vue', 'tailwindcss', 'eslint-comments', 'i18n-json', 'prettier'],
	// Default rules for any file we lint
	rules: {
		/**
		 * Force prettier formatting
		 */
		'prettier/prettier': 'error',
		/**
		 * Disallow the use of console
		 * https://eslint.org/docs/rules/no-console
		 */
		'no-console': 'off',

		/**
		 * Disallow Reassignment of Function Parameters
		 * https://eslint.org/docs/rules/no-param-reassign
		 */
		'no-param-reassign': ['error', { props: false }],

		/** Disallows unnecessary return await
		 * https://eslint.org/docs/rules/no-return-await
		 */
		'no-return-await': 'error',

		/**
		 * Disallow using an async function as a Promise executor
		 * https://eslint.org/docs/rules/no-async-promise-executor
		 */
		'no-async-promise-executor': 'error',

		/**
		 * Disallow await inside of loops
		 * https://eslint.org/docs/rules/no-await-in-loop
		 */
		'no-await-in-loop': 'error',

		/**
		 * Disallow assignments that can lead to race conditions due to
		 * usage of await or yield
		 * https://eslint.org/docs/rules/require-atomic-updates
		 */
		'require-atomic-updates': 'error',

		/**
		 * Disallow async functions which have no await expression
		 * https://eslint.org/docs/rules/require-await
		 */
		'require-await': 'error',

		/**
		 * Require or disallow named function expressions
		 * https://eslint.org/docs/rules/func-names
		 */
		'func-names': 'off',

		/**
		 * Disallow enforcement of consistent linebreak style
		 * https://eslint.org/docs/rules/func-names
		 */
		'linebreak-style': 'off',

		/**
		 * Allow ES6 classes to override methods without using this
		 * https://eslint.org/docs/rules/class-methods-use-this
		 */
		'class-methods-use-this': 'off',
	},
	settings: {
		tailwindcss: {
			config: 'tailwind.js',
		},
	},
};
