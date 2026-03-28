import IsLoadingWithTextAuto from 'app/components/isLoadingWithText/isLoadingWithTextAuto';
import UbDropDown from 'app/components/UbDropDown';
import UBInput from 'app/components/UBInput/UBInput';
import { Buttons } from 'app/constants';
import React, { memo, useEffect, useRef } from 'react';
import { Subscriber, MessageNames } from 'services/messageService';
import type { BroadcastMessage } from 'services/messageService';
import styled from 'styled-components/macro';
import { under } from 'utils/formatters';

import {
  IConstructiveModalElement,
  IFinanceMethod,
} from '../../containers/FinanceMethods';

interface Props {
  initialData: IFinanceMethod;
  modalFields: IConstructiveModalElement[];
  onSubmit: (data: Record<string, unknown>) => void;
  onCancel: () => void;
}

/**
 * Dynamic form modal that generates input/dropdown fields from a modalFields config.
 * Used in the Finance Methods page to create or edit payment method records.
 *
 * @example
 * ```tsx
 * <ConstructiveModal
 *   initialData={method}
 *   modalFields={fields}
 *   onSubmit={(data) => dispatch(saveMethod(data))}
 *   onCancel={() => setOpen(false)}
 * />
 * ```
 */
function ConstructiveModal(props: Props) {
  const { initialData, modalFields, onSubmit, onCancel } = props;
  const sendingData = useRef<IFinanceMethod>({ ...initialData });
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
  const handleChange = (data: {
    field: keyof IFinanceMethod;
    value: string | boolean;
  }) => {
    (sendingData.current as unknown as Record<string, string | boolean>)[data.field] = data.value;
  };
  const handleSubmit = () => {
    const dataToSend: Record<string, unknown> = {};
    for (const key in sendingData.current) {
      if (Object.prototype.hasOwnProperty.call(sendingData.current, key)) {
        const element = (sendingData.current as unknown as Record<string, string | boolean>)[key];
        if (element === 'true') {
          dataToSend[under(key)] = true;
        }
        if (element === 'false') {
          dataToSend[under(key)] = false;
        } else {
          dataToSend[under(key)] = element;
        }
      }
    }
    onSubmit(dataToSend);
  };
  return (
    <ModalWrapper>
      {modalFields.map((item, index) => {
        if (item.type === 'dropDown') {
          return (
            <div className="dropDown" key={index}>
              <UbDropDown
                initialValue={initialData[item.field] as string | boolean}
                // @ts-expect-error — UbDropDown options type is more restrictive than item.options
                options={item.options}
                label={item.name}
                aria-label={item.name}
                onSelect={val => {
                  handleChange({ field: item.field, value: val });
                }}
              />
            </div>
          );
        }
        return (
          <UBInput
          key={index+'input'}
            id={`cm-field-${item.field}`}
            // @ts-expect-error — dynamic field access type
            initialValue={initialData[item.field]}
            label={item.name}
            properties={{
              className: 'narrow',
              disabled: item.editable === true ? false : true,
            }}
            onChange={(val: string) => {
              handleChange({ field: item.field, value: val });
            }}
          />
        );
      })}
      <div className="buttonsWrapper">
        <IsLoadingWithTextAuto
          text="Cancel"
          className={Buttons.BlackButton}
          loadingId="ConstructiveModalCancel"
          onClick={onCancel}
        />
        <IsLoadingWithTextAuto
          text="Submit"
          className={Buttons.SkyBlueButton}
          loadingId="ConstructiveModalSubmit"
          onClick={handleSubmit}
        />
      </div>
    </ModalWrapper>
  );
}

export default memo(ConstructiveModal);
const ModalWrapper = styled.div`
  min-width: 300px;
  min-height: 300px;
  padding: 48px 42px 22px 42px;
  .buttonsWrapper {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 26px;
  }
`;
