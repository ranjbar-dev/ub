import PopupModal from 'components/materialModal/modal';
import React, { FC, useEffect, useRef, useState } from 'react';
import { MessageNames, Subscriber } from 'services/message_service';
import { TwoFaAndEmailBody } from './twoFaBody';
import { ISubmitTwoFaAndEmailCode } from './types';
interface PopupProps{
  onSubmit:(data:ISubmitTwoFaAndEmailCode)=>void,
  onClose:()=>void
}
export const TwofaAndEmailVerificationPopup:FC<PopupProps>=({onSubmit,onClose})=> {
  const requiredVerificationData = useRef<any>({});
const [isOpen, setIsOpen] = useState(false);
const handleSubmit=(data:any)=>{
  onSubmit(data);
  setIsOpen(false);
};
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (
        message.name === MessageNames.OPEN_TWOFA_AND_EMAILCODE_POPUP
      ) {
        requiredVerificationData.current = message.payload;
        setIsOpen(true);
      }if (message.name === MessageNames.CLOSE_MODAL) {
        setIsOpen(false);
      }

    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <PopupModal
    isOpen={isOpen}
    onClose={() => {
      setIsOpen(false);
      onClose();
    }}
  >
<TwoFaAndEmailBody requiredData={requiredVerificationData.current} onFinalSubmit={handleSubmit} />
  </PopupModal>

  );
};

