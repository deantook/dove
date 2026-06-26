import 'package:flutter/material.dart';

enum ConversationStatus { yourTurn, theirTurn }

class Conversation {
  const Conversation({
    required this.id,
    required this.name,
    required this.lastMessage,
    required this.timeLabel,
    required this.status,
    required this.avatarColors,
  });

  final String id;
  final String name;
  final String lastMessage;
  final String timeLabel;
  final ConversationStatus status;
  final List<Color> avatarColors;
}

const mockConversations = [
  Conversation(
    id: '1',
    name: 'Emma',
    lastMessage: '那家咖啡馆听起来不错，周末有空吗？',
    timeLabel: '2 分钟前',
    status: ConversationStatus.yourTurn,
    avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
  ),
  Conversation(
    id: '2',
    name: 'Sophia',
    lastMessage: '哈哈，我也超爱徒步！',
    timeLabel: '昨天',
    status: ConversationStatus.yourTurn,
    avatarColors: [Color(0xFFA8DADC), Color(0xFF457B9D)],
  ),
  Conversation(
    id: '3',
    name: 'Olivia',
    lastMessage: '期待下次见面 ☺️',
    timeLabel: '周二',
    status: ConversationStatus.theirTurn,
    avatarColors: [Color(0xFFD4A574), Color(0xFF8B6914)],
  ),
  Conversation(
    id: '4',
    name: 'Ava',
    lastMessage: '刚看完你推荐的那本书',
    timeLabel: '3 天前',
    status: ConversationStatus.theirTurn,
    avatarColors: [Color(0xFFB8B8D1), Color(0xFF6B6B8D)],
  ),
  Conversation(
    id: '5',
    name: 'Mia',
    lastMessage: '晚安，做个好梦',
    timeLabel: '1 周前',
    status: ConversationStatus.theirTurn,
    avatarColors: [Color(0xFF95D5B2), Color(0xFF40916C)],
  ),
];
