/**
 *
 * PlaceHolder
 *
 */
import { Button, InputAdornment, OutlinedInput } from '@material-ui/core';
import { Buttons } from 'app/constants';
import { translations } from 'locales/i18n';
import React, { memo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components/macro';

interface Props {}

export const PlaceHolder = memo((props: Props) => {
  const { t } = useTranslation();

  const details = {
    adminComments: [],
    adminStatus: 'pending',
    amount: '0.00090000',
    autoTransfer: false,
    country: 'Algeria',
    createdAt: '2020-06-17 13:45:47',
    currencyCode: 'ZEC',
    currencyName: 'zec',
    fee: '0.00010000',
    fromAddress: '',
    id: 321,
    ip: '',
    level: 'level4',
    name: 'mohsenp',
    status: 'cancel',
    tag: '',
    toAddress: 't1a6KGCxRcAE3grePP6mCra7osjsfPXduiE',
    totalAmount: '0.00100000',
    txId: '',
    type: 'withdraw',
    userEmail: 'mhsnprvr@gmail.com',
    userId: 113,
  };
  const rowData = {
    amount: '0.00100000',
    createdAt: '2020-04-23',
    currencyCode: 'ZEC',
    currencyId: 12,
    fromAddress: '',
    id: 321,
    status: 'cancel',
    toAddress: 't1a6KGCxRcAE3grePP6mCra7osjsfPXduiE',
    txId: '',
    type: 'withdraw',
    updatedAt: '2020-06-17',
    userEmail: 'mhsnprvr@gmail.com',
    userId: 113,
  };
  const DetailRow = useCallback(propes => {
    const { title, value, last, small } = propes;
    return (
      <div
        className={`detailRow ${last === true ? 'last' : ''} ${
          small === true ? 'small' : ''
        }`}
      >
        <div className="title">{title}</div>
        <div className="value">{value}</div>
      </div>
    );
  }, []);
  return (
    <Wrapper>
      <div className="wrapperTitle">
        {t(translations.CommonTitles.WithdrawalRequest())}
      </div>
      <div className="content">
        <div className="detailsContainer">
          <div className="detailRowsContainer">
            <DetailRow
              title={t(translations.CommonTitles.UserId())}
              value={details.userId}
            />
            <DetailRow
              title={t(translations.CommonTitles.Name())}
              value={details.name}
            />
            <DetailRow
              title={t(translations.CommonTitles.Email())}
              value={details.userEmail}
            />
            <DetailRow
              title={t(translations.CommonTitles.Country())}
              value={details.country}
            />
            <DetailRow
              title={t(translations.CommonTitles.TrustLevel())}
              value={details.level}
            />
            <DetailRow
              title={t(translations.CommonTitles.RequestDate())}
              value={details.createdAt}
            />
            <DetailRow
              title={t(translations.CommonTitles.LastResponseDate())}
              value={''}
            />
            <DetailRow
              title={t(translations.CommonTitles.IPAddress())}
              value={details.ip}
            />
            <DetailRow
              title={t(translations.CommonTitles.Method())}
              value={details.currencyName}
            />
            <DetailRow
              title={t(translations.CommonTitles.Amount())}
              value={details.amount}
            />
            <DetailRow
              title={t(translations.CommonTitles.WalletAddress())}
              value={details.fromAddress}
              last={true}
            />
          </div>
          <div className="commentsContainer"></div>
        </div>
        <div className="actionsContainer">
          <div className="inp1">
            <OutlinedInput
              fullWidth
              margin="dense"
              endAdornment={
                <InputAdornment position="end">
                  {details.currencyCode}
                </InputAdornment>
              }
            />
          </div>
          <div className="inp2">
            <div className="fee">{t(translations.CommonTitles.Fee())}</div>
            <OutlinedInput
              margin="dense"
              fullWidth
              endAdornment={
                <InputAdornment position="end">
                  {details.currencyCode}
                </InputAdornment>
              }
            />
          </div>
          <div className="autoTransfer"></div>
          <div className="sendMoney">
            <Button
              color="primary"
              className="sendMoneyButton"
              variant="contained"
            >
              {t(translations.CommonTitles.SendMoney())}
            </Button>
          </div>
          <div className="rejectCancel">
            <Button className={Buttons.RedButton}>
              {t(translations.CommonTitles.Reject())}
            </Button>
            <Button className={Buttons.BlackButton}>
              {t(translations.CommonTitles.Cancel())}
            </Button>
          </div>
          <div className="actionDetailRows">
            <DetailRow
              title={t(translations.CommonTitles.TotalDeposit())}
              value={details.totalAmount}
              small={true}
            />
            <DetailRow
              title={t(translations.CommonTitles.TotalWithdraw())}
              value={details.totalAmount}
              small={true}
            />
            <DetailRow
              title={t(translations.CommonTitles.TotalCommissions())}
              value={details.totalAmount}
              small={true}
            />
            <DetailRow
              title={t(translations.CommonTitles.TotalBalances())}
              value={details.totalAmount}
              small={true}
            />
            <DetailRow
              title={t(translations.CommonTitles.TotalOntrade())}
              value={details.totalAmount}
              small={true}
              last={true}
            />
          </div>
          <div className="bottomActions">
            <Button className={Buttons.BlackButton}>
              {t(translations.CommonTitles.SetAsUncheck())}
            </Button>
            <Button className={Buttons.RedButton}>
              {t(translations.CommonTitles.NeedsRecheck())}
            </Button>
            <Button className={Buttons.SkyBlueButton}>
              {t(translations.CommonTitles.MoveToAccepts())}
            </Button>
          </div>
        </div>
      </div>
    </Wrapper>
  );
});

const Wrapper = styled.div`
  width: 100%;
  height: 100%;
  padding-top: 70px;
  display: flex;
  flex-direction: column;
  .wrapperTitle {
    width: 100%;
    height: 40px;
    display: flex;
    align-items: center;
    background: rgb(213, 217, 230);
    padding: 0 10px;
    font-size: 14px;
    font-weight: 600;
  }
  .content {
    display: flex;
    width: 100%;
    flex: 1;

    .detailsContainer {
      background: rgb(245, 246, 248);
      flex: 5;
    }
    .actionsContainer {
      flex: 3;
      display: flex;
      flex-direction: column;
      align-items: center;
      .inp1,
      .inp2 {
        width: 60%;
        margin: 15px 0;
      }
      .inp2 {
        display: flex;
        align-items: center;
        .fee {
          min-width: 50px;
        }
      }
      .sendMoney {
        .sendMoneyButton {
          min-width: 200px;
        }
      }
      .rejectCancel {
        width: 200px;
        display: flex;
        justify-content: space-around;
        margin: 8px 0;
      }
      .bottomActions {
        display: flex;
        justify-content: space-around;
        margin: 8px 0;
        button {
          margin: 0px 8px;
          padding: 0 10px;
        }
      }
    }
    .detailsContainer {
      padding: 10px 20px;
   
      .detailRowsContainer {
      }
    }

    .detailRow {
      display: flex;
      min-height: 35px;
      align-items: center;
      background: white;
      border-bottom: 1px solid #eee;
      &.last {
        border-bottom: none;
      }
      .title {
        min-width: 150px;
        min-width: 150px;
        font-size: 12px;
        padding: 0 12px;
        font-weight: 600;
        color: #6c757e;
      }
      .value {
        flex: 1;
        min-width: 150px;
        min-width: 150px;
        font-size: 12px;
        padding: 0 12px;
        font-weight: 600;
      }
    }
    .actionDetailRows {
      display: flex;
      flex-direction: column;
      align-items: center;
    }
  }
`;
