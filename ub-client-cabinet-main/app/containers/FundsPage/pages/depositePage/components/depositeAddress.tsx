import React, { useState, useEffect, useRef } from 'react';
import styled from 'styles/styled-components';
import { Buttons } from 'containers/App/constants';
import { FormattedMessage } from 'react-intl';
import { Currency } from 'containers/App/types';
import { LocalStorageKeys } from 'services/constants';
import translate from '../../../messages';

import { createStructuredSelector } from 'reselect';
import { makeSelectdepositAndWithDrawData } from 'containers/FundsPage/selectors';
import { useSelector, useDispatch } from 'react-redux';
import { Subscriber, MessageNames } from 'services/message_service';
import { getRawdepositAndWithDrawDataAction } from 'containers/FundsPage/actions';
import attention from 'images/attentionIcon.svg';

import { depositAndWithDrawData } from 'containers/FundsPage/types';

import { QRCode } from 'react-qrcode-logo';

import { CopyToClipboard } from 'utils/formatters';
import { SelectItemWithIcon } from 'containers/FundsPage/components/common/selectItemWithIcon';
import PopupModal from 'components/materialModal/modal';
import CopyIcon from 'images/themedIcons/copyIcon';
import ShowQrIcon from 'images/themedIcons/showQrIcon';
import { Select } from '@material-ui/core';
import { MenuItem } from '@material-ui/core';
import { TextField } from '@material-ui/core';
import { Button } from '@material-ui/core';
import { Tooltip } from '@material-ui/core';
import DepositIcon from 'images/themedIcons/depositIcon';
import Skeleton from '@material-ui/lab/Skeleton';
import ReplayIcon from 'images/themedIcons/replayIcon';
import ExpandMore from 'images/themedIcons/expandMore';
import TransferNetworkSelect from 'containers/FundsPage/components/common/transferNetworkSelect';

const stateSelector = createStructuredSelector({
  depositAndWithDrawData: makeSelectdepositAndWithDrawData(),
});

