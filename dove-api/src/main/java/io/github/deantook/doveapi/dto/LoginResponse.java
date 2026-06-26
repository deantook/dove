package io.github.deantook.doveapi.dto;

import java.util.UUID;

public record LoginResponse(
        String token,
        UUID userId,
        String username,
        boolean created) {
}
