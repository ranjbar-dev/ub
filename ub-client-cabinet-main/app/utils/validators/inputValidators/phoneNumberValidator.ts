export const PhoneNumberValidator = values => {
  const errors = {
    phoneNumber: '',
  };
  if (!values.phoneNumber) {
    errors.phoneNumber = 'Please Enter Phone Number';
  }

  if (values.phoneNumber && values.phoneNumber.length < 5) {
    errors.phoneNumber = 'Phone Number Is Not Valid';
  }
  return errors;
};
