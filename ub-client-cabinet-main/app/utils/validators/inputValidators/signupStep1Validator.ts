export const SignupStep1Validator = values => {
  const errors = {
    hix: '',
    password: '',
    phone: '',
    userName: '',
  };
  if (!values.hix) {
    errors.hix = 'کد HIX الزامی است';
  }

  if (values.hix && values.hix.length < 5) {
    errors.hix = 'کد HIX معتبر نیست';
  }

  if (!values.password) {
    errors.password = 'لطفا کلمه عبور را وارد کنید';
  }

  if (values.password && values.password.length < 4) {
    errors.password = 'حداقل 4 کاراکتر';
  }
  if (!values.phone) {
    errors.phone = 'لطفا تلفن همراه معتبر وارد کنید';
  }

  if (values.phone && values.phone.length < 11) {
    errors.phone = 'تلفن همراه معتبر نیست(مثال:09121234567)';
  }

  if (!values.userName) {
    errors.userName = 'نام کاربری الزامی است';
  }

  if (values.userName && values.userName.length < 5) {
    errors.userName = 'حداقل 5 کاراکتر';
  }

  return errors;
};
