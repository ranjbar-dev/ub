import React, { memo, useState } from 'react';
import styled from 'styles/styled-components';
import { Card, Button } from '@material-ui/core';
import translate from '../../messages';
import { FormattedMessage } from 'react-intl';
import { Buttons } from 'containers/App/constants';

import EyeIcon from 'images/themedIcons/eyeIcon';
import CrossedEyeIcon from 'images/themedIcons/crossedEyeIcon';
import { CurrencyFormater } from 'utils/formatters';
import { LocalStorageKeys } from 'services/constants';
export default memo(function TopInfo (props: {
  data: {
    etimatedBalance: string;
    availableBalance: string;
    inOrder: string;
    btcEtimatedBalance: string;
    btcAvailableBalance: string;
    btcInOrder: string;
  };
}) {
  const data = props.data;
  const [Visible, setVisible] = useState(
    localStorage[LocalStorageKeys.SHOW_TOP_INFO] === 'true' ? true : false,
  );
  const toggleVisible = () => {
    localStorage[LocalStorageKeys.SHOW_TOP_INFO] = (!Visible).toString();
    setVisible(!Visible);
  };
  const dataFormatter = (text: string) => {
    if (Visible === true) {
      return text;
    }
    return '**********';
    // return text
    //   .replace(/[0-9]/g, '*')
    //   .replace(/./g, '*')
    //   .replace(/,/g, '*');
  };
  const Info = (params: {
    title: any;
    data: string;
    btcEqFieldName: string;
  }) => {
    return (
      <div className='info'>
        <span className='title'>{params.title}</span>
        <span className='currency'>{` USD `}</span>
        <span className='data'>
          {dataFormatter(
            CurrencyFormater(
              Number(data[params.data])
                .toFixed(2)
                .toString(),
            ),
          )}
        </span>
        {' / '}
        <span className='currency'>{` BTC `}</span>
        <span className='data'>
          {dataFormatter(CurrencyFormater(data[params.btcEqFieldName]))}
        </span>
      </div>
    );
  };
  return (
    <Wrapper>
      {Info({
        title: <FormattedMessage {...translate.EstimatedBalance} />,
        data: 'etimatedBalance',
        btcEqFieldName: 'btcEtimatedBalance',
      })}
      {Info({
        title: <FormattedMessage {...translate.AvailableBalance} />,
        data: 'availableBalance',
        btcEqFieldName: 'btcAvailableBalance',
      })}
      {Info({
        title: <FormattedMessage {...translate.InOrders} />,
        data: 'inOrder',
        btcEqFieldName: 'btcInOrder',
      })}

      <div className='visibleButtonWrapper'>
        <Button
          onClick={toggleVisible}
          className={`${Buttons.SimpleRoundButton} button`}
        >
          {Visible === true ? <EyeIcon /> : <CrossedEyeIcon />}
        </Button>
      </div>
    </Wrapper>
  );
});
const Wrapper = styled(Card)`
  height: 48px;
  width: 100%;
  min-width: 1000px;
  border-radius: 10px !important;
  box-shadow: none !important;
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  .info {
    flex: 4;
    color: var(--textGrey);
    padding: 0 20px;
    .title {
      font-weight: 600;
    }
    .data {
      font-weight: 600;
      color: var(--blackText);
    }
  }
  .visibleButtonWrapper {
    display: flex;
    justify-content: flex-end;
    padding: 0 20px;
    .button {
      padding: 0px !important;
    }
    img {
      min-height: 30px;
    }
  }
  span {
    font-size: 13px;
  }
`;
