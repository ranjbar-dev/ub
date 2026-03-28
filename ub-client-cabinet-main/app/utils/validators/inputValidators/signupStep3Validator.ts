export const SignupStep3Validator = values => {
  const errors = {
    email: '',
    drugStoreName: '',

  };
  if (!values.email) {
    errors.email = 'آدرس ایمیل الزامی است';
  }

  if (values.email && (values.email.length < 5 || !values.email.includes('@'))) {
    errors.email = 'ایمیل نامعتبر';
  }

  if (!values.drugStoreName) {
    errors.drugStoreName = 'لطفا نام داروخانه را وارد کنید';
  }

  return errors;
};
