package com.catdemo.catprofile;

import com.catdemo.catprofile.entity.CatProfile;
import com.catdemo.catprofile.entity.DeviceAssignment;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Basic tests to verify entity classes and application structure.
 * Full integration tests require PostgreSQL (or Testcontainers).
 */
class CatProfileApplicationTests {

    @Test
    void catProfileEntityCanBeCreated() {
        CatProfile profile = CatProfile.builder()
                .catId(UUID.randomUUID())
                .ownerId("owner-123")
                .name("Whiskers")
                .breed("Persian")
                .ageMonths(24)
                .weightKg(new BigDecimal("4.50"))
                .dietaryRestrictions(List.of("grain-free", "no-fish"))
                .build();

        assertEquals("Whiskers", profile.getName());
        assertEquals("owner-123", profile.getOwnerId());
        assertEquals("Persian", profile.getBreed());
        assertEquals(24, profile.getAgeMonths());
        assertEquals(new BigDecimal("4.50"), profile.getWeightKg());
        assertEquals(2, profile.getDietaryRestrictions().size());
        assertTrue(profile.getDietaryRestrictions().contains("grain-free"));
    }

    @Test
    void deviceAssignmentEntityCanBeCreated() {
        UUID catId = UUID.randomUUID();
        UUID deviceId = UUID.randomUUID();

        DeviceAssignment assignment = DeviceAssignment.builder()
                .assignmentId(UUID.randomUUID())
                .catId(catId)
                .deviceId(deviceId)
                .deviceType("feeder")
                .build();

        assertEquals(catId, assignment.getCatId());
        assertEquals(deviceId, assignment.getDeviceId());
        assertEquals("feeder", assignment.getDeviceType());
    }

    @Test
    void catProfileBuilderHandlesNullableFields() {
        CatProfile profile = CatProfile.builder()
                .ownerId("owner-456")
                .name("Luna")
                .weightKg(new BigDecimal("3.20"))
                .build();

        assertEquals("Luna", profile.getName());
        assertNull(profile.getBreed());
        assertNull(profile.getAgeMonths());
        assertNull(profile.getDietaryRestrictions());
    }
}
