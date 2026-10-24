import 'dart:math' as math;

import 'package:flutter/material.dart';

/// Web/宽屏时限制内容宽度，接近真机观感。
///
/// [expand] 为 false 时用于 bottomNavigationBar：不可纵向撑满，否则会把 body 高度挤成 0。
class AppPageFrame extends StatelessWidget {
  const AppPageFrame({
    super.key,
    required this.child,
    this.maxWidth = 480,
    this.expand = true,
  });

  final Widget child;
  final double maxWidth;
  final bool expand;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final maxW = constraints.maxWidth;
        final width = maxW.isFinite ? math.min(maxW, maxWidth) : maxWidth;
        final framed = SizedBox(
          width: width,
          child: child,
        );

        if (!expand) {
          // 底栏：只占自身高度，水平居中
          return Align(
            alignment: Alignment.bottomCenter,
            heightFactor: 1,
            child: framed,
          );
        }

        // 页面主体：占满可用高度，避免 Column+Expanded 收到无限高度
        final height = constraints.hasBoundedHeight && constraints.maxHeight.isFinite
            ? constraints.maxHeight
            : MediaQuery.sizeOf(context).height;
        return Align(
          alignment: Alignment.topCenter,
          child: SizedBox(
            width: width,
            height: height,
            child: child,
          ),
        );
      },
    );
  }
}
