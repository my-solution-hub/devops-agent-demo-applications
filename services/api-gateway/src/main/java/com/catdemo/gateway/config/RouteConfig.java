package com.catdemo.gateway.config;

import org.springframework.cloud.gateway.route.RouteLocator;
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Additional programmatic route configuration.
 * Primary routes are defined in application.yml.
 * This class can be extended for custom filters or dynamic routing.
 */
@Configuration
public class RouteConfig {

    // Routes are configured declaratively in application.yml via Spring Cloud Gateway.
    // Add programmatic route definitions here if needed for custom logic.
}
