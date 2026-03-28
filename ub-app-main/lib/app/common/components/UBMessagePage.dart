import 'package:flutter/material.dart';

import 'pageContainer.dart';

class UBMessagePage extends StatelessWidget {
  final List<Widget> children;
  final double spaceBetween;

  const UBMessagePage(
      {Key key, @required this.children, this.spaceBetween = 24.0})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    int childNum = children.length;
    return PageContainer(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          for (var i = 0; i < childNum; i++)
            if (i < childNum - 1)
              Container(
                child: Column(
                  children: [
                    children[i],
                    SizedBox(
                      height: spaceBetween,
                    )
                  ],
                ),
              )
            else
              children[i]
        ],
      ),
    );
  }
}
