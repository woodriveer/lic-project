package br.com.woodriver.licauth

import br.com.woodriver.licauth.security.properties.JWTProperties
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.runApplication

@SpringBootApplication
@EnableConfigurationProperties(JWTProperties::class)
class LicAuthApplication

fun main(args: Array<String>) {
    runApplication<LicAuthApplication>(*args)
}
