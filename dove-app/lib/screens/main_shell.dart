import 'package:flutter/material.dart';

import '../widgets/dove_bottom_nav_bar.dart';
import '../widgets/nav_icons/dove_nav_icon.dart';
import 'discover_screen.dart';
import 'likes_screen.dart';
import 'messages_screen.dart';
import 'profile_screen.dart';
import 'standouts_screen.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _currentIndex = 0;

  static const _iconSize = 24.0;

  static final _navItems = [
    DoveNavItem(
      semanticLabel: '发现',
      iconBuilder: (_, color) => DoveNavIcon(
        assetPath: NavIconAssets.discover,
        color: color,
        size: _iconSize,
      ),
    ),
    DoveNavItem(
      semanticLabel: '精选',
      iconBuilder: (_, color) => DoveNavIcon(
        assetPath: NavIconAssets.standouts,
        color: color,
        size: _iconSize,
      ),
    ),
    DoveNavItem(
      semanticLabel: '喜欢',
      iconBuilder: (_, color) => DoveNavIcon(
        assetPath: NavIconAssets.likes,
        color: color,
        size: _iconSize,
      ),
    ),
    DoveNavItem(
      semanticLabel: '消息',
      iconBuilder: (_, color) => DoveNavIcon(
        assetPath: NavIconAssets.messages,
        color: color,
        size: _iconSize,
      ),
    ),
    DoveNavItem(
      semanticLabel: '我的',
      iconBuilder: (_, color) => DoveNavIcon(
        assetPath: NavIconAssets.profile,
        color: color,
        size: _iconSize,
      ),
    ),
  ];

  static const _screens = [
    DiscoverScreen(),
    StandoutsScreen(),
    LikesScreen(),
    MessagesScreen(),
    ProfileScreen(isEmbedded: true),
  ];

  void _onTabSelected(int index) {
    if (index == _currentIndex) return;
    setState(() => _currentIndex = index);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(index: _currentIndex, children: _screens),
      bottomNavigationBar: DoveBottomNavBar(
        currentIndex: _currentIndex,
        items: _navItems,
        onTap: _onTabSelected,
      ),
    );
  }
}
