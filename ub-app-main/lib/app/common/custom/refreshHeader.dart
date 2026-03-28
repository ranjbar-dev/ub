import 'package:flutter/material.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:rive/rive.dart';
import '../components/UBText.dart';
import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/commonConsts.dart';

class UBPullToRefreshHeader extends StatefulWidget {
  final RefreshController controller;
  final String updatingText;
  final String releaseToUpdateText;
  final String beforeUpdateText;
  final String afterUpdateText;
  final bool withUpdateDate;
  final bool showUpdateText;

  const UBPullToRefreshHeader({
    Key key,
    this.controller,
    this.updatingText,
    this.beforeUpdateText,
    this.afterUpdateText,
    this.withUpdateDate,
    this.releaseToUpdateText,
    this.showUpdateText = true,
  }) : super(key: key);
  @override
  State<StatefulWidget> createState() {
    return _UBPullToRefreshHeaderState();
  }
}

class _UBPullToRefreshHeaderState extends State<UBPullToRefreshHeader>
    with TickerProviderStateMixin {
  AnimationController _scaleController;
  RefreshController _refreshController;
  RiveAnimationController _riveAnimationController;
  Artboard _artboard;
  String text;

  //SMITrigger _start;
  //SMITrigger _loading;
  //SMITrigger _finish;

  //void _onRiveInit(Artboard artboard) {
  //  final controller = StateMachineController.fromArtboard(artboard, 'State');
  //  artboard.addController(controller);
  //  _start = controller.findInput<bool>('Start') as SMITrigger;
  //  _loading = controller.findInput<bool>('Loading') as SMITrigger;
  //  _finish = controller.findInput<bool>('Finish') as SMITrigger;
  //}

  void _begin() {
    _artboard?.removeController(_riveAnimationController);
    _artboard
        ?.addController(_riveAnimationController = SimpleAnimation('Start'));
  }

  void _startLoading() {
    _artboard?.removeController(_riveAnimationController);
    _artboard
        ?.addController(_riveAnimationController = SimpleAnimation('Loading'));
  }

  void _finishLoading() {
    _artboard?.removeController(_riveAnimationController);
    _artboard
        ?.addController(_riveAnimationController = SimpleAnimation('Finish'));
  }

  _simpleRiveInit(Artboard artboard) {
    artboard.addController(_riveAnimationController = SimpleAnimation('Start'));
    setState(() {
      _artboard = artboard;
    });
  }

  @override
  void initState() {
    _refreshController = widget.controller;

    _scaleController =
        AnimationController(value: 0.0, vsync: this, upperBound: 1.0);
    _refreshController.headerMode.addListener(() {
      final status = _refreshController.headerStatus;
      String t = this.widget.showUpdateText == false
          ? ''
          : this.widget.beforeUpdateText;
      if (status == RefreshStatus.idle) {
        t = widget.beforeUpdateText;
        _begin();
      } else if (status == RefreshStatus.refreshing) {
        t = widget.updatingText;
        _startLoading();
      } else if (status == RefreshStatus.canRefresh) {
        t = widget.releaseToUpdateText;
        //_startLoading();
      } else if (status == RefreshStatus.completed) {
        t = widget.afterUpdateText;
        if (widget.withUpdateDate == true) {
          final now = DateTime.now();
          t = t + ' ' + now.toString().split(' ')[1].split('.')[0];
        }
        _finishLoading();
      }
      if (widget.showUpdateText == true) {
        setState(() {
          text = t;
        });
      }

      if (_refreshController.headerStatus == RefreshStatus.idle) {
        _scaleController.value = 0.0;
      } else if (_refreshController.headerStatus == RefreshStatus.refreshing) {}
    });
    super.initState();
  }

  @override
  void dispose() {
    _refreshController.dispose();
    _scaleController.dispose();
    _riveAnimationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return CustomHeader(
      height: 65,
      refreshStyle: RefreshStyle.Behind,
      onOffsetChange: (offset) {
        if (_refreshController.headerMode.value != RefreshStatus.refreshing)
          _scaleController.value = offset / 40.0;
      },
      builder: (c, m) {
        return Container(
          child: FadeTransition(
            opacity: _scaleController,
            child: ScaleTransition(
              child: Column(
                children: [
                  vspace12,
                  Container(
                    height: 35.0,
                    child: RiveAnimation.asset(
                      'assets/rive/pullToRefreshLoading.riv',
                      //stateMachines: ['State'],
                      onInit: _simpleRiveInit,
                    ),
                  ),
                  UBText(
                    text: text,
                    color: ColorName.grey80,
                    size: 11.0,
                  )
                ],
              ),
              scale: _scaleController,
            ),
          ),
          alignment: Alignment.center,
        );
      },
    );
  }
}
