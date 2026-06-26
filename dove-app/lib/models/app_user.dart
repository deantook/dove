import 'package:flutter/material.dart';

import '../data/mock_image_urls.dart';
import 'moment.dart';
import 'user_profile.dart';

class AppUser {
  const AppUser({
    required this.name,
    required this.avatarColors,
    required this.profileInfo,
    required this.photos,
    required this.tags,
    this.moments = const [],
    this.isVerified = true,
    this.isCurrentUser = false,
  });

  final String name;
  final List<Color> avatarColors;
  final ProfileInfo profileInfo;
  final List<String> photos;
  final List<String> tags;
  final List<Moment> moments;
  final bool isVerified;
  final bool isCurrentUser;

  static const currentUser = AppUser(
    name: 'Alex',
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    profileInfo: mockProfileInfo,
    photos: [
      'https://randomuser.me/api/portraits/men/32.jpg',
      'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-1476480862126-209bfaa8edc8?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-152499990963-0bbf1c81f0a0?w=600&h=800&fit=crop',
    ],
    tags: ['产品', '徒步', '爵士乐', '咖啡', '黑胶'],
    moments: mockUserMoments,
    isVerified: true,
    isCurrentUser: true,
  );

  factory AppUser.fromName(String name, List<Color> avatarColors) {
    if (name == currentUser.name) return currentUser;
    return AppUser(
      name: name,
      avatarColors: avatarColors,
      profileInfo: _profileInfoFor(name),
      photos: MockImageUrls.profilePhotosFor(name),
      tags: _tagsFor(name),
      moments: _momentsFor(name),
      isVerified: name != '陈朗',
    );
  }
}

ProfileInfo _profileInfoFor(String name) {
  const known = <String, ProfileInfo>{
    'Emma': ProfileInfo(
      prompts: [
        ProfilePrompt(
          title: '关于我',
          content: '烘焙爱好者，周末常泡在厨房。相信美食是最好的开场白。',
        ),
        ProfilePrompt(
          title: '我在找',
          content: '能一起分享生活细节、不回避深度对话的人。',
        ),
      ],
      details: [
        ('年龄', '26'),
        ('城市', '上海'),
        ('职业', '品牌策划'),
      ],
    ),
    '林溪': ProfileInfo(
      prompts: [
        ProfilePrompt(
          title: '关于我',
          content: '独立书店常客，喜欢胶片摄影和 city walk。',
        ),
        ProfilePrompt(
          title: '我在找',
          content: '愿意一起探索城市角落的人。',
        ),
      ],
      details: [
        ('年龄', '27'),
        ('城市', '上海'),
        ('职业', '插画师'),
      ],
    ),
  };

  return known[name] ??
      ProfileInfo(
        prompts: [
          ProfilePrompt(
            title: '关于我',
            content: '正在完善个人资料，期待与你认识。',
          ),
        ],
        details: [
          ('城市', '上海'),
        ],
      );
}

List<String> _tagsFor(String name) {
  const known = <String, List<String>>{
    'Emma': ['烘焙', '美食', '深度对话'],
    '林溪': ['胶片', '书店', 'City Walk'],
  };

  return known[name] ?? const ['上海'];
}

List<Moment> _momentsFor(String name) {
  return [...mockMoments, ...mockUserMoments]
      .where((moment) => moment.authorName == name)
      .toList();
}

const mockUserMoments = [
  Moment(
    id: 'u1',
    authorName: 'Alex',
    content: '周末探店：找到一家很棒的 vinyl 唱片咖啡馆，老板推荐了一张 Coltrane。',
    timeLabel: '3 天前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    hasImage: true,
    likeCount: 15,
  ),
  Moment(
    id: 'u2',
    authorName: 'Alex',
    content: '最近在读《挪威的森林》，有人一起讨论吗？',
    timeLabel: '1 周前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    likeCount: 8,
  ),
  Moment(
    id: 'u3',
    authorName: 'Alex',
    content: '晨跑 8 公里，黄浦边的日出永远看不腻。',
    timeLabel: '1 周前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    hasImage: true,
    likeCount: 22,
  ),
  Moment(
    id: 'u4',
    authorName: 'Alex',
    content: '试做了 sourdough，第三次终于有像样的气孔了 🍞',
    timeLabel: '2 周前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    hasImage: true,
    likeCount: 19,
  ),
  Moment(
    id: 'u5',
    authorName: 'Alex',
    content: '安福路新开的摄影展不错，黑白街拍部分尤其喜欢。',
    timeLabel: '2 周前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    likeCount: 11,
  ),
  Moment(
    id: 'u6',
    authorName: 'Alex',
    content: '下班路上听到街头萨克斯，停下来听了三首，差点错过地铁。',
    timeLabel: '3 周前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    likeCount: 7,
  ),
  Moment(
    id: 'u7',
    authorName: 'Alex',
    content: '整理书架，发现大学时代买的第一本设计书，封面都泛黄了。',
    timeLabel: '1 个月前',
    source: MomentSource.friend,
    avatarColors: [Color(0xFFD4C4B0), Color(0xFF5C7A6B)],
    hasImage: true,
    likeCount: 14,
  ),
];
