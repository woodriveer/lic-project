package br.com.woodriver.licauth.security.configuration

import br.com.woodriver.licauth.security.properties.JWTProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Import
import org.springframework.jdbc.core.JdbcOperations
import org.springframework.security.config.annotation.web.configuration.OAuth2AuthorizationServerConfiguration
import org.springframework.security.oauth2.core.AuthorizationGrantType
import org.springframework.security.oauth2.core.ClientAuthenticationMethod.CLIENT_SECRET_BASIC
import org.springframework.security.oauth2.server.authorization.client.JdbcRegisteredClientRepository
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient
import java.util.UUID.randomUUID

@Configuration
@Import(OAuth2AuthorizationServerConfiguration::class)
class AuthorizationServerConfiguration(
    val jwtProperties: JWTProperties
) {

    @Bean
    fun registeredClientRepository(): JdbcRegisteredClientRepository {
        val registeredClient = RegisteredClient.withId(randomUUID().toString())
            .clientId(jwtProperties.clientId)
            .clientId(jwtProperties.clientSecret)
            .clientAuthenticationMethod(CLIENT_SECRET_BASIC)
            .authorizationGrantType(AuthorizationGrantType.AUTHORIZATION_CODE)
            .authorizationGrantType(AuthorizationGrantType.CLIENT_CREDENTIALS)
            .authorizationGrantType(AuthorizationGrantType.REFRESH_TOKEN)
            .build()

        val jdbc = JdbcRegisteredClientRepository(
            JdbcOperations()
        )
        jdbc.save(registeredClient)
        return jdbc
    }
}