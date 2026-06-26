package io.github.deantook.doveapi.service;

import java.time.Duration;
import java.util.UUID;

import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;

import io.github.deantook.doveapi.config.properties.SessionProperties;

@Service
public class SessionService {

    private static final String SESSION_KEY_PREFIX = "session:";

    private final StringRedisTemplate redisTemplate;
    private final Duration sessionTtl;

    public SessionService(StringRedisTemplate redisTemplate, SessionProperties sessionProperties) {
        this.redisTemplate = redisTemplate;
        this.sessionTtl = sessionProperties.ttl();
    }

    public String createSession(UUID userId) {
        var token = UUID.randomUUID().toString();
        redisTemplate.opsForValue().set(sessionKey(token), userId.toString(), sessionTtl);
        return token;
    }

    public void revokeSession(String token) {
        redisTemplate.delete(sessionKey(token));
    }

    private static String sessionKey(String token) {
        return SESSION_KEY_PREFIX + token;
    }
}
