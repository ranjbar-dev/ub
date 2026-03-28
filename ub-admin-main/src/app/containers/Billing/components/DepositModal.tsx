import SaveOutlinedIcon from '@material-ui/icons/SaveOutlined';
import { RowClickedEvent } from 'ag-grid-community';
import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import { Buttons } from 'app/constants';
import EditDropDown from 'app/containers/UserDetails/components/EditDropDown';
import React, { memo, useMemo, useRef, useEffect, useState } from 'react';
import { Subscriber, MessageNames, BroadcastMessage } from 'services/messageService';

import { Payment, DepositSaveData } from '../types';

interface Props {
  row: RowClickedEvent;
  onCancel: Function;
  onSave: (data: DepositSaveData) => void;
}

function DepositModal(props: Props) {
  const fields = useMemo(
    () => [
      { name: 'Id No', field: 'id', disabled: true },
      { name: 'Method', field: 'currencyCode' },
      { name: 'Amount', field: 'amount' },
      { name: 'FromAddress', field: 'fromAddress' },
      { name: 'ToAddress', field: 'toAddress' },
      { name: 'Transaction Id', field: 'txId' },
      { name: 'Creation Date', field: 'createdAt', disabled: true },
      { name: 'LastUpdate', field: 'updatedAt', disabled: true },
      { name: 'Status', field: 'status', dropDown: true },
    ],
    [],
  );
  const { row, onSave, onCancel } = props;
  const { data: rowData }: { data: Payment } = row;
  const dataToSend = useRef(rowData);
  const [Num, setsetNum] = useState(0);
  const handleDataSend = () => {
    const dts = dataToSend.current;
    const data: DepositSaveData = {
      amount: dts.amount,
      from_address: dts.fromAddress,
      id: dts.id,
      should_deposit: dts.should_deposit ?? false,
      status: dts.status,
      to_address: dts.toAddress,
      tx_id: dts.txId,
    };
    onSave(data);
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
      if (message.name === MessageNames.CLOSE_POPUP) {
        onCancel();
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <div
      style={{
        width: '500px',
        height: ' 410px',
        background: 'white',
        borderRadius: '5px',
        boxShadow: '2px 6px 6px 11px rgba(0, 0, 0, 0.02)',
        padding: '20px 24px',
      }}
    >
      <div className="fields">
        {fields.map((item, index) => {
          return (
            <div
              className="fieldRow"
              style={{
                display: 'flex',
                alignItems: 'center',
                margin: '10px 0 ',
              }}
              key={item.field}
            >
              <div
                className="title"
                style={{ minWidth: '110px', fontSize: '13px' }}
              >
                {item.name}
              </div>
              <div className="input" style={{ flex: '1' }}>
                {item.disabled === true ? (
                  <span style={{ fontSize: '13px', fontWeight: 600 }}>
                    {dataToSend.current[item.field]}
                  </span>
                ) : item.dropDown === true ? (
                  <EditDropDown
                    onSelect={(e: string) => {
                      dataToSend.current.status = e as Payment['status'];
                    }}
                    initialValue={{ name: rowData.status, id: rowData.status }}
                    options={[
                      {
                        name: 'All',
                        id: 'all',
                      },
                      {
                        name: 'Created',
                        id: 'created',
                      },
                      {
                        name: 'In progress',
                        id: 'in_progress',
                      },
                      {
                        name: 'completed',
                        id: 'completed',
                      },
                      {
                        name: 'failed',
                        id: 'failed',
                      },
                      {
                        name: 'Cancel',
                        id: 'cancel',
                      },
                      {
                        name: 'Reject',
                        id: 'reject',
                      },
                    ]}
                  />
                ) : (
                  <input
                    style={{ width: '100%', minHeight: '28px' }}
                    value={dataToSend.current[item.field] as string}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      dataToSend.current[item.field] = e.target.value;
                      setsetNum(Num + 1);
                    }}
                  />
                )}
              </div>
            </div>
          );
        })}
      </div>
      <div
        className="buttons"
        style={{
          display: 'flex',
          position: 'absolute',
          bottom: '16px',
          right: '16px',
        }}
      >
        <IsLoadingWithTextAuto
          text="Cancel"
          style={{
            margin: '0 7px',
          }}
          className={Buttons.BlackButton}
          loadingId="CancelButton"
          onClick={() => {
            props.onCancel();
          }}
        />
        <IsLoadingWithTextAuto
          text="Save"
          loadingId="DepositModalSaveButton"
          style={{
            margin: '0 7px',
          }}
          onClick={() => {
            dataToSend.current.should_deposit = false;
            handleDataSend();
          }}
          className={Buttons.SkyBlueButton}
          icon={<SaveOutlinedIcon />}
        />
        <IsLoadingWithTextAuto
          text="Save And Deposit"
          loadingId="DepositModalSaveAndDepositButton"
          style={{
            margin: '0 7px',
          }}
          className={Buttons.SkyBlueButton}
          onClick={() => {
            dataToSend.current.should_deposit = true;
            handleDataSend();
          }}
        />
      </div>
    </div>
  );
}

export default memo(DepositModal);
