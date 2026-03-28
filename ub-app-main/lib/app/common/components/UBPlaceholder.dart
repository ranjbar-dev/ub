import 'dart:math';

import 'package:flutter/material.dart';

class UBPlaceholder extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    List colors = [
      Colors.red,
      Colors.green,
      Colors.yellow,
      Colors.amber,
      Colors.blue,
      Colors.cyan,
      Colors.deepOrange,
      Colors.deepPurple,
      Colors.lightBlue,
      Colors.pink,
      Colors.tealAccent,
      Colors.grey,
    ];
    Random random = new Random();
    return Placeholder(
      color: colors[random.nextInt(colors.length)],
    );
  }
}
