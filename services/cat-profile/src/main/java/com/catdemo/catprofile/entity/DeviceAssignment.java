package com.catdemo.catprofile.entity;

import jakarta.persistence.*;
import lombok.*;

import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "device_assignments", uniqueConstraints = {
        @UniqueConstraint(columnNames = {"cat_id", "device_type"})
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class DeviceAssignment {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    @Column(name = "assignment_id", updatable = false, nullable = false)
    private UUID assignmentId;

    @Column(name = "cat_id", nullable = false)
    private UUID catId;

    @Column(name = "device_id", nullable = false)
    private UUID deviceId;

    @Column(name = "device_type", nullable = false, length = 50)
    private String deviceType;

    @Column(name = "assigned_at", nullable = false, updatable = false)
    private Instant assignedAt;

    @PrePersist
    protected void onCreate() {
        this.assignedAt = Instant.now();
    }
}
