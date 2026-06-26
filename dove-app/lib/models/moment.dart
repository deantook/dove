import 'package:flutter/material.dart';

import '../data/mock_image_urls.dart';

enum MomentSource { stranger, friend }

class Moment {
  const Moment({
    required this.id,
    required this.authorName,
    required this.content,
    required this.timeLabel,
    required this.source,
    required this.avatarColors,
    this.hasImage = false,
    this.imageUrl,
    this.likeCount = 0,
  });

  final String id;
  final String authorName;
  final String content;
  final String timeLabel;
  final MomentSource source;
  final List<Color> avatarColors;
  final bool hasImage;
  final String? imageUrl;
  final int likeCount;

  String? get resolvedImageUrl =>
      imageUrl ?? (hasImage ? MockImageUrls.momentFor(id) : null);
}

const mockMoments = [
  Moment(
    id: 's1',
    authorName: '林溪',
    content: '周末去了城西那家独立书店，意外发现一本绝版摄影集，待了一整个下午。',
    timeLabel: '2 小时前',
    source: MomentSource.stranger,
    avatarColors: [Color(0xFFB8B8D1), Color(0xFF6B6B8D)],
    hasImage: true,
    likeCount: 24,
  ),
  Moment(
    id: 's2',
    authorName: '陈朗',
    content: '今天滑板的第三个后空翻终于成了，摔了俩礼拜也算值了。',
    timeLabel: '5 小时前',
    source: MomentSource.stranger,
    avatarColors: [Color(0xFFA8DADC), Color(0xFF457B9D)],
    likeCount: 18,
  ),
  Moment(
    id: 's3',
    authorName: 'Zoe',
    content: '有人一起去看下周的爵士乐 live 吗？想找同好。',
    timeLabel: '昨天',
    source: MomentSource.stranger,
    avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
    likeCount: 9,
  ),
  Moment(
    id: 'f1',
    authorName: 'Emma',
    content: '刚烤的 sourdough，外脆里软，谁要分一块 🍞',
    timeLabel: '1 小时前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
    hasImage: true,
    likeCount: 12,
  ),
  Moment(
    id: 'f2',
    authorName: '阿杰',
    content: '徒步回来的 sunset，值了。下次一起？',
    timeLabel: '3 小时前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFF95D5B2), Color(0xFF40916C)],
    hasImage: true,
    likeCount: 31,
  ),
  Moment(
    id: 'f3',
    authorName: 'Mia',
    content: '新书到了，今晚开读。最近迷上推理小说。',
    timeLabel: '昨天',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4A574), Color(0xFF8B6914)],
    likeCount: 6,
  ),
];
