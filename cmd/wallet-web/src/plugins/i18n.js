import { nextTick } from 'vue';
import { createI18n } from 'vue-i18n';
import { setDocumentLang, setDocumentTitle } from '@/mixins/i18n/document';

const loadedLanguages = [];

const i18n = createI18n({
  globalInjection: true,
  legacy: false,
  locale: process.env.VUE_APP_I18N_FALLBACK_LOCALE || 'en',
  fallbackLocale: process.env.VUE_APP_I18N_FALLBACK_LOCALE || 'en',
  messages: {},
});

// This function updates i18n locale, loads new locale's messages and sets document properties accordingly
export async function updateI18nLocale(locale) {
  if (loadedLanguages.length > 0 && i18n.locale === locale) {
    return;
  }

  // If the language was already loaded
  if (loadedLanguages.includes(locale)) {
    if (i18n.mode === 'legacy') {
      i18n.global.locale = locale;
    } else {
      i18n.global.locale.value = locale;
    }
    setDocumentLang(locale);
    setDocumentTitle(i18n.global.t('App.title'));
    return;
  }

  // If the language hasn't been loaded yet
  const messages = await import(
    /* webpackChunkName: "locale-[request]" */ `@/translations/${locale}.js`
  );
  i18n.global.setLocaleMessage(locale, messages.default);
  setDocumentLang(locale);
  setDocumentTitle(i18n.global.t('App.title'));
  if (i18n.mode === 'legacy') {
    i18n.global.locale = locale;
  } else {
    i18n.global.locale.value = locale;
  }
  loadedLanguages.push(locale);
  return nextTick();
}

export default i18n;
