import 'package:flutter/material.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:supercharged/supercharged.dart';

import '../../../generated/colors.gen.dart';
import 'UBPullToRefresh.dart';
import 'UBScrollColumnExpandable.dart';
import 'UBoops.dart';

class PullToRefreshConfig {
  final RefreshController controller;
  final Function onRefreshLoading;
  final bool isLoading;
  final OopsVariant oopsVariant;
  final String updatingText;
  final String releaseToUpdateText;
  final String beforeUpdateText;
  final String afterUpdateText;
  final bool withUpdateDate;
  final bool showUpdateText;

  PullToRefreshConfig({
    this.oopsVariant,
    this.controller,
    this.onRefreshLoading,
    this.isLoading,
    this.showUpdateText = true,
    this.releaseToUpdateText = 'Release to update',
    this.afterUpdateText = 'Updated ',
    this.withUpdateDate = false,
    this.updatingText = 'Updating...',
    this.beforeUpdateText = 'Pull to update...',
  });
}

class UBScrollBar extends StatelessWidget {
  final int itemCount;
  final PullToRefreshConfig pullToRefreshConfig;
  final Widget Function(BuildContext context, int index) builder;
  final ScrollController scrollController;
  final Axis scrollDirection;

  const UBScrollBar({
    Key key,
    @required this.itemCount,
    @required this.builder,
    this.scrollController,
    this.pullToRefreshConfig,
    this.scrollDirection = Axis.vertical,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    ScrollController sController = scrollController ?? ScrollController();

    if (pullToRefreshConfig != null && pullToRefreshConfig.isLoading == false) {
      Future.delayed(100.milliseconds)
          .then((value) => pullToRefreshConfig.controller.refreshCompleted());
    }
    final child = itemCount == 0
        ? UBScrollColumnExpandable(
            children: [
              Expanded(
                child: UBoops(
                  variant: pullToRefreshConfig.oopsVariant,
                ),
              ),
            ],
          )
        : ListView.builder(
            scrollDirection: scrollDirection,
            shrinkWrap: true,
            controller: sController,
            itemCount: itemCount,
            itemBuilder: builder,
          );
    return RawScrollbar(
      controller: sController,
      fadeDuration: 200.milliseconds,
      radius: const Radius.circular(12),
      thumbColor: ColorName.grey23,
      thickness: 4.0,
      child: pullToRefreshConfig == null
          ? child
          : UBPullToRefresh(
              controller: pullToRefreshConfig.controller,
              onRefresh: pullToRefreshConfig.onRefreshLoading,
              afterUpdateText: pullToRefreshConfig.afterUpdateText,
              beforeUpdateText: pullToRefreshConfig.beforeUpdateText,
              releaseToUpdateText: pullToRefreshConfig.releaseToUpdateText,
              showUpdateText: pullToRefreshConfig.showUpdateText,
              updatingText: pullToRefreshConfig.updatingText,
              withUpdateDate: pullToRefreshConfig.withUpdateDate,
              child: child,
            ),
    );
  }
}
