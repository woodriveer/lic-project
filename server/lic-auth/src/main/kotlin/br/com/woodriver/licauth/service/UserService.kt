package br.com.woodriver.licauth.service

import br.com.woodriver.licauth.repository.UserRepository
import org.springframework.security.core.authority.SimpleGrantedAuthority
import org.springframework.security.core.userdetails.User
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.stereotype.Service

@Service
class UserService(private val userRepository: UserRepository): UserDetailsService {
    override fun loadUserByUsername(email: String): UserDetails {
        val user = userRepository.findUserByEmail(email)
        val authority = SimpleGrantedAuthority(user.role.name)
        return User(user.email, user.password, arrayListOf(authority))
    }
}