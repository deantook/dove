import 'package:flutter/material.dart';

import '../models/like_profile.dart';
import '../widgets/segment_tab_bar.dart';
import 'likes/widgets/like_profile_list.dart';

class LikesScreen extends StatefulWidget {
  const LikesScreen({super.key});

  @override
  State<LikesScreen> createState() => _LikesScreenState();
}

class _LikesScreenState extends State<LikesScreen> {
  int _selectedTab = 0;

  static const _tabs = ['谁喜欢了我', '我喜欢了谁'];

  List<LikeProfile> _profilesFor(LikeDirection direction) {
    return mockLikeProfiles
        .where((profile) => profile.direction == direction)
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    final likedMe = _profilesFor(LikeDirection.likedMe);
    final iLiked = _profilesFor(LikeDirection.iLiked);

    return SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 16, 20, 0),
            child: SegmentTabBar(
              tabs: _tabs,
              selectedIndex: _selectedTab,
              onChanged: (index) => setState(() => _selectedTab = index),
            ),
          ),
          const SizedBox(height: 16),
          Expanded(
            child: IndexedStack(
              index: _selectedTab,
              children: [
                LikeProfileList(
                  profiles: likedMe,
                  emptyMessage: '还没有人喜欢你，继续探索吧',
                ),
                LikeProfileList(
                  profiles: iLiked,
                  emptyMessage: '你还没有表达喜欢',
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
