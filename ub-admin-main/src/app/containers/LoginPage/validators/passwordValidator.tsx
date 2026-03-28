import { MessageService, MessageNames } from 'services/messageService';

export const PasswordValidator = (properties: {
  uniqueInputId: string;
  value: string;
  errors: string[];
}) => {
  let value = properties.value;
  let inputId = properties.uniqueInputId;
  const sendError = (error: string | null) => {
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: inputId,
      payload: error,
    });
  };
  if (value == null || value.length === 0) {
    sendError(properties.errors[0]);
    return false;
  } else if (value.length < 12) {
    sendError(properties.errors[1]);
    return false;
  } else if (!/[A-Z]/.test(value) || !/[a-z]/.test(value) || !/[0-9]/.test(value)) {
    sendError(properties.errors[1]);
    return false;
  } else {
    sendError(null);
    return true;
  }
};
