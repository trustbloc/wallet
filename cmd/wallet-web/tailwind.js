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
        chatelle: '#c7c3c8',
        dark: '#190c21',
        light: '#b6b7c7',
        medium: '#6c6d7c',
        mobster: '#7b6f7f',
        mountainMist: '#949095',
        softWhite: '#f4f1f5',
        white: '#ffffff',
      },
      primary: {
        blue: '#2883fb',
        coral: '#ed6765',
        pink: '#e74577',
        purple: '#8a35b7',
      },
      gray:{
        light: '#f4f1f5',
      }
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
      }
    }),
    gradients: theme => ({
      cinnabar: `linear-gradient(-180deg, ${theme('gradientColorStops.cinnabar.DEFAULT')} 0%, ${theme('gradientColorStops.darkTan')} 100%)`,
      cinnabarLight: `linear-gradient(-180deg, ${theme('gradientColorStops.cinnabar.light')} 0%, ${theme('gradientColorStops.darkTan')} 100%)`,
      dark: `linear-gradient(-135deg, ${theme('gradientColorStops.haiti')} 0%, ${theme('gradientColorStops.ebony')} 100%)`,
      full: `linear-gradient(-225deg, ${theme('gradientColorStops.apricot')} 0%, ${theme('gradientColorStops.hibiscus.dark')} 26%, ${theme('gradientColorStops.copperfield')} 47%, ${theme('gradientColorStops.vividViolet.DEFAULT')} 66%, ${theme('gradientColorStops.blueViolet')} 83%, ${theme('gradientColorStops.curiousBlue')} 100%)`,
      iron: `linear-gradient(-180deg, ${theme('gradientColorStops.iron')} 0%, ${theme('gradientColorStops.eastSide')} 100%)`,
      lavender: `linear-gradient(-180deg, ${theme('gradientColorStops.lavender')} 0%, ${theme('gradientColorStops.jagger')} 100%)`,
      pink: `linear-gradient(-135deg, ${theme('gradientColorStops.hibiscus.DEFAULT')} 0%, ${theme('gradientColorStops.grape')} 100%)`,
      purple: `linear-gradient(-135deg, ${theme('gradientColorStops.vividViolet.dark')} 0%, ${theme('gradientColorStops.jagger')} 100%)`,
      moonRaker: `linear-gradient(-180deg, ${theme('gradientColorStops.neutrals.white')} 0%, ${theme('gradientColorStops.moonRaker')} 100%)`,
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
      15: '3.75',
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
    },
    extend: {
      backgroundImage: theme => ({
        'gradient-dark': theme('gradients.dark'),
        'gradient-full': theme('gradients.full'),
        'gradient-pink': theme('gradients.pink'),
        'gradient-purple': theme('gradients.purple'),
        'onboarding': "url('~@/assets/img/onboarding-bg-lg.svg')",
        'flare': "url('~@/assets/img/onboarding-flare-lg.png')",
      }),
      height: {
        lg: '46px',
        xl: '468px',
      },
      width: {
        lg: '550px',
        xl: '896px',
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    plugin(function({ addBase, theme }) {
      addBase({
        'h1': { fontSize: theme('fontSize.9xl'), fontWeight: theme('fontWeight.bold') },
        'h2': { fontSize: theme('fontSize.8xl'), fontWeight: theme('fontWeight.bold') },
        'h3': { fontSize: theme('fontSize.7xl'), fontWeight: theme('fontWeight.bold') },
        'h4': { fontSize: theme('fontSize.6xl'), fontWeight: theme('fontWeight.bold') },
        'h5': { fontSize: theme('fontSize.5xl'), fontWeight: theme('fontWeight.bold') },
        'h6': { fontSize: theme('fontSize.4xl'), fontWeight: theme('fontWeight.bold') },
      })
    }),
    plugin(function({ addComponents, theme }) {
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
          }
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
          }
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
            border: `1px solid ${theme('colors.neutrals.mountainMist')}`,
          },
          '&:focus': {
            boxShadow: `0px 0px 0px 2px rgba(138, 53, 183, 0.7)`, // TODO: replace with shadow from theme once defined
          }
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
          }
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
        }
      }

      addComponents(buttons)
    }),
  ],
}
