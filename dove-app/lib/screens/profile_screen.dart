import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../models/app_user.dart';
import '../../providers/auth_provider.dart';
import '../../theme/app_colors.dart';
import 'profile/widgets/profile_collapsing_header.dart';
import 'profile/widgets/profile_info_tab.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({
    super.key,
    this.user = AppUser.currentUser,
    this.isEmbedded = false,
  });

  /// 底部 Tab 内嵌的「我的」页。
  final AppUser user;
  final bool isEmbedded;

  bool get isOwnProfile => isEmbedded && user.isCurrentUser;

  static Route<void> route({required AppUser user}) {
    return MaterialPageRoute<void>(
      fullscreenDialog: false,
      builder: (_) => ProfileScreen(user: user, isEmbedded: false),
    );
  }

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  static const _collapseStartThreshold = 48.0;
  static const _collapseDistance = 120.0;

  final _scrollController = ScrollController();

  double _scrollOffset = 0;

  double get _collapseProgress => _computeCollapseProgress(_scrollController);

  double _computeCollapseProgress(ScrollController controller) {
    if (!controller.hasClients) return 0;

    final position = controller.position;
    final requiredScroll = _collapseStartThreshold + _collapseDistance;

    if (position.maxScrollExtent < requiredScroll) return 0;
    if (_scrollOffset <= _collapseStartThreshold) return 0;

    return ((_scrollOffset - _collapseStartThreshold) / _collapseDistance)
        .clamp(0.0, 1.0);
  }

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController
      ..removeListener(_onScroll)
      ..dispose();
    super.dispose();
  }

  void _onScroll() {
    final offset = _scrollController.offset;
    if (offset == _scrollOffset) return;
    setState(() => _scrollOffset = offset);
  }

  void _showSettingsMenu(BuildContext context) {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        decoration: BoxDecoration(
          color: AppColors.cardBackground,
          borderRadius: const BorderRadius.vertical(
            top: Radius.circular(24),
          ),
        ),
        padding: const EdgeInsets.symmetric(vertical: 20),
        child: SafeArea(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: AppColors.divider,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 20),
              ListTile(
                leading: const Icon(
                  Icons.edit,
                  color: AppColors.textPrimary,
                ),
                title: const Text('编辑资料'),
                onTap: () {
                  Navigator.pop(context);
                },
              ),
              ListTile(
                leading: const Icon(
                  Icons.settings,
                  color: AppColors.textPrimary,
                ),
                title: const Text('设置'),
                onTap: () {
                  Navigator.pop(context);
                },
              ),
              const Divider(),
              ListTile(
                leading: const Icon(
                  Icons.logout,
                  color: Colors.red,
                ),
                title: const Text(
                  '退出登录',
                  style: TextStyle(color: Colors.red),
                ),
                onTap: () {
                  Navigator.pop(context);
                  _showLogoutConfirm(context);
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showLogoutConfirm(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('退出登录'),
        content: const Text('确定要退出登录吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              context.read<AuthProvider>().logout();
            },
            child: const Text(
              '退出',
              style: TextStyle(color: Colors.red),
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final canPop = Navigator.of(context).canPop();
    final isOwnProfile = widget.isOwnProfile;

    return PopScope(
      canPop: !widget.isEmbedded || canPop,
      child: Scaffold(
        backgroundColor: AppColors.background,
        body: SafeArea(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              ProfileCollapsingHeader(
                collapseProgress: _collapseProgress,
                userName: widget.user.name,
                isVerified: widget.user.isVerified,
                tags: widget.user.tags,
                showSettings: isOwnProfile,
                showBackButton: !widget.isEmbedded && canPop,
                onBackTap: canPop ? () => Navigator.of(context).pop() : null,
                onSettingsTap: isOwnProfile ? () => _showSettingsMenu(context) : null,
                avatarColors: widget.user.avatarColors,
              ),
              Expanded(
                child: ProfileInfoTab(
                  scrollController: _scrollController,
                  profileInfo: widget.user.profileInfo,
                  photos: widget.user.photos,
                  avatarColors: widget.user.avatarColors,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
