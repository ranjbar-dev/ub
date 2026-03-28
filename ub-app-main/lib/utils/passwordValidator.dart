String validatePassword(String password) {
  bool hasUppercase = password.contains(new RegExp(r'[A-Z]'));
  bool hasDigits = password.contains(new RegExp(r'[0-9]'));
  bool hasLowercase = password.contains(new RegExp(r'[a-z]'));
  bool hasSpecialCharacters =
      password.contains(new RegExp(r'[!@#$%^&*(),.?":{}|<>]'));
  bool hasMinLength = password.length >= 8;
  if (password == null || password.isEmpty) {
    return '*Required';
  }
  if (!hasUppercase) {
    return 'At least 1 Uppercase char';
  }
  if (!hasDigits) {
    return 'At least 1 digit';
  }
  if (!hasLowercase) {
    return 'At least 1 lowercase';
  }
  if (!hasSpecialCharacters) {
    return 'At least 1 special character( @#\$...)';
  }
  if (!hasMinLength) {
    return 'At least 8 character';
  }
  return '';
}
