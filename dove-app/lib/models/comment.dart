import 'package:flutter/material.dart';

class Comment {
  const Comment({
    required this.id,
    required this.authorName,
    required this.content,
    required this.timeLabel,
    required this.avatarColors,
  });

  final String id;
  final String authorName;
  final String content;
  final String timeLabel;
  final List<Color> avatarColors;
}

const currentUserAvatarColors = [Color(0xFFD4C4B0), Color(0xFF5C7A6B)];

final mockCommentsByMomentId = <String, List<Comment>>{
  's1': [
    Comment(
      id: 'c1',
      authorName: 'Emma',
      content: '那家书店叫什么？我也想去看看。',
      timeLabel: '1 小时前',
      avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
    ),
    Comment(
      id: 'c2',
      authorName: '阿杰',
      content: '绝版摄影集太羡慕了！',
      timeLabel: '45 分钟前',
      avatarColors: [Color(0xFF95D5B2), Color(0xFF40916C)],
    ),
  ],
  'f1': [
    Comment(
      id: 'c3',
      authorName: 'Mia',
      content: '留一块给我！',
      timeLabel: '30 分钟前',
      avatarColors: [Color(0xFFD4A574), Color(0xFF8B6914)],
    ),
  ],
  'f2': [
    Comment(
      id: 'c4',
      authorName: 'Emma',
      content: '太美了，下次一起！',
      timeLabel: '2 小时前',
      avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
    ),
    Comment(
      id: 'c5',
      authorName: 'Alex',
      content: '算我一个 👋',
      timeLabel: '1 小时前',
      avatarColors: currentUserAvatarColors,
    ),
  ],
};
