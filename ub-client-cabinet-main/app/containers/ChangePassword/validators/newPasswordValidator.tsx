import React from 'react';
import translate from './messages';
import { MessageService, MessageNames } from 'services/message_service';
import { FormattedMessage } from 'react-intl';
const strongRegex = new RegExp(
  '^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})',
);
export const NewPasswordValidator = (properties: {
  uniqueInputId: string;
  value: string;
  compareTo?: string;
  compareId?: string;
}) => {
  const inputId = properties.uniqueInputId;
  const { compareTo, compareId, value } = properties;
  const sendError = (error: any, isLong?: boolean) => {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: error,
      additional: isLong,
    });
  };
  const sendErrorToCompare = (error: any, isLong?: boolean) => {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: compareId,
      payload: error,
      additional: isLong,
    });
  };
  if (value == null || value.length === 0) {
    sendError(<FormattedMessage {...translate.passwordIsRequired} />);
    return false;
  } else if (value.length < 8) {
    sendError(<FormattedMessage {...translate.minimum8Character} />);
    return false;
  } else if (!strongRegex.test(value)) {
    sendError(<FormattedMessage {...translate.strongPasswordError} />, true);
    return false;
  } else if (value !== compareTo && compareTo !== '') {
    sendError(<FormattedMessage {...translate.oldAndNewPassword} />, true);
    return false;
  } else {
    sendError(null);
    if (compareTo !== '') {
      sendErrorToCompare(null);
    }
    return true;
  }
};
