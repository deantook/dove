import 'package:flutter/material.dart';

enum LikeDirection { likedMe, iLiked }

class LikeProfile {
  const LikeProfile({
    required this.id,
    required this.name,
    required this.subtitle,
    required this.timeLabel,
    required this.direction,
    required this.avatarColors,
  });

  final String id;
  final String name;
  final String subtitle;
  final String timeLabel;
  final LikeDirection direction;
  final List<Color> avatarColors;
}

const mockLikeProfiles = [
  LikeProfile(
    id: 'l1',
    name: 'Emma',
    subtitle: '赞了你的 Prompt',
    timeLabel: '2 小时前',
    direction: LikeDirection.likedMe,
    avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
  ),
  LikeProfile(
    id: 'l2',
    name: 'Sophia',
    subtitle: '赞了你的照片',
    timeLabel: '5 小时前',
    direction: LikeDirection.likedMe,
    avatarColors: [Color(0xFFA8DADC), Color(0xFF457B9D)],
  ),
  LikeProfile(
    id: 'l3',
    name: 'Olivia',
    subtitle: '赞了你的 Prompt',
    timeLabel: '昨天',
    direction: LikeDirection.likedMe,
    avatarColors: [Color(0xFFD4A574), Color(0xFF8B6914)],
  ),
  LikeProfile(
    id: 'l4',
    name: '林溪',
    subtitle: '你赞了对方的 Prompt',
    timeLabel: '3 小时前',
    direction: LikeDirection.iLiked,
    avatarColors: [Color(0xFFB8B8D1), Color(0xFF6B6B8D)],
  ),
  LikeProfile(
    id: 'l5',
    name: 'Zoe',
    subtitle: '你赞了对方的照片',
    timeLabel: '昨天',
    direction: LikeDirection.iLiked,
    avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
  ),
  LikeProfile(
    id: 'l6',
    name: 'Mia',
    subtitle: '你赞了对方的 Prompt',
    timeLabel: '2 天前',
    direction: LikeDirection.iLiked,
    avatarColors: [Color(0xFF95D5B2), Color(0xFF40916C)],
  ),
];
