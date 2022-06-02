package br.com.woodriver.licauth.usecase

import br.com.woodriver.licauth.domain.User
import br.com.woodriver.licauth.repository.UserRepository
import org.springframework.stereotype.Component

@Component
class CreateUserUseCase(private val userRepository: UserRepository) {

    fun execute(userEntity: User): User {
        return userRepository.save(userEntity)
    }
}