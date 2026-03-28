import React, { useState } from 'react';
import { IconButton } from '@material-ui/core';
import FilledStar from 'images/themedIcons/filledStar';
import EmptyStar from 'images/themedIcons/emptyStar';
import Anime from 'react-anime';
import { LocalStorageKeys } from 'services/constants';

export default function FavIcon(props: { onFavChange: Function }) {
  const [FavActive, setFavActive] = useState(false);
  const handleFavClick = () => {
    props.onFavChange(!FavActive);
    setFavActive(!FavActive);
    localStorage[LocalStorageKeys.SHOW_FAVS] = !FavActive.toString();
  };
  return (
    <>
      <Anime
        duration={200}
        easing="easeOutCirc"
        scale={[0.1, 1]}
        opacity={[0, 1]}
      >
        <IconButton disableRipple onClick={handleFavClick}>
          {FavActive === true ? (
            <FilledStar size="20" />
          ) : (
            <EmptyStar size="20" />
          )}
        </IconButton>
      </Anime>
    </>
  );
}
