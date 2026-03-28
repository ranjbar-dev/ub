export enum PasswordErrors {
  Min8 = 'Min8',
  NoNumber = 'NoNumber',
  NoUppercase = 'NoUppercase',
  NoSpecialCharacter = 'NoSpicialCharacter',
  NoError = 'NoError',
}
const hasNumberReg = /\d/;
const hasUppercaseReg = /[A-Z]/;
const hasSpicialCharacterReg = /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;
export const passwordContentValidator = ({
  password,
}: {
  password: string;
}) => {
  if (password.length < 8) {
    return PasswordErrors.Min8;
  }
  if (!hasNumberReg.test(password)) {
    return PasswordErrors.NoNumber;
  }
  if (!hasUppercaseReg.test(password)) {
    return PasswordErrors.NoUppercase;
  }
  if (!hasSpicialCharacterReg.test(password)) {
    return PasswordErrors.NoSpecialCharacter;
  }

  return PasswordErrors.NoError;
};
