import { loginData } from '../../hooks/useFormValidation';

export const LoginValidator = (values: loginData): loginData => {
  const errors: loginData = {
    email: '',
    password: '',
  };
  if (!values.email) {
    errors.email = 'enter email';
  }

  if (values.email && values.email.length < 5) {
    errors.email = 'email is not valid';
  }

  if (!values.password) {
    errors.password = 'enter password';
  }

  if (values.password && values.password.length < 4) {
    errors.password = 'minimum 4 character';
  }



  return errors;
};
