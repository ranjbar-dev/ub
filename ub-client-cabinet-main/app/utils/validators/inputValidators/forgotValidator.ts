export const ForgotValidator = values => {
  const errors = {
    hix: '',
  };
  if (!values.hix) {
    errors.hix = 'لطفا نام کاربری را وارد کنید';
  }

  if (values.hix && values.hix.length < 5) {
    errors.hix = 'نام کاربری معتبر نیست';
  }
  return errors;
};
