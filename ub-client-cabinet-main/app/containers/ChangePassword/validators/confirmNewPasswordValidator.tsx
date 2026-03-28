import React from 'react';
import translate from './messages';
import { MessageService, MessageNames } from 'services/message_service';
import { FormattedMessage } from 'react-intl';

export const ConfirmNewPasswordValidator = (properties: {
  uniqueInputId: string;
  value1: string;
  value2: string;
}) => {
  const { value1, value2 } = properties;
  const inputId = properties.uniqueInputId;
  const sendError = (error: any) => {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: error,
    });
  };

  if (value1 !== value2) {
    sendError(<FormattedMessage {...translate.oldAndNewPassword} />);
    return false;
  } else {
    sendError(null);
    return true;
  }
};
