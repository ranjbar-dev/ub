import { render, fireEvent } from '@testing-library/react';
import i18next from 'i18next';
import { i18n, translations } from 'locales/i18n';
import React from 'react';
import { ThemeProvider, DefaultTheme } from 'styled-components';
import { themes } from 'styles/theme/themes';

import { LanguageSwitch } from '..';


const renderLanguageSwitch = (theme?: DefaultTheme) =>
  render(
    <ThemeProvider theme={theme || themes.light}>
      <LanguageSwitch />
    </ThemeProvider>,
  );
describe('<LanguageSwitch />', () => {
  it('should have 2 radio buttons', () => {
    const languageSwitch = renderLanguageSwitch();
    expect(languageSwitch.queryAllByRole('radio', { hidden: true }).length).toBe(2);
  });

  it('should have translated FormLabel ', async () => {
    const t = await i18n;

    i18next.changeLanguage('en');

    let languageSwitch = renderLanguageSwitch();
    let label = languageSwitch.queryByText(
      t(translations.i18nFeature.selectLanguage()),
    );
    expect(label).toBeInTheDocument();

    languageSwitch.unmount();
    i18next.changeLanguage('de');

    languageSwitch = renderLanguageSwitch();
    label = languageSwitch.queryByText(
      t(translations.i18nFeature.selectLanguage()),
    );
    expect(label).toBeInTheDocument();
  });

  it('should change language on click', async () => {
    i18next.changeLanguage('en');
    const languageSwitch = renderLanguageSwitch();
    const radio2 = languageSwitch.queryAllByRole('radio', { hidden: true })[1];
    fireEvent.click(radio2);
    expect(i18next.language).toEqual('de');
  });
});
