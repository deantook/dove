import 'package:flutter/material.dart';

import '../../models/conversation.dart';
import '../../theme/app_colors.dart';
import 'messages/widgets/conversation_tile.dart';

class MessagesScreen extends StatelessWidget {
  const MessagesScreen({super.key});

  List<Conversation> _filter(ConversationStatus status) {
    return mockConversations
        .where((conversation) => conversation.status == status)
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    final yourTurn = _filter(ConversationStatus.yourTurn);
    final theirTurn = _filter(ConversationStatus.theirTurn);

    return SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: ListView(
              padding: const EdgeInsets.fromLTRB(20, 16, 20, 24),
              children: [
                if (yourTurn.isNotEmpty) ...[
                  const ConversationSectionHeader(title: '等你回复'),
                  ...yourTurn.map(
                    (conversation) => ConversationTile(
                      key: ValueKey(conversation.id),
                      conversation: conversation,
                      onTap: () {},
                    ),
                  ),
                ],
                if (yourTurn.isNotEmpty && theirTurn.isNotEmpty)
                  const Padding(
                    padding: EdgeInsets.symmetric(vertical: 8),
                    child: Divider(
                      height: 1,
                      color: AppColors.divider,
                    ),
                  ),
                if (theirTurn.isNotEmpty) ...[
                  const ConversationSectionHeader(title: '进行中'),
                  ...theirTurn.map(
                    (conversation) => ConversationTile(
                      key: ValueKey(conversation.id),
                      conversation: conversation,
                      onTap: () {},
                    ),
                  ),
                ],
                const SizedBox(height: 8),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
