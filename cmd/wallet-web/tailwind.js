/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const plugin = require('tailwindcss/plugin');

module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    colors: {
      neutrals: {
        black: '#000000',
        bonjour: '#eeeaee',
        chatelle: '#c7c3c8',
        dark: '#190c21',
        light: '#b6b7c7',
        medium: '#6c6d7c',
        mischka: '#dcdadd',
        mobster: '#7b6f7f',
        mountainMist: {
          light: '#949095',
          dark: '#8d8a8e',
        },
        softWhite: '#f4f1f5',
        selago: '#f4f1f5',
        thistle: '#dbd7dc',
        white: '#ffffff',
        whiteLilac: '#e4e1e5',
      },
      primary: {
        blue: '#2883fb',
        coral: '#ed6765',
        pink: '#e74577',
        purple: {
          DEFAULT: '#8a35b7',
          hashita: '#8a5d8a',
        },
        valencia: '#d24343',
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
      cinnabar: `linear-gradient(-180deg, ${theme(
        'gradientColorStops.cinnabar.DEFAULT'
      )} 0%, ${theme('gradientColorStops.darkTan')} 100%)`,
      cinnabarLight: `linear-gradient(-180deg, ${theme(
        'gradientColorStops.cinnabar.light'
      )} 0%, ${theme('gradientColorStops.darkTan')} 100%)`,
      dark: `linear-gradient(-135deg, ${theme('gradientColorStops.haiti')} 0%, ${theme(
        'gradientColorStops.ebony'
      )} 100%)`,
      full: `linear-gradient(-225deg, ${theme('gradientColorStops.apricot')} 0%, ${theme(
        'gradientColorStops.hibiscus.dark'
      )} 26%, ${theme('gradientColorStops.copperfield')} 47%, ${theme(
        'gradientColorStops.vividViolet.DEFAULT'
      )} 66%, ${theme('gradientColorStops.blueViolet')} 83%, ${theme(
        'gradientColorStops.curiousBlue'
      )} 100%)`,
      iron: `linear-gradient(-180deg, ${theme('gradientColorStops.iron')} 0%, ${theme(
        'gradientColorStops.eastSide'
      )} 100%)`,
      lavender: `linear-gradient(-180deg, ${theme('gradientColorStops.lavender')} 0%, ${theme(
        'gradientColorStops.jagger'
      )} 100%)`,
      pink: `linear-gradient(-135deg, ${theme('gradientColorStops.hibiscus.DEFAULT')} 0%, ${theme(
        'gradientColorStops.grape'
      )} 100%)`,
      purple: `linear-gradient(-135deg, ${theme('gradientColorStops.vividViolet.dark')} 0%, ${theme(
        'gradientColorStops.jagger'
      )} 100%)`,
      moonRaker: `linear-gradient(-180deg, ${theme(
        'gradientColorStops.neutrals.white'
      )} 0%, ${theme('gradientColorStops.moonRaker')} 100%)`,
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
        'gradient-dark': theme('gradients.dark'),
        'gradient-full': theme('gradients.full'),
        'gradient-pink': theme('gradients.pink'),
        'gradient-purple': theme('gradients.purple'),
        onboarding: "url('~@/assets/img/onboarding-bg-lg.svg')",
        'onboarding-sm': "url('~@/assets/img/onboarding-bg-sm.svg')",
        'onboarding-flare-lg': "url('~@/assets/img/onboarding-flare-lg.png')",
        'header-mobile-flare': "url('~@/assets/img/header-mobile-flare.png')",
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
      backgroundImage: ['focus-within', 'hover'],
      borderRadius: ['focus'],
      gradientColorStops: ['focus-within'],
    },
  },
  plugins: [
    plugin(function ({ addBase, theme }) {
      addBase({
        h1: { fontSize: theme('fontSize.9xl'), fontWeight: theme('fontWeight.bold') },
        h2: { fontSize: theme('fontSize.8xl'), fontWeight: theme('fontWeight.bold') },
        h3: { fontSize: theme('fontSize.7xl'), fontWeight: theme('fontWeight.bold') },
        h4: { fontSize: theme('fontSize.6xl'), fontWeight: theme('fontWeight.bold') },
        h5: { fontSize: theme('fontSize.5xl'), fontWeight: theme('fontWeight.bold') },
        h6: { fontSize: theme('fontSize.4xl'), fontWeight: theme('fontWeight.bold') },
      });
    }),
    plugin(function ({ addComponents, theme }) {
      const buttons = {
        '.btn-primary': {
          height: theme('spacing.11'),
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
            border: `2px solid ${theme('colors.neutrals.white')}`, // TODO: replace with border from theme once defined
            boxShadow: `0px 0px 0px 2px rgba(138, 53, 183, 0.7), 0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          },
        },
        '.btn-inverse': {
          height: theme('spacing.11'),
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
          height: theme('spacing.11'),
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
          height: theme('spacing.11'),
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
          height: theme('spacing.11'),
          padding: `0 ${theme('spacing.8')}`,
          borderRadius: theme('borderRadius.lg'),
          backgroundColor: theme('colors.neutrals.mobster'),
          boxShadow: `0px 1px 2px 0px rgba(0, 0, 0, 0.1)`, // TODO: replace with shadow from theme once defined
          color: theme('colors.neutrals.white'),
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
        },
        '.btn-loading': {
          height: theme('spacing.11'),
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

          '& input': {
            height: theme('spacing.15'),
            width: 342,
            padding: `26px ${theme('spacing.16')} 10px ${theme('spacing.4')}`,
            backgroundColor: theme('colors.neutrals.bonjour'),
            border: 'none',
            borderBottom: `2px solid ${theme('colors.neutrals.mountainMist.light')}`,
            borderRadius: `6px 6px 0px 0px`,
            color: theme('colors.neutrals.dark'),

            '&:focus': {
              borderColor: theme('colors.primary.purple.DEFAULT'),
            },

            '&:invalid': {
              borderColor: theme('colors.primary.valencia'),
              '& ~ img': {
                visibility: 'visible',
              },
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

          '& .input-helper': {
            color: theme('colors.neutrals.medium'),
            fontSize: theme('fontSize.sm'),
            display: 'flex',
            alignItems: 'center',
            marginTop: 5,
          },

          '& .fader': {
            position: `absolute`,
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

          '& img': {
            visibility: 'hidden',
            position: `absolute`,
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

      const skeletons = {
        '.skeleton': {
          backgroundColor: theme('colors.neutrals.whiteLilac'),
          borderRadius: theme('spacing.3'),
          display: 'flex',
          paddingLeft: theme('spacing.4'),
          paddingTop: theme('spacing.4'),
          paddingRight: theme('spacing.14'),
          paddingBottom: theme('spacing.4'),
          width: 342,
          height: 80,
          flexDirection: 'row',
          [`@media (min-width: ${theme('screens.md')})`]: {
            paddingLeft: theme('spacing.6'),
            paddingTop: theme('spacing.5'),
            paddingRight: theme('spacing.6'),
            paddingBottom: theme('spacing.6'),
            width: 256,
            height: 160,
            flexDirection: 'column',
          },
        },
        '.skeleton-photo': {
          backgroundColor: theme('colors.neutrals.selago'),
          borderRadius: '50%',
          width: theme('spacing.12'),
          height: theme('spacing.12'),
          flexFlow: 'col-1',
        },
        '.skeleton-data': {
          backgroundColor: theme('colors.neutrals.selago'),
          position: 'relative',
          borderRadius: theme('spacing.6'),
          overflow: 'hidden',
          marginTop: theme('spacing.0'),
          marginLeft: theme('spacing.3'),
          [`@media (min-width: ${theme('screens.md')})`]: {
            marginTop: theme('spacing.3'),
            marginLeft: theme('spacing.0'),
          },
          '&::before, &::after': {
            content: '""',
            position: 'absolute',
            left: 0,
            width: 80,
            height: 1,
            backgroundImage: `linear-gradient(-90deg, rgba(255, 255, 255, 0) 15%, #867C8C 50%, rgba(255, 255, 255, 0) 85%)`,
            animation: 'shimmer 1.25s linear infinite',
          },
          '&::before': {
            top: 0,
          },

          '&::after': {
            bottom: 0,
          },

          '& .top': {
            width: '100%',
          },

          '& .bottom': {
            width: '60%',
          },
        },
        '@keyframes shimmer': {
          from: {
            transform: 'translateX(-100%)',
          },
          to: {
            transform: 'translateX(250%)',
          },
        },
      };
      const credentialPreview = {
        '.credentialCard': {
          backgroundColor: theme('colors.neutrals.white'),
          borderRadius: theme('spacing.3'),
          padding: `${theme('spacing.6')} ${theme('spacing.3')} ${theme('spacing.6')} ${theme(
            'spacing.4'
          )}`,
          fontSize: theme('fontSize.base'),
          fontWeight: theme('fontWeight.bold'),
          border: `${theme('spacing.px')} solid ${theme('colors.neutrals.thistle')}`,
          width: 400,
          height: theme('spacing.24'),
        },
        '.credentialHeader': {
          paddingLeft: theme('spacing.4'),
          color: theme('colors.neutrals.dark'),
        },
        '.credentialLogoContainer': {
          width: theme('spacing.12'),
          height: theme('spacing.12'),
        },
        '.credentialArrowContainer': {
          backgroundColor: theme('colors.neutrals.thistle'),
          border: `${theme('spacing.px')} solid ${theme('colors.neutrals.thistle')}`,
          borderRadius: '50%',
          width: theme('spacing.8'),
          height: theme('spacing.8'),
        },
        '.credentialArrowLogo': {
          padding: theme('spacing.1'),
        },
      };

      addComponents({
        ...buttons,
        ...inputs,
        ...links,
        ...toastNotifications,
        ...skeletons,
        ...credentialPreview,
      });
    }),
  ],
};
