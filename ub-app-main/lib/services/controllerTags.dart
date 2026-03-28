class ControllerTags {
  static final ControllerTags _singleton = ControllerTags._internal();
  factory ControllerTags() {
    return _singleton;
  }
  ControllerTags._internal();

  static final login = 'login';
}
