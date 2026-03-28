import 'package:flutter/cupertino.dart';
import 'package:get/get.dart';

// this is just a test,and is not implemented anywhere
const isAuthenticated = false;

class AuthMiddleware extends GetMiddleware {
  RouteSettings redirect(String route) {
    String returnUrl = Uri.encodeFull(route ?? '');
    return !isAuthenticated
        ? RouteSettings(name: "/login?return=" + returnUrl)
        : null;
  }
}
