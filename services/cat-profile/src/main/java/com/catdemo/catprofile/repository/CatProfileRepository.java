package com.catdemo.catprofile.repository;

import com.catdemo.catprofile.entity.CatProfile;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface CatProfileRepository extends JpaRepository<CatProfile, UUID> {

    List<CatProfile> findByOwnerId(String ownerId);

    List<CatProfile> findByNameIgnoreCase(String name);
}
