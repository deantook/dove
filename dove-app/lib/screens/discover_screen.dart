import 'package:flutter/material.dart';

import '../models/discover_user.dart';
import '../theme/app_colors.dart';
import '../theme/app_typography.dart';
import '../live_match/live_match_scope.dart';
import 'discover/discover_filter_screen.dart';
import 'discover/widgets/discover_photo_overlay.dart';
import 'discover/widgets/discover_photo_utils.dart';
import 'discover/widgets/discover_waterfall.dart';
import 'discover/widgets/live_match_card.dart';

class DiscoverScreen extends StatefulWidget {
  const DiscoverScreen({super.key});

  @override
  State<DiscoverScreen> createState() => _DiscoverScreenState();
}

class _DiscoverScreenState extends State<DiscoverScreen> {
  late List<DiscoverUser> _users;
  DiscoverUser? _previewUser;
  Rect? _previewSourceRect;
  String? _removingUserId;
  final _likedUserIds = <String>{};
  final _photoKeys = <String, GlobalKey>{};

  @override
  void initState() {
    super.initState();
    _users = List.of(mockDiscoverUsers);
  }

  GlobalKey _photoKeyFor(String userId) {
    return _photoKeys.putIfAbsent(userId, GlobalKey.new);
  }

  Rect? _resolvePreviewSourceRect() {
    final user = _previewUser;
    if (user == null) return _previewSourceRect;
    return globalRectForKey(_photoKeyFor(user.id)) ?? _previewSourceRect;
  }

  void _openPreview(DiscoverUser user, Rect sourceRect) {
    setState(() {
      _previewUser = user;
      _previewSourceRect = sourceRect;
    });
  }

  void _closePreview() {
    setState(() {
      _previewUser = null;
      _previewSourceRect = null;
    });
  }

  void _removeUser(String userId) {
    setState(() {
      _users.removeWhere((user) => user.id == userId);
      _removingUserId = null;
    });
  }

  void _likeUser(String userId) {
    setState(() => _likedUserIds.add(userId));
  }

  void _dislikeUser(String userId) {
    setState(() => _removingUserId = userId);
  }

  @override
  Widget build(BuildContext context) {
    final previewUser = _previewUser;
    final previewSourceRect = _previewSourceRect;

    return Stack(
      fit: StackFit.expand,
      children: [
        SafeArea(
          child: CustomScrollView(
            slivers: [
              SliverPadding(
                padding: const EdgeInsets.fromLTRB(20, 16, 20, 0),
                sliver: SliverToBoxAdapter(
                  child: ListenableBuilder(
                    listenable: LiveMatchScope.of(context),
                    builder: (context, _) {
                      final match = LiveMatchScope.of(context);
                      return LiveMatchCard(
                        isMatching: match.isActive,
                        onStartTap: () {
                          if (match.isActive) {
                            match.expand();
                          } else {
                            match.start();
                          }
                        },
                      );
                    },
                  ),
                ),
              ),
              SliverPadding(
                padding: const EdgeInsets.fromLTRB(20, 24, 20, 0),
                sliver: SliverToBoxAdapter(
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Text('你可能感兴趣的用户', style: AppTypography.cardHeadline),
                      const Spacer(),
                      _FilterIconButton(
                        onTap: () {
                          Navigator.of(
                            context,
                          ).push(DiscoverFilterScreen.route());
                        },
                      ),
                    ],
                  ),
                ),
              ),
              SliverPadding(
                padding: const EdgeInsets.fromLTRB(20, 16, 20, 0),
                sliver: DiscoverWaterfall(
                  users: _users,
                  previewUserId: previewUser?.id,
                  removingUserId: _removingUserId,
                  photoKeyFor: _photoKeyFor,
                  onUserPhotoTap: _openPreview,
                  onUserRemoved: _removeUser,
                ),
              ),
            ],
          ),
        ),
        if (previewUser != null && previewSourceRect != null)
          DiscoverPhotoOverlay(
            key: ValueKey(previewUser.id),
            user: previewUser,
            sourceRect: previewSourceRect,
            resolveSourceRect: () => _resolvePreviewSourceRect(),
            onClose: _closePreview,
            onDislike: () => _dislikeUser(previewUser.id),
            onLike: () => _likeUser(previewUser.id),
          ),
      ],
    );
  }
}

class _FilterIconButton extends StatelessWidget {
  const _FilterIconButton({required this.onTap});

  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(20),
        child: const Padding(
          padding: EdgeInsets.all(8),
          child: Icon(
            Icons.tune_rounded,
            size: 24,
            color: AppColors.textPrimary,
          ),
        ),
      ),
    );
  }
}
