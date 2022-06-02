package br.com.woodriver.licauth.repository

import br.com.woodriver.licauth.domain.User
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface UserRepository: JpaRepository<User, Long> {
    fun findUserByEmail(email: String): User
}