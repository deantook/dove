import 'package:flutter/material.dart';

import '../data/mock_image_urls.dart';

class DiscoverUser {
  const DiscoverUser({
    required this.id,
    required this.name,
    required this.age,
    required this.city,
    required this.avatarColors,
    this.isOnline = false,
    this.photoUrl,
  });

  final String id;
  final String name;
  final int age;
  final String city;
  final List<Color> avatarColors;
  final bool isOnline;
  final String? photoUrl;

  List<String> get photos => MockImageUrls.discoverPhotosFor(id);

  String get resolvedPhotoUrl =>
      photoUrl ?? (photos.isNotEmpty ? photos.first : '');
}

const mockDiscoverUsers = [
  DiscoverUser(
    id: 'd1',
    name: 'Emma',
    age: 26,
    city: '上海',
    avatarColors: [Color(0xFFE8B4B8), Color(0xFFC9ADA7)],
    isOnline: true,
  ),
  DiscoverUser(
    id: 'd2',
    name: 'Sophia',
    age: 28,
    city: '上海',
    avatarColors: [Color(0xFFA8DADC), Color(0xFF457B9D)],
  ),
  DiscoverUser(
    id: 'd3',
    name: 'Olivia',
    age: 25,
    city: '杭州',
    avatarColors: [Color(0xFFD4A574), Color(0xFF8B6914)],
  ),
  DiscoverUser(
    id: 'd4',
    name: 'Ava',
    age: 27,
    city: '上海',
    avatarColors: [Color(0xFFE8D5B7), Color(0xFFB8956A)],
    isOnline: true,
  ),
  DiscoverUser(
    id: 'd5',
    name: 'Mia',
    age: 24,
    city: '苏州',
    avatarColors: [Color(0xFF95D5B2), Color(0xFF40916C)],
  ),
  DiscoverUser(
    id: 'd6',
    name: '林溪',
    age: 27,
    city: '上海',
    avatarColors: [Color(0xFFB8B8D1), Color(0xFF6B6B8D)],
    isOnline: true,
  ),
  DiscoverUser(
    id: 'd7',
    name: 'Zoe',
    age: 29,
    city: '南京',
    avatarColors: [Color(0xFFF4ACB7), Color(0xFF9D8189)],
  ),
  DiscoverUser(
    id: 'd8',
    name: '阿杰',
    age: 30,
    city: '上海',
    avatarColors: [Color(0xFF90A4AE), Color(0xFF546E7A)],
  ),
];
