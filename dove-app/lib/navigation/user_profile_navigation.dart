import 'package:flutter/material.dart';

import '../models/app_user.dart';
import '../screens/profile_screen.dart';

abstract final class UserProfileNavigation {
  static void open(
    BuildContext context, {
    required String name,
    required List<Color> avatarColors,
  }) {
    final user = AppUser.fromName(name, avatarColors);
    Navigator.of(context, rootNavigator: true).push(
      ProfileScreen.route(user: user),
    );
  }
}
