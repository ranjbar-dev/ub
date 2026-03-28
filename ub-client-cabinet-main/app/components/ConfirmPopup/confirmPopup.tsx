import { Button } from '@material-ui/core';
import PopupModal from 'components/materialModal/modal';
import React, { FC } from 'react';

interface PopupProps {
  isOpen: boolean;
  onClose: () => void;
  onCancelClick: () => void;
  onSubmitClick: () => void;
  title: any;
  submitTitle: any;
  cancelTitle: any;
}

export const ConfirmPopup = ({
  isOpen,
  onClose,
  title,
  cancelTitle,
  submitTitle,
  onSubmitClick,
  onCancelClick,
}: PopupProps) => {
  return (
    <PopupModal isOpen={isOpen} onClose={onClose}>
      <div className='alertWrapper alertConfirmWrapper'>{title}</div>
      <div className='alertButtonsWrapper'>
        <Button onClick={onCancelClick}>{cancelTitle}</Button>
        <div className='separator'></div>
        <Button onClick={onSubmitClick}>{submitTitle}</Button>
      </div>
    </PopupModal>
  );
};
