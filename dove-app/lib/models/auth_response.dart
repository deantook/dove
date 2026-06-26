class AuthResponse {
  final String token;
  final String userId;
  final String username;
  final bool created;

  AuthResponse({
    required this.token,
    required this.userId,
    required this.username,
    required this.created,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      token: json['token'] as String,
      userId: json['userId'] as String,
      username: json['username'] as String,
      created: json['created'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'token': token,
      'userId': userId,
      'username': username,
      'created': created,
    };
  }

  @override
  String toString() {
    return 'AuthResponse(token: $token, userId: $userId, username: $username, created: $created)';
  }
}

class LoginRequest {
  final String username;
  final String password;

  LoginRequest({
    required this.username,
    required this.password,
  });

  Map<String, dynamic> toJson() {
    return {
      'username': username,
      'password': password,
    };
  }
}
