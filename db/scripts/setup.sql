-- SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
-- SPDX-License-Identifier: LicenseRef-NvidiaProprietary
--
-- NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
-- property and proprietary rights in and to this material, related
-- documentation and any modifications thereto. Any use, reproduction,
-- disclosure or distribution of this material and related documentation
-- without an express license agreement from NVIDIA CORPORATION or
-- its affiliates is strictly prohibited.

-- Connect as postgres user to template1 DB to run this script
-- Example: PGPASSWORD=postgres psql -U postgres -p 30432 -d template1 < scripts/setup.sql
-- Create extension for pg_trgm
CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- Create Forge DB and user
CREATE DATABASE forge WITH ENCODING 'UTF8';
-- Password should be changed before use in environments deployed in Cloud
CREATE USER forge WITH PASSWORD 'forge';
-- Grant all privileges on Forge DB to Forge user
GRANT ALL PRIVILEGES ON DATABASE forge TO forge;