export default function DepositeAddress () {
  const currencies: Currency[] = localStorage[LocalStorageKeys.CURRENCIES]
    ? JSON.parse(localStorage[LocalStorageKeys.CURRENCIES])
    : [];
  const { depositAndWithDrawData } = useSelector(stateSelector);
  const [IsLoadingData, setIsLoadingData] = useState(false);

  const [IsTooltipOpen, setIsTooltipOpen] = useState(false);
  const [PageData, setPageData]: [depositAndWithDrawData, any] = useState(
    depositAndWithDrawData,
  );

  const [SelectedAddress, setSelectedAddress] = useState(null);

  const [SelectedIndex, setSelectedIndex] = useState(0);
  const dispatch = useDispatch();

  useEffect(() => {
    for (let i = 0; i < currencies.length; i++) {
      if (currencies[i].code === depositAndWithDrawData.balance?.code) {
        setSelectedIndex(i);
      }
    }
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_PAGE_LOADING) {
        setIsLoadingData(message.payload);
      }
      if (message.name === MessageNames.SET_DEPOSIT_PAGE_DATA) {
        setPageData(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
      localStorage.removeItem(LocalStorageKeys.SELECTED_COIN);
    };
  }, []);
  const handleSelectChange = (e: any) => {
    if (SelectedIndex !== e.target.value) {
      setSelectedIndex(e.target.value);
      setPageData({ ...PageData, walletAddress: '' });
      setSelectedAddress(null);
      dispatch(
        getRawdepositAndWithDrawDataAction({
          code: currencies[e.target.value].code,
          type: 'deposit',
          fromCoinChange: true,
        }),
      );
    }
  };
  const handleCopyClick = () => {
    setIsTooltipOpen(true);
    setTimeout(() => {
      setIsTooltipOpen(false);
    }, 1500);
    CopyToClipboard(
      (SelectedAddress ??
        PageData.walletAddress?.replace('bitcoincash:', '')) ||
        '',
    );
  };
  const showQrCode = () => {
    setIsQrOpen(true);
  };
  const handleNetworkSelect = e => {
    if (e.address) {
      setSelectedAddress(e.address);
      return;
    }
    //@ts-ignore
    setSelectedAddress(PageData.walletAddress);
  };
  const [IsQrOpen, setIsQrOpen] = useState(false);
  return (
    <Wrapper>
      <PopupModal
        isOpen={IsQrOpen}
        onClose={() => {
          setIsQrOpen(false);
        }}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            padding: '72px 90px 24px 90px',
          }}
        >
          <QRCode
            size={300}
            value={
              SelectedAddress ??
              PageData.walletAddress?.replace('bitcoincash:', '')
            }
          />
          <div style={{ marginTop: '48px' }}>
            <span
              style={{ fontFamily: 'Open Sans', color: 'var(--blackText)' }}
            >
              Deposit Address :{' '}
            </span>
            <span
              style={{ fontFamily: 'Open Sans', color: 'var(--blackText)' }}
            >
              {SelectedAddress ??
                PageData.walletAddress?.replace('bitcoincash:', '')}
            </span>
          </div>
        </div>
      </PopupModal>
      <div className='dataView'>
        <div className='blueTitle'>
          <FormattedMessage
            {...translate.Pleaseselectanycointogetdepositaddress}
          />
        </div>
        <div className='select'>
          <Select
            IconComponent={ExpandMore}
            MenuProps={{
              getContentAnchorEl: null,
              anchorOrigin: {
                vertical: 'bottom',
                horizontal: 'left',
              },
            }}
            className='select'
            margin='dense'
            fullWidth
            variant='outlined'
            value={SelectedIndex}
            onChange={handleSelectChange}
          >
            {currencies.map((item: Currency, index: number) => {
              return (
                <MenuItem key={'currencyList' + index} value={index}>
                  <SelectItemWithIcon item={item} />
                </MenuItem>
              );
            })}
          </Select>
        </div>
        <div className='descriptions pt1vh'>
          <div className='descreiptionsTitle'>
            <span className='rest'>
              <span>
                {currencies[SelectedIndex]
                  ? currencies[SelectedIndex].code
                  : ''}{' '}
                Deposit Info:{' '}
              </span>
            </span>
          </div>
        </div>
        <div className='coinDescs'>
          {/* <div>
            <img src={attention} alt='' />
          </div> */}
          <div className='descs'>
            {PageData.depositComments && (
              <ul>
                {PageData.depositComments.map((item, index) => {
                  return <li key={index}>{item}</li>;
                })}
              </ul>
            )}
          </div>
        </div>
        <div className='descriptions'>
          <div className='descreiptionsTitle'>
            <span className='rest'>
              <FormattedMessage {...translate.sendyour} />{' '}
            </span>
            <span className='coin'>
              {currencies[SelectedIndex] ? currencies[SelectedIndex].code : ''}{' '}
            </span>
            <span className='rest'>
              <FormattedMessage {...translate.tothisaddress} />
            </span>
          </div>
        </div>
        <div className='addressInput'>
          <TextField
            fullWidth
            value={
              SelectedAddress ??
              PageData.walletAddress?.replace('bitcoincash:', '')
            }
            variant='outlined'
          />
          {(SelectedAddress ?? PageData.walletAddress) === '' && (
            <Skeleton animation='wave' className='skeleton' />
          )}
        </div>

        <div className='buttons'>
          <div className='greyButtons'>
            <Tooltip
              PopperProps={{
                disablePortal: true,
              }}
              placement='top'
              open={IsTooltipOpen}
              disableFocusListener
              disableHoverListener
              disableTouchListener
              title={
                <FormattedMessage {...translate.addressCopiedToClipboard} />
              }
            >
              <Button
                onClick={handleCopyClick}
                endIcon={<CopyIcon />}
                className={Buttons.SimpleGreyButton}
              >
                <FormattedMessage {...translate.CopyAddress} />
              </Button>
            </Tooltip>

            <Button
              onClick={showQrCode}
              endIcon={<ShowQrIcon />}
              className={Buttons.SimpleGreyButton}
            >
              <FormattedMessage {...translate.ShowQRcode} />
            </Button>
          </div>
          <div className='blueButton'>
            <Button
              endIcon={<ReplayIcon />}
              color='primary'
              className={Buttons.TransParentRoundButton}
            >
              <FormattedMessage {...translate.Getnewaddress} />
            </Button>
          </div>
        </div>
        {PageData.otherNetworksConfigsAndAddresses &&
          PageData.otherNetworksConfigsAndAddresses.length > 0 && (
            <TransferNetworkSelect
              mainNetwork={PageData.mainNetwork}
              otherNetworks={PageData.otherNetworksConfigsAndAddresses}
              completedNetworkName={PageData.completedNetworkName}
              type='deposit'
              onNetworkSelect={handleNetworkSelect}
            />
          )}
        <div className='iconWrapper'>
          <DepositIcon />
        </div>
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 2vw;
  height: 100%;
  @media screen and(min-width:1600px) {
    padding: 24px 80px;
  }
  @media screen and (max-height: 870px) {
    padding-top: 8px;
    padding-bottom: 8px;
  }
  .dataView {
    height: 100%;
    flex: 5;
    display: flex;
    flex-direction: column;
    min-width: 465px;
    max-height: 675px;
    .blueTitle {
      flex: 0;
    }
    .select {
      flex: 1.8;
    }
    .descriptions {
      flex: 0;
    }
    .coinDescs {
      flex: 2;
      img {
        margin-top: 2px;
      }
    }
    .addressInput {
      flex: 1.5;
      margin-top: 1vh;
      position: relative;
      input {
        font-weight: 600;
      }
      .skeleton {
        position: absolute;
        width: 90%;
        top: 10px;
        left: 10px;
      }
    }
    .buttons {
      flex: 1;
      padding-top: 1vh;
    }
    .iconWrapper {
      flex: 10;
      display: flex;
      align-items: center;
      min-height: 182px;
    }
  }
  .iconWrapper {
    align-self: center;
  }
  .MuiInputBase-root {
    max-height: 40px;
    &.select {
      margin-top: 8px;
      margin-bottom: -5px;
    }
  }
  .MuiInputLabel-shrink {
    transform: translate(13px, 1.5px) scale(0.75) !important;
    transform-origin: top left;
  }
  /* legend {
    width: 30px !important;
  } */

  .blueTitle {
    color: var(--textBlue);
    font-weight: 600;
  }
  .coinWrapper {
    img {
      height: 25px;
      width: 25px;
      border: 2px solid #c1c1c1;
      border-radius: 50px;
      padding: 1px;
      min-width: 25px;
      min-height: 25px;
    }
    span {
      margin: 0 8px;
      font-weight: 600;
      font-size: 13px;
    }
  }
  .descriptions {
    font-weight: 700;
    color: var(--blackText);
  }
  .coinDescs {
    display: flex;
    align-items: start;
    .descs {
      padding-top: 5px;
      font-weight: 600;
      color: var(--textGrey);
      font-size: 13px !important;
      ul {
        padding-left: 32px;
        margin: 0;
        margin-bottom: 12px;
      }
    }
  }
  .addressInput {
    /* pointer-events: none; */
    legend {
      width: 0px !important;
    }
    .MuiInputLabel-shrink {
      transform: translate(15px, -6.5px) scale(0.75) !important;
    }
  }
  .buttons {
    display: flex;

    .greyButtons {
      display: flex;
      justify-content: space-between;
      flex: 7;
      padding: 0 !important;

      .simpleGreyButton {
        min-width: 145px !important;
        margin-right: 12px;
      }
    }
    .blueButton {
      text-align: end;
      min-width: 160px;
      opacity: 0;
      pointer-events: none;
      .${Buttons.TransParentRoundButton} {
        height: 32px !important;
        padding: 0 12px !important;
      }
    }
    span {
      font-size: 13px !important;
      font-weight: 600;
      line-height: 1px;
    }
    .MuiButton-endIcon {
      margin-left: 5px !important;
    }
  }
`;
