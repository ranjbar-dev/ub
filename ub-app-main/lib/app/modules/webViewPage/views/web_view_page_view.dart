import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../controllers/web_view_page_controller.dart';
import 'dart:io';
import 'package:webview_flutter/webview_flutter.dart';

class WebViewPageView extends GetView<WebViewPageController> {
  final String url;
  final String title;

  WebViewPageView({this.title, this.url});
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          title ?? '',
          style: whiteBold14,
        ),
        centerTitle: true,
      ),
      body: WebViewWidget(
        url: url,
      ),
    );
  }
}

class WebViewWidget extends StatefulWidget {
  final String url;

  const WebViewWidget({Key key, this.url}) : super(key: key);
  @override
  WebViewWidgetState createState() => WebViewWidgetState();
}

class WebViewWidgetState extends State<WebViewWidget> {
  @override
  void initState() {
    super.initState();
    // Enable hybrid composition.
    if (Platform.isAndroid) WebView.platform = SurfaceAndroidWebView();
  }

  @override
  Widget build(BuildContext context) {
    return WebView(
      initialUrl: widget.url,
      javascriptMode: JavascriptMode.unrestricted,
    );
  }
}
