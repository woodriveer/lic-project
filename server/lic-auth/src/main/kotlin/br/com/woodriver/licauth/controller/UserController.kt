package br.com.woodriver.licauth.controller

import br.com.woodriver.licauth.controller.helper.toDomain
import br.com.woodriver.licauth.controller.request.UserRequest
import br.com.woodriver.licauth.domain.User
import br.com.woodriver.licauth.usecase.CreateUserUseCase
import br.com.woodriver.licauth.utils.logger
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RestController

@RestController
class UserController(val createUserUseCase: CreateUserUseCase) {

    val logger = logger<UserController>()

    @PostMapping(value = ["/signup"])
    fun signup(@RequestBody body: UserRequest) : User {
        logger.info("Starting to create [User={}]", body.email)
        return createUserUseCase.execute(body.toDomain())
            .apply {
                logger.info("Done to create user with [ID={}]", this.id)
            }
    }
}
