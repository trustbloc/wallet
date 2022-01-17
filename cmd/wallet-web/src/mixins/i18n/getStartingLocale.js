import supportedLocales from '@/config/supportedLocales';
import getBrowserLocale from '@/mixins/i18n/getBrowserLocale';
import router from '@/router/index';

// Returns locale configuration. By default, try VUE_APP_I18N_LOCALE. As fallback, use 'en'.
function getMappedLocale(locale = process.env.VUE_APP_I18N_LOCALE || 'en') {
  return supportedLocales.find((loc) => loc.id === locale);
}

export default function getStartingLocale() {
  // Get locale parameter form the URL
  const localeUrlParam = router.options.history.location
    .replaceAll(/^\//gi, '')
    .replace(/\/.*$/gi, '');
  // If locale parameter is set, check if it is amongst the supported locales and return it.
  if (localeUrlParam && supportedLocales.find((loc) => loc.id === localeUrlParam)) {
    return getMappedLocale(localeUrlParam);
  }
  // If no locale parameter is set in the URL, use the browser default.
  else {
    const browserLocale = getBrowserLocale({ countryCodeOnly: true });
    return getMappedLocale(browserLocale);
  }
}
