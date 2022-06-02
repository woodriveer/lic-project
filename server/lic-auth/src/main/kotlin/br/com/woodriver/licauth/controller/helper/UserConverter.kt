package br.com.woodriver.licauth.controller.helper

import br.com.woodriver.licauth.controller.request.UserRequest
import br.com.woodriver.licauth.domain.User
import br.com.woodriver.licauth.domain.User.Role.USER


fun UserRequest.toDomain() = User(
    email = email,
    password = password,
    role = USER
)