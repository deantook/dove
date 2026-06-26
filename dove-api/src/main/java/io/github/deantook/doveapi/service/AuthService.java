package io.github.deantook.doveapi.service;

import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import io.github.deantook.doveapi.dto.LoginRequest;
import io.github.deantook.doveapi.dto.LoginResponse;
import io.github.deantook.doveapi.entity.User;
import io.github.deantook.doveapi.exception.InvalidCredentialsException;
import io.github.deantook.doveapi.repository.UserRepository;

@Service
public class AuthService {

    private final UserRepository userRepository;
    private final SessionService sessionService;
    private final PasswordEncoder passwordEncoder;

    public AuthService(
            UserRepository userRepository,
            SessionService sessionService,
            PasswordEncoder passwordEncoder) {
        this.userRepository = userRepository;
        this.sessionService = sessionService;
        this.passwordEncoder = passwordEncoder;
    }

    @Transactional
    public LoginResponse login(LoginRequest request) {
        var username = request.username().trim();
        var password = request.password();

        return userRepository.findByUsername(username)
                .map(user -> authenticateExistingUser(user, password))
                .orElseGet(() -> registerAndLogin(username, password));
    }

    private LoginResponse authenticateExistingUser(User user, String password) {
        if (!passwordEncoder.matches(password, user.getPasswordHash())) {
            throw new InvalidCredentialsException();
        }

        var token = sessionService.createSession(user.getId());
        return new LoginResponse(token, user.getId(), user.getUsername(), false);
    }

    private LoginResponse registerAndLogin(String username, String password) {
        var user = userRepository.save(new User(username, passwordEncoder.encode(password)));
        var token = sessionService.createSession(user.getId());
        return new LoginResponse(token, user.getId(), user.getUsername(), true);
    }
}
