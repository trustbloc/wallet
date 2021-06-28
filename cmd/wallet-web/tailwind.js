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
       }),
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
