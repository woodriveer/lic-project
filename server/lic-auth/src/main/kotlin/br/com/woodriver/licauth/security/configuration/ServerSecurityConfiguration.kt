package br.com.woodriver.licauth.security.configuration

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.authentication.dao.DaoAuthenticationProvider
import org.springframework.security.config.annotation.method.configuration.EnableGlobalMethodSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer

@Configuration
@EnableWebSecurity
@EnableGlobalMethodSecurity(prePostEnabled = true, proxyTargetClass = true)
class ServerSecurityConfiguration(
    val userDetailsService: UserDetailsService
): WebMvcConfigurer {

    @Bean
    fun authenticationProvider(): DaoAuthenticationProvider {
        val provider = DaoAuthenticationProvider()
        provider.setPasswordEncoder(passwordEncoder())
        provider.setUserDetailsService(userDetailsService)
        return provider
    }

    @Bean
    fun passwordEncoder(): PasswordEncoder = BCryptPasswordEncoder()
}