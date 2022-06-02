package br.com.woodriver.licauth.security.properties

import org.springframework.boot.context.properties.ConfigurationProperties

@ConfigurationProperties(value = "jwt")
class JWTProperties(
    val clientId: String,
    val clientSecret: String,
    val accessTokenValiditySeconds: Int,
    val refreshTokenValiditySeconds: Int
)