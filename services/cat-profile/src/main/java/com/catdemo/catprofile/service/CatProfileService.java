package com.catdemo.catprofile.service;

import com.catdemo.catprofile.dto.CreateCatProfileRequest;
import com.catdemo.catprofile.dto.UpdateCatProfileRequest;
import com.catdemo.catprofile.entity.CatProfile;
import com.catdemo.catprofile.exception.ResourceNotFoundException;
import com.catdemo.catprofile.repository.CatProfileRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class CatProfileService {

    private final CatProfileRepository catProfileRepository;

    public CatProfile createCatProfile(CreateCatProfileRequest request) {
        CatProfile catProfile = CatProfile.builder()
                .name(request.getName())
                .ownerId(request.getOwnerId())
                .breed(request.getBreed())
                .ageMonths(request.getAgeMonths())
                .weightKg(request.getWeightKg())
                .dietaryRestrictions(request.getDietaryRestrictions())
                .build();
        return catProfileRepository.save(catProfile);
    }

    public CatProfile getCatProfile(UUID id) {
        return catProfileRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Cat profile not found with id: " + id));
    }

    public CatProfile updateCatProfile(UUID id, UpdateCatProfileRequest request) {
        CatProfile existing = catProfileRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Cat profile not found with id: " + id));

        existing.setName(request.getName());
        existing.setBreed(request.getBreed());
        existing.setAgeMonths(request.getAgeMonths());
        existing.setWeightKg(request.getWeightKg());
        existing.setDietaryRestrictions(request.getDietaryRestrictions());

        return catProfileRepository.save(existing);
    }

    public List<CatProfile> listCatProfiles(String ownerId) {
        if (ownerId != null && !ownerId.isBlank()) {
            return catProfileRepository.findByOwnerId(ownerId);
        }
        return catProfileRepository.findAll();
    }
}
