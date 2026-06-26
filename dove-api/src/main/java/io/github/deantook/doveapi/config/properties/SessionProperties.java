package io.github.deantook.doveapi.config.properties;

import java.time.Duration;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "dove.session")
public record SessionProperties(Duration ttl) {
}
