package br.com.woodriver.licauth.domain

import br.com.woodriver.licauth.domain.User.Role.USER
import jakarta.persistence.*

@Entity
@Table(name = "LicUser",
    uniqueConstraints = [
        UniqueConstraint(
            columnNames = ["email"]
        )
    ])
data class User(
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    val id: Long = 0L,
    val email: String = "",
    val password: String = "",
    val role: Role = USER
) {
    enum class Role {
        ADMIN, USER
    }
}