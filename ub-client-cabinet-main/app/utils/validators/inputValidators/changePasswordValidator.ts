import { ChangePasswordData } from 'containers/ChangePassword/constants';

// var strongRegex = new RegExp(
//   '^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})',
// );

export const ChangePasswordValidator = (
  values: ChangePasswordData,
): ChangePasswordData => {
  const errors: ChangePasswordData = {
    oldPassword: '',
    newPassword: '',
    confirmNewPassword: '',
  };
  if (!values.oldPassword) {
    errors.oldPassword = 'enter old password';
  }

  if (!values.newPassword) {
    errors.newPassword = 'enter new password';
  }

  if (!values.confirmNewPassword) {
    errors.confirmNewPassword = 'please cunfirm your new password';
  }
  if (values.newPassword !== values.confirmNewPassword) {
    errors.confirmNewPassword = 'old and new password must be equal';
    errors.newPassword = 'old and new password must be equal';
  }
  // if (values.newPassword) {
  //   if (!strongRegex.test(values.newPassword)) {
  //     errors.newPassword =
  //       'new password must be at least 8 characters containing one special character one number and one capital letter';
  //   }
  // }

  return errors;
};
