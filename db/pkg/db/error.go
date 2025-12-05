/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */


package db

import (
	"errors"
)

var (
	// ErrDoesNotExist is raised a DB query fails to find the requested entity
	ErrDoesNotExist = errors.New("the requested entity does not exist")
	// ErrDBError is a generalized error to expose to the user when unexpected errors occur when communicating with DB
	ErrDBError = errors.New("error communicating with data store")
	// ErrInvalidValue is raised when a value to be stored in DB is invalid
	ErrInvalidValue = errors.New("provided value is invalid")
	// ErrInvalidParams is raised when a function is called with invalid set of parameters
	ErrInvalidParams = errors.New("provided params are invalid or conflicting")

	// ErrXactAdvisoryLockFailed indicates that the transaction advisory lock could not be taken
	ErrXactAdvisoryLockFailed = errors.New("unable to take transaction advisory lock")
	// ErrSessionAdvisoryLockFailed indicates that the session advisory lock could not be taken
	ErrSessionAdvisoryLockFailed = errors.New("unable to take session advisory lock")
	// ErrSessionAdvisoryLockUnlockFailed indicates that the session advisory lock could not be released.
	ErrSessionAdvisoryLockUnlockFailed = errors.New("unable to release session advisory lock or lock was not held by this session")
)
