import React, { useState, useEffect, memo, useRef } from 'react';
import styled from 'styles/styled-components';
import { Buttons } from 'containers/App/constants';
import { Select, Tooltip, Zoom } from '@material-ui/core';
import { FormattedMessage, injectIntl } from 'react-intl';
import { MenuItem } from '@material-ui/core';
import { Currency } from 'containers/App/types';
import { LocalStorageKeys } from 'services/constants';
import translate from '../../../messages';
import { FormControl } from '@material-ui/core';
import { createStructuredSelector } from 'reselect';
import {
  makeSelectdepositAndWithDrawData,
  makeSelectFormerWithdrawAddresses,
} from 'containers/FundsPage/selectors';
import { useSelector, useDispatch } from 'react-redux';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import {
  getRawdepositAndWithDrawDataAction,
  getFormerWithdrawAddressesAction,
  preWithdrawAction,
} from 'containers/FundsPage/actions';
import attention from 'images/attentionIcon.svg';

import {
  depositAndWithDrawData,
  WithdrawModel,
} from 'containers/FundsPage/types';
import { TextField } from '@material-ui/core';
import { Button } from '@material-ui/core';

import DataRow from 'components/dataRow';
import * as Decimal from 'utils/decimal.js';
import { SelectItemWithIcon } from 'containers/FundsPage/components/common/selectItemWithIcon';
import { WithdrawAddress } from 'containers/AddressManagementPage/types';
import PopupModal from 'components/materialModal/modal';
import G2fa from './2faAndEmailAlert';
import { toast } from 'components/Customized/react-toastify/cjs/react-toastify';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import AddressInputWithdraw from './addressInputWithdraw';
import ExpandMore from 'images/themedIcons/expandMore';
import {
  CurrencyFormater,
  formatCurrencyWithMaxFraction,
} from 'utils/formatters';
import TransferNetworkSelect from '../../../components/common/transferNetworkSelect';
import { currencyMap } from 'utils/sharedData';
const stateSelector = createStructuredSelector({
  depositAndWithDrawData: makeSelectdepositAndWithDrawData(),
  formerWithdrawAddresses: makeSelectFormerWithdrawAddresses(),
  //  userData: makeSelectUserData(),
});
let amountEntryTimeout;
let wData: any = {};
let formerAdresses;
const WithdrawalAddress = ({ intl }) => {
  const currencies: Currency[] = localStorage[LocalStorageKeys.CURRENCIES]
    ? JSON.parse(localStorage[LocalStorageKeys.CURRENCIES])
    : [];
  const { depositAndWithDrawData, formerWithdrawAddresses } = useSelector(
    stateSelector,
  );
  formerAdresses = formerWithdrawAddresses;
  const [PageData, setPageData]: [depositAndWithDrawData, any] = useState(
    depositAndWithDrawData,
  );

  const [IsLoadingData, setIsLoadingData] = useState(false);
  const [FormerAddresses, setFormerAddresses] = useState(
    formerWithdrawAddresses,
  );
  const [SelectedAddress, setSelectedAddress] = useState('');
  const [SelectedAmount, setSelectedAmount] = useState('');
  const [Is2faOpen, setIs2faOpen] = useState(false);
  const [IsSelectAddressOpen, setIsSelectAddressOpen] = useState(false);
  const [IsWithdrawing, setIsWithdrawing] = useState(false);
  const [CanSubmit, setCanSubmit] = useState(false);
  const [YouWillGet, setYouWillGet] = useState({ show: false, value: 0 });

  const requiredVerificationData = useRef<any>({});

  const [ToolTipProps, setToolTipProps]: [any, any] = useState({
    isOpen: false,
    message: <FormattedMessage {...translate.Totalbalance} />,
  });

  const [selectedNetwork, setSelectedNetwork] = useState(PageData.mainNetwork);

  const { networksConfigsAndAddresses } = PageData;
  let fee = PageData.balance?.fee;
  if (networksConfigsAndAddresses) {
    const config = networksConfigsAndAddresses.find(
      config => config.code === selectedNetwork,
    );
    if (config?.fee) {
      fee = config.fee;
    }
  }
  let youGet = 0.0;
  if (fee && SelectedAmount) {
    youGet = +Decimal(SelectedAmount)
      .sub(Decimal(fee))
      .toString();
  }

  useEffect(() => {
    setSelectedNetwork(PageData.mainNetwork);
    return () => {};
  }, [PageData.mainNetwork]);

  const [SelectedIndex, setSelectedIndex] = useState(0);
  const dispatch = useDispatch();

  useEffect(() => {
    let foundCurrency = false;
    for (let i = 0; i < currencies.length; i++) {
      if (currencies[i].code === depositAndWithDrawData.balance?.code) {
        foundCurrency = true;
        dispatch(
          getFormerWithdrawAddressesAction({ code: currencies[i].code }),
        );
        setSelectedIndex(i);
      }
    }
    if (foundCurrency === false) {
      dispatch(getFormerWithdrawAddressesAction({ code: 'BTC' }));
    }
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_PAGE_LOADING) {
        if (message.payload === true) {
          setIsLoadingData(message.payload);
          return;
        }
        setTimeout(() => {
          setIsLoadingData(message.payload);
        }, 200);
      } else if (message.name === MessageNames.SET_DEPOSIT_PAGE_DATA) {
        setPageData(message.payload);
        wData = message.payload;
      } else if (
        message.name === MessageNames.OPEN_WITHDRAW_VERIFICATION_POPUP
      ) {
        requiredVerificationData.current = message.payload;
        setIs2faOpen(true);
      } else if (message.name === MessageNames.SETLOADING) {
        setIsWithdrawing(message.payload);
      } else if (message.name === MessageNames.ADD_DATA_ROW_TO_WITHDRAWS) {
        if (PageData.balance) {
          const data = wData.balance ? wData : PageData;
          data.balance.availableAmount = Decimal(+data.balance.availableAmount)
            .sub(+message.payload.amount)
            .toString();
          setPageData({ ...data });
        }
      }
      if (message.name === MessageNames.SET_FORMER_WITHDRAW_ADDRESSES) {
        setFormerAddresses(message.payload);
      }
      // if (message.name === MessageNames.ADDITIONAL_ACTION) {
      //   if (message.payload.title === 'removeAddButton') {
      //     let former = [...FormerAddresses];
      //     former.unshift(message.payload.value);
      //     setFormerAddresses(former);
      //   }
      // }
    });
    return () => {
      wData = {};
      formerAdresses = [];
      Subscription.unsubscribe();
      localStorage.removeItem(LocalStorageKeys.SELECTED_COIN);
    };
  }, []);
  const handleSelectChange = (e: any) => {
    setSelectedIndex(e.target.value);
    setSelectedAddress('');
    setYouWillGet({ show: false, value: 0 });
    dispatch(
      getRawdepositAndWithDrawDataAction({
        code: currencies[e.target.value].code,
        type: 'withdraw',
        fromCoinChange: true,
      }),
    );
    dispatch(
      getFormerWithdrawAddressesAction({
        code: currencies[e.target.value].code,
      }),
    );
    MessageService.send({
      name: MessageNames.ADDITIONAL_ACTION,
      payload: {
        title: 'removeAddButton',
      },
    });
  };

  const handleSubmit = () => {
    if (SelectedAddress === '') {
      toast.warn('please enter or select the Address to withdraw');
      return;
    } else if (SelectedAmount === '') {
      toast.warn('please enter amount to withdraw');
      return;
    }

    const sendingData: WithdrawModel = {
      label: '',
      code: currencies[SelectedIndex].code,
      ...(PageData.otherNetworksConfigsAndAddresses &&
        PageData.otherNetworksConfigsAndAddresses.length > 0 && {
          network: selectedNetwork,
        }),
      amount: SelectedAmount,
      address: SelectedAddress,
    };

    dispatch(preWithdrawAction(sendingData));
  };

  const handleAmountChange = (e: any) => {
    const value = e;
    if (isNaN(value)) {
      return;
    }
    clearTimeout(amountEntryTimeout);
    amountEntryTimeout = setTimeout(() => {
      if (value === '') {
        setToolTipProps({
          ...ToolTipProps,
          isOpen: false,
          message: <FormattedMessage {...translate.maximumWithdrawAmountIs} />,
        });
        setYouWillGet({ ...YouWillGet, show: false });
        setCanSubmit(false);
      } else if (Number(value) > Number(PageData.balance?.availableAmount)) {
        setToolTipProps({
          ...ToolTipProps,
          isOpen: true,
          message: (
            <div className='tooltipMessage'>
              <FormattedMessage {...translate.maximumWithdrawAmountIs} />
              {PageData.balance?.availableAmount}{' '}
              {currencies[SelectedIndex].code}
            </div>
          ),
        });
        setYouWillGet({ ...YouWillGet, show: false });
        setCanSubmit(false);
      } else if (Number(value) < Number(PageData.balance?.minimumWithdraw)) {
        setToolTipProps({
          ...ToolTipProps,
          isOpen: true,
          message: (
            <div className='tooltipMessage'>
              <FormattedMessage {...translate.minimumWithdrawAmountIs} />
              {PageData.balance?.minimumWithdraw}{' '}
              {currencies[SelectedIndex].code}
            </div>
          ),
        });
        setYouWillGet({ ...YouWillGet, show: false });
        setCanSubmit(false);
      } else {
        setToolTipProps({
          ...ToolTipProps,
          isOpen: false,
          message: (
            <div className='tooltipMessage'>
              <FormattedMessage {...translate.minimumWithdrawAmountIs} />
              {PageData.balance?.minimumWithdraw}
            </div>
          ),
        });
        setYouWillGet({
          ...YouWillGet,
          show: true,
          value: youGet,
        });
        setCanSubmit(true);
      }
    }, 500);
    setSelectedAmount(value);
  };

  const handleCloseAddress = () => {
    setIsSelectAddressOpen(false);
    setFormerAddresses(formerAdresses);
  };

  const handleOpenAddress = () => {
    setIsSelectAddressOpen(true);
  };

  const handleNetworkSelect = e => {
    setSelectedNetwork(e.value);
  };

  const coinFraction =
    currencyMap({ code: currencies[SelectedIndex].code })?.showDigits ?? 8;

  const testOpen = () => {
    MessageService.send({
      name: MessageNames.OPEN_WITHDRAW_VERIFICATION_POPUP,
      payload: { needEmailCode: true, need2fa: true },
    });
  };

  return (
    <Wrapper>
      {/* <button onClick={testOpen}>open</button> */}
      <PopupModal
        isOpen={Is2faOpen}
        onClose={() => {
          setIs2faOpen(false);
        }}
      >
        <G2fa
          label=''
          network={selectedNetwork}
          intl={intl}
          code={currencies[SelectedIndex].code}
          amount={SelectedAmount}
          address={SelectedAddress}
          requiredData={requiredVerificationData.current}
          onClose={() => {
            setIs2faOpen(false);
          }}
        />
      </PopupModal>

      <div className='dataView'>
        <div className='blueTitle'>
          <FormattedMessage {...translate.Pleaseselectanycointowithdraw} />
        </div>
        <div className='select'>
          <FormControl style={{ width: '100%' }}>
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
                    <SelectItemWithIcon item={item}></SelectItemWithIcon>
                  </MenuItem>
                );
              })}
            </Select>
          </FormControl>
        </div>
        <div className='dataColumn'>
          <div className='wrapp'>
            <DataRow
              isLoading={IsLoadingData}
              title={<FormattedMessage {...translate.Totalbalance} />}
              value={
                CurrencyFormater(PageData.balance?.totalAmount + '') +
                ' ' +
                currencies[SelectedIndex].code
              }
            ></DataRow>
            <DataRow
              hideWhenSmall={
                PageData.otherNetworksConfigsAndAddresses &&
                PageData.otherNetworksConfigsAndAddresses.length > 0
              }
              isLoading={IsLoadingData}
              title={<FormattedMessage {...translate.Inorder} />}
              value={
                (PageData.balance
                  ? CurrencyFormater(PageData.balance.inOrderAmount)
                  : '') +
                ' ' +
                currencies[SelectedIndex].code
              }
            ></DataRow>
            <DataRow
              isLoading={IsLoadingData}
              title={<FormattedMessage {...translate.AvailableBalance} />}
              value={
                CurrencyFormater(
                  PageData.balance ? PageData.balance.availableAmount : '',
                ) +
                ' ' +
                currencies[SelectedIndex].code
              }
              boldValue={true}
            ></DataRow>
          </div>
        </div>
        <div className='descriptions'>
          <div className='descreiptionsTitle'>
            <span className='coin'>{currencies[SelectedIndex].code} </span>
            <span className='rest'>
              <FormattedMessage {...translate.withdrawalAddress} />
            </span>
          </div>
        </div>
        <div className='selectionWrapper'>
          <div className='selection'>
            <Select
              disabled={formerWithdrawAddresses.length === 0}
              IconComponent={ExpandMore}
              MenuProps={{
                getContentAnchorEl: null,
                anchorOrigin: {
                  vertical: 'bottom',
                  horizontal: 'left',
                },
              }}
              className='noLabel'
              fullWidth
              variant='outlined'
              open={IsSelectAddressOpen}
              onClose={handleCloseAddress}
              onOpen={handleOpenAddress}
              value=''
              onChange={(e: any) => {
                setSelectedAddress(
                  formerWithdrawAddresses[e.target.value].address,
                );
              }}
            >
              {FormerAddresses.map((item: WithdrawAddress, index: number) => {
                return (
                  <MenuItem key={'currencyList' + index} value={index}>
                    <div className='badge'>
                      {item.code} - {item.label}
                    </div>{' '}
                    {' ' + item.address}
                  </MenuItem>
                );
              })}
            </Select>
          </div>
          <div className={`input `}>
            <AddressInputWithdraw
              code={currencies[SelectedIndex].code}
              selectedAddress={SelectedAddress}
              formerWithdrawAddresses={formerWithdrawAddresses}
              network={
                PageData.otherNetworksConfigsAndAddresses &&
                PageData.otherNetworksConfigsAndAddresses.length > 0
                  ? selectedNetwork
                  : null
              }
              onAddressExists={item => {
                setFormerAddresses([...item]);
                setIsSelectAddressOpen(true);
              }}
              onChange={e => {
                setSelectedAddress(e);
              }}
            />
          </div>
        </div>
        {PageData.otherNetworksConfigsAndAddresses &&
          PageData.otherNetworksConfigsAndAddresses.length > 0 && (
            <TransferNetworkSelect
              mainNetwork={PageData.mainNetwork}
              otherNetworks={PageData.otherNetworksConfigsAndAddresses}
              completedNetworkName={PageData.completedNetworkName}
              type='withdraw'
              onNetworkSelect={handleNetworkSelect}
            />
          )}
        <div className='coinDescs pt1vh'>
          <div className='descs'>
            {PageData.withdrawComments && (
              <ul>
                {PageData.withdrawComments.map((item, index) => {
                  return <li key={index}>{item}</li>;
                })}
              </ul>
            )}
            {/* <FormattedMessage {...translate.dontWithdraw} /> */}
          </div>
        </div>
        <div className='amountWrapper'>
          <Tooltip
            PopperProps={{
              disablePortal: true,
            }}
            arrow
            TransitionComponent={Zoom}
            placement='top'
            open={ToolTipProps.isOpen}
            disableFocusListener
            disableHoverListener
            disableTouchListener
            className='amountTooltip'
            title={ToolTipProps.message}
          >
            <TextField
              className='noLabel'
              placeholder={intl.formatMessage({
                id: 'containers.FundsPage.amount',
                defaultMessage: 'Amount',
              })}
              fullWidth
              variant='outlined'
              value={SelectedAmount}
              onChange={e => handleAmountChange(e.target.value)}
            />
          </Tooltip>
          <div className='valueWrapper'>
            <span className='available'>
              <FormattedMessage {...translate.Available} />
            </span>
            <span className='value'>
              {CurrencyFormater(
                PageData.balance ? PageData.balance.availableAmount : '',
              )}{' '}
              {currencies[SelectedIndex].code}
            </span>
            <span className='buttonWrapper'>
              <Button
                className='allButton'
                onClick={() => {
                  handleAmountChange(PageData.balance?.availableAmount);
                }}
              >
                <FormattedMessage {...translate.all} />
              </Button>
            </span>
          </div>
        </div>
        <div className='dataColumn'>
          <div className='wrap'>
            <DataRow
              isLoading={IsLoadingData}
              title={<FormattedMessage {...translate.Minimumwithdrawal} />}
              value={
                PageData.balance?.minimumWithdraw +
                ' ' +
                currencies[SelectedIndex].code
              }
              small={true}
            ></DataRow>
            <DataRow
              isLoading={IsLoadingData}
              title={<FormattedMessage {...translate.TransactionFee} />}
              value={
                Number(
                  formatCurrencyWithMaxFraction(fee, coinFraction),
                ).toFixed(coinFraction) +
                ' ' +
                currencies[SelectedIndex].code
              }
              boldValue={true}
            ></DataRow>
            {YouWillGet.show === true && (
              <DataRow
                className='shimmer'
                isLoading={IsLoadingData}
                title={<FormattedMessage {...translate.YouWillGet} />}
                value={
                  Number(
                    formatCurrencyWithMaxFraction(
                      youGet.toString(),
                      coinFraction,
                    ),
                  ).toFixed(coinFraction) +
                  ' ' +
                  currencies[SelectedIndex].code
                }
                boldValue={true}
              ></DataRow>
            )}
          </div>
        </div>
        <div className='submitWrapper'>
          <Button
            fullWidth
            disabled={!CanSubmit || SelectedAddress.length == 0}
            onClick={handleSubmit}
            color='primary'
            variant='contained'
          >
            <IsLoadingWithText
              isLoading={IsWithdrawing}
              text={<FormattedMessage {...translate.submit} />}
            />
          </Button>
        </div>
      </div>
    </Wrapper>
  );
};
export default injectIntl(memo(WithdrawalAddress));
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  padding: 24px 2vw;
  height: 100%;
  --heightBreakPoint: 870px;
  @media screen and (min-width: 1600px) {
    padding: 24px 80px;
  }
  @media screen and (max-height: 870px) {
    padding-top: 8px;
    padding-bottom: 8px;
  }
  .dataView {
    height: 100%;
    min-width: 435px;
    display: flex;
    flex-direction: column;
    max-height: 675px;
    overflow-y: auto;
    .blueTitle {
      /* @media screen and (max-height: 870px) {
        margin-bottom: -10px;
      } */
      flex: 1;
    }
    .select {
      flex: 2.4;
      /* margin-bottom: 12px;
      @media screen and (max-height: 870px) {
        margin-bottom: 0.5vh;
      } */
    }
    .dataColumn {
      padding-top: 12px;
      flex: 6;
      /* margin-bottom: 48px;
      .title {
        margin-bottom: 0.9vh;
      }
      @media screen and (max-height: 870px) {
        margin-bottom: 2vh;
        .title {
          margin-bottom: 0vh;
        }
      } */
    }
    .descriptions {
      flex: 1.3;
      /* margin-bottom: 12px;
      @media screen and (max-height: 870px) {
        margin-bottom: 0;
      } */
    }
    .selectionWrapper {
      flex: 2;
    }
    .coinDescs {
      flex: 3;
      display: flex;
      align-items: center;
      padding-top: 0;
      img {
        margin-top: 2px;
      }
      /* margin-bottom: 48px;
      @media screen and (max-height: 870px) {
        margin-bottom: 2vh;
      } */
    }
    .amountWrapper {
      flex: 2;

      .MuiTooltip-tooltipPlacementTop {
        background: var(--orange) !important;
      }
      .MuiTooltip-arrow {
        color: var(--orange) !important;
        left: 23px !important;
      }
      .tooltipMessage {
        font-size: 12px;
        font-weight: 600;
        span {
          font-size: 12px;
          font-weight: 600;
        }
      }
      .MuiTooltip-popper {
        transform: translate3d(20px, -42px, 0px) !important;
      }
    }
    .submitWrapper {
      flex: 8;
    }
  }
  .iconWrapper {
    flex: 5;
    align-self: center;
    padding: 90px 0;
  }
  .MuiInputBase-root {
    max-height: 40px;
    &.select {
      margin-top: 8px;
      margin-bottom: -5px;
    }
  }
  .MuiInputLabel-formControl {
    left: 11px;
    color: var(--textGrey);
  }
  /* legend {
    width: 30px !important;
  } */
  .noLabel {
    legend {
      width: 0px !important;
    }
  }
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

    span {
      font-size: 11px !important;
    }
    .descs {
      padding-top: 2px;
      font-weight: 600;
      color: var(--textGrey);
      ul {
        padding-left: 32px;
        margin: 0;
        margin-bottom: 12px;
      }
    }
  }
  .addressInput {
    pointer-events: none;
    margin-bottom: 24px;
    legend {
      width: 90px !important;
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
    }
    .blueButton {
      flex: 4;
      text-align: end;
      .${Buttons.TransParentRoundButton} {
        height: 32px !important;
        padding: 0 12px !important;
      }
    }
    span {
      font-size: 13px !important;
      font-weight: 600;
    }
    .MuiButton-endIcon {
      margin-left: 5px !important;
    }
  }

  .selectionWrapper {
    position: relative;
    .input {
      position: absolute;
      z-index: 2;
      width: 90%;
      top: -2px;
      left: 10px;

      .MuiInput-underline:before {
        display: none;
      }
      .MuiInput-underline:after {
        display: none;
      }
      .MuiInputBase-root {
        max-height: 31px;
        min-height: 31px;
        border-radius: 10px;
        margin-top: 9px;
      }
      .MuiFormControl-fullWidth {
        background: transparent;
      }
      &.long {
        width: 97%;
        border-radius: 10px;
        top: 1px;
        background: white;
        height: 37px;
        .MuiInputBase-input {
          padding-top: 0px !important;
        }
      }
    }
    .selection {
      .MuiInputBase-input {
        opacity: 0;
        padding: 0px;
        min-height: 40px;
      }
    }
  }
  .amountWrapper {
    position: relative;
    input {
      font-weight: 700;
      &::placeholder {
        font-weight: 500;
        color: var(--textGrey) !important;
        font-size: 13px !important;
      }
    }
    .valueWrapper {
      position: absolute;
      top: 8px;
      right: 5px;
      span {
        font-size: 12px !important;
      }
      .available {
        span {
          font-weight: 600;
          color: var(--textGrey);
        }
      }
      .value {
        margin: 2px;
        color: var(--textBlue);
      }
      .buttonWrapper {
        .allButton {
          max-height: 17px;
          padding: 0 !important;
          background: var(--lightBlue);
          min-width: 45px;
          margin: 0 4px;
          margin-top: -2px;
          span {
            color: var(--textBlue);
            font-weight: 600;
            font-size: 11px !important;
          }
        }
      }
    }
  }
  .MuiInputBase-input {
    font-size: 14px;
  }
  .loadingCircle {
    top: 8px !important;
  }
`;
