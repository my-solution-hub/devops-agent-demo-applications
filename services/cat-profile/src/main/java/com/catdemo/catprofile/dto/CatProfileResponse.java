package com.catdemo.catprofile.dto;

import com.catdemo.catprofile.entity.CatProfile;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.List;
import java.util.UUID;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class CatProfileResponse {

    private UUID catId;
    private String ownerId;
    private String name;
    private String breed;
    private Integer ageMonths;
    private BigDecimal weightKg;
    private List<String> dietaryRestrictions;
    private Instant createdAt;
    private Instant updatedAt;

    public static CatProfileResponse fromEntity(CatProfile entity) {
        return CatProfileResponse.builder()
                .catId(entity.getCatId())
                .ownerId(entity.getOwnerId())
                .name(entity.getName())
                .breed(entity.getBreed())
                .ageMonths(entity.getAgeMonths())
                .weightKg(entity.getWeightKg())
                .dietaryRestrictions(entity.getDietaryRestrictions())
                .createdAt(entity.getCreatedAt())
                .updatedAt(entity.getUpdatedAt())
                .build();
    }
}
