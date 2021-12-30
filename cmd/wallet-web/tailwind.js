/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const plugin = require('tailwindcss/plugin');

module.exports = {
  purge: {
    // Uncomment the following line to enable purging in dev mode (e.g. for testing)
    // enabled: true,
    content: ['./src/**/*.html', './src/**/*.vue', '../../test/**/*.html'],
    // Safe-listing several classes which we have to supply dynamic strings to (`bg-gradient-${color}`, etc.)
    options: { safelist: [/bg-gradient-.*/, /ring-primary-.*/] },
  },
  darkMode: false, // or 'media' or 'class'
  theme: {
    colors: {
      transparent: 'transparent',
      neutrals: {
        black: '#000000',
        bonjour: '#eeeaee',
        chatelle: '#c7c3c8',
        dark: '#190c21',
        light: '#b6b7c7',
        lilacSoft: '#edeaee',
        magnolia: '#fbf8fc',
        medium: '#6c6d7c',
        moist: '#f5f4f6',
        mischka: '#dcdadd',
        mobster: '#7b6f7f',
        mountainMist: {
          light: '#949095',
          dark: '#8d8a8e',
        },
        softWhite: '#f4f1f5',
        thistle: '#dbd7dc',
        victorianPewter: '#81838a',
        white: '#ffffff',
        whiteLilac: '#e4e1e5',
      },
      primary: {
        blue: '#2883fb',
        boatBlue: {
          DEFAULT: '#2b5283',
          dark: '#1e3a5d',
          light: '#386aa9',
        },
        cobalt: {
          DEFAULT: '#313283',
          dark: '#23245e',
          light: '#3f40a8',
        },
        coral: '#ed6765',
        green: {
          DEFAULT: '#277B3E',
          dark: '#1b542a',
          light: '#33a252',
        },
        pink: '#e74577',
        purple: {
          DEFAULT: '#8a35b7',
          hashita: '#8a5d8a',
        },
        red: {
          DEFAULT: '#ae0002',
          dark: '#7b0001',
          light: '#e10003',
        },
        valencia: '#d24343',
        vampire: '#b51313',
      },
      gray: {
        light: '#f4f1f5',
      },
    },
    fontFamily: {
      sans: [
        'ui-sans-serif',
        'system-ui',
        '-apple-system',
        'BlinkMacSystemFont',
        '"Segoe UI"',
        'Roboto',
        '"Helvetica Neue"',
        'Arial',
        '"Noto Sans"',
        'sans-serif',
        '"Apple Color Emoji"',
        '"Segoe UI Emoji"',
        '"Segoe UI Symbol"',
        '"Noto Color Emoji"',
      ],
    },
    fontSize: {
      xs: ['0.75rem', { lineHeight: '1.125rem' }],
      sm: ['0.875rem', { lineHeight: '1.3125rem' }],
      base: ['1rem', { lineHeight: '1.5rem' }],
      lg: ['1.125rem', { lineHeight: '1.6875rem' }],
      xl: ['1.25rem', { lineHeight: '1.875rem' }],
      '2xl': ['1.375rem', { lineHeight: '2.0625rem' }],
      '3xl': ['1.5rem', { lineHeight: '2.25rem' }],
      '4xl': ['1.625rem', { lineHeight: '2.4375rem' }],
      '5xl': ['1.75rem', { lineHeight: '2.625rem' }],
      '6xl': ['1.875rem', { lineHeight: '2.8125rem' }],
      '7xl': ['2rem', { lineHeight: '3rem' }],
      '8xl': ['2.125rem', { lineHeight: '3rem' }],
      '9xl': ['2.25rem', { lineHeight: '3.1875rem' }],
    },
    fontWeight: {
      normal: '400',
      bold: '700',
    },
    gradientColorStops: (theme) => ({
      ...theme('colors'),
      apricot: '#ec857c',
      blueViolet: '#5d5cbd',
      cinnabar: {
        DEFAULT: '#e84c4c',
        light: '#ff6666',
      },
      copperfield: '#df7f75',
      curiousBlue: '#2f97d9',
      darkTan: '#670f0f',
      ebony: '#100716',
      eastSide: '#c082cc',
      grape: '#49205e',
      haiti: '#261131',
      harold: '#a095a5',
      hummer: '#9e95a3',
      hibiscus: {
        DEFAULT: '#cd3a67',
        dark: '#cc3566',
      },
      iron: '#e8e8e8',
      jagger: '#360b4c',
      lavender: '#b964d3',
      moonRaker: '#eec6f6',
      vividViolet: {
        DEFAULT: '#90399e',
        dark: '#8631a0',
      },
    }),
    gradients: (theme) => ({
      boatBlue: `linear-gradient(-225deg, ${theme('colors.primary.boatBlue.dark')} 0%, ${theme(
        'colors.primary.boatBlue.light'
      )} 100%)`,
      cinnabar: `linear-gradient(${theme('gradientColorStops.cinnabar.DEFAULT')} 0%, ${theme(
        'gradientColorStops.darkTan'
      )} 100%)`,
      cinnabarLight: `linear-gradient(${theme('gradientColorStops.cinnabar.light')} 0%, ${theme(
        'gradientColorStops.darkTan'
      )} 100%)`,
      cobalt: `linear-gradient(-225deg, ${theme('colors.primary.cobalt.dark')} 0%, ${theme(
        'colors.primary.cobalt.light'
      )} 100%)`,
      dark: `linear-gradient(-135deg, ${theme('gradientColorStops.haiti')} 0%, ${theme(
        'gradientColorStops.ebony'
      )} 100%)`,
      full: `linear-gradient(153deg, ${theme('gradientColorStops.apricot')} 0%, ${theme(
        'gradientColorStops.hibiscus.dark'
      )} 26%, ${theme('gradientColorStops.copperfield')} 47%, ${theme(
        'gradientColorStops.vividViolet.DEFAULT'
      )} 66%, ${theme('gradientColorStops.blueViolet')} 83%, ${theme(
        'gradientColorStops.curiousBlue'
      )} 100%)`,
      green: `linear-gradient(-225deg, ${theme('colors.primary.green.dark')} 0%, ${theme(
        'colors.primary.green.light'
      )} 100%)`,
      haroldLight: `linear-gradient(-135deg, ${theme('gradientColorStops.harold')} 0%, ${theme(
        'gradientColorStops.hummer'
      )} 100%)`,
      iron: `linear-gradient(${theme('gradientColorStops.iron')} 0%, ${theme(
        'gradientColorStops.eastSide'
      )} 100%)`,
      lavender: `linear-gradient(${theme('gradientColorStops.lavender')} 0%, ${theme(
        'gradientColorStops.jagger'
      )} 100%)`,
      moonRaker: `linear-gradient(${theme('gradientColorStops.neutrals.white')} 0%, ${theme(
        'gradientColorStops.moonRaker'
      )} 100%)`,
      pink: `linear-gradient(-135deg, ${theme('gradientColorStops.hibiscus.DEFAULT')} 0%, ${theme(
        'gradientColorStops.grape'
      )} 100%)`,
      purple: `linear-gradient(${theme('gradientColorStops.vividViolet.dark')} 0%, ${theme(
        'gradientColorStops.jagger'
      )} 100%)`,
      red: `linear-gradient(-225deg, ${theme('colors.primary.red.dark')} 0%, ${theme(
        'colors.primary.red.light'
      )} 100%)`,
    }),
    spacing: {
      px: '1px',
      0: '0px',
      1: '0.25rem',
      2: '0.5rem',
      3: '0.75rem',
      4: '1rem',
      5: '1.25rem',
      6: '1.5rem',
      8: '2rem',
      10: '2.5rem',
      11: '2.75rem',
      12: '3rem',
      13: '3.25rem',
      14: '3.5rem',
      15: '3.75rem',
      16: '4rem',
      18: '4.5rem',
      20: '5rem',
      22: '5.5rem',
      24: '6rem',
      32: '8rem',
      40: '10rem',
      48: '12rem',
      56: '14rem',
      64: '16rem',
      80: '20rem',
    },
    extend: {
      backgroundImage: (theme) => ({
        'gradient-boatBlue': theme('gradients.boatBlue'),
        'gradient-cobalt': theme('gradients.cobalt'),
        'gradient-dark': theme('gradients.dark'),
        'gradient-harold': theme('gradients.haroldLight'),
        'gradient-full': theme('gradients.full'),
        'gradient-green': theme('gradients.green'),
        'gradient-pink': theme('gradients.pink'),
        'gradient-purple': theme('gradients.purple'),
        'gradient-red': theme('gradients.red'),
        onboarding: "url('~@/assets/img/onboarding-bg-lg.svg')",
        'onboarding-sm': "url('~@/assets/img/onboarding-bg-sm.svg')",
        'onboarding-flare-lg': "url('~@/assets/img/onboarding-flare-lg.png')",
      }),
      boxShadow: (theme) => ({
        'inner-outline-blue': `inset 0px 2px 0px 0px ${theme(
          'colors.primary.blue'
        )}, inset 0px -2px 0px 0px ${theme('colors.primary.blue')}`,
        'main-container': '0px 0px 40px 0px rgba(25, 12, 33, 0.1)',
      }),
      height: {
        lg: '46px',
        xl: '468px',
      },
      inset: {
        '-1': '-1px',
      },
      minWidth: {
        80: '20rem',
      },
      width: {
        lg: '550px',
        xl: '896px',
      },
      screens: {
        xs: '350px',
      },
    },
  },
  variants: {
    extend: {
      backgroundImage: ['focus-within', 'hover', 'disabled'],
      borderRadius: ['focus'],
      gradientColorStops: ['focus-within'],
    },
  },
  plugins: [
    plugin(function ({ addBase, theme }) {
      addBase({
        h1: { fontSize: theme('fontSize.9xl'), fontWeight: theme('fontWeight.bold'), margin: 0 },
        h2: { fontSize: theme('fontSize.8xl'), fontWeight: theme('fontWeight.bold'), margin: 0 },
        h3: { fontSize: theme('fontSize.7xl'), fontWeight: theme('fontWeight.bold'), margin: 0 },
        h4: { fontSize: theme('fontSize.6xl'), fontWeight: theme('fontWeight.bold'), margin: 0 },
        h5: { fontSize: theme('fontSize.5xl'), fontWeight: theme('fontWeight.bold'), margin: 0 },
        h6: { fontSize: theme('fontSize.4xl'), fontWeight: theme('fontWeight.bold'), margin: 0 },
      });
    }),
    plugin(function ({ addComponents, theme }) {
      const buttons = {
        '.btn-primary': {
          minHeight: theme('spacing.11'),
          padding: `0 ${theme('spacing.8')}`,
          borderRadius: theme('borderRadius.lg'),
          backgroundImage: theme('gradients.purple'),
          boxShadow: `0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          color: theme('colors.neutrals.white'),
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
          transition: 'all 200ms cubic-bezier(0.4, 0, 1, 1)',
          '&:hover': {
            backgroundImage: theme('gradients.lavender'),
          },
          '&:focus': {
            boxShadow: `0px 0px 0px 2px ${theme(
              'colors.neutrals.white'
            )}, 0px 0px 0px 4px rgba(138, 53, 183, 0.7), 0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          },
          '&:disabled': {
            backgroundImage: 'none',
            backgroundColor: theme('colors.neutrals.mobster'),
            color: theme('colors.neutrals.white'),
            cursor: 'not-allowed',
          },
        },
        '.btn-inverse': {
          minHeight: theme('spacing.11'),
          padding: `0 ${theme('spacing.8')}`,
          borderRadius: theme('borderRadius.lg'),
          backgroundImage: theme('gradients.moonRaker'),
          boxShadow: `0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          color: theme('colors.neutrals.dark'),
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
          transition: 'all 200ms cubic-bezier(0.4, 0, 1, 1)',
          '&:hover': {
            backgroundImage: theme('gradients.iron'),
          },
          '&:focus': {
            border: `2px solid ${theme('colors.neutrals.dark')}`, // TODO: replace with border from theme once defined
            boxShadow: `0px 0px 0px 2px rgba(227, 173, 255, 0.7), 0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          },
        },
        '.btn-outline': {
          minHeight: theme('spacing.11'),
          padding: `0 ${theme('spacing.8')}`,
          border: `1px solid ${theme('colors.neutrals.chatelle')}`, // TODO: replace with border from theme once defined
          borderRadius: theme('borderRadius.lg'),
          backgroundColor: theme('colors.neutrals.white'),
          color: theme('colors.neutrals.medium'),
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
          transition: 'all 200ms cubic-bezier(0.4, 0, 1, 1)',
          '&:hover': {
            border: `1px solid ${theme('colors.neutrals.mountainMist.light')}`,
          },
          '&:focus': {
            boxShadow: `0px 0px 0px 2px rgba(138, 53, 183, 0.7)`, // TODO: replace with shadow from theme once defined
          },
        },
        '.btn-danger': {
          minHeight: theme('spacing.11'),
          padding: `0 ${theme('spacing.8')}`,
          borderRadius: theme('borderRadius.lg'),
          backgroundImage: theme('gradients.cinnabar'),
          boxShadow: `0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          color: theme('colors.neutrals.white'),
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
          transition: 'all 200ms cubic-bezier(0.4, 0, 1, 1)',
          '&:hover': {
            backgroundImage: theme('gradients.cinnabarLight'),
          },
          '&:focus': {
            border: `2px solid ${theme('colors.neutrals.white')}`, // TODO: replace with border from theme once defined
            boxShadow: `0px 0px 0px 2px rgba(174, 49, 49, 0.7), 0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          },
        },
        '.btn-disabled': {
          minHeight: theme('spacing.11'),
          padding: `0 ${theme('spacing.8')}`,
          borderRadius: theme('borderRadius.lg'),
          backgroundColor: theme('colors.neutrals.mobster'),
          boxShadow: `0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          color: theme('colors.neutrals.white'),
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
          cursor: 'not-allowed',
        },
        '.btn-loading': {
          minHeight: theme('spacing.11'),
          padding: `0 ${theme('spacing.8')}`,
          borderRadius: theme('borderRadius.lg'),
          backgroundImage: theme('gradients.purple'),
          boxShadow: `0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          color: theme('colors.neutrals.white'),
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
        },
      };

      const inputs = {
        '.input-container': {
          position: 'relative',
          marginBottom: theme('spacing.6'),
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'flex-start',
          alignItems: 'flex-start',
          width: '100%',

          '& input': {
            height: theme('spacing.15'),
            padding: `26px ${theme('spacing.16')} 10px ${theme('spacing.4')}`,
            backgroundColor: theme('colors.neutrals.bonjour'),
            width: '100%',
            border: 'none',
            borderBottom: `2px solid ${theme('colors.neutrals.mountainMist.light')}`,
            borderRadius: `6px 6px 0px 0px`,
            color: theme('colors.neutrals.dark'),
            outline: 'none',
            marginBottom: 5,

            '&:focus': {
              borderColor: theme('colors.primary.purple.DEFAULT'),
            },

            '&:invalid:not(:placeholder-shown)': {
              borderColor: theme('colors.primary.valencia'),
            },
          },

          '& input::placeholder': {
            visibility: 'hidden',
          },

          '& input:focus::placeholder': {
            visibility: 'visible',
            color: theme('colors.neutrals.medium'),
            fontSize: theme('fontSize.base'),
          },

          '& input:focus + .input-label, input:not(:focus):not(:placeholder-shown) + .input-label':
            {
              top: 5,
              fontSize: theme('fontSize.sm'),
              marginBottom: theme('spacing.8'),
            },

          '& .input-label': {
            fontSize: theme('fontSize.base'),
            fontWeight: theme('fontWeight.bold'),
            color: theme('colors.neutrals.dark'),
            position: 'absolute',
            paddingLeft: theme('spacing.4'),
            transition: `top .2s`,
            top: 17,
          },

          '& .input-word-limit': {
            position: 'absolute',
            height: theme('spacing.14'),
            right: theme('spacing.4'),
            color: theme('colors.neutrals.medium'),
            fontSize: theme('fontSize.sm'),
            top: 1,
            display: 'flex',
            justifyContent: 'flex-end',
            alignItems: 'center',
          },

          '& .input-helper': {
            color: theme('colors.neutrals.medium'),
            fontSize: theme('fontSize.sm'),
            display: 'flex',
            alignItems: 'center',
          },

          '& .fader': {
            position: 'absolute',
            top: 1,
            right: theme('spacing.16'),
            background: `linear-gradient(90deg, rgba(238, 234, 238, 0) 0%, rgb(238, 234, 238) 100%)`,
            borderRadius: `0px 4px 0px 0px`,
            height: theme('spacing.14'),
            width: 68,
            display: 'flex',
            justifyContent: 'flex-end',
            alignItems: 'center',
          },

          '& .danger-icon, .checkmark-icon': {
            position: 'absolute',
            top: 18,
            right: theme('spacing.3'),
          },
        },
      };

      const links = {
        '.underline-blue': {
          borderBottom: `2px solid ${theme('colors.primary.blue')}`,
          paddingBottom: 2,
        },
        '.underline-white': {
          paddingBottom: 2,
          '&:not(:focus):not(:focus-within):hover': {
            borderBottom: `1px solid ${theme('colors.neutrals.white')}`,
            paddingBottom: 1,
          },
        },
      };

      const toastNotifications = {
        '.error-notification': {
          backgroundColor: theme('colors.primary.valencia'),
          top: theme('spacing.8'),
          position: 'absolute',
          borderRadius: theme('spacing.3'),
          zIndex: 10,
          maxWidth: '48rem',
          display: 'flex',
          verticalAlign: 'middle',
          flexDirection: 'row',
          justifyContent: 'flex-start',
          alignItems: 'center',
          padding: `${theme('spacing.4')} ${theme('spacing.5')}`,
          color: theme('colors.neutrals.softWhite'),
        },
      };

      const credentialPreview = {
        '.nocredentialCard': {
          width: 342,
          height: 152,
          [`@media (min-width: ${theme('screens.md')})`]: {
            width: 384,
          },
        },
      };

      addComponents({
        ...buttons,
        ...inputs,
        ...links,
        ...toastNotifications,
        ...credentialPreview,
      });
    }),
  ],
};
