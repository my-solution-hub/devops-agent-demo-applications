package com.catdemo.catprofile.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class CreateCatProfileRequest {

    @NotBlank(message = "must not be blank")
    private String name;

    @NotBlank(message = "must not be blank")
    private String ownerId;

    private String breed;

    @Positive(message = "must be positive")
    private Integer ageMonths;

    @NotNull(message = "must not be null")
    @Positive(message = "must be positive")
    private BigDecimal weightKg;

    private List<String> dietaryRestrictions;
}
