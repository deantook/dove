import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:dove/main.dart';

void main() {
  testWidgets('底部导航切换页面', (WidgetTester tester) async {
    await tester.pumpWidget(const DoveApp());

    expect(find.text('实时匹配'), findsOneWidget);
    expect(find.text('推荐对象'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('nav_tab_3')));
    await tester.pumpAndSettle();

    expect(find.text('Emma'), findsOneWidget);
    expect(find.text('等你回复'), findsWidgets);
  });
}
