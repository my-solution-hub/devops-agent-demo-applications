package com.catdemo.catprofile.entity;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

/**
 * JPA AttributeConverter that converts between List&lt;String&gt; and a comma-separated
 * String for database storage. This approach is compatible with both PostgreSQL TEXT[]
 * (via native queries) and H2 (for testing).
 *
 * For PostgreSQL production use, the column is defined as TEXT[] and Hibernate's native
 * array support handles the mapping via @JdbcTypeCode(SqlTypes.ARRAY).
 * This converter serves as a fallback documentation of the conversion logic.
 */
@Converter
public class StringListConverter implements AttributeConverter<List<String>, String[]> {

    @Override
    public String[] convertToDatabaseColumn(List<String> attribute) {
        if (attribute == null || attribute.isEmpty()) {
            return new String[0];
        }
        return attribute.toArray(new String[0]);
    }

    @Override
    public List<String> convertToEntityAttribute(String[] dbData) {
        if (dbData == null || dbData.length == 0) {
            return new ArrayList<>();
        }
        return new ArrayList<>(Arrays.asList(dbData));
    }
}
