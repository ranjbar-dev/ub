import React from 'react';
import translate from './messages';
import { MessageService, MessageNames } from 'services/message_service';
import { FormattedMessage } from 'react-intl';
const emailtext = new RegExp(
  '^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:.[a-zA-Z0-9-]+)*$',
);
export const EmailValidator = (properties: {
  uniqueInputId: string;
  value: string;
}) => {
  const value = properties.value;
  const inputId = properties.uniqueInputId;
  const sendError = (error: any) => {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: error,
    });
  };
  if (value == null || value.length === 0) {
    sendError(<FormattedMessage {...translate.emailIsRequired} />);
    return false;
  } else if (!emailtext.test(value) || !value.includes('.')) {
    sendError(<FormattedMessage {...translate.emailIsNotValid} />);
    return false;
  } else {
    sendError(null);
    return true;
  }
};
