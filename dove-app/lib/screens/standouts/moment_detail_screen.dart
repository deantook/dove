import 'package:flutter/material.dart';

import '../../models/app_user.dart';
import '../../models/comment.dart';
import '../../models/moment.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_typography.dart';
import '../../widgets/moment_image.dart';
import '../standouts/widgets/comment_section.dart';
import '../standouts/widgets/moment_avatar.dart';

class MomentDetailScreen extends StatefulWidget {
  const MomentDetailScreen({
    super.key,
    required this.moment,
  });

  final Moment moment;

  static Route<void> route(Moment moment) {
    return MaterialPageRoute<void>(
      builder: (_) => MomentDetailScreen(moment: moment),
    );
  }

  @override
  State<MomentDetailScreen> createState() => _MomentDetailScreenState();
}

class _MomentDetailScreenState extends State<MomentDetailScreen> {
  late List<Comment> _comments;
  final _commentController = TextEditingController();
  final _commentFocusNode = FocusNode();
  int _commentIdCounter = 1000;

  Moment get moment => widget.moment;

  @override
  void initState() {
    super.initState();
    _comments = List<Comment>.from(
      mockCommentsByMomentId[moment.id] ?? const [],
    );
  }

  @override
  void dispose() {
    _commentController.dispose();
    _commentFocusNode.dispose();
    super.dispose();
  }

  void _submitComment() {
    final text = _commentController.text.trim();
    if (text.isEmpty) return;

    setState(() {
      _comments.add(
        Comment(
          id: 'new_${_commentIdCounter++}',
          authorName: AppUser.currentUser.name,
          content: text,
          timeLabel: '刚刚',
          avatarColors: AppUser.currentUser.avatarColors,
        ),
      );
      _commentController.clear();
    });
    _commentFocusNode.unfocus();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_new_rounded, size: 20),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Text(moment.authorName, style: AppTypography.cardHeadline),
        centerTitle: true,
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.fromLTRB(20, 8, 20, 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      MomentAvatar(
                        name: moment.authorName,
                        colors: moment.avatarColors,
                        size: 48,
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              moment.authorName,
                              style: AppTypography.cardHeadline,
                            ),
                            Text(moment.timeLabel, style: AppTypography.caption),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),
                  Text(moment.content, style: AppTypography.body),
                  if (moment.resolvedImageUrl != null) ...[
                    const SizedBox(height: 24),
                    MomentImage(
                      imageUrl: moment.resolvedImageUrl!,
                      aspectRatio: 4 / 3,
                      fallbackColors: [
                        moment.avatarColors.first.withValues(alpha: 0.6),
                        moment.avatarColors.last.withValues(alpha: 0.8),
                      ],
                    ),
                  ],
                  const SizedBox(height: 24),
                  Row(
                    children: [
                      Icon(
                        Icons.favorite_border_rounded,
                        size: 22,
                        color: AppColors.textTertiary,
                      ),
                      const SizedBox(width: 6),
                      Text(
                        '${moment.likeCount} 人觉得很赞',
                        style: AppTypography.caption,
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),
                  Text(
                    '评论 ${_comments.length}',
                    style: AppTypography.cardHeadline,
                  ),
                  const SizedBox(height: 16),
                  if (_comments.isEmpty)
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: 24),
                      child: Center(
                        child: Text(
                          '还没有评论，来说第一句吧',
                          style: AppTypography.body.copyWith(
                            color: AppColors.textSecondary,
                          ),
                        ),
                      ),
                    )
                  else
                    ..._comments.map(
                      (comment) => CommentTile(
                        key: ValueKey(comment.id),
                        comment: comment,
                      ),
                    ),
                ],
              ),
            ),
          ),
          CommentInputBar(
            controller: _commentController,
            focusNode: _commentFocusNode,
            onSubmit: _submitComment,
          ),
        ],
      ),
    );
  }
}
