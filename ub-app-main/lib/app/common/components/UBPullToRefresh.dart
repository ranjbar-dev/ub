import 'package:flutter/material.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import '../custom/refreshHeader.dart';

class UBPullToRefresh extends StatelessWidget {
  final Widget child;
  final RefreshController controller;
  final Function onRefresh;
  final String updatingText;
  final String releaseToUpdateText;
  final String beforeUpdateText;
  final String afterUpdateText;
  final bool withUpdateDate;
  final bool showUpdateText;

  const UBPullToRefresh({
    Key key,
    this.child,
    this.controller,
    this.onRefresh,
    this.afterUpdateText,
    this.withUpdateDate,
    this.updatingText,
    this.beforeUpdateText,
    this.releaseToUpdateText,
    this.showUpdateText,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    RefreshController _controller = controller;
    if (_controller == null) {
      _controller = RefreshController();
    }
    return SmartRefresher(
        header: UBPullToRefreshHeader(
          controller: _controller,
          afterUpdateText: afterUpdateText,
          beforeUpdateText: beforeUpdateText,
          updatingText: updatingText,
          withUpdateDate: withUpdateDate,
          releaseToUpdateText: releaseToUpdateText,
          showUpdateText: showUpdateText,
        ),
        //header: const MaterialClassicHeader(
        //  height: 60,
        //  backgroundColor: ColorName.white,
        //  color: ColorName.primaryBlue,
        //),
        controller: _controller,
        onRefresh: onRefresh,
        child: child);
  }
}
