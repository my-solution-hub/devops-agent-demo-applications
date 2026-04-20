package com.catdemo.catprofile.repository;

import com.catdemo.catprofile.entity.DeviceAssignment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface DeviceAssignmentRepository extends JpaRepository<DeviceAssignment, UUID> {

    List<DeviceAssignment> findByCatId(UUID catId);
}
