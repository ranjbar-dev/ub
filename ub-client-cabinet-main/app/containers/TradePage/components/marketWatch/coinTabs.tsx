import React, { useState, useEffect } from 'react';
import styled from 'styles/styled-components';
import { Tabs, Tab } from '@material-ui/core';
import {
  MarketWatchSubscriber,
} from 'services/message_service';

import EmptyStar from 'images/themedIcons/emptyStar';
import FilledStar from 'images/themedIcons/filledStar';
import { LocalStorageKeys } from 'services/constants';
//   <FavIcon onFavChange={handleFavChange} />
let tabs: any[] = ['Favs', 'All'];
export default function CoinTabs (props: {
  onTabChange: Function;
  //  onFavChange: Function;
  subject: string;
}) {
  const [activeIndex, setactiveIndex] = React.useState(
    localStorage[LocalStorageKeys.FAV_COIN] == 'Favs' ? 0 : 1,
  );
  const [CTabs, setCTabs]: [any, any] = useState(['Favs', 'All']);
  const handleChange = (event, newactiveIndex) => {
    localStorage[LocalStorageKeys.FAV_COIN] = CTabs[newactiveIndex];
    setactiveIndex(newactiveIndex);
    props.onTabChange(CTabs[newactiveIndex]);
  };

  useEffect(() => {
    const MarketWatchSubscription = MarketWatchSubscriber.subscribe(
      (message: any) => {
        const payload: any = message.payload;
        const code: string = payload.name.split('-')[1];
        if (!tabs.includes(code)) {
          tabs.push(code);
          setCTabs([...tabs]);
          if (localStorage[LocalStorageKeys.FAV_COIN] === code) {
            setactiveIndex(tabs.length - 1);
          }
        }
      },
    );
    return () => {
      tabs = ['Favs', 'All'];
      MarketWatchSubscription.unsubscribe();
    };
  }, []);
  //  const handleFavChange = (fav: boolean) => {
  //    props.onFavChange(fav);
  //  };
  return (
    <div>
      {CTabs.length > 0 ? (
        <Wrapper>
          {/*<div className="fav">
            {useMemo(
              () => (
                             <FavIcon onFavChange={handleFavChange} />
              ),
              [],
            )}
          </div>*/}
          <Tabs
            value={activeIndex}
            onChange={handleChange}
            indicatorColor='primary'
            textColor='primary'
          >
            {CTabs.map((item: string, index: number) => {
              return (
                <Tab
                  key={'coinT' + index}
                  disableRipple={true}
                  className='typeTab'
                  label={
                    <span>
                      {item === 'Favs' ? (
                        activeIndex === 0 ? (
                          <FilledStar size='20' />
                        ) : (
                          <EmptyStar size='20' />
                        )
                      ) : (
                        item
                      )}
                    </span>
                  }
                />
              );
            })}
          </Tabs>
        </Wrapper>
      ) : (
        <div></div>
      )}
    </div>
  );
}
const Wrapper = styled.div`
  .fav {
    margin: 0 6px;
    svg {
      width: 16px;
    }
  }
  display: flex;
  .MuiTabs-root {
    min-height: 25px !important;
    min-width: 100%;
  }
  --tabWidth: 40px;
  .MuiTabs-indicator {
    /*min-width: var(--tabWidth) !important;*/
    background-color: var(--textBlue) !important;
  }
  .typeTab {
    max-width: var(--tabWidth);
    min-width: var(--tabWidth);
    min-height: 22px;

    .MuiTab-wrapper {
      max-width: var(--tabWidth);
    }
    span {
      transition: color 0.3s;
      color: var(--textGrey) !important;
      font-weight: 600;
      font-size: 12px !important;
    }
    &.Mui-selected {
      span {
        color: var(--textBlue) !important;
      }
    }
  }
  border-bottom: 1px solid var(--lightGrey);
  button {
    background-color: var(--white) !important;
    padding: 0;
  }
  .MuiTab-root {
    min-width: unset !important;
    max-width: fit-content;
    padding: 0 5px;
  }
`;
