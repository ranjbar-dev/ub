import React, { useState, memo } from 'react';
import { IconButton } from '@material-ui/core';
import night from 'images/Ndark.svg';
import light from 'images/NLight.svg';
import { Themes } from 'containers/App/constants';
import { LocalStorageKeys } from 'services/constants';
import { MessageNames, MessageService } from 'services/message_service';
import Anime from 'react-anime';
const ThemeToggler = () => {
  const [Theme, setTheme] = useState(
    localStorage[LocalStorageKeys.Theme] &&
      localStorage[LocalStorageKeys.Theme] != ''
      ? localStorage[LocalStorageKeys.Theme]
      : Themes.DARK,
  );

  return (
    <IconButton
      onClick={() => {
        Theme === Themes.LIGHT ? setTheme(Themes.DARK) : setTheme(Themes.LIGHT);
        const app = document.querySelector('body');
        const hTml = document.querySelector('html');
        if (app) {
          app.classList.toggle('darkTheme');
          hTml?.classList.toggle('htmldark');
          localStorage[LocalStorageKeys.Theme] = app.classList.contains(
            Themes.DARK,
          )
            ? Themes.DARK
            : Themes.LIGHT;
          MessageService.send({
            name: MessageNames.CHANGE_THEME,
            payload: Theme === Themes.LIGHT ? Themes.DARK : Themes.LIGHT,
          });
        }
      }}
      className='headerButton'
      size='small'
    >
      <Anime duration={1000} rotate={[0, 360]}>
        <img
          style={{ maxWidth: '20px', width: '20px' }}
          src={Theme === Themes.LIGHT ? light : night}
          alt=''
        />{' '}
      </Anime>
    </IconButton>
  );
};
export default memo(ThemeToggler);
