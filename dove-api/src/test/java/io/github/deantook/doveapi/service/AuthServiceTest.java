package io.github.deantook.doveapi.service;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.util.Optional;
import java.util.UUID;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.password.PasswordEncoder;

import io.github.deantook.doveapi.dto.LoginRequest;
import io.github.deantook.doveapi.entity.User;
import io.github.deantook.doveapi.exception.InvalidCredentialsException;
import io.github.deantook.doveapi.repository.UserRepository;

@ExtendWith(MockitoExtension.class)
class AuthServiceTest {

    @Mock
    private UserRepository userRepository;

    @Mock
    private SessionService sessionService;

    @Mock
    private PasswordEncoder passwordEncoder;

    @InjectMocks
    private AuthService authService;

    @Test
    void loginCreatesUserWhenMissing() {
        var userId = UUID.randomUUID();
        var savedUser = org.mockito.Mockito.mock(User.class);
        when(savedUser.getId()).thenReturn(userId);
        when(savedUser.getUsername()).thenReturn("new-user");

        when(userRepository.findByUsername("new-user")).thenReturn(Optional.empty());
        when(passwordEncoder.encode("123")).thenReturn("hashed");
        when(userRepository.save(any(User.class))).thenReturn(savedUser);
        when(sessionService.createSession(userId)).thenReturn("token-1");

        var response = authService.login(new LoginRequest("new-user", "123"));

        assertThat(response.created()).isTrue();
        assertThat(response.token()).isEqualTo("token-1");
        assertThat(response.userId()).isEqualTo(userId);
        verify(userRepository).save(any(User.class));
    }

    @Test
    void loginRejectsWrongPasswordForExistingUser() {
        var user = org.mockito.Mockito.mock(User.class);
        when(user.getPasswordHash()).thenReturn("hashed");
        when(userRepository.findByUsername("existing")).thenReturn(Optional.of(user));
        when(passwordEncoder.matches("wrong", "hashed")).thenReturn(false);

        assertThatThrownBy(() -> authService.login(new LoginRequest("existing", "wrong")))
                .isInstanceOf(InvalidCredentialsException.class);

        verify(sessionService, never()).createSession(any());
    }

    @Test
    void loginReturnsSessionForExistingUser() {
        var userId = UUID.randomUUID();
        var user = org.mockito.Mockito.mock(User.class);
        when(user.getId()).thenReturn(userId);
        when(user.getUsername()).thenReturn("existing");
        when(user.getPasswordHash()).thenReturn("hashed");
        when(userRepository.findByUsername("existing")).thenReturn(Optional.of(user));
        when(passwordEncoder.matches("right", "hashed")).thenReturn(true);
        when(sessionService.createSession(userId)).thenReturn("token-2");

        var response = authService.login(new LoginRequest("existing", "right"));

        assertThat(response.created()).isFalse();
        assertThat(response.token()).isEqualTo("token-2");
        verify(userRepository, never()).save(any());
    }

    @Test
    void loginTrimsUsername() {
        var savedUser = org.mockito.Mockito.mock(User.class);
        when(savedUser.getId()).thenReturn(UUID.randomUUID());
        when(savedUser.getUsername()).thenReturn("trimmed");

        when(userRepository.findByUsername("trimmed")).thenReturn(Optional.empty());
        when(passwordEncoder.encode("pw")).thenReturn("hashed");
        when(userRepository.save(any(User.class))).thenReturn(savedUser);
        when(sessionService.createSession(any())).thenReturn("token-3");

        authService.login(new LoginRequest("  trimmed  ", "pw"));

        ArgumentCaptor<User> captor = ArgumentCaptor.forClass(User.class);
        verify(userRepository).save(captor.capture());
        assertThat(captor.getValue().getUsername()).isEqualTo("trimmed");
    }
}
