import '../models/user_profile.dart';

/// Mock 数据用的真实头像与动态配图 URL。
abstract final class MockImageUrls {
  static const avatars = <String, String>{
    'Alex': 'https://randomuser.me/api/portraits/men/32.jpg',
    'Emma': 'https://randomuser.me/api/portraits/women/65.jpg',
    'Sophia': 'https://randomuser.me/api/portraits/women/44.jpg',
    'Olivia': 'https://randomuser.me/api/portraits/women/68.jpg',
    'Ava': 'https://randomuser.me/api/portraits/women/26.jpg',
    'Mia': 'https://randomuser.me/api/portraits/women/32.jpg',
    '林溪': 'https://randomuser.me/api/portraits/women/81.jpg',
    '陈朗': 'https://randomuser.me/api/portraits/men/22.jpg',
    'Zoe': 'https://randomuser.me/api/portraits/women/17.jpg',
    '阿杰': 'https://randomuser.me/api/portraits/men/45.jpg',
  };

  static const moments = <String, String>{
    // 独立书店 · 摄影集
    's1':
        'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=800&h=600&fit=crop',
    // sourdough 面包
    'f1':
        'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=800&h=600&fit=crop',
    // 徒步 sunset
    'f2':
        'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=800&h=600&fit=crop',
    // vinyl 唱片咖啡馆
    'u1':
        'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=800&h=600&fit=crop',
    // 晨跑日出
    'u3':
        'https://images.unsplash.com/photo-1476480862126-209bfaa8edc8?w=800&h=600&fit=crop',
    // sourdough 气孔
    'u4':
        'https://images.unsplash.com/photo-1548998547821-82d53e66d2ed?w=800&h=600&fit=crop',
    // 书架
    'u7':
        'https://images.unsplash.com/photo-152499990963-0bbf1c81f0a0?w=800&h=600&fit=crop',
  };

  static const discoverPhotos = <String, String>{
    'd1':
        'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=600&h=450&fit=crop',
    'd2':
        'https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?w=600&h=450&fit=crop',
    'd3':
        'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=600&h=450&fit=crop',
    'd4':
        'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=600&h=450&fit=crop',
    'd5':
        'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=600&h=450&fit=crop',
    'd6':
        'https://images.unsplash.com/photo-1521577352947-9bb58764b69a?w=600&h=450&fit=crop',
    'd7':
        'https://images.unsplash.com/photo-1525134479668-1bee5c7c6845?w=600&h=450&fit=crop',
    'd8':
        'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=600&h=450&fit=crop',
  };

  static const discoverPhotoGalleries = <String, List<String>>{
    'd1': [
      'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=600&h=450&fit=crop',
    ],
    'd2': [
      'https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=600&h=450&fit=crop',
    ],
    'd3': [
      'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1525134479668-1bee5c7c6845?w=600&h=450&fit=crop',
    ],
    'd4': [
      'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1521577352947-9bb58764b69a?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=600&h=450&fit=crop',
    ],
    'd5': [
      'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?w=600&h=450&fit=crop',
    ],
    'd6': [
      'https://images.unsplash.com/photo-1521577352947-9bb58764b69a?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=600&h=450&fit=crop',
    ],
    'd7': [
      'https://images.unsplash.com/photo-1525134479668-1bee5c7c6845?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=600&h=450&fit=crop',
    ],
    'd8': [
      'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=600&h=450&fit=crop',
      'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600&h=450&fit=crop',
    ],
  };

  static String? avatarFor(String name) => avatars[name];

  static String? momentFor(String momentId) => moments[momentId];

  static String? discoverPhotoFor(String userId) => discoverPhotos[userId];

  static List<String> discoverPhotosFor(String userId) {
    final gallery = discoverPhotoGalleries[userId];
    if (gallery != null && gallery.isNotEmpty) return gallery;

    final cover = discoverPhotos[userId];
    if (cover != null) return [cover];

    return const [];
  }

  static const profilePhotoGalleries = <String, List<String>>{
    'Alex': [
      'https://randomuser.me/api/portraits/men/32.jpg',
      'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-1476480862126-209bfaa8edc8?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-152499990963-0bbf1c81f0a0?w=600&h=800&fit=crop',
    ],
    'Emma': [
      'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=600&h=800&fit=crop',
    ],
    '林溪': [
      'https://images.unsplash.com/photo-1521577352947-9bb58764b69a?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=600&h=800&fit=crop',
      'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=600&h=800&fit=crop',
    ],
  };

  static List<String> profilePhotosFor(String name) {
    final List<String> photos;
    final gallery = profilePhotoGalleries[name];
    if (gallery != null && gallery.isNotEmpty) {
      photos = gallery;
    } else {
      final avatar = avatars[name];
      photos = avatar != null ? [avatar] : const [];
    }
    return photos.take(kMaxProfilePhotos).toList();
  }
}
