/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    colors: {
      neutrals: {
        dark: '#190c21',
        light: '#b6b7c7',
        medium: '#6c6d7c',
        softWhite: '#f4f1f5',
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
    gradientColorStops: {
      apricot: '#ec857c',
      blueViolet: '#5d5cbd',
      copperfield: '#df7f75',
      curiousBlue: '#2f97d9',
      ebony: '#100716',
      grape: '#49205e',
      haiti: '#261131',
      hibiscus: {
        DEFAULT: '#cd3a67',
        dark: '#cc3566',
      },
      jagger: '#360b4c',
      vividViolet: {
        DEFAULT: '#90399e',
        dark: '#8631a0',
      }
    },
    gradients: theme => ({
      dark: `linear-gradient(-135deg, ${theme('gradientColorStops.haiti')} 0%, ${theme('gradientColorStops.ebony')} 100%)`,
      full: `linear-gradient(-225deg, ${theme('gradientColorStops.apricot')} 0%, ${theme('gradientColorStops.hibiscus.dark')} 26%, ${theme('gradientColorStops.copperfield')} 47%, ${theme('gradientColorStops.vividViolet.DEFAULT')} 66%, ${theme('gradientColorStops.blueViolet')} 83%, ${theme('gradientColorStops.curiousBlue')} 100%)`,
      pink: `linear-gradient(-135deg, ${theme('gradientColorStops.hibiscus.DEFAULT')} 0%, ${theme('gradientColorStops.grape')} 100%)`,
      purple: `linear-gradient(-135deg, ${theme('gradientColorStops.vividViolet.dark')} 0%, ${theme('gradientColorStops.jagger')} 100%)`,
    }),
    extend: {
      backgroundImage: theme => ({
        'gradient-dark': theme('gradients.dark'),
        'gradient-full': theme('gradients.full'),
        'gradient-pink': theme('gradients.pink'),
        'gradient-purple': theme('gradients.purple'),
        'onboarding': "url('~@/assets/img/onboarding-bg-lg.svg')",
        'flare': "url('~@/assets/img/onboarding-flare-lg.png')",
      }),
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
