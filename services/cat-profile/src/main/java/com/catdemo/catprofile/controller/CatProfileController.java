package com.catdemo.catprofile.controller;

import com.catdemo.catprofile.dto.CatProfileResponse;
import com.catdemo.catprofile.dto.CreateCatProfileRequest;
import com.catdemo.catprofile.dto.UpdateCatProfileRequest;
import com.catdemo.catprofile.entity.CatProfile;
import com.catdemo.catprofile.service.CatProfileService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/cats")
@RequiredArgsConstructor
public class CatProfileController {

    private final CatProfileService catProfileService;

    @PostMapping
    public ResponseEntity<CatProfileResponse> createCatProfile(@Valid @RequestBody CreateCatProfileRequest request) {
        CatProfile created = catProfileService.createCatProfile(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(CatProfileResponse.fromEntity(created));
    }

    @GetMapping("/{id}")
    public ResponseEntity<CatProfileResponse> getCatProfile(@PathVariable UUID id) {
        CatProfile catProfile = catProfileService.getCatProfile(id);
        return ResponseEntity.ok(CatProfileResponse.fromEntity(catProfile));
    }

    @PutMapping("/{id}")
    public ResponseEntity<CatProfileResponse> updateCatProfile(
            @PathVariable UUID id,
            @Valid @RequestBody UpdateCatProfileRequest request) {
        CatProfile updated = catProfileService.updateCatProfile(id, request);
        return ResponseEntity.ok(CatProfileResponse.fromEntity(updated));
    }

    @GetMapping
    public ResponseEntity<List<CatProfileResponse>> listCatProfiles(
            @RequestParam(name = "owner_id", required = false) String ownerId) {
        List<CatProfile> profiles = catProfileService.listCatProfiles(ownerId);
        List<CatProfileResponse> response = profiles.stream()
                .map(CatProfileResponse::fromEntity)
                .toList();
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}/health")
    public ResponseEntity<Map<String, Object>> getCatHealth(@PathVariable UUID id) {
        // Verify the cat exists (throws 404 if not found)
        catProfileService.getCatProfile(id);
        // Placeholder — will be wired to Health Monitor Service later
        return ResponseEntity.ok(Collections.emptyMap());
    }
}
