package br.com.woodriver.licauth.controller.request

data class UserRequest(
    val email: String = "",
    val password: String = ""
)